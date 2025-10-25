package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"orderservice/config"
	ctrlhttp "orderservice/internal/controller/http"
	"orderservice/internal/di"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	container, err := di.New(logger, ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer container.DB.Close()

	go func() {
		if err := container.Kafka.Start(ctx); err != nil {
			log.Printf("main: kafka controller stopped: %v\n", err)
		}
	}()

	server := ctrlhttp.NewServer(logger, container.Router, cfg.HTTPPort)
	server.Start()

	<-ctx.Done()
	log.Println("main: shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown error: %v", err)
	}

	if err := container.Kafka.Stop(ctx); err != nil {
		log.Printf("kafka stop error: %v", err)
	}

	container.Cache.Close()
}
