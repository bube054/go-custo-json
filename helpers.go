package gocustojson

func IsNewLine(char byte) bool {
	return char == '\n'
}

// IsWhiteSpace reports whether a byte is a whitespace.
func IsWhiteSpace(char byte, AllowExtraWS bool) bool {
	switch char {
	// space, line feed, carriage return, horizontal tab
	case
		'\x20', // space
		'\x0A', // line feed
		'\x0D', // carriage return
		'\x09': // horizontal tab
		return true
	case
		'\x0B',   // line tabulation
		'\x0C',   // form feed
		'\u0085', // next line
		'\u00A0': // no break space
		return AllowExtraWS
	default:
		return false
	}
}
