//go:build go1.18

package json

import (
	"bytes"
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
	var v Value
	const input = `{"foo":{"bar":1,"baz":[1,2,3.14],"200":null}}`
	i := ParseString(Default, input)
	parseVal(i, &v)
	assert.NoError(t, i.Error)
	assert.Equal(t, `{foo: {bar: 1, baz: [1, 2, 3.14], 200: null}}`, v.String())

	buf := new(bytes.Buffer)
	s := NewStream(Default, buf, 1024)
	v.Write(s)
	require.NoError(t, s.Flush())
	require.Equal(t, input, buf.String(), "encoded value should equal to input")
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

type ValType byte

const (
	ValStr ValType = iota
	ValInt
	ValFloat
	ValNull
	ValObj
	ValArr
	ValBool
)

// Value represents any json value as sum type.
type Value struct {
	Type  ValType
	Str   string
	Int   int64
	Float float64
	Key   string
	Bool  bool
	Child []Value
}

// Write json representation of Value to Stream.
func (v Value) Write(s *Stream) {
	if v.Key != "" {
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
	if v.Key != "" {
		b.WriteString(v.Key)
		b.WriteString(": ")
	}
	switch v.Type {
	case ValStr:
		b.WriteString(`"` + v.Str + `"'`)
	case ValFloat:
		b.WriteString(strconv.FormatFloat(v.Float, 'f', -1, 64))
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
		if strings.Contains(n.String(), ".") {
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
			v.Child = append(v.Child, elem)
			return i.Error == nil
		})
	case Array:
		v.Type = ValArr
		return i.ReadArrayCB(func(i *Iterator) bool {
			var elem Value
			if !parseVal(i, &elem) {
				return false
			}
			v.Child = append(v.Child, elem)
			return i.Error == nil
		})
	default:
		panic(i.WhatIsNext())
	}
	return i.Error == nil
}
