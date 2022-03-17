//go:build !appengine && !purego

package jx

import "unsafe"

type sliceType struct {
	Ptr unsafe.Pointer
	Len uintptr
	Cap uintptr
}

type strType struct {
	Ptr unsafe.Pointer
	Len uintptr
}

//go:noescape
//go:linkname noescape runtime.noescape
func noescape(unsafe.Pointer) unsafe.Pointer
