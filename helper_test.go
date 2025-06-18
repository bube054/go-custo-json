package gocustojson

import "testing"

func TestIsNewLine(t *testing.T) {
	var tests = []struct {
		p1       byte
		expected bool
	}{
		{p1: '\n', expected: true},
		{p1: '\t', expected: false},
	}

	for _, test := range tests {
		t.Run(string(test.p1), func(t *testing.T) {
			got := IsNewLine(test.p1)

			if got != test.expected {
				t.Errorf("got %t, expected %t", got, test.expected)
			}
		})
	}
}
