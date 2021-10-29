package jx

import (
	"bytes"
	hexEnc "encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_read_string(t *testing.T) {
	badInputs := []string{
		``,
		`null`,
		`"`,
		`"\"`,
		`"\\\"`,
		"\"\n\"",
		`"\U0001f64f"`,
		`"\uD83D\u00"`,
	}
	for i := 0; i < 32; i++ {
		// control characters are invalid
		badInputs = append(badInputs, string([]byte{'"', byte(i), '"'}))
	}

	for _, input := range badInputs {
		i := ReadString(input)
		_, err := i.Str()
		assert.Error(t, err, "input: %q", input)
	}

	goodInputs := []struct {
		input       string
		expectValue string
	}{
		{`""`, ""},
		{`"a"`, "a"},
		{`"Iñtërnâtiônàlizætiøn,💝🐹🌇⛔"`, "Iñtërnâtiônàlizætiøn,💝🐹🌇⛔"},
		{`"\uD83D"`, string([]byte{239, 191, 189})},
		{`"\uD83D\\"`, string([]byte{239, 191, 189, '\\'})},
		{`"\uD83D\ub000"`, string([]byte{239, 191, 189, 235, 128, 128})},
		{`"\uD83D\ude04"`, "😄"},
		{`"\uDEADBEEF"`, string([]byte{239, 191, 189, 66, 69, 69, 70})},
		{`"hel\"lo"`, `hel"lo`},
		{`"hel\\\/lo"`, `hel\/lo`},
		{`"hel\\blo"`, `hel\blo`},
		{`"hel\\\blo"`, "hel\\\blo"},
		{`"hel\\nlo"`, `hel\nlo`},
		{`"hel\\\nlo"`, "hel\\\nlo"},
		{`"hel\\tlo"`, `hel\tlo`},
		{`"hel\\flo"`, `hel\flo`},
		{`"hel\\\flo"`, "hel\\\flo"},
		{`"hel\\\rlo"`, "hel\\\rlo"},
		{`"hel\\\tlo"`, "hel\\\tlo"},
		{`"\u4e2d\u6587"`, "中文"},
		{`"\ud83d\udc4a"`, "\xf0\x9f\x91\x8a"},
	}

	for _, tc := range goodInputs {
		testReadString(t, tc.input, tc.expectValue, false, "json.Unmarshal", json.Unmarshal)

		i := ReadString(tc.input)
		s, err := i.Str()
		assert.NoError(t, err)
		assert.Equal(t, tc.expectValue, s)
	}
}

func testReadString(t *testing.T, input string, expectValue string, expectError bool, marshalerName string, marshaler func([]byte, interface{}) error) {
	t.Helper()
	var value string
	err := marshaler([]byte(input), &value)
	if expectError != (err != nil) {
		t.Errorf("%q: %s: expected error %v, got %v", input, marshalerName, expectError, err)
		return
	}
	if value != expectValue {
		t.Errorf("%q: %s: expected %q, got %q", input, marshalerName, expectValue, value)
		return
	}
}

func TestIter_Str(t *testing.T) {
	for _, tt := range []struct {
		Name  string
		Input string
	}{
		{Name: "\\x00", Input: "\x00"},
		{Name: "\\x00TrailingSpace", Input: "\x00 "},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			s := NewWriter(buf, 128)
			t.Logf("%v", []rune(tt.Input))

			s.Str(tt.Input)
			require.NoError(t, s.Flush())
			t.Logf("%v", []rune(buf.String()))

			// Check `encoding/json` compatibility.
			var gotStd string
			requireCompat(t, buf.Bytes(), tt.Input)
			require.NoError(t, json.Unmarshal(buf.Bytes(), &gotStd))
			require.Equal(t, tt.Input, gotStd)

			i := ReadBytes(buf.Bytes())
			got, err := i.Str()
			require.NoError(t, err)
			require.Equal(t, tt.Input, got, "%s\n%s", buf, hexEnc.Dump(buf.Bytes()))
		})
	}
}
