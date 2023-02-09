package cache

import (
	"main/repo"
	"sync"
)

type CacheInterface interface {
	Set(value *repo.Record)
	Get(key string) *repo.Record
	Delete(key string)
}

type InMemoCache struct {
	cache map[string]*repo.Record
	mu    sync.Mutex
}

func NewInMemoCache() *InMemoCache {
	return &InMemoCache{cache: make(map[string]*repo.Record)}
}

func (c *InMemoCache) Set(value *repo.Record) {
	c.mu.Lock()
	c.cache[value.ID] = value
	c.mu.Unlock()
}

func (c *InMemoCache) Get(key string) *repo.Record {
	c.mu.Lock()
	res, ok := c.cache[key]
	c.mu.Unlock()
	if !ok {
		return nil
	}
	return res
}

func (c *InMemoCache) Delete(key string) {
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
}
