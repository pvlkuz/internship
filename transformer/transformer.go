package transformer

import (
	"encoding/base64"
	"fmt"
	"io"
)

type Transformer interface {
	Transform(in io.Reader, ioinput bool) (string, error)
}

type CaesarTransformer struct {
	Shift int
}

func NewCaesarTransformer(shift int) *CaesarTransformer {
	return &CaesarTransformer{Shift: shift}
}

func (t *CaesarTransformer) Transform(in io.Reader, ioinput bool) (string, error) {
	f, err := io.ReadAll(in)
	if err != nil {
		return "", err
	}
	if ioinput {
		f = f[:len(f)-1]
	}
	rns := []rune(string(f))
	for i := 0; i < len(rns); i++ {
		r := int(rns[i]) + t.Shift
		switch {
		case r > 'z':
			rns[i] = rune(r - 26)
		case r < 'a':
			rns[i] = rune(r + 26)
		default:
			rns[i] = rune(r)
		}
	}
	result := string(rns)

	return result, nil
}

type ReverseTransformer struct{}

func NewReverseTransformer() *ReverseTransformer {
	return &ReverseTransformer{}
}

func (t *ReverseTransformer) Transform(in io.Reader, ioinput bool) (string, error) {
	f, err := io.ReadAll(in)
	if err != nil {
		return "", err
	}
	if ioinput {
		f = f[:len(f)-1]
	}
	rns := []rune(string(f))
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	result := string(rns)
	return result, nil
}

type Base64Transformer struct{}

func NewBase64Transformer() *Base64Transformer {
	return &Base64Transformer{}
}

func (t *Base64Transformer) Transform(in io.Reader, ioinput bool) (string, error) {
	f, err := io.ReadAll(in)
	if err != nil {
		return "", err
	}
	if ioinput {
		f = f[:len(f)-1]
	}
	result := base64.StdEncoding.EncodeToString(f)
	return result, nil
}

func BasicTransform(in io.Reader, out io.Writer, caesaarShift int, base64Use bool, ioinput bool) error {
	var tr Transformer
	switch {
	case base64Use:
		tr = NewBase64Transformer()
	case caesaarShift != 0:
		tr = NewCaesarTransformer(caesaarShift)
	default:
		tr = NewReverseTransformer()
	}
	result, err := tr.Transform(in, ioinput)
	if err != nil {
		return fmt.Errorf("TRANSFORMER error: %w", err)
	}

	_, err = out.Write([]byte(result))
	if err != nil {
		return fmt.Errorf("write output error: %w", err)
	}
	return nil
}
