package jx

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func testDecoderNum(t *testing.T, num func(*Decoder) (Num, error)) {
	t.Run("Cases", func(t *testing.T) {
		runTestCases(t, testNumbers, func(t *testing.T, d *Decoder) error {
			if _, err := num(d); err != nil {
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

	t.Run("StrEscape", func(t *testing.T) {
		a := require.New(t)
		for _, s := range []string{
			`"\u0030"`,                   // hex 30 = dec 48 = '0'
			`"\u002d\u0031\u0030\u0030"`, // "-100", but escaped
		} {
			_, err := num(DecodeStr(s))
			a.NoErrorf(err, "input: %q", s)
		}
	})
	t.Run("Positive", func(t *testing.T) {
		a := require.New(t)
		for _, s := range []string{
			`100`,
			`100.0`,
			`-100.0`,
			`-100`,
			`"-100"`,
			`"-100.0"`,
		} {
			v, err := num(DecodeStr(s))
			a.NoErrorf(err, "input: %q", s)
			a.Equalf(s, v.String(), "input: %q", s)
		}
	})
	t.Run("Negative", func(t *testing.T) {
		a := require.New(t)
		for _, s := range []string{
			`1.00.0`,
			`"-100`,
			`""`,
			`"-100.0.0"`,
			"false",
			`"false"`,
		} {
			_, err := num(DecodeStr(s))
			a.Errorf(err, "input: %q", s)
		}
	})
}

func TestDecoder_Num(t *testing.T) {
	testDecoderNum(t, (*Decoder).Num)
}

func TestDecoder_NumAppend(t *testing.T) {
	testDecoderNum(t, func(d *Decoder) (Num, error) {
		return d.NumAppend(nil)
	})
}
