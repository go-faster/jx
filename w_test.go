package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriter_Reset(t *testing.T) {
	w := GetWriter()
	defer PutWriter(w)

	w.True()
	require.NotEmpty(t, w.Buf)
	w.Reset()
	require.Empty(t, w.Buf)
}

func TestWriter_String(t *testing.T) {
	w := GetWriter()
	defer PutWriter(w)

	w.True()
	require.Equal(t, "true", w.String())
}
