package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_Num(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		for _, s := range []string{
			`100`,
			`100.0`,
			`-100.0`,
			`-100`,
			`"-100"`,
			`"-100.0"`,
		} {
			v, err := DecodeStr(s).Num()
			require.NoError(t, err)
			require.Equal(t, s, v.String())

			v, err = DecodeStr(s).NumAppend(nil)
			require.NoError(t, err)
			require.Equal(t, s, v.String())
		}
	})
	t.Run("Negative", func(t *testing.T) {
		for _, s := range []string{
			`1.00.0`,
			`"-100`,
			`"-100.0.0"`,
			"false",
			`"false"`,
		} {
			_, err := DecodeStr(s).Num()
			require.Error(t, err, s)

			_, err = DecodeStr(s).NumAppend(nil)
			require.Error(t, err, s)
		}
	})
}
