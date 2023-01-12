package transformer

import (
	"encoding/base64"
	"io"
)

type Transformer interface {
	Transform(in io.Reader) (string, error)
}

type CaesarTransformer struct {
	Shift int
}

func NewCaesarTransformer(shift int) *CaesarTransformer {
	return &CaesarTransformer{Shift: shift}
}

func (t *CaesarTransformer) Transform(in io.Reader) (string, error) {
	var result string
	f, err := io.ReadAll(in)
	if err != nil {
		return result, err
	}

	rns := []rune(string(f))
	for i := 0; i < len(rns); i++ {
		r := int(rns[i]) + t.Shift
		if r > 'z' {
			rns[i] = rune(r - 26)
		} else if r < 'a' {
			rns[i] = rune(r + 26)
		} else {
			rns[i] = rune(r)
		}
	}
	result = string(rns)

	return result, nil
}

type ReverseTransformer struct{}

func NewReverseTransformer() *ReverseTransformer {
	return &ReverseTransformer{}
}

func (t *ReverseTransformer) Transform(in io.Reader) (string, error) {
	var result string
	f, err := io.ReadAll(in)
	if err != nil {
		return result, err
	}

	rns := []rune(string(f))
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	result = string(rns)
	return result, nil
}

type Base64Transformer struct{}

func NewBase64Transformer() *Base64Transformer {
	return &Base64Transformer{}
}

func (t *Base64Transformer) Transform(in io.Reader) (string, error) {
	var result string
	f, err := io.ReadAll(in)
	if err != nil {
		return result, err
	}
	result = base64.StdEncoding.EncodeToString(f)
	return result, nil
}
