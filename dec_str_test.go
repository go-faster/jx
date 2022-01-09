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
	for _, s := range []string{
		`"foo`,
		`"\u1d`,
		`"\u$`,
		`"\21412`,
		`"\`,
		`"\u1337`,
		`"\uD834\1`,
		`"\uD834\u3`,
		`"\uD834\`,
		`"\uD834`,
		`"\u07F9`,
	} {
		d := DecodeStr(s)
		_, err := d.Str()
		require.Error(t, err)
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
	for _, size := range []int{
		1, 8, 16, 64, 128, 1024,
	} {
		b.Run(fmt.Sprintf("%db", size), benchmarkDecoderStrBytes(strings.Repeat("a", size)))
	}
}
