package stream

// ObjStart writes object start, performing indentation if needed.
//
// Use Obj as convenience helper for writing objects.
func (e *Encoder[W]) ObjStart() bool {
	if e.comma() || e.w.writeByte('{') {
		return true
	}
	e.begin()
	return e.writeIndent()
}

// FieldStart encodes field name and writes colon.
//
// For non-zero indentation also writes single space after colon.
//
// Use Field as convenience helper for encoding fields.
func (e *Encoder[W]) FieldStart(field string) bool {
	if e.Str(field) || e.w.writeByte(':') {
		return true
	}
	if e.indent > 0 {
		if e.w.writeByte(' ') {
			return true
		}
	}
	if len(e.first) > 0 {
		e.first[e.current()] = true
	}
	return false
}

// Field encodes field start and then invokes callback.
//
// Has ~5ns overhead over FieldStart.
func (e *Encoder[W]) Field(name string, f func(e *Encoder[W]) bool) bool {
	return e.FieldStart(name) ||
		f(e)
}

// ObjEnd writes end of object token, performing indentation if needed.
//
// Use Obj as convenience helper for writing objects.
func (e *Encoder[W]) ObjEnd() bool {
	e.end()
	return e.writeIndent() ||
		e.w.writeByte('}')
}

// ObjEmpty writes empty object.
func (e *Encoder[W]) ObjEmpty() bool {
	return e.comma() ||
		e.ObjStart() ||
		e.ObjEnd()
}

// Obj writes start of object, invokes callback and writes end of object.
//
// If callback is nil, writes empty object.
func (e *Encoder[W]) Obj(f func(e *Encoder[W]) bool) bool {
	if f == nil {
		return e.ObjEmpty()
	}
	return e.ObjStart() ||
		f(e) ||
		e.ObjEnd()
}
