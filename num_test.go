package jx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder_Num(t *testing.T) {
	var e Encoder
	e.Num(Num{
		Format: NumFormatInt,
		Value:  []byte{'1', '2', '3'},
	})
	require.Equal(t, e.String(), "123")
}

func TestNum(t *testing.T) {
	t.Run("ZeroValue", func(t *testing.T) {
		// Zero value is invalid because there is no Nun.Value.
		var v Num
		require.Equal(t, NumFormatInvalid, v.Format)
		require.True(t, v.Format.Invalid())
		require.False(t, v.Zero())
		require.False(t, v.Positive())
		require.False(t, v.Negative())
	})
	t.Run("Integer", func(t *testing.T) {
		v := Num{
			Format: NumFormatInt,
			Value:  []byte{'1', '2', '3'},
		}
		t.Run("Encode", func(t *testing.T) {
			var e Encoder
			e.Num(v)
			require.Equal(t, e.String(), "123")

			n, err := DecodeBytes(e.Bytes()).Int()
			require.NoError(t, err)
			require.Equal(t, 123, n)
		})
		t.Run("Methods", func(t *testing.T) {
			assert.True(t, v.Positive())
			assert.True(t, v.Format.Int())
			assert.False(t, v.Format.Invalid())
			assert.False(t, v.Negative())
			assert.False(t, v.Zero())
			assert.Equal(t, 1, v.Sign())
			assert.Equal(t, "123", v.String())
		})
	})
}

func BenchmarkNum(b *testing.B) {
	b.Run("Integer", func(b *testing.B) {
		v := Num{
			Format: NumFormatInt,
			Value:  []byte{'1', '2', '3', '5', '7'},
		}
		b.Run("Positive", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if !v.Positive() {
					b.Fatal("should be positive")
				}
			}
		})
		b.Run("Zero", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if v.Zero() {
					b.Fatal("should be not zero")
				}
			}
		})
		b.Run("Encode", func(b *testing.B) {
			b.ReportAllocs()
			var e Encoder
			for i := 0; i < b.N; i++ {
				e.Num(v)
				e.Reset()
			}
		})
		b.Run("Format.Invalid", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if v.Format.Invalid() {
					b.Fatal("invalid")
				}
			}
		})
	})
}
