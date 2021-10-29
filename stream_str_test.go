package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStream_WriteStringWithHTMLEscaped(t *testing.T) {
	s := NewStream(Default, nil, 0)
	const data = `<html>Hello\\\n\r\\` + "\n\rWorld\u2028</html>"
	s.StrHTMLEscaped(data)
	require.NoError(t, s.Flush())
	requireCompat(t, s.Buf(), data)
	const expected = `"\u003chtml\u003eHello\\\\\\n\\r\\\\\n\rWorld\u2028\u003c/html\u003e"`
	require.Equal(t, expected, string(s.Buf()))
}

func TestStream_WriteString(t *testing.T) {
	s := NewStream(Default, nil, 0)
	const data = `\nH\tel\tl\ro\\World\r` + "\n\rHello\r\tHi"
	s.Str(data)
	require.NoError(t, s.Flush())
	const expected = `"\\nH\\tel\\tl\\ro\\\\World\\r\n\rHello\r\tHi"`
	require.Equal(t, expected, string(s.Buf()))
	requireCompat(t, s.Buf(), data)
	t.Run("Read", func(t *testing.T) {
		i := NewIter(Default)
		i.ResetBytes(s.Buf())
		s, err := i.Str()
		require.NoError(t, err)
		require.Equal(t, data, s)
	})
}
