package jx

import (
	"bytes"
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
	i := GetDecoder()
	i.ResetBytes([]byte(input))
	err := i.Obj(func(i *Decoder, key string) error {
		return i.Arr(func(i *Decoder) error {
			// Reading "type" field value first.
			var typ string
			if err := i.Capture(func(i *Decoder) error {
				return i.Obj(func(i *Decoder, key string) error {
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
			return i.Obj(func(i *Decoder, key string) error {
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
	it := GetDecoder()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		it.ResetBytes(input)
		if err := it.Capture(func(i *Decoder) error {
			return i.Skip()
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func TestDecoder_Capture(t *testing.T) {
	strs := []string{
		"foo",
		"bar",
		"baz",
	}
	test := func(i *Decoder) func(t *testing.T) {
		return func(t *testing.T) {
			var elems int
			if err := i.Capture(func(i *Decoder) error {
				return i.Arr(func(i *Decoder) error {
					elems++
					return i.Skip()
				})
			}); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, Array, i.Next())
			require.Equal(t, 6, elems)
			t.Run("Nil", func(t *testing.T) {
				require.NoError(t, i.Capture(nil))
				require.Equal(t, Array, i.Next())
			})

			idx := 0
			require.NoError(t, i.Arr(func(d *Decoder) error {
				v, err := d.Str()
				if err != nil {
					return err
				}
				require.Equal(t, strs[idx%len(strs)], v)

				idx++
				return nil
			}))
		}
	}

	var e Encoder
	e.ArrStart()
	for i := 0; i < 6; i++ {
		e.Str(strs[i%len(strs)])
	}
	e.ArrEnd()
	testData := e.Bytes()

	t.Run("Str", test(DecodeBytes(testData)))
	// Check that we get correct result even if buffer smaller than captured data.
	decoder := Decoder{
		reader: bytes.NewReader(testData),
		buf:    make([]byte, 8),
	}
	t.Run("Reader", test(&decoder))
}
