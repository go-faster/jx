package stream

// True writes true.
func (e *Encoder[W]) True() bool {
	return e.comma() || e.w.writeString("true")
}

// False writes false.
func (e *Encoder[W]) False() bool {
	return e.comma() || e.w.writeString("false")
}

// Bool encodes boolean.
func (e *Encoder[W]) Bool(v bool) bool {
	if v {
		return e.True()
	}
	return e.False()
}
