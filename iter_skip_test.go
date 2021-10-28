package jir

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_skip_number_in_array(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `[-0.12, "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	should.Nil(iter.Error)
	should.Equal("stream", iter.Str())
}

func Test_skip_string_in_array(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `["hello", "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	should.Nil(iter.Error)
	should.Equal("stream", iter.Str())
}

func Test_skip_null(t *testing.T) {
	iter := ParseString(Default, `[null , "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if iter.Str() != "stream" {
		t.FailNow()
	}
}

func Test_skip_true(t *testing.T) {
	iter := ParseString(Default, `[true , "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if iter.Str() != "stream" {
		t.FailNow()
	}
}

func Test_skip_false(t *testing.T) {
	iter := ParseString(Default, `[false , "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if iter.Str() != "stream" {
		t.FailNow()
	}
}

func Test_skip_array(t *testing.T) {
	iter := ParseString(Default, `[[1, [2, [3], 4]], "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if iter.Str() != "stream" {
		t.FailNow()
	}
}

func Test_skip_empty_array(t *testing.T) {
	iter := ParseString(Default, `[ [ ], "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if iter.Str() != "stream" {
		t.FailNow()
	}
}

func Test_skip_nested(t *testing.T) {
	iter := ParseString(Default, `[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if iter.Str() != "stream" {
		t.FailNow()
	}
}
