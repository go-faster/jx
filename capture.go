package jx

// Capture call f and then rolls back buffer to state before call.
// Does not work with reader.
func (it *Iter) Capture(f func(i *Iter) error) error {
	if it.reader != nil {
		panic("capture is not supported")
	}
	head, tail, depth := it.head, it.tail, it.depth
	err := f(it)
	it.head, it.tail, it.depth = head, tail, depth
	return err
}
