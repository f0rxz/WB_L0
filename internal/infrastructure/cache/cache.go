package cache

import (
	"orderservice/internal/model"
	"sync"
)

type Cache interface {
	Set(order *model.Order)
	Get(orderUID string) (*model.Order, bool)
	SetupCache(orders []*model.Order)
}

type cache struct {
	data sync.Map
}

func NewCache() Cache {
	return &cache{
		data: sync.Map{},
	}
}

func (c *cache) Set(order *model.Order) {
	c.data.Store(order.OrderUID, order)
}

func (c *cache) Get(orderUID string) (*model.Order, bool) {
	if v, ok := c.data.Load(orderUID); ok {
		return v.(*model.Order), true
	}
	return nil, false
}

func (c *cache) SetupCache(orders []*model.Order) {
	for _, o := range orders {
		c.Set(o)
	}
}
