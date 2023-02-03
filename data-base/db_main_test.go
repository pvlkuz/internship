package database

import (
	"log"
	"main/repo"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db *RecordDB

const connStr = "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

func Test_NewDB(t *testing.T) {
	var err error
	db, err = NewDB(connStr)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
}

func Test_CRUD(t *testing.T) {
	record := &repo.Record{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "321",
		CreatedAt:   time.Now().Unix(),
	}

	m, err := migrate.New("file://.././migration", connStr)
	if err != nil {
		log.Fatalf("failed to migration init: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	err = db.NewRecord(record)
	//assert.Nil(t, err)
	if err != nil {
		t.Fatalf("error in creating new record")
	}
	_, err = db.GetRecord(record.ID)
	if err != nil {
		t.Fatalf("such record cant be read")
	}

	err = db.NewRecord(record)
	record = &repo.Record{
		ID:          uuid.NewString(),
		Type:        record.Type,
		CaesarShift: record.CaesarShift,
		Result:      "54321",
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   time.Now().Unix(),
	}
	err = db.NewRecord(record)
	if err != nil {
		t.Fatalf("error in creating new record")
	}
	_, err = db.GetRecords()
	if err != nil {
		t.Fatalf("multiple record cant be read")
	}

	record = &repo.Record{
		ID:          record.ID,
		Type:        record.Type,
		CaesarShift: record.CaesarShift,
		Result:      "987654321",
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   time.Now().Unix(),
	}
	err = db.UpdateRecord(record)
	if err != nil {
		t.Fatalf("such record cant be updated")
	}

	err = db.DeleteRecord(record.ID)
	if err != nil {
		t.Fatalf("such record cant be deleted")
	}

	err = m.Down()
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}
