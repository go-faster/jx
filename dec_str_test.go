package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Str(t *testing.T) {
	i := GetDecoder()
	i.ResetBytes([]byte(`"hello, world!"`))
	s, err := i.String()
	require.NoError(t, err)
	require.Equal(t, "hello, world!", s)
}
