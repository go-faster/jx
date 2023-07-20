package jx

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
)

func TestReadUint64Invalid(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(",")
	_, err := iter.UInt64()
	should.Error(err)
}

func TestDecoderIntNumbers(t *testing.T) {
	for i := 1; i < 10; i++ { // 10 digits
		var data []byte
		v := 0
		// Produce number like 123456, where 6 is i.
		for j := 1; j <= i; j++ {
			v = v*10 + j
			data = append(data, byte('0')+byte(j))
		}
		// Ensure that buffer length is at least 15, and it
		// has trialing comma.
		for j := 0; j < 15-i; j++ {
			data = append(data, ',')
		}
		s := string(data)
		t.Run("32", func(t *testing.T) {
			decodeStr(t, s, func(t *testing.T, d *Decoder) {
				got, err := d.Int32()
				require.NoError(t, err)
				require.Equal(t, int32(v), got)
			})
		})
		t.Run("64", func(t *testing.T) {
			decodeStr(t, s, func(t *testing.T, d *Decoder) {
				got, err := d.Int64()
				require.NoError(t, err)
				require.Equal(t, int64(v), got)
			})
		})
		t.Run("int", func(t *testing.T) {
			decodeStr(t, s, func(t *testing.T, d *Decoder) {
				got, err := d.Int()
				require.NoError(t, err)
				require.Equal(t, v, got)
			})
		})
		t.Run("uint", func(t *testing.T) {
			decodeStr(t, s, func(t *testing.T, d *Decoder) {
				got, err := d.UInt()
				require.NoError(t, err)
				require.Equal(t, uint(v), got)
			})
		})
		t.Run("uint32", func(t *testing.T) {
			decodeStr(t, s, func(t *testing.T, d *Decoder) {
				got, err := d.UInt32()
				require.NoError(t, err)
				require.Equal(t, uint32(v), got)
			})
		})
	}
}

func TestDecoderIntOverflow(t *testing.T) {
	t.Run("8", func(t *testing.T) {
		for _, s := range []string{
			"18446744073709551617",
			"18446744073709551616",
			"4294967296",
			"-9223372036854775809",
			"-18446744073709551617",
			"-32768",
			"65537",
			"-129",
			"256",
		} {
			t.Run(s, func(t *testing.T) {
				should := require.New(t)
				d := DecodeStr(s)
				v, err := d.Int8()
				should.Error(err, "%v", v)

				d = DecodeStr(s)
				_, err = d.int(8)
				should.Error(err)

				d = DecodeStr(s)
				vu, err := d.UInt8()
				should.Error(err, "%v", vu)

				d = DecodeStr(s)
				_, err = d.uint(8)
				should.Error(err)
			})
		}
		for _, s := range []string{
			"-9223372036854775809",
			"9223372036854775808",
			"2147483648",
			"-2147483649",
			"-4294967295",
			"32768",
			"-32769",
			"65535",
			"65536",
			"-129",
			"128",
			"255",
			"256",
		} {
			t.Run(s, func(t *testing.T) {
				should := require.New(t)
				d := DecodeStr(s)
				v, err := d.Int8()
				should.Error(err, "%v", v)
			})
		}
	})
	t.Run("16", func(t *testing.T) {
		for _, s := range []string{
			"18446744073709551617",
			"18446744073709551616",
			"4294967296",
			"-9223372036854775809",
			"-18446744073709551617",
			"-32769",
			"65537",
		} {
			t.Run(s, func(t *testing.T) {
				should := require.New(t)
				d := DecodeStr(s)
				v, err := d.Int16()
				should.Error(err, "%v", v)

				d = DecodeStr(s)
				_, err = d.int(16)
				should.Error(err)

				d = DecodeStr(s)
				vu, err := d.UInt16()
				should.Error(err, "%v", vu)

				d = DecodeStr(s)
				_, err = d.uint(16)
				should.Error(err)
			})
		}
		for _, s := range []string{
			"-9223372036854775809",
			"9223372036854775808",
			"2147483648",
			"-2147483649",
			"-4294967295",
			"32768",
			"-32769",
			"65535",
			"65536",
		} {
			t.Run(s, func(t *testing.T) {
				should := require.New(t)
				d := DecodeStr(s)
				v, err := d.Int16()
				should.Error(err, "%v", v)
			})
		}
	})
	t.Run("32", func(t *testing.T) {
		for _, s := range []string{
			"18446744073709551617",
			"18446744073709551616",
			"4294967296",
			"-9223372036854775809",
			"-18446744073709551617",
		} {
			t.Run(s, func(t *testing.T) {
				should := require.New(t)
				d := DecodeStr(s)
				v, err := d.Int32()
				should.Error(err, "%v", v)

				d = DecodeStr(s)
				_, err = d.int(32)
				should.Error(err)

				d = DecodeStr(s)
				vu, err := d.UInt32()
				should.Error(err, "%v", vu)

				d = DecodeStr(s)
				_, err = d.uint(32)
				should.Error(err)
			})
		}
		for _, s := range []string{
			"-9223372036854775809",
			"9223372036854775808",
			"2147483648",
			"-2147483649",
			"-4294967295",
		} {
			t.Run(s, func(t *testing.T) {
				should := require.New(t)
				d := DecodeStr(s)
				v, err := d.Int32()
				should.Error(err, "%v", v)
			})
		}
	})
	t.Run("64", func(t *testing.T) {
		for _, s := range []string{
			"18446744073709551617",
			"18446744073709551616",
			"-9223372036854775809",
			"-18446744073709551617",
		} {
			t.Run(s, func(t *testing.T) {
				should := require.New(t)
				d := DecodeStr(s)
				v, err := d.Int64()
				should.Error(err, "%v", v)

				d = DecodeStr(s)
				_, err = d.int(64)
				should.Error(err)

				d = DecodeStr(s)
				vu, err := d.UInt64()
				should.Error(err, "%v", vu)

				d = DecodeStr(s)
				_, err = d.uint(64)
				should.Error(err)
			})
		}
		for _, s := range []string{
			"-9223372036854775809",
			"9223372036854775808",
		} {
			t.Run(s, func(t *testing.T) {
				should := require.New(t)
				d := DecodeStr(s)
				v, err := d.Int64()
				should.Error(err, "%v", v)
			})
		}
	})
}

