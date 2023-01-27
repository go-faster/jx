package stream

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/segmentio/asm/base64"
	"github.com/stretchr/testify/require"
)

func TestEncoder_Base64(t *testing.T) {
	for i, data := range [][]byte{
		[]byte(`1`),
		[]byte(`12`),
		bytes.Repeat([]byte{1}, 256),
		bytes.Repeat([]byte{1}, defaultBufferSize-1),
		bytes.Repeat([]byte{1}, defaultBufferSize),
		bytes.Repeat([]byte{1}, defaultBufferSize+1),
	} {
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			var sb strings.Builder
			e := NewEncoder(&sb)
			e.Base64(data)

			require.NoError(t, e.Close())
			expected := fmt.Sprintf("%q", base64.StdEncoding.EncodeToString(data))
			got := sb.String()
			require.Equal(t, expected, got, "%#q != %#q", expected, got)
		})
	}

	var sb strings.Builder
	e := NewEncoder(&sb)
	e.Base64(nil)
	require.NoError(t, e.Close())
	require.Equal(t, "null", sb.String())
}
