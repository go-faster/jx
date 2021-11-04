package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_Num(t *testing.T) {
	var e Encoder
	e.Num(Num{
		Format: NumFormatInt,
		Value:  []byte{'1', '2', '3'},
	})
	require.Equal(t, e.String(), "123")
}
