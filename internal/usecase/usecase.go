package usecase

import (
	"context"
	"fmt"
	"orderservice/internal/infrastructure/cache"
	"orderservice/internal/infrastructure/repo"
	"orderservice/internal/model"
)

type OrderUsecase interface {
	GetOrder(ctx context.Context, orderUID string) (*model.Order, error)
	CreateOrder(ctx context.Context, ord *model.Order) error
}

type orderUsecase struct {
	repo  repo.Repo
	cache cache.Cache
}

func NewOrderUsecase(r repo.Repo, c cache.Cache) OrderUsecase {
	return &orderUsecase{repo: r, cache: c}
}

func (u *orderUsecase) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
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

func (u *orderUsecase) CreateOrder(ctx context.Context, ord *model.Order) error {
	if err := ord.Validate(); err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	if _, err := u.repo.CreateOrder(ctx, ord); err != nil {
		return err
	}

	u.cache.Set(ord)
	return nil
}
