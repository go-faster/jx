package jx

// Num encodes number.
func (e *Encoder) Num(v Num) {
	if v.Format.Invalid() {
		e.Null()
		return
	}
	if v.Format.Str() {
		e.byte('"')
	}
	e.RawBytes(v.Value)
	if v.Format.Str() {
		e.byte('"')
	}
}
