package database

import (
	"database/sql"
	"log"
	"main/models"
	"sort"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var db *Database
var myTime = time.Date(2023, 3, 3, 12, 0, 0, 0, time.FixedZone("", 0)) // manually set location to Greenwich
var records = []models.Record{
	{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "321",
		//CreatedAt:   myTime,
	},
	{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "54321",
		//CreatedAt:   myTime,
	},
}

const connStr = "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

func Test_NewDB(t *testing.T) {
	var err error
	db, err = NewDB(connStr)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
}

func Test_CreateAndRead(t *testing.T) {
	m, err := migrate.New("file://.././migration", connStr)
	if err != nil {
		log.Fatalf("failed to migration init: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	err = db.CreateRecord(&records[0])
	assert.Nil(t, err)
	result, err := db.GetRecord(records[0].ID)
	assert.Nil(t, err)
	assert.Equal(t, records[0], result)

	err = m.Down()
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}

func Test_ReadAll(t *testing.T) {
	m, err := migrate.New("file://.././migration", connStr)
	if err != nil {
		log.Fatalf("failed to migration init: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	err = db.CreateRecord(&records[0])
	assert.Nil(t, err)
	err = db.CreateRecord(&records[1])
	assert.Nil(t, err)
	results, err := db.GetAllRecords()
	assert.Nil(t, err)

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})

	assert.Equal(t, records, results)

	err = m.Down()
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}

func Test_Update(t *testing.T) {
	m, err := migrate.New("file://.././migration", connStr)
	if err != nil {
		log.Fatalf("failed to migration init: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	record := models.Record{
		ID:          records[0].ID,
		Type:        records[0].Type,
		CaesarShift: records[0].CaesarShift,
		Result:      "987654321",
	}

	err = db.CreateRecord(&record)
	assert.Nil(t, err)
	record.Result = "123456789"

	err = db.UpdateRecord(&record)
	assert.Nil(t, err)

	err = m.Down()
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}

func Test_Delete(t *testing.T) {
	m, err := migrate.New("file://.././migration", connStr)
	if err != nil {
		log.Fatalf("failed to migration init: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	err = db.CreateRecord(&records[0])
	assert.Nil(t, err)
	err = db.DeleteRecord(records[0].ID)
	assert.Nil(t, err)

	_, err = db.GetRecord(records[0].ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	err = m.Down()
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}
