package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_Raw(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		v := `{"foo":   [1, 2, 3, 4, 5]  }}`
		t.Run("Raw", func(t *testing.T) {
			d := DecodeStr(v)
			require.NoError(t, d.Obj(func(d *Decoder, key string) error {
				raw, err := d.Raw()
				require.NoError(t, err)
				t.Logf("%q", raw)
				return err
			}))
		})
		t.Run("RawAppend", func(t *testing.T) {
			d := DecodeStr(v)
			require.NoError(t, d.Obj(func(d *Decoder, key string) error {
				raw, err := d.RawAppend(nil)
				require.NoError(t, err)
				t.Logf("%q", raw)
				return err
			}))
		})
	})
	t.Run("Negative", func(t *testing.T) {
		v := `{"foo":   [1, 2, 3, 4, 5`
		t.Run("Raw", func(t *testing.T) {
			d := DecodeStr(v)
			var called bool
			require.Error(t, d.Obj(func(d *Decoder, key string) error {
				called = true
				raw, err := d.Raw()
				require.Error(t, err)
				require.Nil(t, raw)
				return err
			}))
			require.True(t, called, "should be called")
		})
		t.Run("RawAppend", func(t *testing.T) {
			d := DecodeStr(v)
			var called bool
			require.Error(t, d.Obj(func(d *Decoder, key string) error {
				called = true
				raw, err := d.RawAppend(make([]byte, 10))
				require.Error(t, err)
				require.Nil(t, raw)
				return err
			}))
			require.True(t, called, "should be called")
		})
	})
	t.Run("Reader", func(t *testing.T) {
		d := Decode(errReader{}, 0)
		if _, err := d.Raw(); err == nil {
			t.Error("should fail")
		}
		if _, err := d.RawAppend(nil); err == nil {
			t.Error("should fail")
		}
	})
}

func BenchmarkDecoder_Raw(b *testing.B) {
	data := []byte(`{"foo": [1,2,3,4,5,6,7,8,9,10,11,12,13,14]}`)
	b.ReportAllocs()

	var d Decoder
	for i := 0; i < b.N; i++ {
		d.ResetBytes(data)
		raw, err := d.Raw()
		if err != nil {
			b.Fatal(err)
		}
		if len(raw) == 0 {
			b.Fatal("blank")
		}
	}
}
