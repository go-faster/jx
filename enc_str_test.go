package jx

import (
	"encoding/json"
	"fmt"
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
	}
	for i, tt := range testCases {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			for _, enc := range []struct {
				name string
				enc  func(e *Encoder, input string)
			}{
				{"Str", (*Encoder).Str},
				{"Bytes", func(e *Encoder, input string) {
					e.ByteStr([]byte(tt.input))
				}},
			} {
				enc := enc
				t.Run(enc.name, func(t *testing.T) {
					e := GetEncoder()
					enc.enc(e, tt.input)
					requireCompat(t, e.Bytes(), tt.input)
					t.Run("Decode", func(t *testing.T) {
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
		var e Encoder
		e.Str(v)
		requireCompat(t, e.Bytes(), v)
	})
	t.Run("QuotesObj", func(t *testing.T) {
		const (
			k = "k"
			v = "\"/\""
		)

		var e Encoder
		e.ObjStart()
		e.FieldStart(k)
		e.Str(v)
		e.ObjEnd()
		t.Log(e)

		var target map[string]string
		require.NoError(t, json.Unmarshal(e.Bytes(), &target))
		assert.Equal(t, v, target[k])
		requireCompat(t, e.Bytes(), map[string]string{k: v})
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
		{`<html>Hello\\\n\r\\` + "\n\rW\torld\u2028</html>",
			`"\u003chtml\u003eHello\\\\\\n\\r\\\\\n\rW\torld\u2028\u003c/html\u003e"`},
	}
	for i, tt := range testCases {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			for _, enc := range []struct {
				name string
				enc  func(e *Encoder, input string)
			}{
				{"Str", (*Encoder).StrEscape},
				{"Bytes", func(e *Encoder, input string) {
					e.ByteStrEscape([]byte(tt.input))
				}},
			} {
				enc := enc
				t.Run(enc.name, func(t *testing.T) {
					e := GetEncoder()
					enc.enc(e, tt.input)
					require.Equal(t, tt.expect, string(e.Bytes()))
					requireCompat(t, e.Bytes(), tt.input)
				})
			}
		})
	}
	t.Run("QuotesEscape", func(t *testing.T) {
		const (
			v = "\"/\""
		)
		var e Encoder
		e.StrEscape(v)
		requireCompat(t, e.Bytes(), v)
	})
}
