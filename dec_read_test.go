package jx

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_readAtLeast(t *testing.T) {
	a := require.New(t)
	d := Decode(strings.NewReader("aboba"), 1)
	a.NoError(d.readAtLeast(4))
	a.Equal(d.buf[d.head:d.tail], []byte("abob"))
}

func TestDecoder_consume(t *testing.T) {
	r := errReader{}
	d := Decode(r, 1)
	require.ErrorIs(t, d.consume('"'), r.Err())
}

func TestDecoder_Next(t *testing.T) {
	d := DecodeBytes(nil)
	d.Next()
	d.Next()
	d.Next()

	_, err := d.Str()
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
