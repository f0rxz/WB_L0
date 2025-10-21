package cache

import (
	"context"
	"sync"
	"time"

	"orderservice/internal/model"
)

type Cache interface {
	Set(order *model.Order)
	Get(orderUID string) (*model.Order, bool)
	SetupCache(orders []*model.Order)
	Close()
}

type entry struct {
	order     *model.Order
	expiresAt time.Time
}

type cache struct {
	mu     sync.RWMutex
	data   map[string]entry
	ttl    time.Duration
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewCache() Cache {
	return NewCacheWithTTL(0)
}

func NewCacheWithTTL(ttl time.Duration) Cache {
	ctx, cancel := context.WithCancel(context.Background())
	c := &cache{
		data:   make(map[string]entry),
		ttl:    ttl,
		ctx:    ctx,
		cancel: cancel,
	}

	if ttl > 0 {
		c.Run(c.ctx)
	}

	return c
}

func (c *cache) Set(order *model.Order) {
	var exp time.Time
	if c.ttl > 0 {
		exp = time.Now().Add(c.ttl)
	}
	c.mu.Lock()
	c.data[order.OrderUID] = entry{order: order, expiresAt: exp}
	c.mu.Unlock()
}

func (c *cache) Get(orderUID string) (*model.Order, bool) {
	c.mu.RLock()
	e, ok := c.data[orderUID]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if !e.expiresAt.IsZero() && time.Now().After(e.expiresAt) {
		c.mu.Lock()
		delete(c.data, orderUID)
		c.mu.Unlock()
		return nil, false
	}
	return e.order, true
}

func (c *cache) SetupCache(orders []*model.Order) {
	c.mu.Lock()
	for _, o := range orders {
		var exp time.Time
		if c.ttl > 0 {
			exp = time.Now().Add(c.ttl)
		}
		c.data[o.OrderUID] = entry{order: o, expiresAt: exp}
	}
	c.mu.Unlock()
}

func (c *cache) Close() {
	if c.cancel != nil {
		c.cancel()
		c.wg.Wait()
	}
}
