package stream

// ArrStart writes start of array, performing indentation if needed.
//
// Use Arr as convenience helper for writing arrays.
func (e *Encoder[W]) ArrStart() bool {
	if e.comma() || e.w.writeByte('[') {
		return true
	}
	e.begin()
	return e.writeIndent()
}

// ArrEnd writes end of array, performing indentation if needed.
//
// Use Arr as convenience helper for writing arrays.
func (e *Encoder[W]) ArrEnd() bool {
	e.end()
	return e.writeIndent() ||
		e.w.writeByte(']')
}

// ArrEmpty writes empty array.
func (e *Encoder[W]) ArrEmpty() bool {
	return e.comma() ||
		e.ArrStart() ||
		e.ArrEnd()
}

// Arr writes start of array, invokes callback and writes end of array.
//
// If callback is nil, writes empty array.
func (e *Encoder[W]) Arr(f func(e *Encoder[W]) bool) bool {
	if f == nil {
		return e.ArrEmpty()
	}
	return e.ArrStart() ||
		f(e) ||
		e.ArrEnd()
}
