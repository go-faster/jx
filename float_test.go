package jx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-faster/errors"
)

// epsilon to compare floats.
const epsilon = 1e-6

func TestReadBigFloat(t *testing.T) {
	should := require.New(t)
	r := DecodeStr(`12.3`)
	val, err := r.BigFloat()
	should.NoError(err)
	val64, _ := val.Float64()
	should.Equal(12.3, val64)
}

func TestReadBigInt(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`92233720368547758079223372036854775807`)
	val, err := iter.BigInt()
	should.NoError(err)
	should.NotNil(val)
	should.Equal(`92233720368547758079223372036854775807`, val.String())
}

func TestEncodeInf(t *testing.T) {
	should := require.New(t)
	_, err := json.Marshal(math.Inf(1))
	should.Error(err)
	_, err = json.Marshal(float32(math.Inf(1)))
	should.Error(err)
	_, err = json.Marshal(math.Inf(-1))
	should.Error(err)
}

func TestEncodeNaN(t *testing.T) {
	should := require.New(t)
	_, err := json.Marshal(math.NaN())
	should.Error(err)
	_, err = json.Marshal(float32(math.NaN()))
	should.Error(err)
	_, err = json.Marshal(math.NaN())
	should.Error(err)
}

func TestReadFloat(t *testing.T) {
	inputs := []string{
		`1.1`,
		`1000`,
		`9223372036854775807`,
		`12.3`,
		`-12.3`,
		`720368.54775807`,
		`720368.547758075`,
		`1e1`,
		`1e+1`,
		`1e-1`,
		`1E1`,
		`1E+1`,
		`1E-1`,
		`-1e1`,
		`-1e+1`,
		`-1e-1`,
	}
	for i, input := range inputs {
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			// non-streaming
			t.Run("Float32", func(t *testing.T) {
				should := require.New(t)
				r := DecodeStr(input + ",")
				expected, err := strconv.ParseFloat(input, 32)
				should.NoError(err)
				got, err := r.Float32()
				should.NoError(err)
				should.Equal(float32(expected), got)
			})
			t.Run("Float64", func(t *testing.T) {
				should := require.New(t)
				r := DecodeStr(input + ",")
				expected, err := strconv.ParseFloat(input, 64)
				should.NoError(err)
				got, err := r.Float64()
				should.NoError(err)
				should.Equal(expected, got)
			})
			t.Run("Reader", func(t *testing.T) {
				should := require.New(t)
				iter := Decode(bytes.NewBufferString(input+","), 2)
				expected, err := strconv.ParseFloat(input, 32)
				should.NoError(err)
				got, err := iter.Float32()
				should.NoError(err)
				should.Equal(float32(expected), got)
			})
			t.Run("StdJSONCompliance", func(t *testing.T) {
				should := require.New(t)
				iter := Decode(bytes.NewBufferString(input+","), 2)
				val := float64(0)
				err := json.Unmarshal([]byte(input), &val)
				should.NoError(err)
				got, err := iter.Float64()
				should.NoError(err)
				should.Equal(val, got)
			})
		})
	}
}

func TestWriteFloat32(t *testing.T) {
	vals := []float32{
		0, 1, -1, 99, 0xff, 0xfff, 0xffff, 0xfffff, 0xffffff, 0x4ffffff, 0xfffffff,
		-0x4ffffff, -0xfffffff, 1.2345, 1.23456, 1.234567, 1.001,
	}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			requireCompat(t, func(e *Encoder) {
				e.Float32(val)
			}, val)
		})
	}
	should := require.New(t)
	e := GetEncoder()
	e.Float32(float32(0.0000001))
	should.Equal("1e-7", string(e.Bytes()))
}

func TestWriteFloat64(t *testing.T) {
	vals := []float64{
		0, 1, -1, 99, 0xff, 0xfff, 0xffff, 0xfffff, 0xffffff, 0x4ffffff, 0xfffffff,
		-0x4ffffff, -0xfffffff, 1.2345, 1.23456, 1.234567, 1.001,
	}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			requireCompat(t, func(e *Encoder) {
				e.Float64(val)
			}, val)
		})
	}
	should := require.New(t)
	e := GetEncoder()
	e.Float64(0.0000001)
	should.Equal("1e-7", e.String())
}

func TestEncoder_FloatError(t *testing.T) {
	e := NewStreamingEncoder(io.Discard, -1)
	e.w.stream.setError(errors.New("foo"))

	require.True(t, e.Float32(10))
	require.True(t, e.Float64(10))
	require.Error(t, e.Close())
}

func TestDecoder_FloatEOF(t *testing.T) {
	d := GetDecoder()

	_, err := d.Float32()
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestDecoder_FloatLeadingDot(t *testing.T) {
	v := `.0`
	_, err := DecodeStr(v).Float32()
	require.Error(t, err)
	_, err = DecodeStr(v).Float64()
	require.Error(t, err)
}

func TestDecoder_FloatReaderErr(t *testing.T) {
	getDecoder := func() *Decoder {
		d := Decode(errReader{}, -1)
		d.tail = 1
		d.buf = []byte{'-'}
		return d
	}
	_, err := getDecoder().Float32()
	require.Error(t, err)
	_, err = getDecoder().Float64()
	require.Error(t, err)
}

func TestEncoder_FloatNanInf(t *testing.T) {
	for _, f := range []float64{
		math.NaN(),
		math.Inf(-1),
		math.Inf(1),
	} {
		var e Encoder
		e.ArrStart()
		e.Float64(f)
		e.Float32(float32(f))
		e.ArrEnd()

		d := DecodeBytes(e.Bytes())
		requireElem(t, d)
		require.NoError(t, d.Null())
		requireElem(t, d)
		require.NoError(t, d.Null())
	}
}
