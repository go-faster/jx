package jx

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_Line(t *testing.T) {
	const number = 4

	{
		// Generate object where value is line number.
		var w bytes.Buffer
		w.WriteString("{\n")
		for i := range [number]struct{}{} {
			// Line starts from 1 + heading '{' newline.
			line := i + 2
			w.WriteString(fmt.Sprintf("%q:%d", string(rune('a'+i)), line))
			if i != number-1 {
				w.WriteByte(',')
			}
			w.WriteByte('\n')
		}
		w.WriteByte('}')

		t.Run("Object", testBufferReader(w.String(), func(t *testing.T, d *Decoder) {
			a := require.New(t)
			a.NoError(d.ObjBytes(func(d *Decoder, key []byte) error {
				line := d.Line()
				v, err := d.Int()
				if err != nil {
					return err
				}
				a.Equalf(v, line, "%q:%d", key, v)
				return nil
			}))
		}))
	}

	{
		// Generate array where value is line number.
		var w bytes.Buffer
		w.WriteString("[\n")
		for i := range [number]struct{}{} {
			// Line starts from 1 + heading '[' newline.
			line := i + 2
			w.WriteString(strconv.Itoa(line))
			if i != number-1 {
				w.WriteByte(',')
			}
			w.WriteByte('\n')
		}
		w.WriteByte(']')

		t.Run("Array", testBufferReader(w.String(), func(t *testing.T, d *Decoder) {
			a := require.New(t)
			a.NoError(d.Arr(func(d *Decoder) error {
				line := d.Line()
				v, err := d.Int()
				if err != nil {
					return err
				}
				a.Equal(v, line)
				return nil
			}))
		}))
	}
}
