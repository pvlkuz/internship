package repo

import (
	"fmt"
)

type Record struct {
	ID          string `db:"id"`
	Type        string `db:"transform_type"`
	CaesarShift int    `db:"caesar_shift"`
	Result      string `db:"result"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}

type TransformRequest struct {
	Type        string `json:"type"`
	CaesarShift int    `json:"shift,omitempty"`
	Input       string `json:"input,omitempty"`
}

var (
	ErrType  = fmt.Errorf("expected tranformation type field: reverse/caesar/base64")
	ErrShift = fmt.Errorf("expected shift field (not 0)")
	ErrIn    = fmt.Errorf("expected input field")
)

func (t TransformRequest) Validate() error {
	if t.Type != "reverse" && t.Type != "caesar" && t.Type != "base64" {
		return ErrType
	}

	if t.Type == "caesar" && t.CaesarShift == 0 {
		return ErrShift
	}

	if t.Input == "" {
		return ErrIn
	}

	return nil
}
