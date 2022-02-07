package jx

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func decodeStr(t *testing.T, s string, f func(d *Decoder)) {
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
			Fn:   func() *Decoder { return Decode(bytes.NewBufferString(s), 0) },
		},
		{
			Name: "DecodeSingleByteBuf",
			Fn:   func() *Decoder { return Decode(bytes.NewBufferString(s), 1) },
		},
	} {
		t.Run(d.Name, func(t *testing.T) {
			t.Helper()
			dec := d.Fn()
			f(dec)
		})
	}
}

func TestDecoder_Float(t *testing.T) {
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
	for _, tc := range []struct {
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
		t.Run(tc.String, func(t *testing.T) {
			t.Run("32Str", func(t *testing.T) {
				v, err := DecodeStr(tc.String).Float32()
				require.InEpsilon(t, tc.Value, v, epsilon)
				require.NoError(t, err)
			})
			t.Run("64Str", func(t *testing.T) {
				v, err := DecodeStr(tc.String).Float64()
				require.InEpsilon(t, tc.Value, v, epsilon)
				require.NoError(t, err)
			})
			t.Run("32", func(t *testing.T) {
				decodeStr(t, tc.String, func(d *Decoder) {
					v, err := d.Float32()
					require.NoError(t, err)
					require.InEpsilon(t, tc.Value, v, epsilon)
				})
			})
			t.Run("64", func(t *testing.T) {
				decodeStr(t, tc.String, func(d *Decoder) {
					v, err := d.Float64()
					require.NoError(t, err)
					require.InEpsilon(t, tc.Value, v, epsilon)
				})
			})
		})
	}
}

func BenchmarkDecoder_Float64(b *testing.B) {
	b.Run("Array", func(b *testing.B) {
		run := func(name string, data []byte) {
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
		}

		runTestdataFile("floats.json", b.Fatal, run)
		runTestdataFile("integers.json", b.Fatal, run)
	})
	b.Run("One", func(b *testing.B) {
		run := func(name string, data []byte) {
			b.Run(name, func(b *testing.B) {
				d := GetDecoder()
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					d.ResetBytes(data)

					if _, err := d.Float64(); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
		run("Float", []byte("0.31152244431052484"))
		run("Integer", []byte("31152244431052484"))
		run("Zero", []byte("0"))
	})
}
