package models

import (
	"fmt"
	"time"
)

type Record struct {
	ID          string    `db:"id"`
	Type        string    `db:"transform_type"`
	CaesarShift int       `db:"caesar_shift"`
	Result      string    `db:"result"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type TransformRequest struct {
	Type        string `json:"type"`
	CaesarShift int    `json:"shift,omitempty"`
	Input       string `json:"input,omitempty"`
}

var (
	ErrInvalidType = fmt.Errorf("expected tranformation type field: reverse/caesar/base64")
	ErrShift       = fmt.Errorf("expected shift field (not 0)")
	ErrIn          = fmt.Errorf("expected input field")
)

func (t TransformRequest) Validate() error {
	if t.Type != "reverse" && t.Type != "caesar" && t.Type != "base64" {
		return ErrInvalidType
	}

	if t.Type == "caesar" && t.CaesarShift == 0 {
		return ErrShift
	}

	if t.Input == "" {
		return ErrIn
	}

	return nil
}
