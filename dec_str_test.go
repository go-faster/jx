package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_StrAppend(t *testing.T) {
	s := `"Hello"`
	d := DecodeStr(s)
	var (
		data []byte
		err  error
	)
	data, err = d.StrAppend(data)
	require.NoError(t, err)
	require.Equal(t, "Hello", string(data))
}
