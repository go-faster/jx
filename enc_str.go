package jx

// StrEscape encodes string with html special characters escaping.
func (e *Encoder) StrEscape(v string) {
	e.comma()
	e.w.StrEscape(v)
}

// Str encodes string without html escaping.
//
// Use StrEscape to escape html, this is default for encoding/json and
// should be used by default for untrusted strings.
func (e *Encoder) Str(v string) {
	e.comma()
	e.w.Str(v)
}
