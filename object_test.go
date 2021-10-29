package jx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty_object(t *testing.T) {
	iter := ParseString(`{}`)
	require.NoError(t, iter.Object(func(iter *Iter, field string) error {
		t.Error("should not call")
		return nil
	}))
}

func Test_one_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`{"a": "stream"}`)
	should.NoError(iter.Object(func(iter *Iter, field string) error {
		should.Equal("a", field)
		return iter.Skip()
	}))
}

func Test_write_object(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	s := NewStream(buf, 4096)
	s.SetIdent(2)
	s.ObjStart()
	s.ObjField("hello")
	s.WriteInt(1)
	s.More()
	s.ObjField("world")
	s.WriteInt(2)
	s.ObjEnd()
	should.NoError(s.Flush())
	should.Equal("{\n  \"hello\": 1,\n  \"world\": 2\n}", buf.String())
}
