package jx

import (
	_ "embed"
	"encoding/json"
	"testing"
)

//go:embed testdata/file.json
var data []byte

/*
200000	      8886 ns/op	    4336 B/op	       6 allocs/op
50000	     34244 ns/op	    6744 B/op	      14 allocs/op
*/
func Benchmark_large_file(b *testing.B) {
	b.ReportAllocs()
	iter := Decode(nil, 4096)

	for n := 0; n < b.N; n++ {
		iter.ResetBytes(data)
		if err := iter.Array(func(iter *Decoder) error {
			return iter.Skip()
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_std_large_file(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var result []struct{}
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Error(err)
		}
	}
}
