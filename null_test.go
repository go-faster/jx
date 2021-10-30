package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_write_null(t *testing.T) {
	should := require.New(t)
	e := NewEncoder()
	e.Null()
	should.Equal("null", e.String())
}

func Test_decode_null_array_element(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`[null,"a"]`)
	should.True(iter.Elem())
	should.NoError(iter.Null())
	should.True(iter.Elem())
	s, err := iter.String()
	should.NoError(err)
	should.Equal("a", s)
}

func Test_decode_null_string(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`[null,"a"]`)
	should.True(iter.Elem())
	should.NoError(iter.Null())
	should.True(iter.Elem())
	s, err := iter.String()
	should.NoError(err)
	should.Equal("a", s)
}

func Test_decode_null_skip(t *testing.T) {
	iter := DecodeStr(`[null,"a"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "a" {
		t.FailNow()
	}
}
