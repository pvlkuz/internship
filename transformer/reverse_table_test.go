package transformer

import (
	"strings"
	"testing"
)

type TestReverse struct {
	input, expected string
	ioinput         bool
}

var TestArrayReverse = []TestReverse{
	TestReverse{"abcd", "dcba", false},
	TestReverse{"123456789", "987654321", false},
	TestReverse{"12345", "54321", false},
	TestReverse{"1", "1", false},
	TestReverse{"", "", false},
	TestReverse{"1230", "321", true},
}

func TestTableReverse(t *testing.T) {

	for _, test := range TestArrayReverse {

		result, err := NewReverseTransformer().Transform(strings.NewReader(test.input), test.ioinput)
		if err != nil {
			t.Errorf("Error transforming")
		}

		if result != test.expected {
			t.Errorf("Error: result = %q, expected = %q", result, test.expected)
		}

	}
}
