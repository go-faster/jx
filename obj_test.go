package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty_object(t *testing.T) {
	iter := DecodeStr(`{}`)
	require.NoError(t, iter.Obj(func(iter *Decoder, field string) error {
		t.Error("should not call")
		return nil
	}))
}

func Test_one_field(t *testing.T) {
	should := require.New(t)
	d := DecodeStr(`{"a": "stream"}`)
	should.NoError(d.Obj(func(iter *Decoder, field string) error {
		should.Equal("a", field)
		return iter.Skip()
	}))
}

func Test_write_object(t *testing.T) {
	should := require.New(t)
	e := NewEncoder()
	e.SetIdent(2)
	e.ObjStart()
	e.ObjField("hello")
	e.Int(1)
	e.More()
	e.ObjField("world")
	e.Int(2)
	e.ObjEnd()
	should.Equal("{\n  \"hello\": 1,\n  \"world\": 2\n}", e.String())
}
