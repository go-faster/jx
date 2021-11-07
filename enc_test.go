package jx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_byte_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.byte('1')
	should.Equal("1", string(e.Bytes()))
	should.Equal(1, len(e.buf))
	e.byte('2')
	should.Equal("12", string(e.Bytes()))
	should.Equal(2, len(e.buf))
	e.threeBytes('3', '4', '5')
	should.Equal("12345", string(e.Bytes()))
}

func TestEncoder(t *testing.T) {
	data := `"hello world"`
	buf := []byte(data)
	var e Encoder
	t.Run("Write", func(t *testing.T) {
		e.Reset()
		n, err := e.Write(buf)
		require.NoError(t, err)
		require.Equal(t, n, len(buf))
		require.Equal(t, data, e.String())
	})
	t.Run("Raw", func(t *testing.T) {
		e.Reset()
		e.Raw(buf)
		require.Equal(t, data, e.String())
	})
	t.Run("SetBytes", func(t *testing.T) {
		e.Reset()
		e.SetBytes(buf)
		require.Equal(t, data, e.String())
	})
}

func TestEncoder_Raw_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.RawStr("123")
	should.Equal("123", string(e.Bytes()))
}

func TestEncoder_Str_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.Str("123")
	should.Equal(`"123"`, string(e.Bytes()))
}

func TestEncoder_ArrEmpty(t *testing.T) {
	e := GetEncoder()
	e.ArrEmpty()
	require.Equal(t, "[]", string(e.Bytes()))
}

func TestEncoder_ObjEmpty(t *testing.T) {
	e := GetEncoder()
	e.ObjEmpty()
	require.Equal(t, "{}", string(e.Bytes()))
}

func TestEncoder_Obj(t *testing.T) {
	t.Run("FieldStart", func(t *testing.T) {
		var e Encoder
		e.Obj(func(e *Encoder) {
			e.Field("hello", func(e *Encoder) {
				e.Str("world")
			})
		})
		require.Equal(t, `{"hello":"world"}`, e.String())
	})
	t.Run("Nil", func(t *testing.T) {
		var e Encoder
		e.Obj(nil)
		require.Equal(t, `{}`, e.String())
	})
}

func TestEncoder_Arr(t *testing.T) {
	t.Run("Elem", func(t *testing.T) {
		var e Encoder
		e.Arr(func(e *Encoder) {
			e.Str("world")
		})
		require.Equal(t, `["world"]`, e.String())
	})
	t.Run("Nil", func(t *testing.T) {
		var e Encoder
		e.Arr(nil)
		require.Equal(t, `[]`, e.String())
	})
}

func BenchmarkEncoder_Arr(b *testing.B) {
	b.Run("Manual", func(b *testing.B) {
		b.ReportAllocs()
		var e Encoder
		for i := 0; i < b.N; i++ {
			e.ArrStart()
			e.Null()
			e.ArrEnd()

			e.Reset()
		}
	})
	b.Run("Callback", func(b *testing.B) {
		b.ReportAllocs()
		var e Encoder
		for i := 0; i < b.N; i++ {
			e.Arr(func(e *Encoder) {
				e.Null()
			})

			e.Reset()
		}
	})
}

func BenchmarkEncoder_Field(b *testing.B) {
	for _, fields := range []int{
		1,
		5,
		10,
		100,
	} {
		b.Run(fmt.Sprintf("%d", fields), func(b *testing.B) {
			b.Run("Manual", func(b *testing.B) {
				b.ReportAllocs()
				var e Encoder
				for i := 0; i < b.N; i++ {
					e.ObjStart()
					for j := 0; j < fields; j++ {
						e.FieldStart("field")
						e.Null()
					}
					e.ObjEnd()

					e.Reset()
				}
			})
			b.Run("Callback", func(b *testing.B) {
				b.ReportAllocs()
				var e Encoder
				for i := 0; i < b.N; i++ {
					e.Obj(func(e *Encoder) {
						for j := 0; j < fields; j++ {
							e.Field("field", func(e *Encoder) {
								e.Null()
							})
						}
					})

					e.Reset()
				}
			})
		})
	}
}
