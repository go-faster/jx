package jx

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
)

func fuzzCallback[Int constraints.Signed](decoder func(*Decoder) (Int, error)) func(*testing.T, Int) {
	return func(t *testing.T, expected Int) {
		a := require.New(t)
		buf := make([]byte, 0, 32)
		buf = strconv.AppendInt(buf, int64(expected), 10)

		d := DecodeBytes(buf)
		got, err := decoder(d)
		a.NoError(err)
		a.Equal(expected, got)
	}
}

func FuzzDecoderInt8(f *testing.F)  { f.Fuzz(fuzzCallback[int8]((*Decoder).Int8)) }
func FuzzDecoderInt16(f *testing.F) { f.Fuzz(fuzzCallback[int16]((*Decoder).Int16)) }
func FuzzDecoderInt32(f *testing.F) { f.Fuzz(fuzzCallback[int32]((*Decoder).Int32)) }
func FuzzDecoderInt64(f *testing.F) { f.Fuzz(fuzzCallback[int64]((*Decoder).Int64)) }
