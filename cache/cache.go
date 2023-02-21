package cache

import (
	"container/list"
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

type LruCache struct {
	capacity int
	queue    *list.List
	cache    map[string]*Item
}

type Item struct {
	data *repo.Record
	key  *list.Element
}

func NewLruCache(capacity int) *LruCache {
	return &LruCache{
		capacity: capacity,
		queue:    list.New(),
		cache:    make(map[string]*Item),
	}
}

func (c *LruCache) Set(value *repo.Record) {
	item, ok := c.cache[value.ID]
	if !ok {
		if c.capacity == len(c.cache) {
			back := c.queue.Back()
			c.queue.Remove(back)
			delete(c.cache, back.Value.(string))
		}

		c.cache[value.ID] = &Item{
			data: value,
			key:  c.queue.PushFront(value.ID),
		}
	} else {
		item.data = value
		c.cache[value.ID] = item
		c.queue.MoveToFront(item.key)
	}
}

func (c *LruCache) Get(key string) *repo.Record {
	item, ok := c.cache[key]
	if ok {
		c.queue.MoveToFront(item.key)
		return item.data
	}
	return nil
}

func (c *LruCache) Delete(key string) {
	item, ok := c.cache[key]
	if ok {
		c.queue.Remove(item.key)
		delete(c.cache, key)
	}
}
