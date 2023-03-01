package cache

import (
	"main/repo"
	"main/service"
	"testing"
	"time"

	"github.com/google/uuid"
)

var testcache service.CacheInterface

var records = []repo.Record{
	{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "321",
		CreatedAt:   time.Now().Unix(),
	},
	{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "54321",
		CreatedAt:   time.Now().Unix(),
	},
	{
		ID:          uuid.NewString(),
		Type:        "base64",
		CaesarShift: 0,
		Result:      "Man",
		CreatedAt:   time.Now().Unix(),
	},
	{
		ID:          uuid.NewString(),
		Type:        "caesar",
		CaesarShift: -3,
		Result:      "abc",
		CreatedAt:   time.Now().Unix(),
	},
}

func Test_NewLruCache(t *testing.T) {
	testcache = NewLruCache(3)
}

func Test_LruCache(t *testing.T) {
	testcache.Set(&records[0])
	testcache.Set(&records[0])

	testcache.Set(&records[1])
	testcache.Set(&records[2])
	testcache.Set(&records[3])

	testcache.Get(records[1].ID)
	testcache.Get(records[0].ID)

	testcache.Delete(records[3].ID)
}

func Test_NewInMemoCache(t *testing.T) {
	testcache = NewInMemoCache()
}

func Test_InMemoCache(t *testing.T) {
	testcache.Set(&records[0])
	testcache.Set(&records[1])
	testcache.Set(&records[2])
	testcache.Set(&records[3])

	testcache.Get(records[1].ID)
	testcache.Get(records[0].ID)

	testcache.Delete(records[2].ID)
}

func Test_NewRedisCache(t *testing.T) {
	testcache = NewRedisCache("localhost:6379")
}

func Test_RedisCache(t *testing.T) {
	testcache.Set(&records[0])
	testcache.Get(records[0].ID)
}
