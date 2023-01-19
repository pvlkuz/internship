package transformer

import (
	"bytes"
	"strings"
	"testing"
)

type TestBasic struct {
	input, expected    string
	caesarshift        int
	base64use, ioinput bool
}

var TestArrayBasic = []TestBasic{
	TestBasic{"Man", "TWFu", 0, true, false},
	TestBasic{"Ma", "TWE=", 0, true, false},
	TestBasic{"Ma1", "TWE=", 0, true, true},
	TestBasic{"za", "ab", 1, false, false},
	TestBasic{"abc", "xyz", -3, false, false},
	TestBasic{"za1", "ab", 1, false, true},
	TestBasic{"12345", "54321", 0, false, false},
	TestBasic{"123450", "54321", 0, false, true},
}

func TestTableBasic(t *testing.T) {

	for _, test := range TestArrayBasic {

		buf := new(bytes.Buffer)
		err := BasicTransform(strings.NewReader(test.input), buf, test.caesarshift, test.base64use, test.ioinput)
		if err != nil {
			t.Errorf("Error transforming")
		}
		result := buf.String()

		if result != test.expected {
			t.Errorf("Error: result = %q, expected = %q", result, test.expected)
		}

	}
}
