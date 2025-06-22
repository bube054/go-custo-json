package gocustojson

import (
	"testing"
)

func TestIsNewLine(t *testing.T) {
	var tests = []struct {
		msg      string
		p1       byte
		expected bool
	}{
		{msg: "Valid newline", p1: '\n', expected: true},
		{msg: "Invalid newline", p1: '\t', expected: false},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			got := IsNewLine(test.p1)

			if got != test.expected {
				t.Errorf("got %t, expected %t", got, test.expected)
			}
		})
	}
}

func TestIsWhiteSpace(t *testing.T) {
	var tests = []struct {
		msg      string
		p1       byte
		p2       bool
		expected bool
	}{
		// Standard whitespaces without AllowExtraWS
		{msg: "Valid space, without AllowExtraWS", p1: '\x20', p2: false, expected: true},
		{msg: "Valid line feed, without AllowExtraWS", p1: '\x0A', p2: false, expected: true},
		{msg: "Valid carriage return, without AllowExtraWS", p1: '\x0D', p2: false, expected: true},
		{msg: "Valid horizontal tab, without AllowExtraWS", p1: '\x09', p2: false, expected: true},

		//  Standard whitespaces with AllowExtraWS
		{msg: "Valid space, with AllowExtraWS", p1: '\x20', p2: true, expected: true},
		{msg: "Valid line feed, with AllowExtraWS", p1: '\x0A', p2: true, expected: true},
		{msg: "Valid carriage return, with AllowExtraWS", p1: '\x0D', p2: true, expected: true},
		{msg: "Valid horizontal tab, with AllowExtraWS", p1: '\x09', p2: true, expected: true},

		//  Additional whitespaces without AllowExtraWS
		{msg: "Valid line tabulation, with AllowExtraWS", p1: '\x0B', p2: true, expected: true},
		{msg: "Valid form feed, with AllowExtraWS", p1: '\x0C', p2: true, expected: true},
		{msg: "Valid next line, with AllowExtraWS", p1: '\u0085', p2: true, expected: true},
		{msg: "Valid no break space, with AllowExtraWS", p1: '\u00A0', p2: true, expected: true},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			got := IsWhiteSpace(test.p1, test.p2)

			if got != test.expected {
				t.Errorf("got %t, expected %t", got, test.expected)
			}
		})
	}
}

func TestIs4HexDigits(t *testing.T) {
	var tests = []struct {
		msg      string
		p1       [4]byte
		expected bool
	}{
		// Valid 4 hex digits
		{msg: "Valid hex only digits", p1: [4]byte{48, 48, 51, 49}, expected: true},
		{msg: "Valid hex uppercase", p1: [4]byte{'A', 'B', 'C', 'F'}, expected: true},
		{msg: "Valid hex lowercase", p1: [4]byte{'a', 'b', 'c', 'f'}, expected: true},
		{msg: "Mixed valid hex", p1: [4]byte{'0', 'F', 'a', '9'}, expected: true},

		// Invalid 4 hex digits
		{msg: "Invalid non-hex character", p1: [4]byte{'0', 'G', 'a', '9'}, expected: false},
		{msg: "All invalid", p1: [4]byte{'g', 'h', 'z', '!'}, expected: false},
		{msg: "Edge case: just below range", p1: [4]byte{'/', 'A', 'B', 'C'}, expected: false},
		{msg: "Edge case: just above range", p1: [4]byte{'G', '0', '1', '2'}, expected: false},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			got := Is4HexDigits(test.p1)

			if got != test.expected {
				t.Errorf("got %t, expected %t", got, test.expected)
			}
		})
	}
}

func TestIsJSIdentifier(t *testing.T) {
	var tests = []struct {
		msg      string
		p1       []byte
		expected bool
	}{
		// Valid js indent
		{msg: "Latin letter with accent", p1: []byte("café"), expected: true},
		{msg: "German umlaut", p1: []byte("über"), expected: true},
		{msg: "'ï' is a Unicode letter", p1: []byte("naïve"), expected: true},
		{msg: "Cyrillic (Russian)", p1: []byte("привет"), expected: true},
		{msg: "Greek", p1: []byte("κόσμος"), expected: true},
		{msg: "Chinese (means \"variable\")", p1: []byte("变量"), expected: true},
		{msg: "Arabic (means \"variable\")", p1: []byte("متغير"), expected: true},
		{msg: "Hindi (Devanagari script)", p1: []byte("फ़ाइल"), expected: true},

		// Invalid js indent
		{msg: "starts with a digit", p1: []byte("1variable"), expected: false},
		{msg: "hyphen is not allowed", p1: []byte("var-name"), expected: false},
		{msg: "emoji is *not* a letter", p1: []byte("💻"), expected: false},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			got := IsJSIdentifier(test.p1)

			if got != test.expected {
				t.Errorf("got %t, expected %t", got, test.expected)
			}
		})
	}
}

func TestIsPossibleJSIdentifier(t *testing.T) {
	var tests = []struct {
		msg      string
		p1       byte
		expected bool
	}{
		// Valid possible js indent
		{msg: "dollar sign", p1: '$', expected: true},
		{msg: "underscore", p1: '_', expected: true},
		{msg: "digit", p1: '8', expected: true},
		{msg: "lowercase letter", p1: 'b', expected: true},
		{msg: "upper letter", p1: 'F', expected: true},
		{msg: "unicode letter", p1: 'ï', expected: true},

		// Invalid possible js indent
		{msg: "space", p1: ' ', expected: false},
		{msg: "nothing", p1: 0, expected: false},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			got := IsPossibleJSIdentifier(test.p1)

			if got != test.expected {
				t.Errorf("got %t, expected %t", got, test.expected)
			}
		})
	}
}
