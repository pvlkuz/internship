package repo

type Record struct {
	ID          string `db:"id"`
	Type        string `db:"transform_type"`
	CaesarShift int    `db:"caesar_shift"`
	Result      string `db:"result"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}
