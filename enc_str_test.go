package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_StringEscape(t *testing.T) {
	s := NewEncoder(nil, 0)
	const data = `<html>Hello\\\n\r\\` + "\n\rWorld\u2028</html>"
	s.StringEscape(data)
	require.NoError(t, s.Flush())
	requireCompat(t, s.Bytes(), data)
	const expected = `"\u003chtml\u003eHello\\\\\\n\\r\\\\\n\rWorld\u2028\u003c/html\u003e"`
	require.Equal(t, expected, string(s.Bytes()))
}

func TestEncoder_String(t *testing.T) {
	s := NewEncoder(nil, 0)
	const data = `\nH\tel\tl\ro\\World\r` + "\n\rHello\r\tHi"
	s.String(data)
	require.NoError(t, s.Flush())
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
