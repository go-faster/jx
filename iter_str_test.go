package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Str(t *testing.T) {
	i := Default.GetIter([]byte(`"hello, world!"`))
	s, err := i.Str()
	require.NoError(t, err)
	require.Equal(t, "hello, world!", s)
}
