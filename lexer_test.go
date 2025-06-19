package gocustojson

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	var tests = []struct {
		input    string
		cfg      *Config
		expected []Token
	}{
		// lex white spaces
		{input: "", expected: []Token{NewToken(EOF, nil, 1, 0)}},
		// {p1: '\n', expected: true},
		// {p1: '\t', expected: false},
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			lexer := NewLexer(test.input, test.cfg)
			got := lexer.GenerateTokens()

			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("got %v, expected %v", got, test.expected)
			}
		})
	}
}
