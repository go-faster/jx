package bench

import (
	"bytes"
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson/jwriter"
	fflib "github.com/pquerna/ffjson/fflib/v1"
	"github.com/romshark/jscan"
	"github.com/sugawarayuuta/sonnet"

	"github.com/go-faster/jx"
)

// setupHelloWorld should be called on each "HelloWorld" benchmark.
func setupHelloWorld(b *testing.B) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(helloWorld)))
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
	})
}
