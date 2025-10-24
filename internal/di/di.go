package di

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"orderservice/config"
	ctrlhttp "orderservice/internal/controller/http"
	ctrlkafka "orderservice/internal/controller/kafkacontroller"
	"orderservice/internal/infrastructure/cache"
	"orderservice/internal/infrastructure/repo"
	"orderservice/internal/usecase"
	"orderservice/pkg/connectors"
	"orderservice/pkg/consumer"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type Container struct {
	Config   *config.Config
	DB       *pgxpool.Pool
	Repo     repo.Repo
	Cache    cache.Cache
	Usecase  usecase.OrderUsecase
	Consumer *consumer.Consumer
	Kafka    *ctrlkafka.KafkaController
	Router   http.Handler
}

func New(ctx context.Context, cfg *config.Config) (*Container, error) {
	db, err := connectors.ConnectPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("app: connect postgres: %w", err)
	}

	r := repo.NewRepo(db)
	c := cache.NewCacheWithTTL(cfg.CacheTTL)
	u := usecase.NewOrderUsecase(r, c)

	loadCtx, loadCancel := context.WithTimeout(ctx, 30*time.Second)
	defer loadCancel()
	orders, err := r.GetAllOrders(loadCtx)
	if err != nil {
		log.Printf("di: failed to preload cache, starting with empty cache: %v", err)
	} else {
		c.SetupCache(orders)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBrokers},
		Topic:   cfg.KafkaTopic,
		GroupID: cfg.KafkaGroupID,
	})
	cons := consumer.NewConsumer(reader)

	router := ctrlhttp.NewRouter(u)

	kctrl := ctrlkafka.NewKafkaController(u, cons)

	return &Container{
		Config:   cfg,
		DB:       db,
		Repo:     r,
		Cache:    c,
		Usecase:  u,
		Consumer: cons,
		Kafka:    kctrl,
		Router:   router,
	}, nil
}
