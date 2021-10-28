package jir

// Capture call f and then rolls back buffer to state before call.
// Does not work with reader.
func (it *Iterator) Capture(f func(i *Iterator)) {
	if it.reader != nil {
		panic("capture is not supported")
	}
	head, tail, depth := it.head, it.tail, it.depth
	f(it)
	it.head, it.tail, it.depth = head, tail, depth
}
