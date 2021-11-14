//go:build !go1.17
// +build !go1.17

package bench

import (
	"testing"

	sonic "github.com/bytedance/sonic/encoder"
)

func sonicHelloWorld(b *testing.B) {
	buf := make([]byte, 0, 1024)
	v := &HelloWorld{Message: helloWorldMessage}
	setupHelloWorld(b)
	for i := 0; i < b.N; i++ {
		buf = buf[:0] // reset buffer
		if err := sonic.EncodeInto(&buf, v, 0); err != nil {
			b.Fatal(err)
		}
	}
}
