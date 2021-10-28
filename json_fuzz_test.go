//go:build go1.18

package json

import (
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
