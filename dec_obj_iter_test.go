package jx

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-faster/errors"
)

func TestDecoder_ObjIter(t *testing.T) {
	testIter := func(d *Decoder) error {
		iter, err := d.ObjIter()
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
		if err := iter.Err(); err != nil {
			return err
		}

		// Check for any trialing json.
		if d.head != d.tail {
			if err := d.Skip(); err != io.EOF {
				return errors.Wrap(err, "unexpected trialing data")
			}
		}
		return nil
	}
	for i, s := range testObjs {
		s := s
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			checker := require.Error
			if json.Valid([]byte(s)) {
				checker = require.NoError
			}

			d := DecodeStr(s)
			checker(t, testIter(d), s)
		})
	}
	t.Run("Depth", func(t *testing.T) {
		d := DecodeStr(`{`)
		// Emulate depth
		d.depth = maxDepth
		require.ErrorIs(t, testIter(d), errMaxDepth)
	})
	t.Run("Empty", func(t *testing.T) {
		d := DecodeStr(``)
		require.ErrorIs(t, testIter(d), io.ErrUnexpectedEOF)
	})
	t.Run("Key", testBufferReader(`{"foo":1,"bar":1,"baz":1}`, func(t *testing.T, d *Decoder) {
		a := require.New(t)

		iter, err := d.ObjIter()
		a.NoError(err)

		var r []string
		for iter.Next() {
			r = append(r, string(iter.Key()))
			a.NoError(d.Skip())
		}
		a.False(iter.Next())
		a.NoError(iter.Err())

		a.Equal([]string{"foo", "bar", "baz"}, r)
	}))
}
