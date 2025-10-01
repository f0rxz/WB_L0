package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"orderservice/internal/broker"
	"orderservice/internal/config"
	"orderservice/internal/infrastructure/cache"
	"orderservice/internal/infrastructure/repo"
	"orderservice/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

func main() {
	cfg, err := config.LoadConfig()
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
	c := cache.NewCache()
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

		w.Write(data)
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
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("main: shutting down")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("http server shutdown error: %v", err)
	}

	consumer.Close()
	u.Shutdown(ctx)
}
