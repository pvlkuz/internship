package transformer

import (
	"strings"
	"testing"
)

type TestBase64 struct {
	input, expected string
	ioinput         bool
}

var TestArrayBase64 = []TestBase64{
	TestBase64{"Man", "TWFu", false},
	TestBase64{"Ma", "TWE=", false},
	TestBase64{"M", "TQ==", false},
	TestBase64{"Man1", "TWFu", true},
}

func TestTableBase64(t *testing.T) {

	for _, test := range TestArrayBase64 {

		result, err := NewBase64Transformer().Transform(strings.NewReader(test.input), test.ioinput)
		if err != nil {
			t.Errorf("Error transforming")
		}

		if result != test.expected {
			t.Errorf("Error: result = %q, expected = %q", result, test.expected)
		}

	}
}
