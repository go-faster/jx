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
)

func Test_read_big_float(t *testing.T) {
	should := require.New(t)
	r := DecodeStr(`12.3`)
	val, err := r.BigFloat()
	should.NoError(err)
	val64, _ := val.Float64()
	should.Equal(12.3, val64)
}

func Test_read_big_int(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`92233720368547758079223372036854775807`)
	val, err := iter.BigInt()
	should.NoError(err)
	should.NotNil(val)
	should.Equal(`92233720368547758079223372036854775807`, val.String())
}

func Test_read_number(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`92233720368547758079223372036854775807`)
	val, err := iter.Number()
	should.NoError(err)
	should.Equal(`92233720368547758079223372036854775807`, string(val))
}

func Test_encode_inf(t *testing.T) {
	should := require.New(t)
	_, err := json.Marshal(math.Inf(1))
	should.Error(err)
	_, err = json.Marshal(float32(math.Inf(1)))
	should.Error(err)
	_, err = json.Marshal(math.Inf(-1))
	should.Error(err)
}

func Test_encode_nan(t *testing.T) {
	should := require.New(t)
	_, err := json.Marshal(math.NaN())
	should.Error(err)
	_, err = json.Marshal(float32(math.NaN()))
	should.Error(err)
	_, err = json.Marshal(math.NaN())
	should.Error(err)
}

func Test_read_float(t *testing.T) {
	inputs := []string{
		`1.1`, `1000`, `9223372036854775807`, `12.3`, `-12.3`, `720368.54775807`, `720368.547758075`,
		`1e1`, `1e+1`, `1e-1`, `1E1`, `1E+1`, `1E-1`, `-1e1`, `-1e+1`, `-1e-1`,
	}
	for _, input := range inputs {
		// non-streaming
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			r := DecodeStr(input + ",")
			expected, err := strconv.ParseFloat(input, 32)
			should.NoError(err)
			got, err := r.Float32()
			should.NoError(err)
			should.Equal(float32(expected), got)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			r := DecodeStr(input + ",")
			expected, err := strconv.ParseFloat(input, 64)
			should.NoError(err)
			got, err := r.Float64()
			should.NoError(err)
			should.Equal(expected, got)
		})
		// streaming
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Decode(bytes.NewBufferString(input+","), 2)
			expected, err := strconv.ParseFloat(input, 32)
			should.NoError(err)
			got, err := iter.Float32()
			should.NoError(err)
			should.Equal(float32(expected), got)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Decode(bytes.NewBufferString(input+","), 2)
			val := float64(0)
			err := json.Unmarshal([]byte(input), &val)
			should.NoError(err)
			got, err := iter.Float64()
			should.NoError(err)
			should.Equal(val, got)
		})
	}
}

func Test_write_float32(t *testing.T) {
	vals := []float32{0, 1, -1, 99, 0xff, 0xfff, 0xffff, 0xfffff, 0xffffff, 0x4ffffff, 0xfffffff,
		-0x4ffffff, -0xfffffff, 1.2345, 1.23456, 1.234567, 1.001}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			w := GetEncoder()
			w.Float32(val)
			output, err := json.Marshal(val)
			should.Nil(err)
			should.Equal(output, w.Bytes())
		})
	}
	should := require.New(t)
	e := GetEncoder()
	e.Float32(float32(0.0000001))
	should.Equal("1e-7", string(e.Bytes()))
}

func Test_write_float64(t *testing.T) {
	vals := []float64{0, 1, -1, 99, 0xff, 0xfff, 0xffff, 0xfffff, 0xffffff, 0x4ffffff, 0xfffffff,
		-0x4ffffff, -0xfffffff, 1.2345, 1.23456, 1.234567, 1.001}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			e := GetEncoder()
			e.Float64(val)
			s := strconv.FormatFloat(val, 'f', -1, 64)
			should.Equal(s, string(e.Bytes()))
		})
	}
	should := require.New(t)
	e := GetEncoder()
	e.Float64(0.0000001)
	should.Equal("1e-7", e.String())
}

func TestDecoder_FloatEOF(t *testing.T) {
	d := GetDecoder()

	_, err := d.Float32()
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestEncoder_FloatNanInf(t *testing.T) {
	for _, f := range []float64{
		math.NaN(),
		math.Inf(-1),
		math.Inf(1),
	} {
		var e Encoder
		e.Float64(f)
		e.More()
		e.Float32(float32(f))

		d := DecodeBytes(e.Bytes())
		require.NoError(t, d.Null())
		requireElem(t, d)
		require.NoError(t, d.Null())
	}
}
