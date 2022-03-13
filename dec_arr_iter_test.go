package jx

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_ArrIter(t *testing.T) {
	testIter := func(d *Decoder) error {
		iter, err := d.ArrIter()
		if err != nil {
			return err
		}
		for iter.Next() {
			if err := d.Skip(); err != nil {
				return err
			}
		}
		if iter.Next() {
			panic("BUG")
		}
		return iter.Err()
	}
	for _, s := range testArrs {
		checker := require.Error
		if json.Valid([]byte(s)) {
			checker = require.NoError
		}

		d := DecodeStr(s)
		checker(t, testIter(d), s)
	}
	t.Run("Depth", func(t *testing.T) {
		d := DecodeStr(`[`)
		// Emulate depth
		d.depth = maxDepth
		require.ErrorIs(t, testIter(d), errMaxDepth)
	})
	t.Run("Empty", func(t *testing.T) {
		d := DecodeStr(``)
		require.ErrorIs(t, testIter(d), io.ErrUnexpectedEOF)
	})
}
