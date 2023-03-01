package cache

import (
	"main/repo"
	"main/service"
	"testing"
	"time"

	"github.com/google/uuid"
)

var cache service.Cache

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
	cache = NewLruCache(3)
}

func Test_LruCache(t *testing.T) {
	cache.Set(&records[0])
	cache.Set(&records[0])

	cache.Set(&records[1])
	cache.Set(&records[2])
	cache.Set(&records[3])

	cache.Get(records[1].ID)
	cache.Get(records[0].ID)

	cache.Delete(records[3].ID)
}

func Test_NewInMemoCache(t *testing.T) {
	cache = NewInMemoCache()
}

func Test_InMemoCache(t *testing.T) {
	cache.Set(&records[0])
	cache.Set(&records[1])
	cache.Set(&records[2])
	cache.Set(&records[3])

	cache.Get(records[1].ID)
	cache.Get(records[0].ID)

	cache.Delete(records[2].ID)
}
