package stream

// Null writes null.
func (e *Encoder[W]) Null() bool {
	return e.comma() || e.w.writeString("null")
}
