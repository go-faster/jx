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
		e.Int(3)
		e.ArrEnd()
		e.ArrEnd()

		require.Equal(t, "[1,[2,3]]", e.String())
	})
	t.Run("Object", func(t *testing.T) {
		var e Encoder
		e.ObjStart()
		e.Field("a")
		e.Int(1)
		e.Field("b")
		e.Int(2)
		e.Field("c")
		e.ArrStart()
		e.Int(1)
		e.Int(2)
		e.ArrEnd()
		e.ObjEnd()

		require.Equal(t, `{"a":1,"b":2,"c":[1,2]}`, e.String())
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
