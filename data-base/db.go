package database

import (
	"fmt"
	"main/models"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	queryCreate     = `INSERT INTO records VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	querySingleRead = `SELECT * FROM records WHERE id = $1`
	queryMultiRead  = `SELECT * FROM records`
	queryUpdate     = `UPDATE records SET transform_type = $1, caesar_shift = $2, result = $3, updated_at = $4 WHERE id = $5 RETURNING *`
	queryDelete     = `DELETE FROM records WHERE id = $1`
)

func NewDB(connStr string) (*Database, error) {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening new DB connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("checking DB connection: %w", err)
	}

	return &Database{DB: db}, nil
}

type Database struct {
	*sqlx.DB
}

func (db *Database) CreateRecord(r *models.Record) error {
	err := db.Get(r, queryCreate, r.ID, r.Type, r.CaesarShift, r.Result, r.CreatedAt, r.UpdatedAt)

	if err != nil {
		return fmt.Errorf("creating new record: %w", err)
	}

	return nil
}

func (db *Database) GetRecord(id string) (models.Record, error) {
	var result models.Record

	err := db.Get(&result, querySingleRead, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.Record{}, nil
		}

		return models.Record{}, fmt.Errorf("reading one record: %w", err)
	}

	return result, nil
}

func (db *Database) GetAllRecords() ([]models.Record, error) {
	var result []models.Record
	err := db.Select(&result, queryMultiRead)

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	if err != nil {
		return result, fmt.Errorf("reading all records: %w", err)
	}

	return result, nil
}

func (db *Database) UpdateRecord(r *models.Record) error {
	_, err := db.Exec(queryUpdate, r.Type, r.CaesarShift, r.Result, r.UpdatedAt, r.ID)
	if err != nil {
		return fmt.Errorf("updating record: %w", err)
	}

	return nil
}

func (db *Database) DeleteRecord(id string) error {
	_, err := db.Exec(queryDelete, id)
	if err != nil {
		return fmt.Errorf("deleting record: %w", err)
	}

	return nil
}
