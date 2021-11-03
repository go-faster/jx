package jx

import (
	"io"
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
