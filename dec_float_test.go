package jx

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-faster/errors"
)

func decodeStr(t *testing.T, s string, f func(t *testing.T, d *Decoder)) {
	t.Helper()
	for _, d := range []struct {
		Name string
		Fn   func() *Decoder
	}{
		{
			Name: "DecodeStr",
			Fn:   func() *Decoder { return DecodeStr(s) },
		},
		{
			Name: "DecodeBytes",
			Fn:   func() *Decoder { return DecodeBytes([]byte(s)) },
		},
		{
			Name: "Decode",
			Fn:   func() *Decoder { return Decode(strings.NewReader(s), 0) },
		},
		{
			Name: "DecodeSingleByteBuf",
			Fn:   func() *Decoder { return Decode(strings.NewReader(s), 1) },
		},
	} {
		t.Run(d.Name, func(t *testing.T) {
			t.Helper()
			dec := d.Fn()
			f(t, dec)
		})
	}
}

func TestDecoder_Float(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		runTestCases(t, testNumbers, func(t *testing.T, d *Decoder) error {
			_, err := d.Float64()
			if err != nil {
				return err
			}
			if err := d.Skip(); err != nil {
				if err != io.EOF && err != io.ErrUnexpectedEOF {
					return err
				}
			}
			return nil
		})
	})
	t.Run("Slow", func(t *testing.T) {
		s := `,0.1`
		t.Run("64", func(t *testing.T) {
			d := Decode(bytes.NewBufferString(s), 2)
			requireElem(t, d)
			_, err := d.Float64()
			require.NoError(t, err)
		})
		t.Run("32", func(t *testing.T) {
			d := Decode(bytes.NewBufferString(s), 2)
			requireElem(t, d)
			v, err := d.Float32()
			require.NoError(t, err)
			t.Logf("%f", v)
		})
	})
}

func TestDecoder_BigFloat(t *testing.T) {
	data := []byte{'1'}
	for i := 0; i < 64; i++ {
		data = append(data, '0')
	}
	data = append(data, ".0"...)
	f, err := DecodeBytes(data).BigFloat()
	require.NoError(t, err)
	require.Equal(t, `1e+64`, f.String())
}

func TestDecoder_Float32(t *testing.T) {
	v, err := DecodeStr(`429496739.0`).Float32()
	require.NoError(t, err)
	require.InEpsilon(t, 429496729.0, v, epsilon)
}

func TestDecoder_Float64(t *testing.T) {
	for i, tc := range []struct {
		String string
		Value  float64
	}{
		{
			String: `18446744073709551700.0`,
			Value:  18446744073709551700.0,
		},
		{
			String: `18446744073709551.7000`,
			Value:  18446744073709551.7000,
		},
	} {
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			t.Run("32Str", func(t *testing.T) {
				v, err := DecodeStr(tc.String).Float32()
				require.InEpsilonf(t, tc.Value, v, epsilon, "%v != %v", tc.Value, v)
				require.NoError(t, err)
			})
			t.Run("64Str", func(t *testing.T) {
				v, err := DecodeStr(tc.String).Float64()
				require.InEpsilonf(t, tc.Value, v, epsilon, "%v != %v", tc.Value, v)
				require.NoError(t, err)
			})
			t.Run("32", func(t *testing.T) {
				decodeStr(t, tc.String, func(t *testing.T, d *Decoder) {
					v, err := d.Float32()
					require.NoError(t, err)
					require.InEpsilonf(t, tc.Value, v, epsilon, "%v != %v", tc.Value, v)
				})
			})
			t.Run("64", func(t *testing.T) {
				decodeStr(t, tc.String, func(t *testing.T, d *Decoder) {
					v, err := d.Float64()
					require.NoError(t, err)
					require.InEpsilonf(t, tc.Value, v, epsilon, "%v != %v", tc.Value, v)
				})
			})
		})
	}
}

func TestDecoderFloatUnexpectedChar(t *testing.T) {
	type floatFunc struct {
		name    string
		bitSize int
		fn      func(*Decoder) error
	}
	floatFuncs := []floatFunc{
		{"Float32", 32, decoderOnlyError((*Decoder).Float32)},
		{"Float64", 64, decoderOnlyError((*Decoder).Float64)},
	}

	tests := []struct {
		input       string
		size        int // 0 for any
		errContains string
	}{
		// Leading space.
		{" 10", 0, ""},
		{"   10", 0, ""},
		{" -10", 0, ""},

		// Digit after leading zero.
		{"00", 0, "leading zero: unexpected byte 48 '0' at 1"},
		{"01", 0, "leading zero: unexpected byte 49 '1' at 1"},
		{"-00", 0, "leading zero: unexpected byte 48 '0' at 2"},
		{"-01", 0, "leading zero: unexpected byte 49 '1' at 2"},

		// Double minus.
		{"--10", 0, "unexpected byte 45 '-' at 1"},

		// Leading dot.
		{".0", 0, "unexpected byte 46 '.' at 0"},
		// Leading exponent.
		{"e0", 0, "unexpected byte 101 'e' at 0"},
		{"E0", 0, "unexpected byte 69 'E' at 0"},

		// Non-digit after minus.
		{"-.0", 0, "unexpected byte 46 '.' at 1"},
		{"-e0", 0, "unexpected byte 101 'e' at 1"},
		{"-E0", 0, "unexpected byte 69 'E' at 1"},

		// Unexpected character.
		{"-a", 0, "unexpected byte 97 'a' at 1"},
		{"0a", 0, "unexpected byte 97 'a' at 1"},
		{"0.a", 0, "unexpected byte 97 'a' at 2"},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			check := func(fns []floatFunc) {
				for _, intFn := range fns {
					intFn := intFn
					if tt.size != 0 && tt.size > intFn.bitSize {
						continue
					}
					t.Run(intFn.name, func(t *testing.T) {
						decodeStr(t, tt.input, func(t *testing.T, d *Decoder) {
							a := assert.New(t)

							err := intFn.fn(d)
							if e := tt.errContains; e != "" {
								a.ErrorContains(err, e)
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

			check(floatFuncs)
		})
	}
}

func BenchmarkDecoder_Float64(b *testing.B) {
	for _, file := range []string{
		"floats.json",
		"slow_floats.json",
		"integers.json",
	} {
		runTestdataFile(file, b.Fatal, func(name string, data []byte) {
			b.Run(name, func(b *testing.B) {
				d := GetDecoder()
				cb := func(d *Decoder) error {
					_, err := d.Float64()
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
}
