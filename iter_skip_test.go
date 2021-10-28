package jir

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_skip_number_in_array(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `[-0.12, "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	should.Nil(iter.Error)
	should.Equal("stream", iter.ReadString())
}

func Test_skip_string_in_array(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `["hello", "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	should.Nil(iter.Error)
	should.Equal("stream", iter.ReadString())
}

func Test_skip_null(t *testing.T) {
	iter := ParseString(Default, `[null , "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_true(t *testing.T) {
	iter := ParseString(Default, `[true , "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_false(t *testing.T) {
	iter := ParseString(Default, `[false , "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_array(t *testing.T) {
	iter := ParseString(Default, `[[1, [2, [3], 4]], "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_empty_array(t *testing.T) {
	iter := ParseString(Default, `[ [ ], "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_nested(t *testing.T) {
	iter := ParseString(Default, `[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_and_return_bytes(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`)
	iter.ReadArray()
	skipped := iter.SkipAndReturnBytes()
	should.Equal(`{"a" : [{"stream": "c"}], "d": 102 }`, string(skipped))
}

func Test_skip_and_return_bytes_with_reader(t *testing.T) {
	should := require.New(t)
	iter := Parse(Default, bytes.NewBufferString(`[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`), 4)
	iter.ReadArray()
	skipped := iter.SkipAndReturnBytes()
	should.Equal(`{"a" : [{"stream": "c"}], "d": 102 }`, string(skipped))
}

func Test_append_skip_and_return_bytes_with_reader(t *testing.T) {
	should := require.New(t)
	iter := Parse(Default, bytes.NewBufferString(`[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`), 4)
	iter.ReadArray()
	buf := make([]byte, 0, 1024)
	buf = iter.SkipAndAppendBytes(buf)
	should.Equal(`{"a" : [{"stream": "c"}], "d": 102 }`, string(buf))
}
