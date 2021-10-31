package jx

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_Arr(t *testing.T) {
	t.Run("Blank", func(t *testing.T) {
		d := DecodeStr(`[]`)
		require.NoError(t, d.Arr(nil))
	})
	t.Run("Invalid", func(t *testing.T) {
		d := DecodeStr(`{`)
		require.Error(t, d.Arr(nil))
	})
	t.Run("ErrUnexpectedEOF", func(t *testing.T) {
		d := DecodeStr("")
		require.ErrorIs(t, d.Arr(nil), io.ErrUnexpectedEOF)
	})
	t.Run("ErrUnexpectedEOF", func(t *testing.T) {
		d := DecodeStr("[")
		require.ErrorIs(t, d.Arr(nil), io.ErrUnexpectedEOF)
	})
}

func TestDecoder_Elem(t *testing.T) {
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