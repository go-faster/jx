package jx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_true(t *testing.T) {
	should := require.New(t)
	iter := DecodeString(`true`)
	should.True(iter.Bool())
}

func Test_false(t *testing.T) {
	should := require.New(t)
	iter := DecodeString(`false`)
	should.False(iter.Bool())
}

func Test_write_true_false(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	w := NewEncoder(buf, 4096)
	w.True()
	w.False()
	w.Bool(false)
	should.NoError(w.Flush())
	should.Equal("truefalsefalse", buf.String())
}
