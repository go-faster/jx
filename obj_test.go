package jx

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyObject(t *testing.T) {
	iter := DecodeStr(`{}`)
	require.NoError(t, iter.Obj(func(iter *Decoder, field string) error {
		t.Error("should not call")
		return nil
	}))
}

func TestOneField(t *testing.T) {
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
	e.FieldStart("hello")
	e.Int(1)
	e.FieldStart("world")
	e.Int(2)
	e.FieldStart("obj")
	e.ObjStart()
	e.FieldStart("a")
	e.Str("b")
	e.ObjEnd()
	e.FieldStart("data")
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

	t.Run("Std", func(t *testing.T) {
		b := new(bytes.Buffer)
		enc := json.NewEncoder(b)
		enc.SetIndent("", "  ")
		require.NoError(t, enc.Encode(struct {
			Hello int `json:"hello"`
			World int `json:"world"`
			Obj   struct {
				A string `json:"a"`
			} `json:"obj"`
			Data []int `json:"data"`
		}{
			Hello: 1,
			World: 2,
			Obj: struct {
				A string `json:"a"`
			}{A: "b"},
			Data: []int{1, 2},
		}))

		// Remove trialing newline from expected.
		exp := b.String()
		exp = strings.TrimRight(exp, "\n")

		require.Equal(t, exp, e.String())
	})
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
