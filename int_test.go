package jx

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_read_uint64_invalid(t *testing.T) {
	should := require.New(t)
	iter := ReadString(",")
	_, err := iter.Uint64()
	should.Error(err)
}

func Test_read_int32(t *testing.T) {
	inputs := []string{`1`, `12`, `123`, `1234`, `12345`, `123456`, `2147483647`, `-2147483648`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ReadString(input)
			expected, err := strconv.ParseInt(input, 10, 32)
			should.NoError(err)
			v, err := iter.Int32()
			should.NoError(err)
			should.Equal(int32(expected), v)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Read(bytes.NewBufferString(input), 2)
			expected, err := strconv.ParseInt(input, 10, 32)
			should.NoError(err)
			v, err := iter.Int32()
			should.NoError(err)
			should.Equal(int32(expected), v)
		})
	}
}

func Test_read_int_overflow(t *testing.T) {
	for _, s := range []string{"1234232323232323235678912", "-1234567892323232323212"} {
		t.Run(s, func(t *testing.T) {
			should := require.New(t)
			iter := ReadString(s)
			_, err := iter.Int32()
			should.Error(err)

			iterUint := ReadString(s)
			_, err = iterUint.Uint32()
			should.Error(err)
		})
	}

	for _, s := range []string{"123456789232323232321545111111111111111111111111111111145454545445", "-1234567892323232323212"} {
		t.Run(s, func(t *testing.T) {
			should := require.New(t)
			iter := ReadString(s)
			v, err := iter.Int64()
			should.Error(err, "%v", v)

			iterUint := ReadString(s)
			vu, err := iterUint.Uint64()
			should.Error(err, "%v", vu)
		})
	}
}

func Test_read_int64_overflow(t *testing.T) {
	s := `123456789232323232321545111111111111111111111111111111145454545445`
	iter := ReadString(s)
	_, err := iter.Int64()
	require.Error(t, err)
}

func Test_read_int64(t *testing.T) {
	inputs := []string{`1`, `12`, `123`, `1234`, `12345`, `123456`, `9223372036854775807`, `-9223372036854775808`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ReadString(input)
			expected, err := strconv.ParseInt(input, 10, 64)
			should.NoError(err)
			v, err := iter.Int64()
			should.NoError(err)
			should.Equal(expected, v)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Read(bytes.NewBufferString(input), 2)
			expected, err := strconv.ParseInt(input, 10, 64)
			should.NoError(err)
			v, err := iter.Int64()
			should.NoError(err)
			should.Equal(expected, v)
		})
	}
}

func Test_write_uint32(t *testing.T) {
	vals := []uint32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewWriter(buf, 4096)
			stream.Uint32(val)
			should.NoError(stream.Flush())
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewWriter(buf, 10)
	stream.Raw("a")
	stream.Uint32(0xffffffff) // should clear buffer
	should.NoError(stream.Flush())
	should.Equal("a4294967295", buf.String())
}

func Test_write_int32(t *testing.T) {
	vals := []int32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0x7fffffff, -0x80000000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewWriter(buf, 4096)
			stream.Int32(val)
			should.NoError(stream.Flush())
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewWriter(buf, 11)
	stream.Raw("a")
	stream.Int32(-0x7fffffff) // should clear buffer
	should.NoError(stream.Flush())
	should.Equal("a-2147483647", buf.String())
}

func Test_write_uint64(t *testing.T) {
	vals := []uint64{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff, 0xffffffffff, 0xfffffffffff, 0xffffffffffff, 0xfffffffffffff, 0xffffffffffffff,
		0xfffffffffffffff, 0xffffffffffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewWriter(buf, 4096)
			stream.Uint64(val)
			should.NoError(stream.Flush())
			should.Equal(strconv.FormatUint(val, 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewWriter(buf, 10)
	stream.Raw("a")
	stream.Uint64(0xffffffff) // should clear buffer
	should.NoError(stream.Flush())
	should.Equal("a4294967295", buf.String())
}

func Test_write_int64(t *testing.T) {
	vals := []int64{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff, 0xffffffffff, 0xfffffffffff, 0xffffffffffff, 0xfffffffffffff, 0xffffffffffffff,
		0xfffffffffffffff, 0x7fffffffffffffff, -0x8000000000000000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewWriter(buf, 4096)
			stream.Int64(val)
			should.NoError(stream.Flush())
			should.Equal(strconv.FormatInt(val, 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewWriter(buf, 10)
	stream.Raw("a")
	stream.Int64(0xffffffff) // should clear buffer
	should.NoError(stream.Flush())
	should.Equal("a4294967295", buf.String())
}
