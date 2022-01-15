package jx

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"
	"unicode/utf8"

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
	testStr := func(d *Decoder, input string, valid bool) func(t *testing.T) {
		return func(t *testing.T) {
			t.Cleanup(func() {
				if t.Failed() {
					t.Logf("Input: %q", input)
				}
			})

			_, err := d.Str()
			if valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}
	}
	for i, input := range testStrings {
		valid := json.Valid([]byte(input))

		t.Run(fmt.Sprintf("Test%d", i), func(t *testing.T) {
			t.Run("Buffer", testStr(DecodeStr(input), input, valid))

			r := strings.NewReader(input)
			d := Decode(r, 512)
			t.Run("Reader", testStr(d, input, valid))

			r.Reset(input)
			obr := iotest.OneByteReader(r)
			d.Reset(obr)
			t.Run("OneByteReader", testStr(d, input, valid))
		})
	}
}

func Benchmark_appendRune(b *testing.B) {
	b.ReportAllocs()
	buf := make([]byte, 0, 4)
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf = appendRune(buf, 'f')
	}
}

func benchmarkDecoderStrBytes(str string) func(b *testing.B) {
	return func(b *testing.B) {
		e := GetEncoder()
		e.Str(str)
		data := e.Bytes()

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
	runBench := func(char string) func(b *testing.B) {
		return func(b *testing.B) {
			for _, size := range []int{
				2, 8, 16, 64, 128, 1024,
			} {
				count := utf8.RuneCountInString(char)
				b.Run(fmt.Sprintf("%db", size), benchmarkDecoderStrBytes(strings.Repeat(char, size/count)))
			}
		}
	}

	b.Run("Plain", runBench("a"))
	b.Run("Escaped", runBench("ф"))
}
