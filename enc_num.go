package jx

// Num encodes number.
func (e *Encoder) Num(v Num) {
	e.comma()
	e.w.Num(v)
}
