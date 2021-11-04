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

func TestNum(t *testing.T) {
	t.Run("Integer", func(t *testing.T) {
		v := Num{
			Format: NumFormatInt,
			Value:  []byte{'1', '2', '3'},
		}
		t.Run("Encode", func(t *testing.T) {
			var e Encoder
			e.Num(v)
			require.Equal(t, e.String(), "123")
		})
		t.Run("Methods", func(t *testing.T) {
			require.True(t, v.Positive())
			require.True(t, v.Format.Int())
			require.False(t, v.Format.Invalid())
			require.False(t, v.Negative())
			require.False(t, v.Zero())
			require.Equal(t, 1, v.Sign())
		})
	})
}
