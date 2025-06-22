package gocustojson

import (
	"unicode"
	"unicode/utf8"
)

func IsNewLine(char byte) bool {
	return char == '\n'
}

// IsWhiteSpace reports whether a byte is a whitespace.
func IsWhiteSpace(char byte, AllowExtraWS bool) bool {
	switch char {
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

// IsWhiteSpace reports whether a byte slice is a hex.
func Is4HexDigits(chars [4]byte) bool {
	for _, char := range chars {
		if (char >= 48 && char <= 57) || (char >= 65 && char <= 70) || (char >= 97 && char <= 102) {
			continue
		}

		return false
	}

	return true
}

func IsJSIdentifier(input []byte) bool {
	str := string(input)
	for i, char := range str {
		isDollarOrUnderscore := char == '$' || char == '_'
		if i == 0 {
			if !(unicode.IsLetter(char) || isDollarOrUnderscore) {
				return false
			}
		} else {
			if !(unicode.IsLetter(char) ||
				unicode.IsDigit(char) ||
				isDollarOrUnderscore ||
				isJSCombiningMark(char)) {
				return false
			}
		}
	}
	return utf8.Valid(input)
}

func isJSCombiningMark(r rune) bool {
	return unicode.In(r, unicode.Mn, unicode.Mc)
}

func IsPossibleJSIdentifier(b byte) bool {
	if b >= 48 && b <= 57 {
		return true
	}

	if b == 36 || b == 95 {
		return  true
	}

	if unicode.IsLetter(rune(b)) {
		return true
	}

	return false
}
