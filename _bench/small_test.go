package bench

import (
	"bytes"
	"encoding/json"
	stdv2 "encoding/json/v2"
	"testing"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"github.com/minio/simdjson-go"
	"github.com/romshark/jscan"
	"github.com/stretchr/testify/require"
	"github.com/sugawarayuuta/sonnet"
	"github.com/valyala/fastjson"

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
		Title:   "Foo",
		Titles:  []string{"hello", "world"},
		Price:   40.8,
		Prices:  []float64{-0.1, 0.1},
		Hot:     true,
		Hots:    []bool{true, true, true},
		Author:  author,
		Authors: []SmallAuthor{author, author, author},
		Weights: []int{1, 2, 3},
	}
)

func TestSmall(t *testing.T) {
	data, err := json.Marshal(small)
	if err != nil {
		t.Fatal(err)
	}
	t.Run(JX, func(t *testing.T) {
		var s Small
		d := jx.DecodeBytes(data)
		d.ResetBytes(data)
		if err := s.Decode(d); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, small, &s)
	})
	t.Run(FastJSON, func(t *testing.T) {
		var s Small
		var p fastjson.Parser
		if err := s.DecodeFastJSON(&p, data); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, small, &s)
	})
}

func BenchmarkSmall(b *testing.B) {
	v := small
	b.Run(Encode, func(b *testing.B) {
		b.Run(JX, func(b *testing.B) {
			b.Run("Encoder", func(b *testing.B) {
				setupSmall(b)
				var e jx.Encoder
				for i := 0; i < b.N; i++ {
					e.Reset()
					v.Encode(&e)
				}
			})
			b.Run("Writer", func(b *testing.B) {
				setupSmall(b)
				var w jx.Writer
				for i := 0; i < b.N; i++ {
					w.Reset()
					v.Write(&w)
				}
			})
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
		b.Run(Sonnet, func(b *testing.B) {
			w := new(bytes.Buffer)
			e := sonnet.NewEncoder(w)
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
		b.Run(EasyJSON, func(b *testing.B) {
			data := setupSmall(b)
			var d Small
			for i := 0; i < b.N; i++ {
				d.Reset()
				l := jlexer.Lexer{Data: data}
				d.UnmarshalEasyJSON(&l)
				if err := l.Error(); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(Sonic, sonicDecodeSmall)
		b.Run(Std, func(b *testing.B) {
			data := setupSmall(b)
			var d Small
			for i := 0; i < b.N; i++ {
				d.Reset()
				if err := json.Unmarshal(data, &d); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(Sonnet, func(b *testing.B) {
			data := setupSmall(b)
			var d Small
			for i := 0; i < b.N; i++ {
				d.Reset()
				if err := sonnet.Unmarshal(data, &d); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(JX, func(b *testing.B) {
			data := setupSmall(b)
			var s Small
			d := jx.DecodeBytes(data)
			for i := 0; i < b.N; i++ {
				s.Reset()
				d.ResetBytes(data)
				if err := s.Decode(d); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(FastJSON, func(b *testing.B) {
			data := setupSmall(b)
			var p fastjson.Parser
			var s Small
			for i := 0; i < b.N; i++ {
				s.Reset()
				if err := s.DecodeFastJSON(&p, data); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(StdV2, func(b *testing.B) {
			data := setupSmall(b)
			var s Small
			for i := 0; i < b.N; i++ {
				s.Reset()
				if err := stdv2.Unmarshal(data, &s); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run(Scan, func(b *testing.B) {
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
		b.Run(JScan, func(b *testing.B) {
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
		b.Run(SIMD, func(b *testing.B) {
			if !simdjson.SupportedCPU() {
				b.SkipNow()
			}
			pj := new(simdjson.ParsedJson)
			data := setupSmall(b)
			for i := 0; i < b.N; i++ {
				var err error
				if pj, err = simdjson.Parse(data, pj, simdjson.WithCopyStrings(false)); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(FastJSON, func(b *testing.B) {
			p := new(fastjson.Parser)
			data := setupSmall(b)
			for i := 0; i < b.N; i++ {
				if _, err := p.ParseBytes(data); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}
