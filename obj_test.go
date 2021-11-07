package jx

import (
	"encoding/json"
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

func TestEncoder_SetIdent(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.SetIdent(2)
	e.ObjStart()
	e.Field("hello")
	e.Int(1)
	e.Field("world")
	e.Int(2)
	e.Field("obj")
	e.ObjStart()
	e.Field("a")
	e.Str("b")
	e.ObjEnd()
	e.Field("data")
	e.ArrStart()
	e.Int(1)
	e.Int(2)
	e.ArrEnd()
	e.ObjEnd()
	expected := `{
  "hello": 1,
  "world": 2,
  "obj": {
    "a": "b"
  },
  "data": [
    1,
    2
  ]
}`
	should.Equal(expected, e.String())
}

func TestDecoder_Obj(t *testing.T) {
	// https://github.com/json-iterator/go/issues/549
	b := []byte(`{"\u6D88\u606F":"\u6D88\u606F"}`)

	v := struct {
		Message string `json:"消息"`
	}{}
	require.NoError(t, json.Unmarshal(b, &v))
	require.Equal(t, "消息", v.Message)

	var gotKey, gotVal string
	require.NoError(t, DecodeBytes(b).Obj(func(d *Decoder, key string) error {
		str, err := d.Str()
		if err != nil {
			return err
		}
		gotKey = key
		gotVal = str
		return nil
	}))

	require.Equal(t, v.Message, gotVal)
	require.Equal(t, v.Message, gotKey)
}
