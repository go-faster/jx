package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_comma(t *testing.T) {
	t.Run("Array", func(t *testing.T) {
		var e Encoder
		e.ArrStart()
		e.Int(1)
		e.ArrStart()
		e.Int(2)
		e.RawStr(`3`)
		e.ArrEnd()
		e.ArrEnd()

		require.Equal(t, "[1,[2,3]]", e.String())
	})
	t.Run("Object", func(t *testing.T) {
		var e Encoder
		e.ObjStart()
		e.FieldStart("a")
		e.Int(1)
		e.FieldStart("b")
		e.Int(2)
		e.FieldStart("c")
		e.ArrStart()
		e.Int(1)
		e.Raw([]byte{'2'})
		e.Float32(3.0)
		e.Float64(4.5)
		e.Num(Num{'2', '3'})
		e.True()
		e.False()
		e.Null()
		e.Base64(Raw{1})
		e.Bool(true)
		e.ArrEnd()
		e.ObjEnd()

		require.Equal(t, `{"a":1,"b":2,"c":[1,2,3,4.5,23,true,false,null,"AQ==",true]}`, e.String())
	})
	t.Run("NoPanic", func(t *testing.T) {
		var e Encoder
		e.ObjEnd()
		e.ObjEnd()
		e.ArrEnd()
		e.ArrEnd()
	})
}

func BenchmarkEncoder_comma_overhead(b *testing.B) {
	// Measure overhead of ArrStart + comma + resetComma.
	// BenchmarkEncoder_comma-32 5.057 ns/op 0 B/op	0 allocs/op.
	var e Encoder
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e.ArrStart()
		e.comma()
		e.resetComma()
		e.Reset()
	}
}
