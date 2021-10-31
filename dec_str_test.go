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
	e := &UnexpectedTokenErr{
		Token: 'c',
	}
	s := error(e).Error()
	require.Equal(t, "unexpected byte 99 'c'", s)
}
