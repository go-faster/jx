package jx

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder_Str(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{``},
		{`abcd`},
		{
			`abcd\nH\tel\tl\ro\\World\r` + "\n\rHello\r\tHi",
		},
		{"\x00"},
		{"\x00 "},
		{`"hello, world!"`},

		{strings.Repeat("a", encoderBufSize)},
	}
	for i, tt := range testCases {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			for _, enc := range []struct {
				name string
				enc  func(e *Encoder, input string) bool
			}{
				{"Str", (*Encoder).Str},
				{"Bytes", func(e *Encoder, input string) bool {
					return e.ByteStr([]byte(tt.input))
				}},
			} {
				enc := enc
				t.Run(enc.name, func(t *testing.T) {
					requireCompat(t, func(e *Encoder) {
						enc.enc(e, tt.input)
					}, tt.input)

					t.Run("Decode", func(t *testing.T) {
						e := GetEncoder()
						enc.enc(e, tt.input)

						i := GetDecoder()
						i.ResetBytes(e.Bytes())
						s, err := i.Str()
						require.NoError(t, err)
						require.Equal(t, tt.input, s)
					})
				})
			}
		})
	}
	t.Run("Quotes", func(t *testing.T) {
		const (
			v = "\"/\""
		)
		requireCompat(t, func(e *Encoder) {
			e.StrEscape(v)
		}, v)
	})
	t.Run("QuotesObj", func(t *testing.T) {
		const (
			k = "k"
			v = "\"/\""
		)

		cb := func(e *Encoder) {
			e.ObjStart()
			e.FieldStart(k)
			e.Str(v)
			e.ObjEnd()
			t.Log(e)
		}

		var e Encoder
		cb(&e)

		var target map[string]string
		require.NoError(t, json.Unmarshal(e.Bytes(), &target))
		assert.Equal(t, v, target[k])
		requireCompat(t, cb, map[string]string{k: v})
	})
}

func TestEncoder_StrEscape(t *testing.T) {
	testCases := []struct {
		input, expect string
	}{
		{"Foo", `"Foo"`},
		{"\uFFFD", `"�"`},
		{"a\xc5z", `"a\ufffdz"`},
		{"<f\xed\xa0\x80", `"\u003cf\ufffd\ufffd\ufffd"`},
		{
			`<html>Hello\\\n\r\\` + "\n\rW\torld\u2028</html>",
			`"\u003chtml\u003eHello\\\\\\n\\r\\\\\n\rW\torld\u2028\u003c/html\u003e"`,
		},
	}
	for i, tt := range testCases {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			for _, enc := range []struct {
				name string
				enc  func(e *Encoder, input string) bool
			}{
				{"Str", (*Encoder).StrEscape},
				{"Bytes", func(e *Encoder, input string) bool {
					return e.ByteStrEscape([]byte(tt.input))
				}},
			} {
				enc := enc
				t.Run(enc.name, func(t *testing.T) {
					requireCompat(t, func(e *Encoder) {
						enc.enc(e, tt.input)
					}, tt.input)
				})
			}
		})
	}
	t.Run("QuotesEscape", func(t *testing.T) {
		const (
			v = "\"/\""
		)
		requireCompat(t, func(e *Encoder) {
			e.StrEscape(v)
		}, v)
	})
}
