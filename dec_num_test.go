package jx

import (
	"fmt"
	"io"
	"strconv"
	"strings"
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

func BenchmarkDecoder_Num(b *testing.B) {
	number := strconv.FormatInt(1234567890421, 10)
	// escapeHex escapes the number as a string in \uXXXX format.
	escapeHex := func(number string) string {
		var b strings.Builder
		b.WriteByte('"')
		for _, c := range []byte(number) {
			b.WriteString("\\u00")
			b.WriteString(strconv.FormatInt(int64(c), 16))
		}
		b.WriteByte('"')
		return b.String()
	}

	for _, bt := range []struct {
		name  string
		input string
	}{
		{`Number`, number},
		{`String`, strconv.Quote(number)},
		{`EscapedString`, escapeHex(number)},
	} {
		bt := bt
		b.Run(bt.name, func(b *testing.B) {
			b.Run("Buffer", func(b *testing.B) {
				var (
					input = []byte(bt.input)
					d     = DecodeBytes(input)
					err   error
				)

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					d.ResetBytes(input)
					_, err = d.Num()
				}

				if err != nil {
					b.Fatal(err, bt.input)
				}
			})
			b.Run("Reader", func(b *testing.B) {
				var (
					r   = strings.NewReader(bt.input)
					d   = Decode(r, 512)
					err error
				)

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					r.Reset(bt.input)
					d.Reset(r)
					_, err = d.Num()
				}

				if err != nil {
					b.Fatal(err, bt.input)
				}
			})
		})
	}
}
