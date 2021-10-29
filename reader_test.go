package jx

import (
	"testing"
)

func TestType_String(t *testing.T) {
	met := map[string]bool{}
	for i := Invalid; i <= Object+1; i++ {
		s := i.String()
		if s == "" {
			t.Error("blank")
		}
		if met[s] {
			t.Errorf("met %s", s)
		}
		met[s] = true
	}
	if len(met) != 8 {
		t.Error("unexpected met types")
	}
}
