package jx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty_object(t *testing.T) {
	iter := ParseString(Default, `{}`)
	require.NoError(t, iter.Object(func(iter *Iter, field string) error {
		t.Error("should not call")
		return nil
	}))
}

func Test_one_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `{"a": "stream"}`)
	should.NoError(iter.Object(func(iter *Iter, field string) error {
		should.Equal("a", field)
		return iter.Skip()
	}))
}

func Test_write_object(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Config{IndentionStep: 2}.API(), buf, 4096)
	stream.ObjStart()
	stream.ObjField("hello")
	stream.WriteInt(1)
	stream.More()
	stream.ObjField("world")
	stream.WriteInt(2)
	stream.ObjEnd()
	should.NoError(stream.Flush())
	should.Equal("{\n  \"hello\": 1,\n  \"world\": 2\n}", buf.String())
}
