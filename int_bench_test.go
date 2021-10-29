package jx

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"testing"
)

func Benchmark_encode_int(b *testing.B) {
	stream := NewStream(Default, ioutil.Discard, 64)
	for n := 0; n < b.N; n++ {
		stream.Reset(nil)
		stream.WriteUint64(0xffffffff)
	}
}

func Benchmark_itoa(b *testing.B) {
	for n := 0; n < b.N; n++ {
		strconv.FormatInt(0xffffffff, 10)
	}
}

func Benchmark_int(b *testing.B) {
	iter := NewIterator(Default)
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
