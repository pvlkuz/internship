package database

import (
	"errors"
	"fmt"
	"main/models"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	queryCreate     = `INSERT INTO records VALUES ($1, $2, $3, $4) RETURNING *`
	querySingleRead = `SELECT * FROM records WHERE id = $1`
	queryMultiRead  = `SELECT * FROM records ORDER BY created_at DESC`
	queryUpdate     = `UPDATE records SET transform_type = $1, caesar_shift = $2, result = $3 WHERE id = $4 RETURNING *`
	queryDelete     = `DELETE FROM records WHERE id = $1`
)

type Database struct {
	*sqlx.DB
}

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

func (db *Database) CreateRecord(r *models.Record) error {
	err := db.Get(r, queryCreate, r.ID, r.Type, r.CaesarShift, r.Result)
	if err != nil {
		return fmt.Errorf("creating new record: %w", err)
	}

	return nil
}

func (db *Database) GetRecord(id string) (models.Record, error) {
	var result models.Record

	err := db.Get(&result, querySingleRead, id)
	if err != nil {
		return result, fmt.Errorf("reading one record: %w", err)
	}

	return result, nil
}

func (db *Database) GetAllRecords() ([]models.Record, error) {
	var result []models.Record
	err := db.Select(&result, queryMultiRead)

	if err != nil {
		return result, fmt.Errorf("reading all records: %w", err)
	}

	return result, nil
}

func (db *Database) UpdateRecord(r *models.Record) error {
	_, err := db.Exec(queryUpdate, r.Type, r.CaesarShift, r.Result, r.ID)
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

func (db *Database) MigrateUp(connStr string, path string) error {
	filePath := fmt.Sprintf("file://%s", path)

	m, err := migrate.New(filePath, connStr)
	if err != nil {
		return fmt.Errorf("creating migration err: %w", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("migration up err: %w", err)
	}

	return nil
}

func (db *Database) MigrateDown(connStr string, path string) error {
	filePath := fmt.Sprintf("file://%s", path)

	m, err := migrate.New(filePath, connStr)
	if err != nil {
		return fmt.Errorf("creating migration err: %w", err)
	}

	err = m.Down()
	if err != nil {
		return fmt.Errorf("migration up err: %w", err)
	}

	return nil
}
