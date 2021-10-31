package jx

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_Arr(t *testing.T) {
	t.Run("Blank", func(t *testing.T) {
		d := DecodeStr(`[]`)
		ok, err := d.Elem()
		require.NoError(t, err)
		require.False(t, ok)
	})
	t.Run("Invalid", func(t *testing.T) {
		d := DecodeStr(`{`)
		ok, err := d.Elem()
		require.Error(t, err)
		require.False(t, ok)
	})
	t.Run("EOF", func(t *testing.T) {
		d := DecodeStr("")
		ok, err := d.Elem()
		require.ErrorIs(t, err, io.EOF)
		require.False(t, ok)
	})
	t.Run("ErrUnexpectedEOF", func(t *testing.T) {
		d := DecodeStr("[")
		ok, err := d.Elem()
		require.ErrorIs(t, err, io.ErrUnexpectedEOF)
		require.False(t, ok)
	})
}
