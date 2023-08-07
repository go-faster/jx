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

func TestWriter_Grow(t *testing.T) {
	should := require.New(t)
	e := &Writer{}
	should.Equal(0, len(e.Buf))
	should.Equal(0, cap(e.Buf))
	e.Grow(1024)
	should.Equal(0, len(e.Buf))
	should.Equal(1024, cap(e.Buf))
	e.Grow(512)
	should.Equal(0, len(e.Buf))
	should.Equal(1024, cap(e.Buf))
	e.Grow(4096)
	should.Equal(0, len(e.Buf))
	should.Equal(4096, cap(e.Buf))
}
