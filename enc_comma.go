package jx

// begin should be called before new Array or Object.
func (e *Encoder) begin() {
	e.first = append(e.first, true)
}

// end should be called after Array or Object.
func (e *Encoder) end() {
	if len(e.first) == 0 {
		return
	}
	e.first = e.first[:e.last()]
}

func (e *Encoder) last() int {
	return len(e.first) - 1
}

func (e *Encoder) resetComma() {
	if len(e.first) == 0 {
		return
	}
	e.first[e.last()] = true
}

// comma should be called before any new value.
func (e *Encoder) comma() {
	// Writing commas.
	// 1. Before every field expect first.
	// 2. Before every array element except first.
	if len(e.first) == 0 {
		return
	}
	if e.first[e.last()] {
		e.first[e.last()] = false
		return
	}
	e.byte(',')
	e.writeIndent(0)
}
