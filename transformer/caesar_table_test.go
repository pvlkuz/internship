package transformer

import (
	"strings"
	"testing"
)

type TestCaesar struct {
	input, expected string
	caesarshift     int
	ioinput         bool
}

var TestArrayCaesar = []TestCaesar{
	TestCaesar{"zab", "abc", 1, false},
	TestCaesar{"abc", "xyz", -3, false},
	TestCaesar{"abc", "xyz", -3, false},
	TestCaesar{"abc1", "xyz", -3, true},
}

func TestTableCaesar(t *testing.T) {

	for _, test := range TestArrayCaesar {

		result, err := NewCaesarTransformer(test.caesarshift).Transform(strings.NewReader(test.input), test.ioinput)
		if err != nil {
			t.Errorf("Error transforming")
		}

		if result != test.expected {
			t.Errorf("Error: result = %q, expected = %q", result, test.expected)
		}

	}
}
