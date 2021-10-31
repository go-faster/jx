package jx

import (
	"encoding/json"
	"testing"
)

func BenchmarkFloat(b *testing.B) {
	data := []byte(`1.1`)
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var result float64
			if err := json.Unmarshal([]byte(`1.1`), &result); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("JX", func(b *testing.B) {
		b.ReportAllocs()
		d := GetDecoder()
		for n := 0; n < b.N; n++ {
			d.ResetBytes(data)
			if _, err := d.Float64(); err != nil {
				b.Fatal(err)
			}
		}
	})
}
