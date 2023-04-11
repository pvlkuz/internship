package service

import (
	"log"
	"main/cache"
	database "main/data-base"
	"main/models"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testRecord = []models.Record{
	{
		ID:          "1111",
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "54321",
		CreatedAt:   time.Now(),
	},
	{
		ID:          "2222",
		Type:        "caesar",
		CaesarShift: -3,
		Result:      "xyz",
		CreatedAt:   time.Now(),
	},
}
var requestTable = []models.TransformRequest{
	{Type: "caesar", CaesarShift: -3, Input: "abc"},
	{Type: "reverse", CaesarShift: 0, Input: "54321"},
	{Type: "base64", CaesarShift: 0, Input: "Man"},
}
var resultTable = []string{"xyz", "12345", "TWFu"}

type mockDb struct {
	mock.Mock
}

func (mock *mockDb) CreateRecord(r *models.Record) error {
	return nil
}
func (mock *mockDb) GetRecord(id string) (models.Record, error) {
	result := testRecord[0]
	return result, nil
}
func (mock *mockDb) GetAllRecords() ([]models.Record, error) {
	result := testRecord
	return result, nil
}
func (mock *mockDb) UpdateRecord(r *models.Record) error {
	return nil
}
func (mock *mockDb) DeleteRecord(id string) error {
	return nil
}

type mockCache struct {
	mock.Mock
}

func (mock *mockCache) Set(value *models.Record) {

}
func (mock *mockCache) Get(key string) (*models.Record, bool) {
	return nil, false
}
func (mock *mockCache) Delete(key string) {

}

func Test_NewRecord(t *testing.T) {
	testService := NewService(new(mockDb), new(mockCache))

	for i, test := range requestTable {
		res, err := testService.CreateRecord(test)
		assert.ErrorIs(t, err, nil)
		assert.Equal(t, test.Type, res.Type)
		assert.Equal(t, test.CaesarShift, res.CaesarShift)
		assert.Equal(t, resultTable[i], res.Result)
	}
}

func Test_GetRecord(t *testing.T) {
	testService := NewService(new(mockDb), new(mockCache))

	res, err := testService.GetRecord("1111")
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, testRecord[0], *res)
}

func Test_GetAllRecords(t *testing.T) {
	testService := NewService(new(mockDb), new(mockCache))

	res, err := testService.GetAllRecords()
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, testRecord, res)
}

func Test_UpdateRecord(t *testing.T) {
	testService := NewService(new(mockDb), new(mockCache))

	for _, test := range requestTable {
		testService.UpdateRecord("123", test)
	}

	res, err := testService.UpdateRecord("1111", requestTable[0])
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, testRecord[0], *res)
}

func Test_DeleteRecord(t *testing.T) {
	testService := NewService(new(mockDb), new(mockCache))

	err := testService.DeleteRecord("123")
	assert.ErrorIs(t, err, nil)
}

const connStr = "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

func Benchmark_GetRecord(b *testing.B) {
	db, err := database.NewDB(connStr)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	err = db.MigrateUp(connStr, ".././migration")
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	testcache := new(mockCache)
	s := NewService(db, testcache)
	var ids [20]string
	for i := 0; i < 20; i++ {
		result, _ := s.CreateRecord(requestTable[1])
		ids[i] = result.ID
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.GetRecord(ids[i%10])
	}

	err = db.MigrateDown(connStr, ".././migration")
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}

func Benchmark_GetRecord_WithCache(b *testing.B) {
	db, err := database.NewDB(connStr)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	err = db.MigrateUp(connStr, ".././migration")
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	testcache := cache.NewLruCache(10)
	s := NewService(db, testcache)
	var ids [20]string
	for i := 0; i < 20; i++ {
		result, _ := s.CreateRecord(requestTable[1])
		ids[i] = result.ID
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.GetRecord(ids[i%10])
	}

	err = db.MigrateDown(connStr, ".././migration")
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}
