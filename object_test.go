package jir

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty_object(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `{}`)
	field := iter.ReadField()
	should.Equal("", field)
	iter = ParseString(Default, `{}`)
	iter.ReadObject(func(iter *Iterator, field string) bool {
		should.FailNow("should not call")
		return true
	})
}

func Test_one_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `{"a": "stream"}`)
	field := iter.ReadField()
	should.Equal("a", field)
	value := iter.String()
	should.Equal("stream", value)
	field = iter.ReadField()
	should.Equal("", field)
	iter = ParseString(Default, `{"a": "stream"}`)
	should.True(iter.ReadObject(func(iter *Iterator, field string) bool {
		should.Equal("a", field)
		iter.Skip()
		return true
	}))

}

func Test_two_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `{ "a": "stream" , "c": "d" }`)
	field := iter.ReadField()
	should.Equal("a", field)
	value := iter.String()
	should.Equal("stream", value)
	field = iter.ReadField()
	should.Equal("c", field)
	value = iter.String()
	should.Equal("d", value)
	field = iter.ReadField()
	should.Equal("", field)
	iter = ParseString(Default, `{"field1": "1", "field2": 2}`)
	for field := iter.ReadField(); field != ""; field = iter.ReadField() {
		switch field {
		case "field1":
			iter.String()
		case "field2":
			iter.ReadInt64()
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
