//go:build go1.18

package json

import (
	"testing"
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

func goRecursive(i *Iterator) bool {
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
			return goRecursive(i)
		})
	case Array:
		return i.ReadArrayCB(func(i *Iterator) bool {
			return goRecursive(i)
		})
	default:
		panic(i.WhatIsNext())
	}
	return i.Error == nil
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
		goRecursive(i)
	})
}
