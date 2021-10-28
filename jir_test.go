package jir

import (
	"bytes"
	hexEnc "encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	_ = Default
	_ = Fastest
}

func Test_parseVal(t *testing.T) {
	t.Run("Object", func(t *testing.T) {
		var v Value
		const input = `{"foo":{"bar":1,"baz":[1,2,3.14],"200":null}}`
		i := ParseString(Default, input)
		parseVal(i, &v)
		assert.NoError(t, i.Error)
		assert.Equal(t, `{foo: {bar: 1, baz: [1, 2, f3.14], 200: null}}`, v.String())

		buf := new(bytes.Buffer)
		s := NewStream(Default, buf, 1024)
		v.Write(s)
		require.NoError(t, s.Flush())
		require.Equal(t, input, buf.String(), "encoded value should equal to input")

	})
	t.Run("Inputs", func(t *testing.T) {
		for _, tt := range []struct {
			Input string
		}{
			{Input: "1"},
			{Input: "0.0"},
		} {
			t.Run(tt.Input, func(t *testing.T) {
				var v Value
				input := []byte(tt.Input)
				i := ParseBytes(Default, input)
				parseVal(i, &v)
				if i.Error != nil && i.Error != io.EOF {
					t.Fatal(i.Error)
				}

				buf := new(bytes.Buffer)
				s := NewStream(Default, buf, 1024)
				v.Write(s)
				require.NoError(t, s.Flush())
				require.Equal(t, tt.Input, buf.String(), "encoded value should equal to input")

				var otherValue Value
				i.ResetBytes(buf.Bytes())
				parseVal(i, &otherValue)
				if i.Error != nil && i.Error != io.EOF {
					t.Log(hexEnc.Dump(input))
					t.Log(hexEnc.Dump(buf.Bytes()))
				}
			})

		}
	})
}

type ValType byte

const (
	ValInvalid ValType = iota
	ValStr
	ValInt
	ValFloat
	ValNull
	ValObj
	ValArr
	ValBool
)

// Value represents any json value as sum type.
type Value struct {
	Type   ValType
	Str    string
	Int    int64
	Float  float64
	Key    string
	KeySet bool // Key can be ""
	Bool   bool
	Child  []Value
}

// Write json representation of Value to Stream.
func (v Value) Write(s *Stream) {
	if v.KeySet {
		s.WriteObjectField(v.Key)
	}
	switch v.Type {
	case ValStr:
		s.WriteString(v.Str)
	case ValFloat:
		s.WriteFloat64(v.Float)
	case ValInt:
		s.WriteInt64(v.Int)
	case ValBool:
		s.WriteBool(v.Bool)
	case ValNull:
		s.WriteNil()
	case ValArr:
		s.WriteArrayStart()
		for i, c := range v.Child {
			if i != 0 {
				s.WriteMore()
			}
			c.Write(s)
		}
		s.WriteArrayEnd()
	case ValObj:
		s.WriteObjectStart()
		for i, c := range v.Child {
			if i != 0 {
				s.WriteMore()
			}
			c.Write(s)
		}
		s.WriteObjectEnd()
	default:
		panic(v.Type)
	}
}

func (v Value) String() string {
	var b strings.Builder
	if v.KeySet {
		if v.Key == "" {
			b.WriteString("<blank>")
		}
		b.WriteString(v.Key)
		b.WriteString(": ")
	}
	switch v.Type {
	case ValStr:
		b.WriteString(`"` + v.Str + `"'`)
	case ValFloat:
		b.WriteRune('f')
		b.WriteString(fmt.Sprintf("%v", v.Float))
	case ValInt:
		b.WriteString(strconv.FormatInt(v.Int, 10))
	case ValBool:
		b.WriteString(strconv.FormatBool(v.Bool))
	case ValNull:
		b.WriteString("null")
	case ValArr:
		b.WriteString("[")
		for i, c := range v.Child {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(c.String())
		}
		b.WriteString("]")
	case ValObj:
		b.WriteString("{")
		for i, c := range v.Child {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(c.String())
		}
		b.WriteString("}")
	default:
		panic(v.Type)
	}

	return b.String()
}

func parseVal(i *Iterator, v *Value) bool {
	switch i.WhatIsNext() {
	case Invalid:
		return false
	case Number:
		n := i.ReadNumber()
		idx := strings.Index(n.String(), ".")
		if (idx > 0 && idx != len(n.String())-1) || strings.Contains(n.String(), "e") {
			f, err := n.Float64()
			if err != nil {
				i.ReportError("ReadNumber", err.Error())
				return false
			}
			v.Float = f
			v.Type = ValFloat
		} else {
			f, err := n.Int64()
			if err != nil {
				i.ReportError("ReadNumber", err.Error())
				return false
			}
			v.Int = f
			v.Type = ValInt
		}
	case String:
		v.Str = i.Str()
		v.Type = ValStr
	case Nil:
		i.Null()
		v.Type = ValNull
	case Bool:
		v.Bool = i.Bool()
		v.Type = ValBool
	case Object:
		v.Type = ValObj
		return i.Object(func(i *Iterator, s string) bool {
			var elem Value
			if !parseVal(i, &elem) {
				return false
			}
			elem.Key = s
			elem.KeySet = true
			v.Child = append(v.Child, elem)
			return true
		})
	case Array:
		v.Type = ValArr
		return i.Array(func(i *Iterator) bool {
			var elem Value
			if !parseVal(i, &elem) {
				return false
			}
			v.Child = append(v.Child, elem)
			return true
		})
	default:
		panic(i.WhatIsNext())
	}
	return true
}

// requireCompat fails if `encoding/json` will encode v differently than exp.
func requireCompat(t testing.TB, exp []byte, v interface{}) {
	t.Helper()
	buf, err := json.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, exp, buf)
}
