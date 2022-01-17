package jx

import (
	"fmt"
	"io"
	"strings"
	"testing"

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
}

func Benchmark_appendRune(b *testing.B) {
	b.ReportAllocs()
	buf := make([]byte, 0, 4)
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf = appendRune(buf, 'f')
	}
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
