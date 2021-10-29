package jx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIterator_Capture(t *testing.T) {
	const input = `{
	"objects": [
		{
			"type": "foo",
			"foo": "string"
		},
		{
			"type": "bar",
			"bar": 1000
		}
	]
}`
	i := GetIter()
	i.ResetBytes([]byte(input))
	err := i.Obj(func(i *Iter, key string) error {
		return i.Array(func(i *Iter) error {
			// Reading "type" field value first.
			var typ string
			if err := i.Capture(func(i *Iter) error {
				return i.Obj(func(i *Iter, key string) error {
					switch key {
					case "type":
						s, err := i.Str()
						if err != nil {
							return err
						}
						typ = s
					default:
						return i.Skip()
					}
					return nil
				})
			}); err != nil {
				return err
			}
			// Reading objects depending on type.
			return i.Obj(func(i *Iter, key string) error {
				if key == "type" {
					s, err := i.Str()
					if err != nil {
						return err
					}
					assert.Equal(t, typ, s)
					return nil
				}
				switch typ {
				case "foo":
					_, _ = i.Str()
				case "bar":
					_, err := i.Int()
					return err
				}
				return nil
			})
		})
	})
	require.NoError(t, err)
}

func BenchmarkIterator_Skip(b *testing.B) {
	var input = []byte(`{"type": "foo", "foo": "string"}`)
	it := GetIter()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		it.ResetBytes(input)
		if err := it.Capture(func(i *Iter) error {
			return i.Skip()
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func TestIter_Capture(t *testing.T) {
	i := ParseString(`["foo", "bar", "baz"]`)
	var elems int
	if err := i.Capture(func(i *Iter) error {
		return i.Array(func(i *Iter) error {
			elems++
			return i.Skip()
		})
	}); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, Array, i.Next())
	require.Equal(t, 3, elems)
}
