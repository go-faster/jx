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
	iter := ParseString(Default, ",")
	_, err := iter.Uint64()
	should.Error(err)
}

func Test_read_int32(t *testing.T) {
	inputs := []string{`1`, `12`, `123`, `1234`, `12345`, `123456`, `2147483647`, `-2147483648`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(Default, input)
			expected, err := strconv.ParseInt(input, 10, 32)
			should.NoError(err)
			v, err := iter.Int32()
			should.NoError(err)
			should.Equal(int32(expected), v)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(Default, bytes.NewBufferString(input), 2)
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
			iter := ParseString(Default, s)
			_, err := iter.Int32()
			should.Error(err)

			iterUint := ParseString(Default, s)
			_, err = iterUint.Uint32()
			should.Error(err)
		})
	}

	for _, s := range []string{"123456789232323232321545111111111111111111111111111111145454545445", "-1234567892323232323212"} {
		t.Run(s, func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(Default, s)
			v, err := iter.Int64()
			should.Error(err, "%v", v)

			iterUint := ParseString(Default, s)
			vu, err := iterUint.Uint64()
			should.Error(err, "%v", vu)
		})
	}
}

func Test_read_int64_overflow(t *testing.T) {
	s := `123456789232323232321545111111111111111111111111111111145454545445`
	iter := ParseString(Default, s)
	_, err := iter.Int64()
	require.Error(t, err)
}

func Test_read_int64(t *testing.T) {
	inputs := []string{`1`, `12`, `123`, `1234`, `12345`, `123456`, `9223372036854775807`, `-9223372036854775808`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(Default, input)
			expected, err := strconv.ParseInt(input, 10, 64)
			should.NoError(err)
			v, err := iter.Int64()
			should.NoError(err)
			should.Equal(expected, v)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(Default, bytes.NewBufferString(input), 2)
			expected, err := strconv.ParseInt(input, 10, 64)
			should.NoError(err)
			v, err := iter.Int64()
			should.NoError(err)
			should.Equal(expected, v)
		})
	}
}

func Test_write_uint8(t *testing.T) {
	vals := []uint8{0, 1, 11, 111, 255}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteUint8(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 3)
	stream.WriteRaw("a")
	stream.WriteUint8(100) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a100", buf.String())
}

func Test_write_int8(t *testing.T) {
	vals := []int8{0, 1, -1, 99, 0x7f, -0x80}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteInt8(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 4)
	stream.WriteRaw("a")
	stream.WriteInt8(-100) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a-100", buf.String())
}

func Test_write_uint16(t *testing.T) {
	vals := []uint16{0, 1, 11, 111, 255, 0xfff, 0xffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteUint16(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 5)
	stream.WriteRaw("a")
	stream.WriteUint16(10000) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a10000", buf.String())
}

func Test_write_int16(t *testing.T) {
	vals := []int16{0, 1, 11, 111, 255, 0xfff, 0x7fff, -0x8000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteInt16(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 6)
	stream.WriteRaw("a")
	stream.WriteInt16(-10000) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a-10000", buf.String())
}

func Test_write_uint32(t *testing.T) {
	vals := []uint32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteUint32(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 10)
	stream.WriteRaw("a")
	stream.WriteUint32(0xffffffff) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a4294967295", buf.String())
}

func Test_write_int32(t *testing.T) {
	vals := []int32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0x7fffffff, -0x80000000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteInt32(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 11)
	stream.WriteRaw("a")
	stream.WriteInt32(-0x7fffffff) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
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
			stream := NewStream(Default, buf, 4096)
			stream.WriteUint64(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(val, 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 10)
	stream.WriteRaw("a")
	stream.WriteUint64(0xffffffff) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
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
			stream := NewStream(Default, buf, 4096)
			stream.WriteInt64(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(val, 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 10)
	stream.WriteRaw("a")
	stream.WriteInt64(0xffffffff) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a4294967295", buf.String())
}
