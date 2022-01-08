package jx

import (
	_ "embed"
	"encoding/hex"
	"github.com/go-faster/errors"
	"github.com/stretchr/testify/require"
	"testing"
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
		e.Uint8(o.Severity)
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
			v, err := d.Int32()
			if err != nil {
				return errors.Wrap(err, "severity number")
			}
			o.Severity = byte(v)
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

	t.Run("Encode", func(t *testing.T) {
		var e Encoder
		e.SetIdent(2)
		v.Encode(&e)

		require.JSONEq(t, string(otelEx1), e.String())
	})
}

func BenchmarkOTEL_Decode(b *testing.B) {
	d := GetDecoder()
	b.ReportAllocs()
	b.SetBytes(int64(len(otelEx1)))

	var v OTEL
	for i := 0; i < b.N; i++ {
		v.Reset()
		d.ResetBytes(otelEx1)
		if err := v.Decode(d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOTEL_Encode(b *testing.B) {
	d := DecodeBytes(otelEx1)
	var v OTEL

	require.NoError(b, v.Decode(d))

	b.ReportAllocs()
	e := GetEncoder()
	v.Encode(e)
	b.SetBytes(int64(len(e.Bytes())))

	for i := 0; i < b.N; i++ {
		e.Reset()
		v.Encode(e)
	}
}
