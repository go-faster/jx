package jx

import (
	hexEnc "encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Reset Any value to reuse.
func (v *Any) Reset() {
	v.Type = AnyInvalid
	v.Child = v.Child[:0]
	v.KeyValid = false

	v.Str = ""
	v.Key = ""
}

// Obj calls f for any child that is field if v is AnyObj.
func (v Any) Obj(f func(k string, v Any)) {
	if v.Type != AnyObj {
		return
	}
	for _, c := range v.Child {
		if !c.KeyValid {
			continue
		}
		f(c.Key, c)
	}
}

func TestAny_Read(t *testing.T) {
	t.Run("Obj", func(t *testing.T) {
		var v Any
		const input = `{"foo":{"bar":1,"baz":[1,2,3.14],"200":null,"f":"s","t":true,"":""}}`
		r := DecodeStr(input)
		assert.NoError(t, v.Read(r))
		assert.Equal(t, "{foo: {bar: 1, baz: [1, 2, 3.14], 200: null, f: 's', t: true, <blank>: ''}}", v.String())

		e := GetEncoder()
		e.Any(v)
		require.Equal(t, input, e.String(), "encoded value should equal to input")
	})
	t.Run("Inputs", func(t *testing.T) {
		for _, tt := range []struct {
			Input string
		}{
			{Input: "1"},
			{Input: "0.0"},
		} {
			t.Run(tt.Input, func(t *testing.T) {
				input := []byte(tt.Input)
				r := DecodeBytes(input)
				v, err := r.Any()
				require.NoError(t, err)

				e := GetEncoder()
				v.Write(e)

				var otherValue Any
				r.ResetBytes(e.Bytes())

				if err := otherValue.Read(r); err != nil {
					t.Error(err)
					t.Log(hexEnc.Dump(input))
					t.Log(hexEnc.Dump(e.Bytes()))
				}

				require.True(t, otherValue.Equal(v))
			})
		}
	})
}

func TestAny_Equal(t *testing.T) {
	t.Run("ZeroValues", func(t *testing.T) {
		for _, typ := range []AnyType{
			AnyInvalid,
			AnyStr,
			AnyNumber,
			AnyNull,
			AnyObj,
			AnyArr,
			AnyBool,
		} {
			a := Any{Type: typ}
			t.Run("Equal", func(t *testing.T) {
				b := Any{Type: typ}
				require.True(t, a.Equal(b))
				t.Run("Child", func(t *testing.T) {
					aArr := Any{
						Type:  AnyArr,
						Child: []Any{a},
					}
					bArr := Any{
						Type: AnyArr,
					}
					require.False(t, aArr.Equal(bArr))
					bArr.Child = []Any{b}
					require.True(t, aArr.Equal(bArr))
				})
			})
			t.Run("NotEqual", func(t *testing.T) {
				b := Any{Type: typ + 1}
				require.False(t, a.Equal(b))
				t.Run("Child", func(t *testing.T) {
					aArr := Any{
						Type:  AnyArr,
						Child: []Any{a},
					}
					bArr := Any{
						Type:  AnyArr,
						Child: []Any{b},
					}
					require.False(t, aArr.Equal(bArr))
				})
			})
			t.Run("Keys", func(t *testing.T) {
				a.KeyValid = true
				b := Any{Type: typ, KeyValid: true}
				require.True(t, a.Equal(b))
				b.Key = "1"
				require.False(t, a.Equal(b))
			})

		}
	})
}

func BenchmarkAny(b *testing.B) {
	data := []byte(`[true, null, false, 100, "false"]`)
	r := GetDecoder()

	b.ReportAllocs()
	b.SetBytes(int64(len(data)))

	var v Any
	for i := 0; i < b.N; i++ {
		v.Reset()
		r.ResetBytes(data)
		if err := v.Read(r); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAnyStd(b *testing.B) {
	data := []byte(`[true, null, false, 100, "false"]`)
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))

	var v []interface{}
	for i := 0; i < b.N; i++ {
		v = v[:0]
		if err := json.Unmarshal(data, &v); err != nil {
			b.Fatal(err)
		}
	}
}
