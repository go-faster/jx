package jx

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_read_uint64_invalid(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(",")
	_, err := iter.Uint64()
	should.Error(err)
}

func Test_read_int32(t *testing.T) {
	inputs := []string{`1`, `12`, `123`, `1234`, `12345`, `123456`, `2147483647`, `-2147483648`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := DecodeStr(input)
			expected, err := strconv.ParseInt(input, 10, 32)
			should.NoError(err)
			v, err := iter.Int32()
			should.NoError(err)
			should.Equal(int32(expected), v)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Decode(bytes.NewBufferString(input), 2)
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
			iter := DecodeStr(s)
			_, err := iter.Int32()
			should.Error(err)

			iterUint := DecodeStr(s)
			_, err = iterUint.Uint32()
			should.Error(err)
		})
	}

	for _, s := range []string{"123456789232323232321545111111111111111111111111111111145454545445", "-1234567892323232323212"} {
		t.Run(s, func(t *testing.T) {
			should := require.New(t)
			iter := DecodeStr(s)
			v, err := iter.Int64()
			should.Error(err, "%v", v)

			iterUint := DecodeStr(s)
			vu, err := iterUint.Uint64()
			should.Error(err, "%v", vu)
		})
	}
}

func Test_read_int64_overflow(t *testing.T) {
	s := `123456789232323232321545111111111111111111111111111111145454545445`
	iter := DecodeStr(s)
	_, err := iter.Int64()
	require.Error(t, err)
}

func Test_read_int64(t *testing.T) {
	inputs := []string{`1`, `12`, `123`, `1234`, `12345`, `123456`, `9223372036854775807`, `-9223372036854775808`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := DecodeStr(input)
			expected, err := strconv.ParseInt(input, 10, 64)
			should.NoError(err)
			v, err := iter.Int64()
			should.NoError(err)
			should.Equal(expected, v)
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Decode(bytes.NewBufferString(input), 2)
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
			e := NewEncoder()
			e.Uint32(val)
			should.Equal(strconv.FormatUint(uint64(val), 10), e.String())
		})
	}
	should := require.New(t)
	e := NewEncoder()
	e.Raw("a")
	e.Uint32(0xffffffff) // should clear buffer
	should.Equal("a4294967295", e.String())
}

func Test_write_int32(t *testing.T) {
	vals := []int32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0x7fffffff, -0x80000000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			e := NewEncoder()
			e.Int32(val)
			should.Equal(strconv.FormatInt(int64(val), 10), e.String())
		})
	}
	should := require.New(t)
	e := NewEncoder()
	e.Raw("a")
	e.Int32(-0x7fffffff) // should clear buffer
	should.Equal("a-2147483647", e.String())
}

func Test_write_uint64(t *testing.T) {
	vals := []uint64{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff, 0xffffffffff, 0xfffffffffff, 0xffffffffffff, 0xfffffffffffff, 0xffffffffffffff,
		0xfffffffffffffff, 0xffffffffffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			e := NewEncoder()
			e.Uint64(val)
			should.Equal(strconv.FormatUint(val, 10), e.String())
		})
	}
	should := require.New(t)
	e := NewEncoder()
	e.Raw("a")
	e.Uint64(0xffffffff) // should clear buffer
	should.Equal("a4294967295", e.String())
}

func Test_write_int64(t *testing.T) {
	vals := []int64{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff, 0xffffffffff, 0xfffffffffff, 0xffffffffffff, 0xfffffffffffff, 0xffffffffffffff,
		0xfffffffffffffff, 0x7fffffffffffffff, -0x8000000000000000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			e := NewEncoder()
			e.Int64(val)
			should.Equal(strconv.FormatInt(val, 10), e.String())
		})
	}
	should := require.New(t)
	e := NewEncoder()
	e.Raw("a")
	e.Int64(0xffffffff) // should clear buffer
	should.Equal("a4294967295", e.String())
}

func intPow(n, m int64) int64 {
	if m == 0 {
		return 1
	}
	result := n
	for i := int64(2); i <= m; i++ {
		result *= n
	}
	return result
}

func requireArrEnd(t testing.TB, d *Decoder) {
	t.Helper()
	ok, err := d.Elem()
	require.False(t, ok)
	require.NoError(t, err)
	requireEOF(t, d)
}

func requireElem(t testing.TB, d *Decoder) {
	t.Helper()
	ok, err := d.Elem()
	require.True(t, ok)
	require.NoError(t, err)
}

