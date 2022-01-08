package bench

import (
	"bytes"
	"encoding/json"
	"github.com/romshark/jscan"
	"testing"

	"github.com/mailru/easyjson/jwriter"

	"github.com/go-faster/jx"
)

// setupSmall should be called on each "Small" benchmark.
func setupSmall(b *testing.B) []byte {
	b.Helper()
	b.ReportAllocs()
	data, err := json.Marshal(small)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(data)))
	return data
}

var (
	author = SmallAuthor{
		Name: "json",
		Age:  99,
		Male: true,
	}
	small = &Small{
		BookId:  12125925,
		BookIds: []int{-2147483648, 2147483647},
		Title:   "未来简史-从智人到智神",
		Titles:  []string{"hello", "world"},
		Price:   40.8,
		Prices:  []float64{-0.1, 0.1},
		Hot:     true,
		Hots:    []bool{true, true, true},
		Author:  author,
		Authors: []SmallAuthor{author, author, author},
		Weights: nil,
	}
)

func BenchmarkSmall(b *testing.B) {
	v := small
	b.Run(Encode, func(b *testing.B) {
		b.Run(JX, func(b *testing.B) {
			setupSmall(b)
			var e jx.Encoder
			for i := 0; i < b.N; i++ {
				e.Reset()
				v.Encode(&e)
			}
		})
		b.Run(Std, func(b *testing.B) {
			w := new(bytes.Buffer)
			e := json.NewEncoder(w)
			setupSmall(b)
			for i := 0; i < b.N; i++ {
				w.Reset()
				if err := e.Encode(v); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(Sonic, sonicSmall)
		b.Run(EasyJSON, func(b *testing.B) {
			jw := jwriter.Writer{}
			setupSmall(b)
			for i := 0; i < b.N; i++ {
				jw.Buffer.Buf = jw.Buffer.Buf[:0] // reset
				v.MarshalEasyJSON(&jw)
			}
		})
	})
	b.Run(Decode, func(b *testing.B) {
		b.Run(JX, func(b *testing.B) {
			data := setupSmall(b)
			var d jx.Decoder

			for i := 0; i < b.N; i++ {
				d.ResetBytes(data)
				if err := d.Skip(); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(JSCan, func(b *testing.B) {
			data := string(setupSmall(b))
			for i := 0; i < b.N; i++ {
				r := jscan.Scan(
					jscan.Options{},
					data,
					func(i *jscan.Iterator) bool { return false },
				)
				if r.IsErr() {
					b.Fatal("err")
				}
			}
		})
	})
}
