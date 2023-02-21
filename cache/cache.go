package cache

import (
	"container/list"
	"main/repo"
	"sync"
)

type InMemoCache struct {
	cache map[string]*repo.Record
	mu    sync.Mutex
}

func NewInMemoCache() *InMemoCache {
	return &InMemoCache{cache: make(map[string]*repo.Record)}
}

func (c *InMemoCache) Set(value *repo.Record) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[value.ID] = value
}

func (c *InMemoCache) Get(key string) (*repo.Record, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	res, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return res, true
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
	if ok {
		item.data = value
		c.cache[value.ID] = item
		c.queue.MoveToFront(item.key)
	} else {
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
}

func (c *LruCache) Get(key string) (*repo.Record, bool) {
	item, ok := c.cache[key]
	if ok {
		c.queue.MoveToFront(item.key)
		return item.data, true
	}
	return nil, false
}

func (c *LruCache) Delete(key string) {
	item, ok := c.cache[key]
	if ok {
		c.queue.Remove(item.key)
		delete(c.cache, key)
	}
}