func requireEOF(t testing.TB, d *Decoder) {
	t.Helper()
	require.ErrorIs(t, d.Skip(), io.EOF)
}

func TestDecoder_Int64(t *testing.T) {
	// Generate some diverse numbers.
	for i := int64(0); i < 40; i++ {
		v := int64(3)
		for k := int64(0); k < i; k++ {
			// No special meaning, just trying to make digits more diverse.
			v += i + 1
			v += intPow(10, k) * (k%7 + 1)
		}
		t.Run(fmt.Sprintf("%d", v), func(t *testing.T) {
			e := GetEncoder()
			e.ArrStart()
			e.Int64(v)
			e.More()
			e.Int64(-v)
			e.ArrEnd()

			d := DecodeBytes(e.Bytes())
			requireElem(t, d)
			got, err := d.Int64()
			require.NoError(t, err)
			require.Equal(t, v, got)
			requireElem(t, d)
			got, err = d.Int64()
			require.NoError(t, err)
			require.Equal(t, -v, got)
			requireArrEnd(t, d)
		})
	}
}

func int32Pow(n, m int32) int32 {
	if m == 0 {
		return 1
	}
	result := n
	for i := int32(2); i <= m; i++ {
		result *= n
	}
	return result
}

func TestDecoder_Int32(t *testing.T) {
	// Generate some diverse numbers.
	for i := int32(0); i < 10; i++ {
		v := int32(3)
		for k := int32(0); k < i; k++ {
			// No special meaning, just trying to make digits more diverse.
			v += i + 1
			v += int32Pow(10, k) * (k%7 + 1)
		}
		t.Run(fmt.Sprintf("%d", v), func(t *testing.T) {
			e := GetEncoder()
			e.ArrStart()
			e.Int32(v)
			e.More()
			e.Int32(-v)
			e.ArrEnd()

			d := DecodeBytes(e.Bytes())
			requireElem(t, d)
			got, err := d.Int32()
			require.NoError(t, err)
			require.Equal(t, v, got)
			requireElem(t, d)
			got, err = d.Int32()
			require.NoError(t, err)
			require.Equal(t, -v, got)
			requireArrEnd(t, d)
		})
	}
}

func uintPow(n, m uint64) uint64 {
	if m == 0 {
		return 1
	}
	result := n
	for i := uint64(2); i <= m; i++ {
		result *= n
	}
	return result
}

func TestDecoder_Uint64(t *testing.T) {
	// Generate some diverse numbers.
	for i := uint64(0); i < 28; i++ {
		v := uint64(3)
		for k := uint64(0); k < i; k++ {
			// No special meaning, just trying to make digits more diverse.
			v += i + 1
			v += uintPow(10, k) * (k%7 + 1)
		}
		t.Run(fmt.Sprintf("%d", v), func(t *testing.T) {
			e := GetEncoder()
			e.ArrStart()
			e.Uint64(v)
			e.ArrEnd()

			d := GetDecoder()
			d.ResetBytes(e.Bytes())
			requireElem(t, d)
			got, err := d.Uint64()
			require.NoError(t, err)
			require.Equal(t, v, got)
			requireArrEnd(t, d)
		})
	}
}

func uint32Pow(n, m uint32) uint32 {
	if m == 0 {
		return 1
	}
	result := n
	for i := uint32(2); i <= m; i++ {
		result *= n
	}
	return result
}

func TestDecoder_Uint32(t *testing.T) {
	// Generate some diverse numbers.
	for i := uint32(0); i < 20; i++ {
		v := uint32(3)
		for k := uint32(0); k < i; k++ {
			// No special meaning, just trying to make digits more diverse.
			v += i + 1
			v += uint32Pow(10, k) * (k%7 + 1)
		}
		t.Run(fmt.Sprintf("%d", v), func(t *testing.T) {
			e := GetEncoder()
			e.ArrStart()
			e.Uint32(v)
			e.More()
			e.Uint(uint(v))
			e.ArrEnd()

			d := GetDecoder()
			d.ResetBytes(e.Bytes())
			requireElem(t, d)
			got, err := d.Uint32()
			require.NoError(t, err)
			require.Equal(t, v, got)
			requireElem(t, d)
			gotUint, err := d.Uint()
			require.NoError(t, err)
			require.Equal(t, uint(v), gotUint)
			requireArrEnd(t, d)
		})
	}
}
