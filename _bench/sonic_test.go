//go:build !go1.24

package bench

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/encoder"
)

func sonicHelloWorld(b *testing.B) {
	buf := make([]byte, 0, 1024)
	v := &HelloWorld{Message: helloWorldMessage}
	setupHelloWorld(b)
	for i := 0; i < b.N; i++ {
		buf = buf[:0] // reset buffer
		if err := encoder.EncodeInto(&buf, v, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func sonicSmall(b *testing.B) {
	buf := make([]byte, 0, 1024)
	setupSmall(b)
	for i := 0; i < b.N; i++ {
		buf = buf[:0] // reset buffer
		if err := encoder.EncodeInto(&buf, small, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func sonicDecodeSmall(b *testing.B) {
	data := string(setupSmall(b))
	var v Small
	for i := 0; i < b.N; i++ {
		if err := sonic.UnmarshalString(data, &v); err != nil {
			b.Fatal(err)
		}
	}
}
