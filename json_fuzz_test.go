//go:build go1.18

package jir

import (
	"bytes"
	hexe "encoding/hex"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzValid(f *testing.F) {
	f.Add("{}")
	f.Add(`{"foo": "bar"}`)
	f.Add(``)
	f.Add(`"foo"`)
	f.Add(`"{"`)
	f.Add(`"{}"`)
	f.Fuzz(func(t *testing.T, queryStr string) {
		Valid([]byte(queryStr))
	})
}

func iterDown(i *Iterator, count *int) bool {
	*count++
	switch i.WhatIsNext() {
	case Invalid:
		return false
	case Number:
		_ = i.ReadNumber()
	case String:
		_ = i.ReadString()
	case Nil:
		i.ReadNil()
	case Bool:
		i.ReadBool()
	case Object:
		return i.ReadObjectCB(func(i *Iterator, s string) bool {
			return iterDown(i, count)
		})
	case Array:
		return i.ReadArrayCB(func(i *Iterator) bool {
			return iterDown(i, count)
		})
	default:
		panic(i.WhatIsNext())
	}
	return i.Error == nil
}

func Test_iterDown(t *testing.T) {
	var count int
	i := ParseString(Default, `{"foo": {"bar": 1, "baz": [1, 2, 3]}}`)
	iterDown(i, &count)
	assert.NoError(t, i.Error)
	assert.Equal(t, 7, count)
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
					t.Log(hexe.Dump(input))
					t.Log(hexe.Dump(buf.Bytes()))
				}
			})

		}
	})
}

func FuzzIter(f *testing.F) {
	f.Add([]byte("{}"))
	f.Add([]byte(`"foo"`))
	f.Add([]byte(`123"`))
	f.Add([]byte(`null`))
	f.Add([]byte(`{"foo": {"bar": 1, "baz": [1, 2, 3]}}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		i := Default.Iterator(data)
		defer Default.PutIterator(i)
		var count int
		iterDown(i, &count)
	})
}

func FuzzDecEnc(f *testing.F) {
	f.Add([]byte("{}"))
	f.Add([]byte(`"foo"`))
	f.Add([]byte(`123"`))
	f.Add([]byte(`null`))
	f.Add([]byte(`{"foo": {"bar": 1, "baz": [1, 2, 3]}}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		i := Default.Iterator(nil)
		i.ResetBytes(data)
		defer Default.PutIterator(i)

		// Parsing to v.
		var v Value
		if !parseVal(i, &v) {
			t.Skip()
		}
		if v.Type == ValInvalid {
			t.Skip()
		}
		if i.Error != nil && i.Error != io.EOF {
			t.Skip()
		}
		// Writing v to buf.
		var buf bytes.Buffer
		s := Default.Stream(&buf)
		v.Write(s)
		if err := s.Flush(); err != nil {
			t.Fatal(err)
		}

		// Parsing from buf to new value.
		i.ResetBytes(buf.Bytes())
		var parsed Value
		parseVal(i, &parsed)
		if i.Error != nil && i.Error != io.EOF {
			t.Fatalf("%v:\nBuf:   %s\nValue: %s\nData:  %s",
				i.Error, buf.Bytes(), v, data)
		}
		if !reflect.DeepEqual(parsed, v) {
			t.Fatalf("%v:\nBuf:   %s\nValue: %s != %s \nData:  %s",
				i.Error, buf.Bytes(), parsed, v, data)
		}
		// Writing parsed value to newBuf.
		var newBuf bytes.Buffer
		s.Reset(&newBuf)
		parsed.Write(s)
		if err := s.Flush(); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(newBuf.Bytes(), buf.Bytes()) {
			t.Fatalf("%s != %s", &newBuf, &buf)
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
		v.Str = i.ReadString()
		v.Type = ValStr
	case Nil:
		i.ReadNil()
		v.Type = ValNull
	case Bool:
		v.Bool = i.ReadBool()
		v.Type = ValBool
	case Object:
		v.Type = ValObj
		return i.ReadObjectCB(func(i *Iterator, s string) bool {
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
		return i.ReadArrayCB(func(i *Iterator) bool {
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
