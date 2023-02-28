package cache

import (
	"container/list"
	"errors"
	"fmt"
	"main/repo"
	"sync"
	"time"

	"github.com/go-redis/redis"
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

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		//nolint:exhaustivestruct, exhaustruct
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (r *RedisCache) Set(value *repo.Record) {
	r.client.Set(value.ID, value, 15*time.Second)
}

func (r *RedisCache) Get(key string) (*repo.Record, bool) {
	result := repo.Record{} //nolint:exhaustivestruct, exhaustruct

	get := r.client.Get(key)
	fmt.Println(get.Bytes())

	err := r.client.Get(key).Scan(result)
	fmt.Println(err)
	if errors.Is(err, redis.Nil) {
		return nil, false
	}

	return &result, true
}

func (r *RedisCache) Delete(key string) {

}