func TestReadInt64Overflow(t *testing.T) {
	s := `123456789232323232321545111111111111111111111111111111145454545445`
	iter := DecodeStr(s)
	_, err := iter.Int64()
	require.Error(t, err)
}

func TestReadInt8(t *testing.T) {
	inputs := []string{
		`-127`,
		`-12`,
		`-1`,
		`0`,
		`1`,
		`12`,
		`123`,
		`127`,
	}
	for i, input := range inputs {
		input := input
		t.Run(fmt.Sprintf("Test%d", i+1), createTestCase(input, func(t *testing.T, d *Decoder) error {
			should := require.New(t)
			expected, err := strconv.ParseInt(input, 10, 8)
			should.NoError(err)
			v, err := d.Int8()
			should.NoError(err)
			should.Equal(int8(expected), v)
			return nil
		}))
	}

	{
		input := "[" + strings.Join(inputs, ",") + "]"
		t.Run("Array", createTestCase(input, func(t *testing.T, d *Decoder) error {
			should := require.New(t)
			i := 0

			return d.Arr(func(d *Decoder) error {
				expected, err := strconv.ParseInt(inputs[i], 10, 8)
				should.NoError(err)

				v, err := d.Int8()
				if err != nil {
					return err
				}
				should.Equal(int8(expected), v)

				i++
				return nil
			})
		}))
	}
}

func TestReadInt16(t *testing.T) {
	inputs := []string{
		`-32767`,
		`-32766`,
		`-12`,
		`-1`,
		`0`,
		`1`,
		`12`,
		`123`,
		`1234`,
		`12345`,
		`32767`,
	}
	for i, input := range inputs {
		input := input
		t.Run(fmt.Sprintf("Test%d", i+1), createTestCase(input, func(t *testing.T, d *Decoder) error {
			should := require.New(t)
			expected, err := strconv.ParseInt(input, 10, 16)
			should.NoError(err)
			v, err := d.Int16()
			should.NoError(err)
			should.Equal(int16(expected), v)
			return nil
		}))
	}

	{
		input := "[" + strings.Join(inputs, ",") + "]"
		t.Run("Array", createTestCase(input, func(t *testing.T, d *Decoder) error {
			should := require.New(t)
			i := 0

			return d.Arr(func(d *Decoder) error {
				expected, err := strconv.ParseInt(inputs[i], 10, 16)
				should.NoError(err)

				v, err := d.Int16()
				if err != nil {
					return err
				}
				should.Equal(int16(expected), v)

				i++
				return nil
			})
		}))
	}
}

