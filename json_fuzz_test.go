//go:build go1.18

package json

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	i := ParseString(Default, `{"foo": {"bar": 1, "baz": [1, 2, 3]}}`)
	parseVal(i, &v)
	assert.NoError(t, i.Error)
	assert.Equal(t, `{foo: {bar: 1, baz: [1, 2, 3]}}`, v.String())
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

type Value struct {
	Type  ValType
	Str   string
	Int   int64
	Float float64
	Key   string
	Bool  bool
	Child []Value
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
