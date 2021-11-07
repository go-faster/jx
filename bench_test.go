//go:build !gofuzz && go1.17
// +build !gofuzz,go1.17

package jx

import (
	_ "embed"
	"encoding/json"
	"testing"
)

//go:embed testdata/file.json
var benchData []byte

func Benchmark_large_file(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(benchData)))
	d := Decode(nil, 4096)

	for n := 0; n < b.N; n++ {
		d.ResetBytes(benchData)
		if err := d.Arr(nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValid(b *testing.B) {
	b.Run("JX", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchData)))
		var d Decoder
		for n := 0; n < b.N; n++ {
			d.ResetBytes(benchData)
			if err := d.Validate(); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchData)))

		for n := 0; n < b.N; n++ {
			if !json.Valid(benchData) {
				b.Fatal("invalid")
			}
		}
	})
}

func Benchmark_std_large_file(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var result []struct{}
		err := json.Unmarshal(benchData, &result)
		if err != nil {
			b.Error(err)
		}
	}
}
