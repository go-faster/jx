package jx

import (
	"encoding/json"
	"strconv"
	"testing"
)

func Benchmark_encode_int(b *testing.B) {
	e := NewEncoder()
	for n := 0; n < b.N; n++ {
		e.Reset()
		e.Uint64(0xffffffff)
	}
}

func Benchmark_itoa(b *testing.B) {
	for n := 0; n < b.N; n++ {
		strconv.FormatInt(0xffffffff, 10)
	}
}

func Benchmark_int(b *testing.B) {
	iter := NewDecoder()
	input := []byte(`100`)
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(input)
		iter.Int64()
	}
}

func Benchmark_std_int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := int64(0)
		_ = json.Unmarshal([]byte(`-100`), &result)
	}
}
