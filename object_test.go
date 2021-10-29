package jx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty_object(t *testing.T) {
	iter := ParseString(Default, `{}`)
	require.NoError(t, iter.Object(func(iter *Iterator, field string) error {
		t.Error("should not call")
		return nil
	}))
}

func Test_one_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `{"a": "stream"}`)
	should.NoError(iter.Object(func(iter *Iterator, field string) error {
		should.Equal("a", field)
		return iter.Skip()
	}))

}

func Test_write_object(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Config{IndentionStep: 2}.API(), buf, 4096)
	stream.WriteObjectStart()
	stream.WriteObjectField("hello")
	stream.WriteInt(1)
	stream.WriteMore()
	stream.WriteObjectField("world")
	stream.WriteInt(2)
	stream.WriteObjectEnd()
	require.NoError(t, stream.Flush())
	should.Nil(stream.Error)
	should.Equal("{\n  \"hello\": 1,\n  \"world\": 2\n}", buf.String())
}
