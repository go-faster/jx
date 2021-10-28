package jir

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
	i := Default.Iterator([]byte(input))
	i.Object(func(i *Iterator, key string) bool {
		return i.Array(func(i *Iterator) bool {
			// Reading "type" field value first.
			var typ string
			i.Capture(func(i *Iterator) {
				i.Object(func(i *Iterator, key string) bool {
					switch key {
					case "type":
						typ = i.String()
					default:
						i.Skip()
					}
					return true
				})
			})
			// Reading objects depending on type.
			return i.Object(func(i *Iterator, key string) bool {
				if key == "type" {
					assert.Equal(t, typ, i.String())
					return true
				}
				switch typ {
				case "foo":
					i.String()
				case "bar":
					i.Int()
				}
				return true
			})
		})
	})
	require.NoError(t, i.Error)
}

func BenchmarkIterator_Skip(b *testing.B) {
	var input = []byte(`{"type": "foo", "foo": "string"}`)
	it := Default.Iterator(nil)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		it.ResetBytes(input)
		it.Capture(func(i *Iterator) {
			i.Skip()
		})
	}
}
