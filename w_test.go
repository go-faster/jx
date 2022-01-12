package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriter_Reset(t *testing.T) {
	var w Writer
	w.True()
	require.NotEmpty(t, w.Buf)
	w.Reset()
	require.Empty(t, w.Buf)
}

func TestWriter_String(t *testing.T) {
	w := Writer{Buf: []byte(`true`)}
	require.Equal(t, "true", w.String())
}
