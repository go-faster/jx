package jx

import (
	_ "embed"
	"encoding/json"
	"testing"
)

//go:embed testdata/file.json
var data []byte

func Benchmark_large_file(b *testing.B) {
	b.ReportAllocs()
	d := Decode(nil, 4096)

	for n := 0; n < b.N; n++ {
		d.ResetBytes(data)
		if err := d.Arr(nil); err != nil {
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
