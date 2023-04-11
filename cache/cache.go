package cache

import (
	"container/list"
	"main/models"
	"sync"
)

type InMemoCache struct {
	cache map[string]*models.Record
	mu    sync.Mutex
}

func NewInMemoCache() *InMemoCache {
	return &InMemoCache{
		cache: make(map[string]*models.Record),
		mu:    sync.Mutex{},
	}
}

func (c *InMemoCache) Set(value *models.Record) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[value.ID] = value
}

func (c *InMemoCache) Get(key string) (*models.Record, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	res, ok := c.cache[key]

	return res, ok
}

func (c *InMemoCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

type LruCache struct {
	capacity int
	queue    *list.List
	cache    map[string]*Item
	mu       sync.Mutex
}

type Item struct {
	data *models.Record
	key  *list.Element
}

func NewLruCache(capacity int) *LruCache {
	return &LruCache{
		capacity: capacity,
		queue:    list.New(),
		cache:    make(map[string]*Item),
		mu:       sync.Mutex{},
	}
}

func (c *LruCache) Set(value *models.Record) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.cache[value.ID]
	if ok {
		item.data = value
		c.cache[value.ID] = item
		c.queue.MoveToFront(item.key)

		return
	}

	if c.capacity == len(c.cache) {
		back := c.queue.Back()
		c.queue.Remove(back)
		key, _ := back.Value.(string)
		delete(c.cache, key)
	}

	c.cache[value.ID] = &Item{
		data: value,
		key:  c.queue.PushFront(value.ID),
	}
}

func (c *LruCache) Get(key string) (*models.Record, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item.key)

	return item.data, true
}

func (c *LruCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.cache[key]
	if ok {
		c.queue.Remove(item.key)
		delete(c.cache, key)
	}
}
