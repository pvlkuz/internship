package database

import (
	"main/repo"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
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
		return err
	}
	return nil
}

func (db *RecordDB) GetRecord(id string) (repo.Record, error) {
	var r repo.Record
	err := db.Get(&r, QuerySingleRead, id)
	if err != nil {
		return repo.Record{}, err
	}
	return r, nil
}

func (db *RecordDB) GetRecords() ([]repo.Record, error) {
	var r []repo.Record
	err := db.Select(&r, QueryMultiRead)
	if err != nil {
		return []repo.Record{}, err
	}
	return r, err
}

func (db *RecordDB) UpdateRecord(r *repo.Record) error {
	err := db.Get(r, QueryUpdate, r.Type, r.CaesarShift, r.Result, r.UpdatedAt, r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *RecordDB) DeleteRecord(id string) error {
	_, err := db.Exec(QueryDelete, id)
	if err != nil {
		return err
	}
	return nil
}
