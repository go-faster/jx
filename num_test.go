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
	t.Run("String", func(t *testing.T) {
		for _, cc := range []struct {
			Name   string
			String string
			Value  Num
		}{
			{
				Name:   "Int",
				String: "-12",
				Value: Num{
					Format: NumFormatInt,
					Value:  []byte("-12"),
				},
			},
			{
				Name:   "IntStr",
				String: `"-12"`,
				Value: Num{
					Format: NumFormatIntStr,
					Value:  []byte("-12"),
				},
			},
		} {
			t.Run(cc.Name, func(t *testing.T) {
				require.Equal(t, cc.String, cc.Value.String())
			})
		}
	})
	t.Run("ZeroValue", func(t *testing.T) {
		// Zero value is invalid because there is no Num.Value.
		var v Num
		require.Equal(t, NumFormatInvalid, v.Format)
		require.True(t, v.Format.Invalid())
		require.False(t, v.Zero())
		require.False(t, v.Positive())
		require.False(t, v.Negative())
		require.Equal(t, "<invalid>", v.String())
	})
	t.Run("Integer", func(t *testing.T) {
		t.Run("Int", func(t *testing.T) {
			v := Num{
				Format: NumFormatInt,
				Value:  []byte{'1', '2', '3'},
			}
			t.Run("Methods", func(t *testing.T) {
				assert.True(t, v.Positive())
				assert.True(t, v.Format.Int())
				assert.False(t, v.Format.Invalid())
				assert.False(t, v.Negative())
				assert.False(t, v.Zero())
				assert.Equal(t, 1, v.Sign())
				assert.Equal(t, "123", v.String())
			})
			t.Run("Encode", func(t *testing.T) {
				var e Encoder
				e.Num(v)
				require.Equal(t, e.String(), "123")

				n, err := DecodeBytes(e.Bytes()).Int()
				require.NoError(t, err)
				require.Equal(t, 123, n)
			})
		})
		t.Run("FloatAsInt", func(t *testing.T) {
			t.Run("Positive", func(t *testing.T) {
				v := Num{
					Format: NumFormatFloat,
					Value:  []byte{'1', '2', '3', '.', '0'},
				}
				n, err := v.Int()
				require.NoError(t, err)
				require.Equal(t, 123, n)
			})
			t.Run("Negative", func(t *testing.T) {
				v := Num{
					Format: NumFormatFloat,
					Value:  []byte{'1', '2', '3', '.', '0', '0', '1'},
				}
				_, err := v.Int()
				require.Error(t, err)
			})
		})
		t.Run("Decode", func(t *testing.T) {
			n, err := DecodeStr("12345").NumTo(Num{})
			require.NoError(t, err)
			require.Equal(t, NumFormatInt, n.Format)
			require.Equal(t, "12345", n.String())
		})
	})
	t.Run("Float", func(t *testing.T) {
		const (
			s = `1.23`
			f = 1.23
		)
		v := Num{
			Format: NumFormatFloat,
			Value:  []byte(s),
		}
		t.Run("Encode", func(t *testing.T) {
			var e Encoder
			e.Num(v)
			require.Equal(t, e.String(), s)

			n, err := DecodeBytes(e.Bytes()).Float64()
			require.NoError(t, err)
			require.InEpsilon(t, f, n, epsilon)
		})
		t.Run("Decode", func(t *testing.T) {
			n, err := DecodeStr(s).NumTo(Num{})
			require.NoError(t, err)
			require.Equal(t, NumFormatFloat, n.Format)
			require.Equal(t, s, n.String())
		})
		t.Run("Methods", func(t *testing.T) {
			assert.True(t, v.Positive())
			assert.True(t, v.Format.Float())
			assert.False(t, v.Format.Invalid())
			assert.False(t, v.Negative())
			assert.False(t, v.Zero())
			assert.Equal(t, 1, v.Sign())
			assert.Equal(t, s, v.String())
		})
	})
}

func BenchmarkNum(b *testing.B) {
	b.Run("FloatAsInt", func(b *testing.B) {
		b.Run("Integer", func(b *testing.B) {
			v := Num{
				Format: NumFormatFloat,
				Value:  []byte{'1', '2', '3', '5', '7', '.', '0'},
			}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := v.Int(); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
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
