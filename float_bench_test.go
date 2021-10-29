package jx

import (
	"encoding/json"
	"testing"
)

func Benchmark_json_float(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := float64(0)
		_ = json.Unmarshal([]byte(`1.1`), &result)
	}
}
