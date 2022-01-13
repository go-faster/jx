package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkip_number_in_array(t *testing.T) {
	var err error
	a := require.New(t)
	d := DecodeStr(`[-0.12, "stream"]`)
	_, err = d.Elem()
	a.NoError(err)
	err = d.Skip()
	a.NoError(err)
	_, err = d.Elem()
	a.NoError(err)
	if s, _ := d.Str(); s != "stream" {
		t.FailNow()
	}
}

func TestSkip_string_in_array(t *testing.T) {
	d := DecodeStr(`["hello", "stream"]`)
	d.Elem()
	d.Skip()
	d.Elem()
	if s, _ := d.Str(); s != "stream" {
		t.FailNow()
	}
}

func TestSkip_null(t *testing.T) {
	d := DecodeStr(`[null , "stream"]`)
	d.Elem()
	d.Skip()
	d.Elem()
	if s, _ := d.Str(); s != "stream" {
		t.FailNow()
	}
}

func TestSkip_true(t *testing.T) {
	d := DecodeStr(`[true , "stream"]`)
	d.Elem()
	d.Skip()
	d.Elem()
	if s, _ := d.Str(); s != "stream" {
		t.FailNow()
	}
}

func TestSkip_false(t *testing.T) {
	d := DecodeStr(`[false , "stream"]`)
	d.Elem()
	d.Skip()
	d.Elem()
	if s, _ := d.Str(); s != "stream" {
		t.FailNow()
	}
}

func TestSkip_array(t *testing.T) {
	d := DecodeStr(`[[1, [2, [3], 4]], "stream"]`)
	d.Elem()
	d.Skip()
	d.Elem()
	if s, _ := d.Str(); s != "stream" {
		t.FailNow()
	}
}

func TestSkip_empty_array(t *testing.T) {
	d := DecodeStr(`[ [ ], "stream"]`)
	d.Elem()
	d.Skip()
	d.Elem()
	if s, _ := d.Str(); s != "stream" {
		t.FailNow()
	}
}

func TestSkip_nested(t *testing.T) {
	d := DecodeStr(`[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`)
	if _, err := d.Elem(); err != nil {
		t.Fatal(err)
	}
	require.NoError(t, d.Skip())
	if _, err := d.Elem(); err != nil {
		t.Fatal(err)
	}
	s, err := d.Str()
	require.NoError(t, err)
	require.Equal(t, "stream", s)
}

func TestSkip_simple_nested(t *testing.T) {
	d := DecodeStr(`["foo", "bar", "baz"]`)
	require.NoError(t, d.Skip())
}

func TestDecoder_Bool(t *testing.T) {
	for _, s := range []string{
		"tru",
		"fals",
		"",
		"nope",
	} {
		d := DecodeStr(s)
		v, err := d.Bool()
		require.False(t, v)
		require.Error(t, err)
	}
}

func TestDecoder_Null(t *testing.T) {
	for _, s := range []string{
		"",
		"nope",
		"nul",
		"nil",
	} {
		d := DecodeStr(s)
		require.Error(t, d.Null())
	}
}

func TestDecoder_Skip(t *testing.T) {
	for _, s := range []string{
		"",
		"nope",
		"nul",
		"nil",
		"tru",
		"fals",
		"1.2.3",
	} {
		d := DecodeStr(s)
		require.Error(t, d.Skip())
	}
}
