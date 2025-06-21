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
		{msg: "Valid hex only digits", p1: [4]byte{48, 48, 51, 49}, expected: true},
		{msg: "Valid hex uppercase", p1: [4]byte{'A', 'B', 'C', 'F'}, expected: true},
		{msg: "Valid hex lowercase", p1: [4]byte{'a', 'b', 'c', 'f'}, expected: true},
		{msg: "Mixed valid hex", p1: [4]byte{'0', 'F', 'a', '9'}, expected: true},
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
