package jx

// Bool reads a json object as Bool
func (d *Decoder) Bool() (bool, error) {
	if err := d.skipSpace(); err != nil {
		return false, err
	}

	var buf [4]byte
	if err := d.readExact4(&buf); err != nil {
		return false, err
	}

	switch string(buf[:]) {
	case "true":
		return true, nil
	case "fals":
		c, err := d.byte()
		if err != nil {
			return false, err
		}
		if c != 'e' {
			return false, badToken(c)
		}
		return false, nil
	default:
		switch c := buf[0]; c {
		case 't':
			const encodedTrue = 't' | 'r'<<8 | 'u'<<16 | 'e'<<24
			return false, findInvalidToken4(buf, encodedTrue)
		case 'f':
			const encodedAlse = 'a' | 'l'<<8 | 's'<<16 | 'e'<<24
			return false, findInvalidToken4(buf, encodedAlse)
		default:
			return false, badToken(c)
		}
	}
}
