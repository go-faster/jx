package jx

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-faster/errors"
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

func TestDecoderIntUnexpectedChar(t *testing.T) {
	type intFunc struct {
		name    string
		bitSize int
		fn      func(*Decoder) error
	}
	signed := []intFunc{
		{"Int", strconv.IntSize, intDecoderOnlyError((*Decoder).Int)},
		{"Int8", 8, intDecoderOnlyError((*Decoder).Int8)},
		{"Int16", 16, intDecoderOnlyError((*Decoder).Int16)},
		{"Int32", 32, intDecoderOnlyError((*Decoder).Int32)},
		{"Int64", 64, intDecoderOnlyError((*Decoder).Int64)},
	}
	unsigned := []intFunc{
		{"UInt", strconv.IntSize, intDecoderOnlyError((*Decoder).UInt)},
		{"UInt8", 8, intDecoderOnlyError((*Decoder).UInt8)},
		{"UInt16", 16, intDecoderOnlyError((*Decoder).UInt16)},
		{"UInt32", 32, intDecoderOnlyError((*Decoder).UInt32)},
		{"UInt64", 64, intDecoderOnlyError((*Decoder).UInt64)},
	}

	tests := []struct {
		input     string
		unsigned  bool
		size      int // 0 for any
		errString string
	}{
		// Leading space.
		{" 10", true, 0, ""},
		{"   10", true, 0, ""},
		{" -10", false, 0, ""},

		// Space in the middle.
		{"- 10", false, 0, "unexpected byte 32 ' ' at 1"},

		// Digit after leading zero.
		{"00", true, 0, "digit after leading zero: unexpected byte 48 '0' at 1"},
		{"01", true, 0, "digit after leading zero: unexpected byte 49 '1' at 1"},

		// Unexpected character.
		// 8 bits.
		{"0a0", true, 0, "unexpected byte 97 'a' at 1"},
		{"1a00000000000", true, 0, "unexpected byte 97 'a' at 1"},
		{"10a0000000000", true, 0, "unexpected byte 97 'a' at 2"},
		{"100a000000000", true, 0, "unexpected byte 97 'a' at 3"},
		// 16 bits.
		{"1000a00000000", true, 16, "unexpected byte 97 'a' at 4"},
		{"10000a0000000", true, 16, "unexpected byte 97 'a' at 5"},
		// 32 bits.
		{"100000a000000", true, 32, "unexpected byte 97 'a' at 6"},
		{"1000000a00000", true, 32, "unexpected byte 97 'a' at 7"},
		{"10000000a0000", true, 32, "unexpected byte 97 'a' at 8"},
		{"100000000a000", true, 32, "unexpected byte 97 'a' at 9"},
		{"1000000000a00", true, 32, "unexpected byte 97 'a' at 10"},
		// 64 bits.
		{"10000000000a0", true, 64, "unexpected byte 97 'a' at 11"},

		// Dot in integer.
		// 8 bits.
		{"0.0", true, 0, "unexpected floating point character: unexpected byte 46 '.' at 1"},
		{"1.00000000000", true, 0, "unexpected floating point character: unexpected byte 46 '.' at 1"},
		{"10.0000000000", true, 0, "unexpected floating point character: unexpected byte 46 '.' at 2"},
		{"100.000000000", true, 0, "unexpected floating point character: unexpected byte 46 '.' at 3"},
		// 16 bits.
		{"1000.00000000", true, 16, "unexpected floating point character: unexpected byte 46 '.' at 4"},
		{"10000.0000000", true, 16, "unexpected floating point character: unexpected byte 46 '.' at 5"},
		// 32 bits.
		{"100000.000000", true, 32, "unexpected floating point character: unexpected byte 46 '.' at 6"},
		{"1000000.00000", true, 32, "unexpected floating point character: unexpected byte 46 '.' at 7"},
		{"10000000.0000", true, 32, "unexpected floating point character: unexpected byte 46 '.' at 8"},
		{"100000000.000", true, 32, "unexpected floating point character: unexpected byte 46 '.' at 9"},
		{"1000000000.00", true, 32, "unexpected floating point character: unexpected byte 46 '.' at 10"},
		// 64 bits.
		{"10000000000.0", true, 64, "unexpected floating point character: unexpected byte 46 '.' at 11"},

		// Exp in integer.
		{"0e0", true, 0, "unexpected floating point character: unexpected byte 101 'e' at 1"},
		{"0E0", true, 0, "unexpected floating point character: unexpected byte 69 'E' at 1"},
		{"0e-0", true, 0, "unexpected floating point character: unexpected byte 101 'e' at 1"},
		{"0e+0", true, 0, "unexpected floating point character: unexpected byte 101 'e' at 1"},

		{"1e0", true, 0, "unexpected floating point character: unexpected byte 101 'e' at 1"},
		{"1E0", true, 0, "unexpected floating point character: unexpected byte 69 'E' at 1"},
		{"1e-0", true, 0, "unexpected floating point character: unexpected byte 101 'e' at 1"},
		{"1e+0", true, 0, "unexpected floating point character: unexpected byte 101 'e' at 1"},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			check := func(fns []intFunc) {
				for _, intFn := range fns {
					intFn := intFn
					if tt.size != 0 && tt.size > intFn.bitSize {
						continue
					}
					t.Run(intFn.name, func(t *testing.T) {
						decodeStr(t, tt.input, func(t *testing.T, d *Decoder) {
							a := assert.New(t)

							err := intFn.fn(d)
							if e := tt.errString; e != "" {
								a.EqualError(err, e)
								v, ok := errors.Into[*badTokenErr](err)
								if !ok {
									return
								}
								a.Equal(v.Token, tt.input[v.Offset])
								return
							}
							a.NoError(err)
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