func TestReadInt32(t *testing.T) {
	inputs := []string{
		`-12`,
		`-1`,
		`0`,
		`1`,
		`12`,
		`123`,
		`1234`,
		`12345`,
		`123456`,
		`1234567`,
		`12345678`,
		`123456789`,
		`1234567890`,
		`2147483647`,
		`-2147483648`,
	}
	for i, input := range inputs {
		input := input
		t.Run(fmt.Sprintf("Test%d", i+1), createTestCase(input, func(t *testing.T, d *Decoder) error {
			should := require.New(t)
			expected, err := strconv.ParseInt(input, 10, 32)
			should.NoError(err)
			v, err := d.Int32()
			should.NoError(err)
			should.Equal(int32(expected), v)
			return nil
		}))
	}

	{
		input := "[" + strings.Join(inputs, ",") + "]"
		t.Run("Array", createTestCase(input, func(t *testing.T, d *Decoder) error {
			should := require.New(t)
			i := 0

			return d.Arr(func(d *Decoder) error {
				expected, err := strconv.ParseInt(inputs[i], 10, 32)
				should.NoError(err)

				v, err := d.Int32()
				if err != nil {
					return err
				}
				should.Equal(int32(expected), v)

				i++
				return nil
			})
		}))
	}
}

func TestReadInt64(t *testing.T) {
	inputs := []string{
		`-12`,
		`-1`,
		`0`,
		`1`,
		`12`,
		`123`,
		`1234`,
		`12345`,
		`123456`,
		`1234567`,
		`12345678`,
		`123456789`,
		`1234567890`,
		`12345678901`,
		`9223372036854775807`,
		`-9223372036854775808`,
		"-9223372036854775808\r",
	}
	for i, input := range inputs {
		input := input
		t.Run(fmt.Sprintf("Test%d", i+1), createTestCase(input, func(t *testing.T, d *Decoder) error {
			should := require.New(t)
			expected, err := strconv.ParseInt(strings.Trim(input, "\r"), 10, 64)
			should.NoError(err)
			v, err := d.Int64()
			should.NoError(err)
			should.Equal(expected, v)
			return nil
		}))
	}

	{
		input := "[" + strings.Join(inputs, ",") + "]"
		t.Run("Array", createTestCase(input, func(t *testing.T, d *Decoder) error {
			should := require.New(t)
			i := 0

			return d.Arr(func(d *Decoder) error {
				expected, err := strconv.ParseInt(strings.Trim(inputs[i], "\r"), 10, 64)
				should.NoError(err)

				v, err := d.Int64()
				if err != nil {
					return err
				}
				should.Equal(expected, v)

				i++
				return nil
			})
		}))
	}
}

func TestWriteUint32(t *testing.T) {
	vals := []uint32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			requireCompat(t, func(e *Encoder) {
				e.UInt32(val)
			}, val)
		})
	}
	should := require.New(t)
	e := GetEncoder()
	e.RawStr("a")
	e.UInt32(0xffffffff) // should clear buffer
	should.Equal("a4294967295", e.String())
}

func TestWriteInt32(t *testing.T) {
	vals := []int32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0x7fffffff, -0x80000000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			requireCompat(t, func(e *Encoder) {
				e.Int32(val)
			}, val)
		})
	}
	should := require.New(t)
	e := GetEncoder()
	e.RawStr("a")
	e.Int32(-0x7fffffff) // should clear buffer
	should.Equal("a-2147483647", e.String())
}

func TestWriteUint64(t *testing.T) {
	vals := []uint64{
		0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff, 0xffffffffff, 0xfffffffffff, 0xffffffffffff, 0xfffffffffffff, 0xffffffffffffff,
		0xfffffffffffffff, 0xffffffffffffffff,
	}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			requireCompat(t, func(e *Encoder) {
				e.UInt64(val)
			}, val)
		})
	}
	should := require.New(t)
	e := GetEncoder()
	e.RawStr("a")
	e.UInt64(0xffffffff) // should clear buffer
	should.Equal("a4294967295", e.String())
}

