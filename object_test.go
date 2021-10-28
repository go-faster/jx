package jir

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty_object(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `{}`)
	field := iter.Field()
	should.Equal("", field)
	iter = ParseString(Default, `{}`)
	iter.Object(func(iter *Iterator, field string) bool {
		should.FailNow("should not call")
		return true
	})
}

func Test_one_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `{"a": "stream"}`)
	field := iter.Field()
	should.Equal("a", field)
	value := iter.Str()
	should.Equal("stream", value)
	field = iter.Field()
	should.Equal("", field)
	iter = ParseString(Default, `{"a": "stream"}`)
	should.True(iter.Object(func(iter *Iterator, field string) bool {
		should.Equal("a", field)
		iter.Skip()
		return true
	}))

}

func Test_two_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `{ "a": "stream" , "c": "d" }`)
	field := iter.Field()
	should.Equal("a", field)
	value := iter.Str()
	should.Equal("stream", value)
	field = iter.Field()
	should.Equal("c", field)
	value = iter.Str()
	should.Equal("d", value)
	field = iter.Field()
	should.Equal("", field)
	iter = ParseString(Default, `{"field1": "1", "field2": 2}`)
	for field := iter.Field(); field != ""; field = iter.Field() {
		switch field {
		case "field1":
			iter.Str()
		case "field2":
			iter.Int64()
		default:
			iter.ReportError("bind object", "unexpected field")
		}
	}
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
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("{\n  \"hello\": 1,\n  \"world\": 2\n}", buf.String())
}
