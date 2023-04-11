package database

import (
	"database/sql"
	"main/models"
	"sort"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const connStr = "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

var records = []models.Record{
	{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "321",
	},
	{
		ID:          uuid.NewString(),
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "54321",
	},
}

func Test_CreateAndRead(t *testing.T) {
	db, err := NewDB(connStr)
	assert.Nil(t, err)

	err = db.MigrateUp(connStr, ".././migration")
	assert.Nil(t, err)

	err = db.CreateRecord(&records[0])
	assert.Nil(t, err)

	result, err := db.GetRecord(records[0].ID)
	assert.Nil(t, err)
	assert.Equal(t, records[0], result)

	err = db.MigrateDown(connStr, ".././migration")
	assert.Nil(t, err)
}

func Test_ReadAll(t *testing.T) {
	db, err := NewDB(connStr)
	assert.Nil(t, err)

	err = db.MigrateUp(connStr, ".././migration")
	assert.Nil(t, err)

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

	err = db.MigrateDown(connStr, ".././migration")
	assert.Nil(t, err)
}

func Test_Update(t *testing.T) {
	db, err := NewDB(connStr)
	assert.Nil(t, err)

	err = db.MigrateUp(connStr, ".././migration")
	assert.Nil(t, err)

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

	err = db.MigrateDown(connStr, ".././migration")
	assert.Nil(t, err)
}

func Test_Delete(t *testing.T) {
	db, err := NewDB(connStr)
	assert.Nil(t, err)

	err = db.MigrateUp(connStr, ".././migration")
	assert.Nil(t, err)

	err = db.CreateRecord(&records[0])
	assert.Nil(t, err)
	err = db.DeleteRecord(records[0].ID)
	assert.Nil(t, err)

	_, err = db.GetRecord(records[0].ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	err = db.MigrateDown(connStr, ".././migration")
	assert.Nil(t, err)
}
