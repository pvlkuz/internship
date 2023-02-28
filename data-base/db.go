package database

import (
	"main/repo"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // DB driver
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "error opening new DB connection")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "error checking DB connection")
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
	// if err != nil {
	// 	return err
	// }

	// return nil

	return errors.Wrap(err, "error creating new record")
}

func (db *RecordDB) GetRecord(id string) (repo.Record, error) {
	var result repo.Record

	err := db.Get(&result, QuerySingleRead, id)
	if err != nil {
		return repo.Record{}, errors.Wrap(err, "error reading one record")
	}

	return result, nil
}

func (db *RecordDB) GetAllRecords() ([]repo.Record, error) {
	var result []repo.Record
	err := db.Select(&result, QueryMultiRead)
	// if err != nil {
	// 	return []repo.Record{}, err
	// }
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})

	return result, errors.Wrap(err, "error reading all records")
}

func (db *RecordDB) UpdateRecord(r *repo.Record) error {
	err := db.Get(r, QueryUpdate, r.Type, r.CaesarShift, r.Result, r.UpdatedAt, r.ID)
	if err != nil {
		return errors.Wrap(err, "error updating record")
	}

	return nil
}

func (db *RecordDB) DeleteRecord(id string) error {
	_, err := db.Exec(QueryDelete, id)
	if err != nil {
		return errors.Wrap(err, "error deleting record")
	}

	return nil
}
