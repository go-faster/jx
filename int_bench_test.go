package jx

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"
)

func BenchmarkEncoder_Int(b *testing.B) {
	const v = 0xffffff
	b.Run("Strconv", func(b *testing.B) {
		b.ReportAllocs()
		var buf []byte
		for i := 0; i < b.N; i++ {
			buf = buf[:0]
			buf = strconv.AppendInt(buf, v, 10)
		}
	})
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		buf := new(bytes.Buffer)
		e := json.NewEncoder(buf)
		for i := 0; i < b.N; i++ {
			buf.Reset()
			if err := e.Encode(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("JX", func(b *testing.B) {
		b.ReportAllocs()
		e := GetEncoder()
		for i := 0; i < b.N; i++ {
			e.Reset()
			e.UInt64(v)
		}
	})
}

func BenchmarkDecoder_Int64(b *testing.B) {
	input := []byte(`100`)
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := int64(0)
			if err := json.Unmarshal(input, &result); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("JX", func(b *testing.B) {
		b.ReportAllocs()
		d := GetDecoder()
		for i := 0; i < b.N; i++ {
			d.ResetBytes(input)
			if _, err := d.Int64(); err != nil {
				b.Fatal(err)
			}
		}
	})
}
