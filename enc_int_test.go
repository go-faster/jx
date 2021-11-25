package jx

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_Uint64(t *testing.T) {
	test := func(i uint64) (string, func(t *testing.T)) {
		return fmt.Sprintf("Test%d", i), func(t *testing.T) {
			enc := GetEncoder()
			enc.Reset()
			enc.Uint64(i)
			require.Equal(t, enc.String(), strconv.FormatUint(i, 10))
		}
	}
	overflows := func(i uint64) bool {
		result := i * 10
		return result/10 != i
	}

	t.Run(test(0))
	t.Run(test(1337))
	t.Run(test(math.MaxUint64))
	for i := uint64(1); !overflows(i); i *= 10 {
		t.Run(test(i))
		t.Run(test(i + 1))
	}
}

func TestEncoder_Uint32(t *testing.T) {
	test := func(i uint32) (string, func(t *testing.T)) {
		return fmt.Sprintf("Test%d", i), func(t *testing.T) {
			enc := GetEncoder()
			enc.Reset()
			enc.Uint32(i)
			require.Equal(t, enc.String(), strconv.FormatUint(uint64(i), 10))
		}
	}
	overflows := func(i uint32) bool {
		result := i * 10
		return result/10 != i
	}

	t.Run(test(0))
	t.Run(test(1337))
	t.Run(test(math.MaxUint32))
	for i := uint32(1); !overflows(i); i *= 10 {
		t.Run(test(i))
		t.Run(test(i + 1))
	}
}

func TestEncoder_Uint16(t *testing.T) {
	test := func(i uint16) (string, func(t *testing.T)) {
		return fmt.Sprintf("Test%d", i), func(t *testing.T) {
			enc := GetEncoder()
			enc.Reset()
			enc.Uint16(i)
			require.Equal(t, enc.String(), strconv.FormatUint(uint64(i), 10))
		}
	}
	overflows := func(i uint16) bool {
		result := i * 10
		return result/10 != i
	}

	t.Run(test(0))
	t.Run(test(1337))
	t.Run(test(math.MaxUint16))
	for i := uint16(1); !overflows(i); i *= 10 {
		t.Run(test(i))
		t.Run(test(i + 1))
	}
}

func TestEncoder_Uint8(t *testing.T) {
	test := func(i uint8) (string, func(t *testing.T)) {
		return fmt.Sprintf("Test%d", i), func(t *testing.T) {
			enc := GetEncoder()
			enc.Reset()
			enc.Uint8(i)
			require.Equal(t, enc.String(), strconv.FormatUint(uint64(i), 10))
		}
	}
	overflows := func(i uint8) bool {
		result := i * 10
		return result/10 != i
	}

	t.Run(test(0))
	t.Run(test(237))
	t.Run(test(math.MaxUint8))
	for i := uint8(1); !overflows(i); i *= 10 {
		t.Run(test(i))
		t.Run(test(i + 1))
	}
}

func TestEncoder_Int16(t *testing.T) {
	test := func(i int16) (string, func(t *testing.T)) {
		return fmt.Sprintf("Test%d", i), func(t *testing.T) {
			enc := GetEncoder()
			enc.Reset()
			enc.Int16(i)
			require.Equal(t, enc.String(), strconv.FormatInt(int64(i), 10))
		}
	}
	overflows := func(i int16) bool {
		result := i * 10
		return result/10 != i
	}

	t.Run(test(0))
	t.Run(test(-13))
	t.Run(test(math.MaxInt16))
	for i := int16(1); !overflows(i); i *= 10 {
		t.Run(test(i))
		t.Run(test(i + 1))
	}
}

func TestEncoder_Int8(t *testing.T) {
	test := func(i int8) (string, func(t *testing.T)) {
		return fmt.Sprintf("Test%d", i), func(t *testing.T) {
			enc := GetEncoder()
			enc.Reset()
			enc.Int8(i)
			require.Equal(t, enc.String(), strconv.FormatInt(int64(i), 10))
		}
	}
	overflows := func(i int8) bool {
		result := i * 10
		return result/10 != i
	}

	t.Run(test(0))
	t.Run(test(-13))
	t.Run(test(math.MaxInt8))
	for i := int8(1); !overflows(i); i *= 10 {
		t.Run(test(i))
		t.Run(test(i + 1))
	}
}
