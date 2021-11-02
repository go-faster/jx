package jx

import (
	"bytes"
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
	t.Run("Invalid", func(t *testing.T) {
		for _, s := range []string{
			``,
			`-`,
			`-.`,
			`.`,
			`.-`,
			`00`,
			`.00`,
			`00.1`,
		} {
			t.Run(s, func(t *testing.T) {
				t.Run("64", func(t *testing.T) {
					decodeStr(t, s, func(d *Decoder) {
						_, err := d.Float64()
						require.Error(t, err, s)
					})
				})
				t.Run("32", func(t *testing.T) {
					decodeStr(t, s, func(d *Decoder) {
						_, err := d.Float32()
						require.Error(t, err, s)
					})
				})
			})
		}
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
	require.InEpsilon(t, 429496729.0, v, 1e-6)
}

func TestDecoder_Float64(t *testing.T) {
	t.Logf("%v", maxFloat64 < uint64SafeToMultiple10)
	t.Logf("%v", (uint64SafeToMultiple10*10) > maxFloat64)
	t.Run("64Str", func(t *testing.T) {
		v, err := DecodeStr(`0.184467440737095517`).Float64()
		require.InEpsilon(t, .184467440737095517, v, 1e-6)
		require.NoError(t, err)
	})
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
			const epsilon = 1e-6
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
