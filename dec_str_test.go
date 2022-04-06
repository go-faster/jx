package jx

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoder_StrAppend(t *testing.T) {
	s := `"Hello"`
	d := DecodeStr(s)
	var (
		data []byte
		err  error
	)
	data, err = d.StrAppend(data)
	require.NoError(t, err)
	require.Equal(t, "Hello", string(data))

	_, err = d.StrAppend(data)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestUnexpectedTokenErr_Error(t *testing.T) {
	e := &badTokenErr{
		Token: 'c',
	}
	s := error(e).Error()
	require.Equal(t, "unexpected byte 99 'c'", s)
}

func TestDecoder_Str(t *testing.T) {
	runTestCases(t, testStrings, func(t *testing.T, d *Decoder) error {
		_, err := d.Str()
		return err
	})

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
		i := DecodeStr(input)
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

		i := DecodeStr(tc.input)
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

func TestDecoder_strSlow(t *testing.T) {
	r := errReader{}
	d := Decode(r, 1)
	_, err := d.strSlow(value{})
	require.ErrorIs(t, err, r.Err())
}

func Benchmark_appendRune(b *testing.B) {
	b.ReportAllocs()
	buf := make([]byte, 0, 4)
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf = appendRune(buf, 'f')
	}
}

func BenchmarkDecoder_escapedChar(b *testing.B) {
	bench := func(char byte, data []byte) func(b *testing.B) {
		return func(b *testing.B) {
			d := DecodeBytes(data)
			v := value{buf: make([]byte, 0, 16)}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				d.ResetBytes(data)
				v.buf = v.buf[:0]
				if _, err := d.escapedChar(v, char); err != nil {
					b.Fatal(err)
				}
			}
		}
	}
	b.Run("Unicode", bench('u', []byte(`000c`)))
	b.Run("Newline", bench('n', nil))
}

func benchmarkDecoderStrBytes(data []byte) func(b *testing.B) {
	return func(b *testing.B) {
		d := GetDecoder()

		b.SetBytes(int64(len(data)))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			d.ResetBytes(data)
			if _, err := d.StrBytes(); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkDecoder_StrBytes(b *testing.B) {
	runBench := func(char string, maxCount int) func(b *testing.B) {
		return func(b *testing.B) {
			for _, count := range []int{
				1, 8, 16, 64, 128, 1024,
			} {
				if maxCount > 0 && count >= maxCount {
					break
				}
				e := GetEncoder()
				str := strings.Repeat(char, count)
				e.StrEscape(str)
				data := e.Bytes()

				b.Run(fmt.Sprintf("%db", len(data)-2), benchmarkDecoderStrBytes(data))
			}
		}
	}

	b.Run("Plain", runBench("a", -1))
	b.Run("EscapedNewline", runBench("\n", -1))
	b.Run("EscapedUnicode", runBench("\f", -1))
	b.Run("Mixed", runBench("aaaa\naaaa\faaaaa", 64))
}
