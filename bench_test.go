//go:build !gofuzz && go1.17
// +build !gofuzz,go1.17

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

func BenchmarkValid(b *testing.B) {
	b.Run("JX", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(data)))
		var d Decoder
		for n := 0; n < b.N; n++ {
			d.ResetBytes(data)
			if err := d.Validate(); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(data)))

		for n := 0; n < b.N; n++ {
			if !json.Valid(data) {
				b.Fatal("invalid")
			}
		}
	})
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
