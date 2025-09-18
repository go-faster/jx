package bench

import (
	"bytes"
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson/jwriter"
	"github.com/minio/simdjson-go"
	fflib "github.com/pquerna/ffjson/fflib/v1"
	"github.com/romshark/jscan"
	"github.com/sugawarayuuta/sonnet"
	"github.com/valyala/fastjson"

	"github.com/go-faster/jx"
)

// setupHelloWorld should be called on each "HelloWorld" benchmark.
func setupHelloWorld(b *testing.B) []byte {
	b.Helper()
	b.ReportAllocs()
	data := []byte(helloWorld)
	b.SetBytes(int64(len(data)))
	return data
}

func BenchmarkHelloWorld(b *testing.B) {
	v := &HelloWorld{Message: helloWorldMessage}
	b.Run(Encode, func(b *testing.B) {
		b.Run(JX, func(b *testing.B) {
			b.Run("Encoder", func(b *testing.B) {
				setupHelloWorld(b)
				var e jx.Encoder
				for i := 0; i < b.N; i++ {
					e.Reset()
					v.Encode(&e)
				}
			})
			b.Run("Writer", func(b *testing.B) {
				setupHelloWorld(b)
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
			setupHelloWorld(b)
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
			setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				w.Reset()
				if err := e.Encode(v); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(Sonic, sonicHelloWorld)
		b.Run(JSONIter, func(b *testing.B) {
			s := jsoniter.NewStream(jsoniter.ConfigFastest, nil, 1024)
			setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				s.SetBuffer(s.Buffer()[:0]) // reset buffer
				v.EncodeIter(s)
			}
		})
		b.Run(EasyJSON, func(b *testing.B) {
			jw := jwriter.Writer{}
			setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				jw.Buffer.Buf = jw.Buffer.Buf[:0] // reset
				v.MarshalEasyJSON(&jw)
			}
		})
		b.Run(FFJSON, func(b *testing.B) {
			var buf fflib.EncodingBuffer = new(fflib.Buffer)
			v := &HelloWorldFFJSON{Message: helloWorldMessage}
			setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				buf.Reset()
				if err := v.MarshalJSONBuf(buf); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(Baseline, func(b *testing.B) {
			setupHelloWorld(b)
			buf := new(bytes.Buffer)
			for i := 0; i < b.N; i++ {
				buf.Reset()
				buf.WriteString(helloWorld)
			}
		})
	})
	b.Run(Scan, func(b *testing.B) {
		b.Run(JX, func(b *testing.B) {
			setupHelloWorld(b)
			var d jx.Decoder
			data := []byte(helloWorld)
			for i := 0; i < b.N; i++ {
				d.ResetBytes(data)
				if err := d.Skip(); err != nil {
					b.Fatal()
				}
			}
		})
		b.Run(JScan, func(b *testing.B) {
			setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				r := jscan.Scan(
					jscan.Options{},
					helloWorld,
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
			setupHelloWorld(b)
			pj := new(simdjson.ParsedJson)
			data := setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				var err error
				if pj, err = simdjson.Parse(data, pj, simdjson.WithCopyStrings(false)); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(FastJSON, func(b *testing.B) {
			p := new(fastjson.Parser)
			data := setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				if _, err := p.ParseBytes(data); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run(Decode, func(b *testing.B) {
		b.Run(JX, func(b *testing.B) {
			d := new(jx.Decoder)
			data := setupHelloWorld(b)
			var v HelloWorld
			for i := 0; i < b.N; i++ {
				d.ResetBytes(data)
				if err := v.Decode(d); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(SIMD, func(b *testing.B) {
			if !simdjson.SupportedCPU() {
				b.SkipNow()
			}
			pj := new(simdjson.ParsedJson)
			data := setupHelloWorld(b)
			var v HelloWorld
			for i := 0; i < b.N; i++ {
				var err error
				if pj, err = v.DecodeSIMD(data, pj); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(Std, func(b *testing.B) {
			data := setupHelloWorld(b)
			var v HelloWorld
			for i := 0; i < b.N; i++ {
				if err := json.Unmarshal(data, &v); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(FastJSON, func(b *testing.B) {
			p := new(fastjson.Parser)
			data := setupHelloWorld(b)
			var v HelloWorld
			for i := 0; i < b.N; i++ {
				if err := v.DecodeFastJSON(p, data); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}
