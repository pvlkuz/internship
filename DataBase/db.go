package database

import (
	"main/repo"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB() (*RecordDB, error) {
	connStr := "postgresql://postgres:password@database:5432/postgres?sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &RecordDB{
		DB: db,
	}, nil
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
	err := db.Get(r, `INSERT INTO records VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`, r.Id, r.Type, r.CaesarShift, r.Result, r.Created_at, r.Updated_at)
	if err != nil {
		return err
	}
	return nil
}

func (db *RecordDB) GetRecord(id string) (repo.Record, error) {
	var r repo.Record
	err := db.Get(&r, `SELECT * FROM records WHERE id = $1`, id)
	if err != nil {
		return repo.Record{}, err
	}
	return r, nil
}

func (db *RecordDB) GetRecords() ([]repo.Record, error) {
	var r []repo.Record
	err := db.Select(&r, `SELECT * FROM records`)
	if err != nil {
		return []repo.Record{}, err
	}
	return r, err
}

func (db *RecordDB) UpdateRecord(r *repo.Record) error {
	err := db.Get(r, `UPDATE records SET transform_type = $1, caesar_shift = $2, result = $3, created_at = $4, updated_at = $5 WHERE id = $6 RETURNING *`,
		r.Type, r.CaesarShift, r.Result, r.Created_at, r.Updated_at, r.Id)
	if err != nil {
		return err
	}
	return nil
}

func (db *RecordDB) DeleteRecord(id string) error {
	_, err := db.Exec(`DELETE FROM records WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
