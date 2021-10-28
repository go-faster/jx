package jir

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_write_null(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 4096)
	stream.WriteNil()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("null", buf.String())
}

func Test_decode_null_object_field(t *testing.T) {
	iter := ParseString(Default, `[null,"a"]`)
	iter.Elem()
	if iter.ReadField() != "" {
		t.FailNow()
	}
	iter.Elem()
	if iter.String() != "a" {
		t.FailNow()
	}
}

func Test_decode_null_array_element(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `[null,"a"]`)
	should.True(iter.Elem())
	should.True(iter.ReadNil())
	should.True(iter.Elem())
	should.Equal("a", iter.String())
}

func Test_decode_null_string(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `[null,"a"]`)
	should.True(iter.Elem())
	should.Equal("", iter.String())
	should.True(iter.Elem())
	should.Equal("a", iter.String())
}

func Test_decode_null_skip(t *testing.T) {
	iter := ParseString(Default, `[null,"a"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if iter.String() != "a" {
		t.FailNow()
	}
}
