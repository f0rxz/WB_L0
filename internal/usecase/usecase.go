package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"orderservice/internal/infrastructure/cache"
	"orderservice/internal/infrastructure/repo"
	"orderservice/internal/model"
)

type Usecase interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	GetOrder(ctx context.Context, orderUID string) (*model.Order, error)
	HandleKafkaMessage(ctx context.Context, key, value []byte) error
}

type usecase struct {
	repo   repo.Repo
	cache  cache.Cache
	cancel context.CancelFunc
}

func NewUsecase(r repo.Repo, c cache.Cache) Usecase {
	return &usecase{
		repo:  r,
		cache: c,
	}
}

func (u *usecase) Start(ctx context.Context) error {
	loadCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	orders, err := u.repo.GetAllOrders(loadCtx)
	if err != nil {
		return fmt.Errorf("usecase: failed to load orders for cache: %w", err)
	}
	u.cache.SetupCache(orders)
	return nil
}

func (u *usecase) Shutdown(ctx context.Context) error {
	if u.cancel != nil {
		u.cancel()
	}
	return nil
}

func (u *usecase) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	if o, ok := u.cache.Get(orderUID); ok {
		return o, nil
	}

	o, err := u.repo.GetOrderByID(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	u.cache.Set(o)
	return o, nil
}

func (u *usecase) HandleKafkaMessage(ctx context.Context, key, value []byte) error {
	var ord model.Order
	if err := json.Unmarshal(value, &ord); err != nil {
		return err
	}

	wctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := u.repo.CreateOrder(wctx, &ord)
	if err != nil {
		return err
	}

	u.cache.Set(&ord)
	log.Printf("usecase: order %s processed from kafka\n", ord.OrderUID)
	return nil
}
