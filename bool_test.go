package jx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_true(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `true`)
	should.True(iter.Bool())
}

func Test_false(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `false`)
	should.False(iter.Bool())
}

func Test_write_true_false(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 4096)
	stream.True()
	stream.False()
	stream.Bool(false)
	should.NoError(stream.Flush())
	should.Equal("truefalsefalse", buf.String())
}
