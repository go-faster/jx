package jx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_read_big_float(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `12.3`)
	val, err := iter.BigFloat()
	should.NoError(err)
	val64, _ := val.Float64()
	should.Equal(12.3, val64)
}

func Test_read_big_int(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `92233720368547758079223372036854775807`)
	val, err := iter.BigInt()
	should.NoError(err)
	should.NotNil(val)
	should.Equal(`92233720368547758079223372036854775807`, val.String())
}

func Test_read_number(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, `92233720368547758079223372036854775807`)
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
			iter := ParseString(Default, input+",")
			expected, err := strconv.ParseFloat(input, 32)
			should.NoError(err)
			got, err := iter.Float32()
			should.NoError(err)
			should.Equal(float32(expected), got)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(Default, input+",")
			expected, err := strconv.ParseFloat(input, 64)
			should.NoError(err)
			got, err := iter.Float64()
			should.NoError(err)
			should.Equal(expected, got)
		})
		// streaming
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(Default, bytes.NewBufferString(input+","), 2)
			expected, err := strconv.ParseFloat(input, 32)
			should.NoError(err)
			got, err := iter.Float32()
			should.NoError(err)
			should.Equal(float32(expected), got)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(Default, bytes.NewBufferString(input+","), 2)
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
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteFloat32Lossy(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			output, err := json.Marshal(val)
			should.Nil(err)
			should.Equal(string(output), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 10)
	stream.WriteRaw("abcdefg")
	stream.WriteFloat32Lossy(1.123456)
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("abcdefg1.123456", buf.String())

	stream = NewStream(Default, nil, 0)
	stream.WriteFloat32(float32(0.0000001))
	should.Equal("1e-07", string(stream.Buffer()))
}

func Test_write_float64(t *testing.T) {
	vals := []float64{0, 1, -1, 99, 0xff, 0xfff, 0xffff, 0xfffff, 0xffffff, 0x4ffffff, 0xfffffff,
		-0x4ffffff, -0xfffffff, 1.2345, 1.23456, 1.234567, 1.001}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteFloat64(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			s := strconv.FormatFloat(val, 'f', -1, 64)
			if !strings.Contains(s, ".") {
				s += ".0"
			}
			should.Equal(s, buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 10)
	stream.WriteRaw("abcdefg")
	stream.WriteFloat64Lossy(1.123456)
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("abcdefg1.123456", buf.String())

	stream = NewStream(Default, nil, 0)
	stream.WriteFloat64(0.0000001)
	should.Equal("1e-07", string(stream.Buffer()))
}
