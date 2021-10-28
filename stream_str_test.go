package jir

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStream_WriteStringWithHTMLEscaped(t *testing.T) {
	s := NewStream(Default, nil, 0)
	const data = `<html>Привет, мир!</html>`
	s.WriteStringWithHTMLEscaped(data)
	require.NoError(t, s.Error)
	const expected = `"\u003chtml\u003eПривет, мир!\u003c/html\u003e"`
	require.Equal(t, expected, string(s.Buffer()))
}
