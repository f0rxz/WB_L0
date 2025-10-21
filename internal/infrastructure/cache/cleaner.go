package cache

import (
	"context"
	"time"
)

func (c *cache) Run(ctx context.Context) {
	if c.ttl <= 0 {
		return
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		ticker := time.NewTicker(c.ttl)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				c.mu.Lock()
				for k, e := range c.data {
					if !e.expiresAt.IsZero() && now.After(e.expiresAt) {
						delete(c.data, k)
					}
				}
				c.mu.Unlock()
			}
		}
	}()
}
