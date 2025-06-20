package gocustojson

import "testing"

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
