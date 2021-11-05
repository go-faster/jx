package jx

// Num encodes number.
func (e *Encoder) Num(v Num) {
	if len(v) == 0 {
		e.Null()
		return
	}
	e.RawBytes(v)
}
