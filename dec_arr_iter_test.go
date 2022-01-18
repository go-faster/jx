package jx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_ArrIter(t *testing.T) {
	testIter := func(d *Decoder) error {
		iter, err := d.ArrIter()
		if err != nil {
			return err
		}
		for {
			ok, err := iter.Next()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}

			if err := d.Skip(); err != nil {
				return err
			}
		}
	}
	for _, s := range testArrs {
		checker := require.Error
		if json.Valid([]byte(s)) {
			checker = require.NoError
		}

		d := DecodeStr(s)
		checker(t, testIter(d), s)
	}
}
