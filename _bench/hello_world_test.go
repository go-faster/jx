package bench

import (
	"bytes"
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-faster/jx"
)

// HelloWorld case.
//
// Example:
//	{"message": "Hello, world!"}
type HelloWorld struct {
	Message string `json:"message"`
}

const (
	helloWorldField   = "message"
	helloWorldMessage = "Hello, world!"
	helloWorld        = `{"message": "Hello, world!"}`
)

// setupHelloWorld should be called on each "HelloWorld" benchmark.
func setupHelloWorld(b *testing.B) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(helloWorld)))
}

// Common benchmark names.
const (
	// Encode is name for encoding benchmarks.
	Encode = "Encode"
	// Decode is name for decoding benchmarks.
	Decode = "Decode"
	// JX is name for benchmarks related to go-faster/jx package.
	JX = "jx"
	// Std is name for benchmarks related to encoding/json.
	Std = "std"
	// Sonic is name for benchmarks related to bytedance/sonic package.
	Sonic = "sonic"
	// JSONIter for json-iterator/go.
	JSONIter = "json-iterator"
)

func BenchmarkHelloWorld(b *testing.B) {
	b.Run(Encode, func(b *testing.B) {
		b.Run(JX, func(b *testing.B) {
			setupHelloWorld(b)
			var e jx.Encoder
			for i := 0; i < b.N; i++ {
				e.Reset()
				e.ObjStart()
				e.FieldStart(helloWorldField)
				e.Str(helloWorldMessage)
				e.ObjEnd()
			}
		})
		b.Run(Std, func(b *testing.B) {
			w := new(bytes.Buffer)
			e := json.NewEncoder(w)
			v := &HelloWorld{Message: helloWorldMessage}
			setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				w.Reset()
				if err := e.Encode(v); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(Sonic, func(b *testing.B) {
			sonicHelloWorld(b)
		})
		b.Run(JSONIter, func(b *testing.B) {
			e := jsoniter.NewStream(jsoniter.ConfigFastest, nil, 1024)
			setupHelloWorld(b)
			for i := 0; i < b.N; i++ {
				e.SetBuffer(e.Buffer()[:0]) // reset buffer
				e.WriteObjectStart()
				e.WriteObjectField(helloWorldField)
				e.WriteString(helloWorldMessage)
				e.WriteObjectEnd()
			}
		})
	})
}
