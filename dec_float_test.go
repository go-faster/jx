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
