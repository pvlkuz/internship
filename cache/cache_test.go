package cache

import (
	"main/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var records = []models.Record{
	{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "321",
		CreatedAt:   time.Now(),
	},
	{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "54321",
		CreatedAt:   time.Now(),
	},
	{
		ID:          uuid.NewString(),
		Type:        "base64",
		CaesarShift: 0,
		Result:      "Man",
		CreatedAt:   time.Now(),
	},
	{
		ID:          uuid.NewString(),
		Type:        "caesar",
		CaesarShift: -3,
		Result:      "abc",
		CreatedAt:   time.Now(),
	},
}

func Test_LruSetAndGet(t *testing.T) {
	cache := NewLruCache(3)

	cache.Set(&records[0])

	result, _ := cache.Get(records[0].ID)
	assert.Equal(t, records[0], *result)
}

func Test_LruDelete(t *testing.T) {
	cache := NewLruCache(3)

	cache.Set(&records[0])
	cache.Delete(records[0].ID)

	_, ok := cache.Get(records[0].ID)
	assert.Equal(t, false, ok)
}

func Test_LruCheckCapacity(t *testing.T) {
	cache := NewLruCache(3)

	cache.Set(&records[0])
	cache.Set(&records[1])
	cache.Set(&records[2])
	cache.Set(&records[3])

	// Check that 1st value is auto-deleted
	_, ok := cache.Get(records[0].ID)
	assert.Equal(t, false, ok)

	// Check that 2nd value is still there
	result, ok := cache.Get(records[1].ID)
	assert.Equal(t, true, ok)
	assert.Equal(t, records[1], *result)

}

func Test_InMemoSetAndGet(t *testing.T) {
	cache := NewInMemoCache()

	cache.Set(&records[0])

	result, _ := cache.Get(records[0].ID)
	assert.Equal(t, records[0], *result)
}

func Test_InMemoDelete(t *testing.T) {
	cache := NewInMemoCache()

	cache.Set(&records[0])
	cache.Delete(records[0].ID)

	_, ok := cache.Get(records[0].ID)
	assert.Equal(t, false, ok)
}
