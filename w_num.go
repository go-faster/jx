package jx

// Num encodes number.
func (w *Writer) Num(v Num) {
	if len(v) == 0 {
		w.Null()
		return
	}
	w.Raw(v)
}