func TestWriteInt64(t *testing.T) {
	vals := []int64{
		0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff, 0xffffffffff, 0xfffffffffff, 0xffffffffffff, 0xfffffffffffff, 0xffffffffffffff,
		0xfffffffffffffff, 0x7fffffffffffffff, -0x8000000000000000,
	}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			requireCompat(t, func(e *Encoder) {
				e.Int64(val)
			}, val)
		})
	}
	should := require.New(t)
	e := GetEncoder()
	e.RawStr("a")
	e.Int64(0xffffffff) // should clear buffer
	should.Equal("a4294967295", e.String())
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
	var values []int64
	values = append(values, 0, math.MaxInt64, math.MinInt64)
	for i := int64(0); i < 28; i++ {
		v := int64(3)
		for k := int64(0); k < i; k++ {
			v += i + 1
			v += intPow(10, k) * (k%7 + 1)
			values = append(values, v)
		}
	}
	for _, v := range values {
		t.Run(fmt.Sprintf("%d", v), func(t *testing.T) {
			e := GetEncoder()
			e.ArrStart()
			e.Int64(v)
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

func TestDecoder_Int32(t *testing.T) {
	var values []int32
	values = append(values, 0, math.MaxInt32, math.MinInt32)
	for i := int32(0); i < 28; i++ {
		v := int32(3)
		for k := int32(0); k < i; k++ {
			v += i + 1
			v += intPow(10, k) * (k%7 + 1)
			values = append(values, v)
		}
	}
	for _, v := range values {
		t.Run(fmt.Sprintf("%d", v), func(t *testing.T) {
			e := GetEncoder()
			e.ArrStart()
			e.Int32(v)
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

func TestDecoder_UInt64(t *testing.T) {
	// Generate some diverse numbers.
	var values []uint64
	values = append(values, 0, math.MaxUint64)
	for i := uint64(0); i < 28; i++ {
		v := uint64(3)
		for k := uint64(0); k < i; k++ {
			v += i + 1
			v += intPow(10, k) * (k%7 + 1)
			values = append(values, v)
		}
	}
	for _, v := range values {
		t.Run(fmt.Sprintf("%d", v), func(t *testing.T) {
			e := GetEncoder()
			e.ArrStart()
			e.UInt64(v)
			e.ArrEnd()

			d := GetDecoder()
			d.ResetBytes(e.Bytes())
			requireElem(t, d)
			got, err := d.UInt64()
			require.NoError(t, err)
			require.Equal(t, v, got)
			requireArrEnd(t, d)
		})
	}
}

func TestDecoder_UInt32(t *testing.T) {
	var values []uint32
	values = append(values, 0, math.MaxUint32)
	for i := uint32(0); i < 28; i++ {
		v := uint32(3)
		for k := uint32(0); k < i; k++ {
			// No special meaning, just trying to make digits more diverse.
			v += i + 1
			v += intPow(10, k) * (k%7 + 1)
			values = append(values, v)
		}
	}
	for _, v := range values {
		t.Run(fmt.Sprintf("%d", v), func(t *testing.T) {
			e := GetEncoder()
			e.ArrStart()
			e.UInt32(v)
			e.UInt(uint(v))
			e.ArrEnd()

			d := GetDecoder()
			d.ResetBytes(e.Bytes())
			requireElem(t, d)
			got, err := d.UInt32()
			require.NoError(t, err)
			require.Equal(t, v, got)
			requireElem(t, d)
			gotUInt, err := d.UInt()
			require.NoError(t, err)
			require.Equal(t, uint(v), gotUInt)
			requireArrEnd(t, d)
		})
	}
}

func TestIntEOF(t *testing.T) {
	t.Run("Start", func(t *testing.T) {
		d := GetDecoder()
		var err error
		_, err = d.Int()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.Int64()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.Int32()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.Int16()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.Int8()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.UInt()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.UInt64()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.UInt32()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.UInt16()
		assert.Error(t, err, io.ErrUnexpectedEOF)
		_, err = d.UInt8()
		assert.Error(t, err, io.ErrUnexpectedEOF)
	})
	t.Run("Minus", func(t *testing.T) {
		d := DecodeStr(`-`)
		var err error
		_, err = d.Int()
		assert.Error(t, err, io.ErrUnexpectedEOF)

		d = DecodeStr(`-`)
		_, err = d.Int64()
		assert.Error(t, err, io.ErrUnexpectedEOF)

		d = DecodeStr(`-`)
		_, err = d.Int32()
		assert.Error(t, err, io.ErrUnexpectedEOF)

		d = DecodeStr(`-`)
		_, err = d.Int16()
		assert.Error(t, err, io.ErrUnexpectedEOF)

		d = DecodeStr(`-`)
		_, err = d.Int8()
		assert.Error(t, err, io.ErrUnexpectedEOF)
	})
}

func intPow[Int constraints.Integer](n, m Int) Int {
	if m == 0 {
		return 1
	}
	result := n
	for i := Int(2); i <= m; i++ {
		result *= n
	}
	return result
}
