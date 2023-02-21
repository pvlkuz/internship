package service

import (
	"fmt"
	"log"
	"main/cache"
	database "main/data-base"
	"main/repo"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (mock *MockDB) NewRecord(r *repo.Record) error {
	// args := mock.Called()
	// result :=args.Get(0)
	return nil
}
func (mock *MockDB) GetRecord(id string) (repo.Record, error) {
	result := repo.Record{
		ID:          "1111",
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "54321",
		CreatedAt:   time.Now().Unix(),
	}
	return result, nil
}
func (mock *MockDB) GetRecords() ([]repo.Record, error) {
	result := []repo.Record{
		{
			ID:          uuid.NewString(),
			Type:        "reverse",
			CaesarShift: 0,
			Result:      "54321",
			CreatedAt:   time.Now().Unix(),
		},
		{
			ID:          uuid.NewString(),
			Type:        "caesar",
			CaesarShift: -3,
			Result:      "xyz",
			CreatedAt:   time.Now().Unix(),
		},
	}
	return result, nil
}
func (mock *MockDB) UpdateRecord(r *repo.Record) error {
	return nil
}
func (mock *MockDB) DeleteRecord(id string) error {
	return nil
}

type MockCache struct {
	mock.Mock
}

func (mock *MockCache) Set(value *repo.Record) {

}
func (mock *MockCache) Get(key string) (*repo.Record, bool) {
	return nil, false
}
func (mock *MockCache) Delete(key string) {

}

var TestService ServiceInterface

func Test_NewService(t *testing.T) {
	TestService = NewService(new(MockDB), new(MockCache))
}

var NewRecordRequestTable = []TransformRequest{
	TransformRequest{Type: "caesar", CaesarShift: -3, Input: "abc"},
	TransformRequest{Type: "reverse", CaesarShift: 0, Input: "54321"},
	TransformRequest{Type: "base64", CaesarShift: 0, Input: "Man"},
}
var NewRecordResultTable = []string{
	"xyz", "12345", "TWFu",
}

func Test_NewRecord(t *testing.T) {
	for _, test := range NewRecordRequestTable {
		TestService.NewRecord(test)
	}
}

func Test_GetRecord(t *testing.T) {
	TestService.GetRecord("123")
}

func Test_GetRecords(t *testing.T) {
	TestService.GetRecords()
}

func Test_UpdateRecord(t *testing.T) {
	for _, test := range NewRecordRequestTable {
		TestService.UpdateRecord("123", test)
	}
}

func Test_DeleteRecord(t *testing.T) {
	TestService.DeleteRecord("123")
}

const connStr = "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

func Benchmark_GetRecord(b *testing.B) {
	m, err := migrate.New("file://.././migration", connStr)
	if err != nil {
		log.Fatalf("failed to migration init: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		log.Print(fmt.Errorf("failed to migrate up: %s", err.Error()))
		return
	}
	db, err := database.NewDB(connStr)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	//testcache := cache.NewLruCache(10)
	// testcache := cache.NewInMemoCache()
	testcache := new(MockCache)
	s := NewService(db, testcache)
	var ids [20]string
	for i := 0; i < 20; i++ {
		result := s.NewRecord(NewRecordRequestTable[1])
		ids[i] = result.ID
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.GetRecord(ids[i%10])
	}

	err = m.Down()
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}

func Benchmark_GetRecord_WithCache(b *testing.B) {
	m, err := migrate.New("file://.././migration", connStr)
	if err != nil {
		log.Fatalf("failed to migration init: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		log.Print(fmt.Errorf("failed to migrate up: %s", err.Error()))
		return
	}
	db, err := database.NewDB(connStr)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	testcache := cache.NewLruCache(10)
	// testcache := cache.NewInMemoCache()
	//testcache := new(MockCache)
	s := NewService(db, testcache)
	var ids [20]string
	for i := 0; i < 20; i++ {
		result := s.NewRecord(NewRecordRequestTable[1])
		ids[i] = result.ID
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.GetRecord(ids[i%10])
	}

	err = m.Down()
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}
