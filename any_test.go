package jx

import (
	"bytes"
	hexEnc "encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAny_Read(t *testing.T) {
	t.Run("Obj", func(t *testing.T) {
		var v Any
		const input = `{"foo":{"bar":1,"baz":[1,2,3.14],"200":null}}`
		r := ReadString(input)
		assert.NoError(t, v.Read(r))
		assert.Equal(t, `{foo: {bar: 1, baz: [1, 2, f3.14], 200: null}}`, v.String())

		buf := new(bytes.Buffer)
		w := NewWriter(buf, 1024)
		require.NoError(t, w.Any(v))
		require.NoError(t, w.Flush())
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
				var v Any
				input := []byte(tt.Input)
				r := ReadBytes(input)
				require.NoError(t, v.Read(r))

				buf := new(bytes.Buffer)
				s := NewWriter(buf, 1024)
				require.NoError(t, v.Write(s))
				require.NoError(t, s.Flush())
				require.Equal(t, tt.Input, buf.String(), "encoded value should equal to input")

				var otherValue Any
				r.ResetBytes(buf.Bytes())

				if err := otherValue.Read(r); err != nil {
					t.Error(err)
					t.Log(hexEnc.Dump(input))
					t.Log(hexEnc.Dump(buf.Bytes()))
				}
			})
		}
	})
}

func BenchmarkAny(b *testing.B) {
	data := []byte(`[true, null, false, 100, "false"]`)
	r := NewReader()

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
