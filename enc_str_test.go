package jx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder_StringEscape(t *testing.T) {
	s := GetEncoder()
	const data = `<html>Hello\\\n\r\\` + "\n\rW\torld\u2028</html>"
	s.StrEscape(data)
	requireCompat(t, s.Bytes(), data)
	const expected = `"\u003chtml\u003eHello\\\\\\n\\r\\\\\n\rW\torld\u2028\u003c/html\u003e"`
	require.Equal(t, expected, string(s.Bytes()))
}

func TestEncoder_String(t *testing.T) {
	testCases := []struct {
		name, input, expect string
	}{
		{`Empty`, ``, `""`},
		{`Simple`, `abcd`, `"abcd"`},
		{
			`Escape`,
			`abcd\nH\tel\tl\ro\\World\r` + "\n\rHello\r\tHi",
			`"abcd\\nH\\tel\\tl\\ro\\\\World\\r\n\rHello\r\tHi"`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			e := GetEncoder()
			e.Str(tt.input)
			require.Equal(t, tt.expect, string(e.Bytes()))
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
	t.Run("StrEscapeFast", func(t *testing.T) {
		e := GetEncoder()
		e.StrEscape("Foo")
		require.Equal(t, `"Foo"`, e.String())
	})
	t.Run("StrEscapeBad", func(t *testing.T) {
		e := GetEncoder()
		e.StrEscape("\uFFFD")
		require.Equal(t, `"�"`, e.String())
		v, err := DecodeBytes(e.Bytes()).Str()
		require.NoError(t, err)
		require.Equal(t, "�", v)
	})
	t.Run("BadUnicode", func(t *testing.T) {
		e := GetEncoder()
		e.StrEscape("a\xc5z")
		require.Equal(t, `"a\ufffdz"`, e.String())
		v, err := DecodeBytes(e.Bytes()).Str()
		require.NoError(t, err)
		require.Equal(t, "a�z", v)
	})
	t.Run("Emoji", func(t *testing.T) {
		e := GetEncoder()
		e.Str(string([]byte{240, 159, 144, 152}))
		v, err := DecodeBytes(e.Bytes()).Str()
		require.NoError(t, err)
		require.Equal(t, "🐘", v)
	})
	t.Run("BadUnicodeAfterSafeEscape", func(t *testing.T) {
		e := GetEncoder()
		e.StrEscape("<f\xed\xa0\x80")
		require.Equal(t, `"\u003cf\ufffd\ufffd\ufffd"`, e.String())
	})
	t.Run("Quotes", func(t *testing.T) {
		const (
			v = "\"/\""
		)
		var e Encoder
		e.Str(v)
		requireCompat(t, e.Bytes(), v)
	})
	t.Run("QuotesEscape", func(t *testing.T) {
		const (
			v = "\"/\""
		)
		var e Encoder
		e.StrEscape(v)
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
