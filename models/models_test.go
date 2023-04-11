package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	var requestTable = []TransformRequest{
		{Type: "caesar", CaesarShift: -3, Input: "abc"},
		{Type: "revers", CaesarShift: 0, Input: "54321"},
		{Type: "caesar", Input: "abc"},
		{Type: "reverse", CaesarShift: 0, Input: ""},
	}

	err := requestTable[0].Validate()
	assert.Nil(t, err)

	err = requestTable[1].Validate()
	assert.Equal(t, err.Error(), "expected tranformation type field: reverse/caesar/base64")

	err = requestTable[2].Validate()
	assert.Equal(t, err.Error(), "expected shift field (not 0)")

	err = requestTable[3].Validate()
	assert.Equal(t, err.Error(), "expected input field")
}
