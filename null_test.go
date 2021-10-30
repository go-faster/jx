package jx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_write_null(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewEncoder(buf, 4096)
	stream.Null()
	should.NoError(stream.Flush())
	should.Equal("null", buf.String())
}

func Test_decode_null_array_element(t *testing.T) {
	should := require.New(t)
	iter := DecodeString(`[null,"a"]`)
	should.True(iter.Elem())
	should.NoError(iter.Null())
	should.True(iter.Elem())
	s, err := iter.String()
	should.NoError(err)
	should.Equal("a", s)
}

func Test_decode_null_string(t *testing.T) {
	should := require.New(t)
	iter := DecodeString(`[null,"a"]`)
	should.True(iter.Elem())
	should.NoError(iter.Null())
	should.True(iter.Elem())
	s, err := iter.String()
	should.NoError(err)
	should.Equal("a", s)
}

func Test_decode_null_skip(t *testing.T) {
	iter := DecodeString(`[null,"a"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "a" {
		t.FailNow()
	}
}
