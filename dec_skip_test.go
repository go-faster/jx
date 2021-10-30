package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_skip_number_in_array(t *testing.T) {
	iter := DecodeStr(`[-0.12, "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "stream" {
		t.FailNow()
	}
}

func Test_skip_string_in_array(t *testing.T) {
	iter := DecodeStr(`["hello", "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "stream" {
		t.FailNow()
	}
}

func Test_skip_null(t *testing.T) {
	iter := DecodeStr(`[null , "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "stream" {
		t.FailNow()
	}
}

func Test_skip_true(t *testing.T) {
	iter := DecodeStr(`[true , "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "stream" {
		t.FailNow()
	}
}

func Test_skip_false(t *testing.T) {
	iter := DecodeStr(`[false , "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "stream" {
		t.FailNow()
	}
}

func Test_skip_array(t *testing.T) {
	iter := DecodeStr(`[[1, [2, [3], 4]], "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "stream" {
		t.FailNow()
	}
}

func Test_skip_empty_array(t *testing.T) {
	iter := DecodeStr(`[ [ ], "stream"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.String(); s != "stream" {
		t.FailNow()
	}
}

func Test_skip_nested(t *testing.T) {
	iter := DecodeStr(`[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`)
	if _, err := iter.Elem(); err != nil {
		t.Fatal(err)
	}
	require.NoError(t, iter.Skip())
	if _, err := iter.Elem(); err != nil {
		t.Fatal(err)
	}
	s, err := iter.String()
	require.NoError(t, err)
	require.Equal(t, "stream", s)
}

func Test_skip_simple_nested(t *testing.T) {
	iter := DecodeStr(`["foo", "bar", "baz"]`)
	require.NoError(t, iter.Skip())
}
