package audit

import (
	"context"
	"sync"
	"time"
)

// TTLCache wraps a Scorer and memoises results by URL.
type TTLCache struct {
	inner Scorer
	ttl   time.Duration
	mu    sync.Mutex
	items map[string]cacheEntry
}

type cacheEntry struct {
	result Result
	expiry time.Time
}

func NewTTLCache(inner Scorer, ttl time.Duration) *TTLCache {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &TTLCache{
		inner: inner,
		ttl:   ttl,
		items: make(map[string]cacheEntry),
	}
}

func (c *TTLCache) Score(ctx context.Context, pageURL string) (Result, error) {
	now := time.Now()
	c.mu.Lock()
	if e, ok := c.items[pageURL]; ok && now.Before(e.expiry) {
		res := e.result
		c.mu.Unlock()
		return res, nil
	}
	c.mu.Unlock()

	res, err := c.inner.Score(ctx, pageURL)
	if err != nil {
		return Result{}, err
	}

	c.mu.Lock()
	c.items[pageURL] = cacheEntry{result: res, expiry: now.Add(c.ttl)}
	c.mu.Unlock()
	return res, nil
}
