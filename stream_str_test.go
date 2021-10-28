package jir

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStream_WriteStringWithHTMLEscaped(t *testing.T) {
	s := NewStream(Default, nil, 0)
	const data = `<html>Hello\\\n\r\\` + "\n\rWorld\u2028</html>"
	s.WriteStringWithHTMLEscaped(data)
	require.NoError(t, s.Error)
	const expected = `"\u003chtml\u003eHello\\\\\\n\\r\\\\\n\rWorld\u2028\u003c/html\u003e"`
	require.Equal(t, expected, string(s.Buffer()))
}
