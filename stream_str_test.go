package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStream_WriteStringWithHTMLEscaped(t *testing.T) {
	s := NewStream(Default, nil, 0)
	const data = `<html>Hello\\\n\r\\` + "\n\rWorld\u2028</html>"
	s.WriteStringWithHTMLEscaped(data)
	require.NoError(t, s.Error)
	requireCompat(t, s.Buffer(), data)
	const expected = `"\u003chtml\u003eHello\\\\\\n\\r\\\\\n\rWorld\u2028\u003c/html\u003e"`
	require.Equal(t, expected, string(s.Buffer()))
}

func TestStream_WriteString(t *testing.T) {
	s := NewStream(Default, nil, 0)
	const data = `\nH\tel\tl\ro\\World\r` + "\n\rHello\r\tHi"
	s.WriteString(data)
	require.NoError(t, s.Error)
	const expected = `"\\nH\\tel\\tl\\ro\\\\World\\r\n\rHello\r\tHi"`
	require.Equal(t, expected, string(s.Buffer()))
	requireCompat(t, s.Buffer(), data)
	t.Run("Read", func(t *testing.T) {
		i := NewIter(Default)
		i.ResetBytes(s.Buffer())
		s, err := i.Str()
		require.NoError(t, err)
		require.Equal(t, data, s)
	})
}
