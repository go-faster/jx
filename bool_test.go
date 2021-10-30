package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_true(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`true`)
	should.True(iter.Bool())
}

func Test_false(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`false`)
	should.False(iter.Bool())
}

func Test_write_true_false(t *testing.T) {
	should := require.New(t)
	w := NewEncoder()
	w.True()
	w.False()
	w.Bool(false)
	should.Equal("truefalsefalse", string(w.Bytes()))
}
