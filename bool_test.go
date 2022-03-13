package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrue(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`true`)
	should.True(iter.Bool())
}

func TestFalse(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`false`)
	should.False(iter.Bool())
}

func TestWriteTrueFalse(t *testing.T) {
	should := require.New(t)
	w := GetEncoder()
	w.Bool(true)
	w.Bool(false)
	w.Bool(false)
	should.Equal("truefalsefalse", string(w.Bytes()))
}
