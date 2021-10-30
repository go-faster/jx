package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_invalid_float(t *testing.T) {
	inputs := []string{
		`1.e1`, // dot without following digit
		`1.`,   // dot can not be the last char
		``,     // empty number
		`01`,   // extra leading zero
		`-`,    // negative without digit
		`--`,   // double negative
		`--2`,  // double negative
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			should := require.New(t)
			iter := DecodeStr(input + ",")
			should.Error(iter.Skip())
			iter = DecodeStr(input + ",")
			_, err := iter.Float64()
			should.Error(err)
			iter = DecodeStr(input + ",")
			_, err = iter.Float32()
			should.Error(err)
		})
	}
}

func TestValid_positive(t *testing.T) {
	for _, tt := range []struct {
		Name  string
		Value string
	}{
		// https://github.com/json-iterator/go/issues/520
		{Name: "Number", Value: "1"},

		{Name: "BlankObj", Value: "{}"},
		{Name: "BlankArr", Value: "[]"},
		{Name: "Example", Value: `{"menu": {
  "id": "file",
  "value": "File",
  "popup": {
    "menuitem": [
      {"value": "New", "onclick": "CreateNewDoc()"},
      {"value": "Open", "onclick": "OpenDoc()"},
      {"value": "Close", "onclick": "CloseDoc()"}
    ]
  }
}}`},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			require.True(t, Valid([]byte(tt.Value)), "should be valid")
		})
	}
}

func TestValid_negative(t *testing.T) {
	for _, tt := range []struct {
		Name  string
		Value string
	}{
		{Name: "Blank", Value: ""},
		{Name: "InvalidCharacters", Value: "foo"},
		{Name: "NotClosed", Value: "{"},
		{Name: "NotClosed2", Value: "{{}"},
		{Name: "NotOpened", Value: "{}}"},
		{Name: "NotOpenedArr", Value: "{[}]"},
		{Name: "NotOpenedArr2", Value: "{}]"},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			require.False(t, Valid([]byte(tt.Value)), "should be invalid")
		})
	}
}
