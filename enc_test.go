package jx

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testEncoderModes(t *testing.T, cb func(*Encoder), expected string) {
	t.Run("Buffer", func(t *testing.T) {
		e := GetEncoder()
		cb(e)
		require.Equal(t, expected, e.String())
	})
	t.Run("Writer", func(t *testing.T) {
		var sb strings.Builder
		e := NewStreamingEncoder(&sb, -1)
		cb(e)
		require.NoError(t, e.Close())
		require.Equal(t, expected, sb.String())
	})
}

// requireCompat fails if `encoding/json` will encode v differently than exp.
func requireCompat(t *testing.T, cb func(*Encoder), v any) {
	t.Helper()
	buf, err := json.Marshal(v)
	require.NoError(t, err, "json.Marshal(%#v)", v)
	testEncoderModes(t, cb, string(buf))
}

func TestEncoderGrow(t *testing.T) {
	should := require.New(t)
	e := &Encoder{}
	should.Equal(0, len(e.Bytes()))
	should.Equal(0, cap(e.Bytes()))
	e.Grow(1024)
	should.Equal(0, len(e.Bytes()))
	should.Equal(1024, cap(e.Bytes()))
	e.Grow(512)
	should.Equal(0, len(e.Bytes()))
	should.Equal(1024, cap(e.Bytes()))
	e.Grow(4096)
	should.Equal(0, len(e.Bytes()))
	should.Equal(4096, cap(e.Bytes()))
}

func TestEncoderByteShouldGrowBuffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.byte('1')
	should.Equal("1", string(e.Bytes()))
	should.Equal(1, len(e.w.Buf))
	e.byte('2')
	should.Equal("12", string(e.Bytes()))
	should.Equal(2, len(e.w.Buf))
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

func TestEncoderRawShouldGrowBuffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.RawStr("123")
	should.Equal("123", string(e.Bytes()))
}

func TestEncoderStrShouldGrowBuffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.Str("123")
	should.Equal(`"123"`, string(e.Bytes()))
}

func TestEncoder_ArrEmpty(t *testing.T) {
	testEncoderModes(t, func(e *Encoder) {
		e.ArrEmpty()
	}, "[]")
}

func TestEncoder_ObjEmpty(t *testing.T) {
	testEncoderModes(t, func(e *Encoder) {
		e.ObjEmpty()
	}, "{}")
}

func TestEncoder_Obj(t *testing.T) {
	t.Run("Field", func(t *testing.T) {
		testEncoderModes(t, func(e *Encoder) {
			e.Obj(func(e *Encoder) {
				e.Field("hello", func(e *Encoder) {
					e.Str("world")
				})
			})
		}, `{"hello":"world"}`)
	})
	t.Run("Nil", func(t *testing.T) {
		testEncoderModes(t, func(e *Encoder) {
			e.Obj(nil)
		}, `{}`)
	})
}

func TestEncoder_Arr(t *testing.T) {
	t.Run("Elem", func(t *testing.T) {
		testEncoderModes(t, func(e *Encoder) {
			e.Arr(func(e *Encoder) {
				e.Str("world")
			})
		}, `["world"]`)
	})
	t.Run("Nil", func(t *testing.T) {
		testEncoderModes(t, func(e *Encoder) {
			e.Arr(nil)
		}, `[]`)
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
