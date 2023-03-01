package cache

import (
	"container/list"
	"context"
	"main/repo"
	"sync"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type InMemoCache struct {
	cache map[string]*repo.Record
	mu    sync.Mutex
}

func NewInMemoCache() *InMemoCache {
	return &InMemoCache{
		cache: make(map[string]*repo.Record),
		mu:    sync.Mutex{},
	}
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
	data *repo.Record
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

func (c *LruCache) Set(value *repo.Record) {
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

func (c *LruCache) Get(key string) (*repo.Record, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.cache[key]
	if ok {
		c.queue.MoveToFront(item.key)

		return item.data, true
	}

	return nil, false
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

type MyRedisCache struct {
	cache *cache.Cache
}

//nolint:exhaustivestruct, exhaustruct
func NewRedisCache(address string) *MyRedisCache {
	redisClient := redis.NewClient(&redis.Options{
		Addr: address,
	})

	mycache := cache.New(&cache.Options{
		Redis: redisClient,
	})

	return &MyRedisCache{
		cache: mycache,
	}
}

func (r *MyRedisCache) Set(value *repo.Record) {
	//nolint:exhaustivestruct, exhaustruct, errcheck
	r.cache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   value.ID,
		Value: value,
		TTL:   20 * time.Second,
	})
}

func (r *MyRedisCache) Get(key string) (*repo.Record, bool) {
	result := repo.Record{} //nolint:exhaustivestruct, exhaustruct

	err := r.cache.Get(context.TODO(), key, &result)
	if err != nil {
		return nil, false
	}

	return &result, true
}

func (r *MyRedisCache) Delete(key string) {

}
