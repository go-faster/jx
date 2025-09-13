//go:build go1.24

package bench

import "testing"

func sonicSkip(b *testing.B) {
	// Sonic tests are skipped because sonic does not support go1.18.
	// Ref:
	//  https://github.com/bytedance/sonic/pull/116
	//  https://github.com/bytedance/sonic/issues/75
	b.Helper()
	b.Skip("not supported on current go version")
}

func sonicHelloWorld(b *testing.B) {
	b.Helper()
	sonicSkip(b)
}

func sonicSmall(b *testing.B) {
	b.Helper()
	sonicSkip(b)
}

func sonicDecodeSmall(b *testing.B) {
	b.Helper()
	sonicSkip(b)
}
