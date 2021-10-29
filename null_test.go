package jx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_write_null(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewWriter(buf, 4096)
	stream.Null()
	should.NoError(stream.Flush())
	should.Equal("null", buf.String())
}

func Test_decode_null_array_element(t *testing.T) {
	should := require.New(t)
	iter := ReadString(`[null,"a"]`)
	should.True(iter.Elem())
	should.NoError(iter.Null())
	should.True(iter.Elem())
	s, err := iter.Str()
	should.NoError(err)
	should.Equal("a", s)
}

func Test_decode_null_string(t *testing.T) {
	should := require.New(t)
	iter := ReadString(`[null,"a"]`)
	should.True(iter.Elem())
	should.NoError(iter.Null())
	should.True(iter.Elem())
	s, err := iter.Str()
	should.NoError(err)
	should.Equal("a", s)
}

func Test_decode_null_skip(t *testing.T) {
	iter := ReadString(`[null,"a"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.Str(); s != "a" {
		t.FailNow()
	}
}
