package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"orderservice/config"
	"orderservice/internal/broker"
	"orderservice/internal/infrastructure/cache"
	"orderservice/internal/infrastructure/repo"
	"orderservice/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := repo.NewRepo(db)
	c := cache.NewCacheWithTTL(cfg.CacheTTL)
	u := usecase.NewUsecase(r, c)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBrokers},
		Topic:   cfg.KafkaTopic,
		GroupID: cfg.KafkaGroupID,
	})
	consumer := broker.NewConsumer(reader)

	if err := u.Start(ctx); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := consumer.Consume(ctx, u.HandleKafkaMessage); err != nil {
			log.Printf("main: consumer stopped: %v\n", err)
		}
	}()

	rtr := chi.NewRouter()
	rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./client/client.html")
		if err != nil {
			http.Error(w, "failed to read client.html", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(data); err != nil {
			log.Printf("http write error: %v", err)
		}
	})

	rtr.Get("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "id")

		order, err := u.GetOrder(r.Context(), orderID)
		if err != nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: rtr,
	}

	go func() {
		log.Println("HTTP server started at ", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("main: shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown error: %v", err)
	}

	if err := consumer.Close(); err != nil {
		log.Printf("consumer close error: %v", err)
	}

	c.Close()

	if err := u.Shutdown(ctx); err != nil {
		log.Printf("usecase shutdown error: %v", err)
	}
}
