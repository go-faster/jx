package jx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder_Num(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		var e Encoder
		e.Num(Num{'1', '2', '3'})
		require.Equal(t, e.String(), "123")
	})
	t.Run("Invalid", func(t *testing.T) {
		var e Encoder
		e.Num(Num{})
		require.Equal(t, e.String(), "null")
	})
}

func TestNum(t *testing.T) {
	t.Run("Format", func(t *testing.T) {
		assert.Equal(t, `-12`, fmt.Sprintf("%d", Num(`"-12.0"`)))
		assert.Equal(t, `-12.000000`, fmt.Sprintf("%f", Num(`"-12.0"`)))
		assert.Equal(t, `%!invalid(Num=f)`, fmt.Sprintf("%f", Num(`f`)))
		assert.Equal(t, `"-12.0"`, fmt.Sprintf("%s", Num(`"-12.0"`)))
		assert.Equal(t, `"-12.0"`, fmt.Sprintf("%v", Num(`"-12.0"`)))
		assert.Equal(t, `%!invalid(Num=f)`, fmt.Sprintf("%d", Num(`f`)))
		assert.Equal(t, `%!invalid(Num="-12.1")`, fmt.Sprintf("%d", Num(`"-12.1"`)))
	})
	t.Run("String", func(t *testing.T) {
		for _, cc := range []struct {
			Name   string
			String string
			Value  Num
		}{
			{
				Name:   "Int",
				String: "-12",
				Value:  Num("-12"),
			},
			{
				Name:   "IntStr",
				String: `"-12"`,
				Value:  Num(`"-12"`),
			},
		} {
			t.Run(cc.Name, func(t *testing.T) {
				require.Equal(t, cc.String, cc.Value.String())
				require.True(t, cc.Value.IsInt())
			})
		}
	})
	t.Run("Str", func(t *testing.T) {
		v := Num{'"', '1', '2', '3', '"'}
		assert.True(t, v.Positive())
		assert.False(t, v.Negative())
		assert.False(t, v.Zero())
		assert.True(t, v.Str())
		assert.Equal(t, 1, v.Sign())
		assert.Equal(t, `"123"`, v.String())
		assert.True(t, v.Equal(v))
		assert.True(t, v.IsInt())
		assert.False(t, v.Equal(Num{}))

		assert.Equal(t, 0, Num{'"'}.Sign())
	})
	t.Run("ZeroValue", func(t *testing.T) {
		// Zero value is invalid because there is no Num.Value.
		var v Num
		require.False(t, v.Zero())
		require.False(t, v.Positive())
		require.False(t, v.Negative())
		require.False(t, v.IsInt())
		require.False(t, v.Str())
		require.Equal(t, "<invalid>", v.String())
	})
	t.Run("IntZero", func(t *testing.T) {
		v := Num{'0'}
		require.True(t, v.Zero())
		require.False(t, v.Positive())
		require.False(t, v.Negative())
		require.True(t, v.IsInt())
		require.False(t, v.Str())
		require.Equal(t, "0", v.String())
	})
	t.Run("FloatZero", func(t *testing.T) {
		v := Num{'-', '0', '.', '0'}
		require.True(t, v.Zero())
		require.False(t, v.Positive())
		require.True(t, v.Negative())
		require.False(t, v.IsInt())
		require.False(t, v.Str())
		require.Equal(t, "-0.0", v.String())
	})
	t.Run("Integer", func(t *testing.T) {
		t.Run("Int", func(t *testing.T) {
			v := Num{'1', '2', '3'}
			t.Run("Methods", func(t *testing.T) {
				assert.True(t, v.Positive())
				assert.False(t, v.Negative())
				assert.False(t, v.Zero())
				assert.False(t, v.Str())
				assert.Equal(t, 1, v.Sign())
				assert.Equal(t, "123", v.String())
				assert.True(t, v.Equal(v))
				assert.True(t, v.IsInt())
				assert.False(t, v.Equal(Num{}))
			})
			t.Run("Write", func(t *testing.T) {
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
				v := Num{'1', '2', '3', '.', '0'}
				n, err := v.Int64()
				require.NoError(t, err)
				require.Equal(t, int64(123), n)

				un, err := v.Uint64()
				require.NoError(t, err)
				require.Equal(t, uint64(123), un)

				f, err := v.Float64()
				require.NoError(t, err)
				require.InEpsilon(t, 123, f, epsilon)
			})
			t.Run("Negative", func(t *testing.T) {
				v := Num{'1', '2', '3', '.', '0', '0', '1'}
				_, err := v.Int64()
				require.Error(t, err)
				_, err = v.Uint64()
				require.Error(t, err)
			})
		})
		t.Run("Decode", func(t *testing.T) {
			n, err := DecodeStr("12345").Num()
			require.NoError(t, err)
			require.Equal(t, "12345", n.String())
		})
	})
	t.Run("Float", func(t *testing.T) {
		const (
			s = `1.23`
			f = 1.23
		)
		v := Num(s)
		t.Run("Write", func(t *testing.T) {
			var e Encoder
			e.Num(v)
			require.Equal(t, e.String(), s)

			n, err := DecodeBytes(e.Bytes()).Float64()
			require.NoError(t, err)
			require.InEpsilon(t, f, n, epsilon)
		})
		t.Run("Decode", func(t *testing.T) {
			n, err := DecodeStr(s).Num()
			require.NoError(t, err)
			require.Equal(t, s, n.String())
		})
		t.Run("Methods", func(t *testing.T) {
			assert.True(t, v.Positive())
			assert.False(t, v.Negative())
			assert.False(t, v.Zero())
			assert.False(t, v.IsInt())
			assert.Equal(t, 1, v.Sign())
			assert.Equal(t, s, v.String())
		})
	})
}

func BenchmarkNum(b *testing.B) {
	b.Run("FloatAsInt", func(b *testing.B) {
		b.Run("Integer", func(b *testing.B) {
			v := Num{'1', '2', '3', '5', '7', '.', '0'}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := v.Int64(); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("Float30Chars", func(b *testing.B) {
			var v Num
			for i := 0; i < 28; i++ {
				v = append(v, '1')
			}
			v = append(v, '.', '0')

			b.Run("AsInt", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					if err := v.floatAsInt(); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("IsInt", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					if v.IsInt() {
						b.Fatal("unexpected")
					}
				}
			})
		})
	})
	b.Run("Integer", func(b *testing.B) {
		v := Num{'1', '2', '3', '5', '7'}
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
		b.Run("Write", func(b *testing.B) {
			b.ReportAllocs()
			var e Encoder
			for i := 0; i < b.N; i++ {
				e.Num(v)
				e.Reset()
			}
		})
	})
}
