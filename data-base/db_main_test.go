package database

import (
	"log"
	"main/models"
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
}

const connStr = "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

func Test_NewTestDB(t *testing.T) {
	var err error
	db, err = NewDB(connStr)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
}

func Test_Database(t *testing.T) {

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

	err = db.CreateRecord(&records[0])
	err = db.CreateRecord(&records[1])
	assert.Nil(t, err)
	results, err := db.GetAllRecords()
	assert.Nil(t, err)
	assert.Equal(t, records, results)

	record := models.Record{
		ID:          records[0].ID,
		Type:        records[0].Type,
		CaesarShift: records[0].CaesarShift,
		Result:      "987654321",
		CreatedAt:   records[0].CreatedAt,
		UpdatedAt:   time.Now(),
	}
	records = append(records, record)
	err = db.UpdateRecord(&records[2])
	assert.Nil(t, err)
	result, err = db.GetRecord(records[2].ID)
	assert.Equal(t, records[2], result)

	err = db.DeleteRecord(records[2].ID)
	assert.Nil(t, err)
	result, err = db.GetRecord(records[2].ID)
	assert.NotNil(t, err)

	err = m.Down()
	if err != nil {
		log.Fatalf("failed to migrate down: %s", err.Error())
	}
}
