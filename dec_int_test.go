package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkDecoder_Int(b *testing.B) {
	data := []byte(`69315063`)
	d := GetDecoder()
	for i := 0; i < b.N; i++ {
		d.ResetBytes(data)
		if _, err := d.Int(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecoder_Uint(b *testing.B) {
	data := []byte(`69315063`)
	d := GetDecoder()
	for i := 0; i < b.N; i++ {
		d.ResetBytes(data)
		if _, err := d.Uint(); err != nil {
			b.Fatal(err)
		}
	}
}

func TestDecoder_int_sizes(t *testing.T) {
	data := []byte(`69315063`)
	d := GetDecoder()
	for _, size := range []int{32, 64} {
		d.ResetBytes(data)
		v, err := d.int(size)
		require.NoError(t, err)
		require.Equal(t, 69315063, v)
	}
}

func TestDecoder_uint_sizes(t *testing.T) {
	data := []byte(`69315063`)
	d := GetDecoder()
	for _, size := range []int{32, 64} {
		d.ResetBytes(data)
		v, err := d.uint(size)
		require.NoError(t, err)
		require.Equal(t, uint(69315063), v)
	}
}
