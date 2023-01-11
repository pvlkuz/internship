package repo

import "github.com/google/uuid"

type Record struct {
	Id          string `db:"id"`
	Type        string `db:"transform_type"`
	CaesarShift int    `db:"caesar_shift"`
	Result      string `db:"result"`
	Created_at  int64  `db:"created_at"`
	Updated_at  int64  `db:"updated_at"`
}

type RecordDB interface {
	NewRecord(r *Record) error
	GetRecord(id uuid.UUID) error
	GetRecords() error
	UpdateRecord(r *Record) error
	DeleteRecord(id uuid.UUID) error
}
