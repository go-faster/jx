package jx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkDecoder_Int(b *testing.B) {
	runTestdataFile("integers.json", b.Fatal, func(name string, data []byte) {
		b.Run(name, func(b *testing.B) {
			d := GetDecoder()
			cb := func(d *Decoder) error {
				_, err := d.Int()
				return err
			}
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				d.ResetBytes(data)

				if err := d.Arr(cb); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

func BenchmarkDecoder_Uint(b *testing.B) {
	runTestdataFile("integers.json", b.Fatal, func(name string, data []byte) {
		b.Run(name, func(b *testing.B) {
			d := GetDecoder()
			cb := func(d *Decoder) error {
				_, err := d.UInt()
				return err
			}
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				d.ResetBytes(data)

				if err := d.Arr(cb); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

func TestDecoderIntSizes(t *testing.T) {
	data := []byte(`69315063`)
	d := GetDecoder()
	for _, size := range []int{32, 64} {
		d.ResetBytes(data)
		v, err := d.int(size)
		require.NoError(t, err)
		require.Equal(t, 69315063, v)
	}
}

func TestDecoderUintSizes(t *testing.T) {
	data := []byte(`69315063`)
	d := GetDecoder()
	for _, size := range []int{32, 64} {
		d.ResetBytes(data)
		v, err := d.uint(size)
		require.NoError(t, err)
		require.Equal(t, uint(69315063), v)
	}
}

func TestDecoderIntError(t *testing.T) {
	r := errReader{}
	get := func() *Decoder {
		return &Decoder{
			buf:    []byte{'1', '2'},
			tail:   2,
			reader: r,
		}
	}
	t.Run("Int8", func(t *testing.T) {
		d := get()
		_, err := d.Int8()
		require.ErrorIs(t, err, r.Err())
	})
	t.Run("Int16", func(t *testing.T) {
		d := get()
		_, err := d.Int16()
		require.ErrorIs(t, err, r.Err())
	})
	t.Run("Int32", func(t *testing.T) {
		d := get()
		_, err := d.Int32()
		require.ErrorIs(t, err, r.Err())
	})
	t.Run("Int64", func(t *testing.T) {
		d := get()
		_, err := d.Int64()
		require.ErrorIs(t, err, r.Err())
	})
}

func intDecoderOnlyError[T any](fn func(*Decoder) (T, error)) func(*Decoder) error {
	return func(d *Decoder) error {
		_, err := fn(d)
		return err
	}
}

func TestDecoderIntUnexpectedSpace(t *testing.T) {
	type intFunc struct {
		name string
		fn   func(*Decoder) error
	}
	signed := []intFunc{
		{"Int", intDecoderOnlyError((*Decoder).Int)},
		{"Int8", intDecoderOnlyError((*Decoder).Int8)},
		{"Int16", intDecoderOnlyError((*Decoder).Int16)},
		{"Int32", intDecoderOnlyError((*Decoder).Int32)},
		{"Int64", intDecoderOnlyError((*Decoder).Int64)},
	}
	unsigned := []intFunc{
		{"UInt", intDecoderOnlyError((*Decoder).UInt)},
		{"UInt8", intDecoderOnlyError((*Decoder).UInt8)},
		{"UInt16", intDecoderOnlyError((*Decoder).UInt16)},
		{"UInt32", intDecoderOnlyError((*Decoder).UInt32)},
		{"UInt64", intDecoderOnlyError((*Decoder).UInt64)},
	}

	tests := []struct {
		input    string
		unsigned bool
		wantErr  bool
	}{
		{" 10", true, false},
		{"   10", true, false},
		{" -10", false, false},

		{"- 10", false, true},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			check := func(fns []intFunc) {
				for _, intFn := range fns {
					intFn := intFn
					t.Run(intFn.name, func(t *testing.T) {
						decodeStr(t, tt.input, func(t *testing.T, d *Decoder) {
							err := intFn.fn(d)
							if tt.wantErr {
								require.Error(t, err)
								return
							}
							require.NoError(t, err)
						})
					})
				}
			}

			check(signed)
			if tt.unsigned {
				check(unsigned)
			}
		})
	}
}
