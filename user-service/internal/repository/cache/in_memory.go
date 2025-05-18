package cache

import (
	"context"
	"sync"
	"time"
)

type memcacheClient struct {
	sync.RWMutex
	data map[string]any
	sets map[string]map[any]struct{}
}

func NewMemcacheClient() Client {
	return &memcacheClient{
		data: make(map[string]any),
		sets: make(map[string]map[any]struct{}),
	}
}

func (c *memcacheClient) Set(ctx context.Context, key string, data any, ttl ...time.Duration) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = data
	return nil
}

func (c *memcacheClient) Get(ctx context.Context, key string) (any, error) {
	c.RLock()
	defer c.RUnlock()

	if value, exists := c.data[key]; exists {
		return value, nil
	}
	return nil, ErrCacheMiss
}

func (c *memcacheClient) AddToSet(ctx context.Context, key string, data ...any) error {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.sets[key]; !exists {
		c.sets[key] = make(map[any]struct{})
	}

	for _, item := range data {
		c.sets[key][item] = struct{}{}
	}
	return nil
}

func (c *memcacheClient) IsDataInSet(ctx context.Context, key string, data any) (bool, error) {
	c.RLock()
	defer c.RUnlock()

	if set, exists := c.sets[key]; exists {
		_, isMember := set[data]
		return isMember, nil
	}
	return false, nil
}
