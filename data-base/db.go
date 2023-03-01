package database

import (
	"fmt"
	"main/repo"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // DB driver
)

const (
	QueryCreate     = `INSERT INTO records VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	QuerySingleRead = `SELECT * FROM records WHERE id = $1`
	QueryMultiRead  = `SELECT * FROM records`
	QueryUpdate     = `UPDATE records SET transform_type = $1, caesar_shift = $2, result = $3, updated_at = $4 WHERE id = $5 RETURNING *`
	QueryDelete     = `DELETE FROM records WHERE id = $1`
)

func NewDB(connStr string) (*RecordDB, error) {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening new DB connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error checking DB connection: %w", err)
	}

	return &RecordDB{DB: db}, nil
}

type RecordDB struct {
	*sqlx.DB
}

func NewRecordDB(db *sqlx.DB) *RecordDB {
	return &RecordDB{
		DB: db,
	}
}
func (db *RecordDB) NewRecord(r *repo.Record) error {
	err := db.Get(r, QueryCreate, r.ID, r.Type, r.CaesarShift, r.Result, r.CreatedAt, r.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating new record: %w", err)
	}

	return nil
}

func (db *RecordDB) GetRecord(id string) (repo.Record, error) {
	var result repo.Record

	err := db.Get(&result, QuerySingleRead, id)
	if err != nil {
		return repo.Record{}, fmt.Errorf("error reading one record: %w", err)
	}

	return result, nil
}

func (db *RecordDB) GetAllRecords() ([]repo.Record, error) {
	var result []repo.Record
	err := db.Select(&result, QueryMultiRead)

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})

	if err != nil {
		return result, fmt.Errorf("error reading all records: %w", err)
	}

	return result, nil
}

func (db *RecordDB) UpdateRecord(r *repo.Record) error {
	err := db.Get(r, QueryUpdate, r.Type, r.CaesarShift, r.Result, r.UpdatedAt, r.ID)
	if err != nil {
		return fmt.Errorf("error updating record: %w", err)
	}

	return nil
}

func (db *RecordDB) DeleteRecord(id string) error {
	_, err := db.Exec(QueryDelete, id)
	if err != nil {
		return fmt.Errorf("error deleting record: %w", err)
	}

	return nil
}
