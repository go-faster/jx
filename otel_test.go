package jx

import (
	_ "embed"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-faster/errors"
)

type Pos struct {
	Start int
	End   int
}

type Bytes struct {
	Buf []byte
	Pos []Pos
}

func (b Bytes) Elem(i int) []byte {
	p := b.Pos[i]
	return b.Buf[p.Start:p.End]
}

func (b Bytes) ForEachBytes(f func(i int, b []byte) error) error {
	for i, p := range b.Pos {
		if err := f(i, b.Buf[p.Start:p.End]); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bytes) Append(v []byte) {
	start := len(b.Buf)
	b.Buf = append(b.Buf, v...)
	end := len(b.Buf)
	b.Pos = append(b.Pos, Pos{Start: start, End: end})
}

func (b *Bytes) Reset() {
	b.Buf = b.Buf[:0]
	b.Pos = b.Pos[:0]
}

type Map struct {
	Keys   Bytes
	Values Bytes
}

func (m *Map) Append(k, v []byte) {
	m.Keys.Append(k)
	m.Values.Append(v)
}

func (m Map) Write(w *Writer) {
	w.ObjStart()
	for i, p := range m.Keys.Pos {
		if i != 0 {
			w.Comma()
		}
		w.FieldStart(string(m.Keys.Buf[p.Start:p.End]))
		w.Raw(m.Values.Elem(i))
	}
	w.ObjEnd()
}

func (m Map) Encode(e *Encoder) {
	e.ObjStart()
	defer e.ObjEnd()
	_ = m.Keys.ForEachBytes(func(i int, b []byte) error {
		e.FieldStart(string(b))
		e.Raw(m.Values.Elem(i))
		return nil
	})
}

func (m *Map) Decode(d *Decoder) error {
	return d.ObjBytes(func(d *Decoder, k []byte) error {
		v, err := d.Raw()
		if err != nil {
			return errors.Wrap(err, "value")
		}
		m.Append(k, v)
		return nil
	})
}

func (m *Map) Reset() {
	m.Keys.Reset()
	m.Values.Reset()
}

// OTEL log model.
//
// See https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/data-model.md#log-data-model
type OTEL struct {
	Timestamp  Num
	Attributes Map
	Resource   Map
	TraceID    [16]byte
	SpanID     [8]byte
	Severity   byte
	Body       Raw
}

func (o *OTEL) Write(w *Writer) {
	w.ObjStart()

	w.RawStr(`"Timestamp":`)
	w.Num(o.Timestamp)

	w.RawStr(`,"Attributes":`)
	o.Attributes.Write(w)

	w.RawStr(`,"Resource":`)
	o.Resource.Write(w)

	{
		// Hex encoding.
		buf := make([]byte, 32) // 32 = 16 * 2
		var n int

		n = hex.Encode(buf, o.TraceID[:])
		w.RawStr(`,"TraceId":`)
		w.Str(string(buf[:n]))

		n = hex.Encode(buf, o.SpanID[:])
		w.RawStr(`,"SpanId":`)
		w.Str(string(buf[:n]))
	}

	if o.Severity > 0 && o.Severity <= 24 {
		w.RawStr(`,"SeverityText":`)
		switch {
		case o.Severity >= 1 && o.Severity <= 4:
			w.RawStr(`"TRACE"`)
		case o.Severity >= 5 && o.Severity <= 8:
			w.RawStr(`"DEBUG"`)
		case o.Severity >= 9 && o.Severity <= 12:
			w.RawStr(`"INFO"`)
		case o.Severity >= 13 && o.Severity <= 16:
			w.RawStr(`"WARN"`)
		case o.Severity >= 17 && o.Severity <= 20:
			w.RawStr(`"ERROR"`)
		case o.Severity >= 21 && o.Severity <= 24:
			w.RawStr(`"FATAL"`)
		}
		w.RawStr(`,"SeverityNumber":`)
		w.UInt8(o.Severity)
	}

	w.RawStr(`,"Body":`)
	w.Raw(o.Body)
	w.ObjEnd()
}

func (o *OTEL) Encode(e *Encoder) {
	e.ObjStart()
	defer e.ObjEnd()

	e.FieldStart("Timestamp")
	e.Num(o.Timestamp)

	e.FieldStart("Attributes")
	o.Attributes.Encode(e)

	e.FieldStart("Resource")
	o.Resource.Encode(e)

	{
		// Hex encoding.
		buf := make([]byte, 32) // 32 = 16 * 2
		var n int

		n = hex.Encode(buf, o.TraceID[:])
		e.FieldStart("TraceId")
		e.Str(string(buf[:n]))

		n = hex.Encode(buf, o.SpanID[:])
		e.FieldStart("SpanId")
		e.Str(string(buf[:n]))
	}

	if o.Severity > 0 && o.Severity <= 24 {
		e.FieldStart("SeverityText")
		switch {
		case o.Severity >= 1 && o.Severity <= 4:
			e.Str("TRACE")
		case o.Severity >= 5 && o.Severity <= 8:
			e.Str("DEBUG")
		case o.Severity >= 9 && o.Severity <= 12:
			e.Str("INFO")
		case o.Severity >= 13 && o.Severity <= 16:
			e.Str("WARN")
		case o.Severity >= 17 && o.Severity <= 20:
			e.Str("ERROR")
		case o.Severity >= 21 && o.Severity <= 24:
			e.Str("FATAL")
		}
		e.FieldStart("SeverityNumber")
		e.UInt8(o.Severity)
	}

	e.FieldStart("Body")
	e.Raw(o.Body)
}

func (o *OTEL) Decode(d *Decoder) error {
	return d.ObjBytes(func(d *Decoder, key []byte) error {
		switch string(key) {
		case "Body":
			v, err := d.RawAppend(o.Body[:0])
			if err != nil {
				return errors.Wrap(err, "body")
			}
			o.Body = v
			return nil
		case "SeverityNumber":
			v, err := d.UInt8()
			if err != nil {
				return errors.Wrap(err, "severity number")
			}
			o.Severity = v
			return nil
		case "SeverityText":
			return d.Skip()
		case "Timestamp":
			v, err := d.NumAppend(o.Timestamp[:0])
			if err != nil {
				return errors.Wrap(err, "timestamp")
			}
			o.Timestamp = v
			return nil
		case "TraceId":
			v, err := d.StrBytes()
			if err != nil {
				return errors.Wrap(err, "trace id")
			}
			if _, err := hex.Decode(o.TraceID[:], v); err != nil {
				return errors.Wrap(err, "trace id decode")
			}
			return nil
		case "SpanId":
			v, err := d.StrBytes()
			if err != nil {
				return errors.Wrap(err, "span id")
			}
			if _, err := hex.Decode(o.SpanID[:], v); err != nil {
				return errors.Wrap(err, "span id decode")
			}
			return nil
		case "Attributes":
			if err := o.Attributes.Decode(d); err != nil {
				return errors.Wrap(err, "attributes")
			}
			return nil
		case "Resource":
			if err := o.Resource.Decode(d); err != nil {
				return errors.Wrap(err, "resource")
			}
			return nil
		default:
			return errors.Errorf("unknown key %q", key)
		}
	})
}

func (o *OTEL) Reset() {
	o.Body = o.Body[:0]
	o.Severity = 0
	o.TraceID = [16]byte{}
	o.SpanID = [8]byte{}
	o.Timestamp = o.Timestamp[:0]
	o.Attributes.Reset()
	o.Resource.Reset()
}

//go:embed testdata/otel_ex_1.json
var otelEx1 []byte

func TestOTELDecode(t *testing.T) {
	d := DecodeBytes(otelEx1)
	var v OTEL
	require.NoError(t, v.Decode(d))

	t.Run("Write", func(t *testing.T) {
		var w Writer
		v.Write(&w)

		require.JSONEq(t, string(otelEx1), w.String())
	})
}

func BenchmarkOTEL(b *testing.B) {
	var v OTEL
	dec := DecodeBytes(otelEx1)
	require.NoError(b, v.Decode(dec))

	b.Run("Decode", func(b *testing.B) {
		d := GetDecoder()
		b.ReportAllocs()
		b.SetBytes(int64(len(otelEx1)))

		var vDec OTEL
		for i := 0; i < b.N; i++ {
			vDec.Reset()
			d.ResetBytes(otelEx1)
			if err := vDec.Decode(d); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Validate", func(b *testing.B) {
		d := GetDecoder()
		b.ReportAllocs()
		b.SetBytes(int64(len(otelEx1)))

		for i := 0; i < b.N; i++ {
			d.ResetBytes(otelEx1)
			if d.Validate() != nil {
				b.Fatal("invalid")
			}
		}
	})
	b.Run("Write", func(b *testing.B) {
		w := GetWriter()
		defer PutWriter(w)
		v.Write(w)

		b.ReportAllocs()
		b.SetBytes(int64(len(w.Buf)))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			w.Reset()
			v.Write(w)
		}
	})
	b.Run("Encode", func(b *testing.B) {
		b.ReportAllocs()
		e := GetEncoder()
		v.Encode(e)
		b.SetBytes(int64(len(e.Bytes())))

		for i := 0; i < b.N; i++ {
			e.Reset()
			v.Encode(e)
		}
	})
}
