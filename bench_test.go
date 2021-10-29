package jir

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
	iter := Parse(Default, nil, 4096)

	for n := 0; n < b.N; n++ {
		iter.ResetBytes(data)
		iter.Array(func(iter *Iterator) bool {
			iter.Skip()
			return true
		})
		if iter.Error != nil {
			b.Error(iter.Error)
		}
	}
}

func Benchmark_std_large_file(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		result := []struct{}{}
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Error(err)
		}
	}
}
