package jx

import (
	"fmt"
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

	testNum := func(inputs []string, cb func(input string, t *testing.T, d *Decoder)) func(t *testing.T) {
		return func(t *testing.T) {
			for i, input := range inputs {
				input := input
				t.Run(fmt.Sprintf("Test%d", i+1), testBufferReader(input, func(t *testing.T, d *Decoder) {
					cb(input, t, d)
				}))
			}
		}
	}

	t.Run("StrEscape", testNum([]string{
		`"\u0030"`,                   // hex 30 = dec 48 = '0'
		`"\u002d\u0031\u0030\u0030"`, // "-100", but escaped
	}, func(input string, t *testing.T, d *Decoder) {
		_, err := num(d)
		require.NoErrorf(t, err, "input: %q", input)
	}))
	t.Run("Positive", testNum([]string{
		`100`,
		`100.0`,
		`-100.0`,
		`-100`,
		`"-100"`,
		`"-100.0"`,
	}, func(input string, t *testing.T, d *Decoder) {
		v, err := num(d)
		require.NoErrorf(t, err, "input: %q", input)
		require.Equalf(t, input, v.String(), "input: %q", input)
	}))
	t.Run("Negative", testNum([]string{
		`1.00.0`,
		`"-100`,
		`""`,
		`"-100.0.0"`,
		"false",
		`"false"`,
	}, func(input string, t *testing.T, d *Decoder) {
		_, err := num(d)
		require.Errorf(t, err, "input: %q", input)
	}))
}

func TestDecoder_Num(t *testing.T) {
	testDecoderNum(t, (*Decoder).Num)
}

func TestDecoder_NumAppend(t *testing.T) {
	testDecoderNum(t, func(d *Decoder) (Num, error) {
		return d.NumAppend(nil)
	})
}
