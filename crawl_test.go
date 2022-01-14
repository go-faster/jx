package jx

import "github.com/go-faster/errors"

func crawlValue(d *Decoder) error {
	switch d.Next() {
	case Null:
		if err := d.Null(); err != nil {
			return err
		}
	case Number:
		if _, err := d.Float64(); err != nil {
			return err
		}
	case Bool:
		if _, err := d.Bool(); err != nil {
			return err
		}
	case String:
		if _, err := d.Str(); err != nil {
			return err
		}
	case Array:
		return d.Arr(crawlValue)
	case Object:
		return d.ObjBytes(func(d *Decoder, key []byte) error {
			return crawlValue(d)
		})
	default:
		return errors.New("invalid token")
	}
	return nil
}
