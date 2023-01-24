package stream

import "github.com/go-faster/jx"

// Num encodes number.
func (e *Encoder[W]) Num(v jx.Num) {
	if len(v) == 0 {
		e.Null()
		return
	}
	e.Raw(v)
}
