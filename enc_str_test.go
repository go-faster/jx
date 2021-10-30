package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_StringEscape(t *testing.T) {
	s := NewEncoder()
	const data = `<html>Hello\\\n\r\\` + "\n\rWorld\u2028</html>"
	s.StrEscape(data)
	requireCompat(t, s.Bytes(), data)
	const expected = `"\u003chtml\u003eHello\\\\\\n\\r\\\\\n\rWorld\u2028\u003c/html\u003e"`
	require.Equal(t, expected, string(s.Bytes()))
}

func TestEncoder_String(t *testing.T) {
	s := NewEncoder()
	const data = `\nH\tel\tl\ro\\World\r` + "\n\rHello\r\tHi"
	s.Str(data)
	const expected = `"\\nH\\tel\\tl\\ro\\\\World\\r\n\rHello\r\tHi"`
	require.Equal(t, expected, string(s.Bytes()))
	requireCompat(t, s.Bytes(), data)
	t.Run("Decode", func(t *testing.T) {
		i := NewDecoder()
		i.ResetBytes(s.Bytes())
		s, err := i.String()
		require.NoError(t, err)
		require.Equal(t, data, s)
	})
}
