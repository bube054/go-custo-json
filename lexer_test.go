package gocustojson

import (
	"reflect"
	"testing"
)

type LexerTest struct {
	msg      string
	input    []byte
	cfg      *Config
	expected Tokens
}

func RunLexerTests(t *testing.T, tests []LexerTest) {
	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			lexer := NewLexer(test.input, test.cfg)
			got := lexer.GenerateTokens()

			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("got %v, expected %v", got, test.expected)
			}
		})
	}
}

func TestLexNothing(t *testing.T) {
	var tests = []LexerTest{
		// lex nothing
		{msg: "Lex nothing", input: []byte(``), expected: []Token{NewToken(EOF, nil, 1, 0, nil)}},
	}

	RunLexerTests(t, tests)
}

func TestLexWhiteSpace(t *testing.T) {
	var tests = []LexerTest{
		// lex standard whitespace without AllowExtraWS
		{msg: "Lex ws, space without AllowExtraWS", input: []byte(" "), expected: []Token{NewToken(WHITESPACE, []byte(" "), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, line feed without AllowExtraWS", input: []byte("\n"), expected: []Token{NewToken(WHITESPACE, []byte("\n"), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, carriage return without AllowExtraWS", input: []byte("\r"), expected: []Token{NewToken(WHITESPACE, []byte("\r"), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, horizontal tab without AllowExtraWS", input: []byte("\t"), expected: []Token{NewToken(WHITESPACE, []byte("\t"), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},

		// lex standard whitespace with AllowExtraWS
		{msg: "Lex ws, space with AllowExtraWS", input: []byte(" "), expected: []Token{NewToken(WHITESPACE, []byte(" "), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, line feed with AllowExtraWS", input: []byte("\n"), expected: []Token{NewToken(WHITESPACE, []byte("\n"), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, carriage return with AllowExtraWS", input: []byte("\r"), expected: []Token{NewToken(WHITESPACE, []byte("\r"), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, horizontal tab with AllowExtraWS", input: []byte("\t"), expected: []Token{NewToken(WHITESPACE, []byte("\t"), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},

		// lex extra whitespace without AllowExtraWS
		{msg: "Lex ws, line tabulation without AllowExtraWS", input: []byte("\v"), expected: []Token{NewToken(ILLEGAL, []byte("\v"), 1, 0, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, form feed without AllowExtraWS", input: []byte("\f"), expected: []Token{NewToken(ILLEGAL, []byte("\f"), 1, 0, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, next line without AllowExtraWS", input: []byte{0x85}, expected: []Token{NewToken(ILLEGAL, []byte{0x85}, 1, 0, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, no-break space without AllowExtraWS", input: []byte{0xA0}, expected: []Token{NewToken(ILLEGAL, []byte{0xA0}, 1, 0, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},

		// lex extra whitespace with AllowExtraWS
		{msg: "Lex ws, line tabulation with AllowExtraWS", input: []byte("\v"), expected: []Token{NewToken(WHITESPACE, []byte("\v"), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, form feed with AllowExtraWS", input: []byte("\f"), expected: []Token{NewToken(WHITESPACE, []byte("\f"), 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, next line with AllowExtraWS", input: []byte{0x85}, expected: []Token{NewToken(WHITESPACE, []byte{0x85}, 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, no-break space with AllowExtraWS", input: []byte{0xA0}, expected: []Token{NewToken(WHITESPACE, []byte{0xA0}, 1, 0, nil), NewToken(EOF, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
	}

	RunLexerTests(t, tests)
}

func TestLexComment(t *testing.T) {
	var tests = []LexerTest{
		// lex line comment without AllowLineComments
		{msg: "Lex valid line comment without AllowLineComments", input: []byte("//\n"), expected: []Token{NewToken(ILLEGAL, []byte("//\n"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(false))},
		{msg: "Lex invalid line comment without AllowLineComments", input: []byte("//"), expected: []Token{NewToken(ILLEGAL, []byte("//"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(false))},

		// lex line comment with AllowLineComments
		{msg: "Lex valid line comment with AllowLineComments", input: []byte("//\n"), expected: []Token{NewToken(LINE_COMMENT, []byte("//\n"), 1, 0, nil), NewToken(EOF, nil, 1, 3, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex invalid line comment with AllowLineComments", input: []byte("//"), expected: []Token{NewToken(ILLEGAL, []byte("//"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex valid line comment with AllowLineComments and text", input: []byte("// This is a line comment\n"), expected: []Token{NewToken(LINE_COMMENT, []byte("// This is a line comment\n"), 1, 0, nil), NewToken(EOF, nil, 1, 26, nil)}, cfg: NewConfig(WithAllowLineComments(true))},

		// lex block comment without AllowBlockComments
		{msg: "Lex valid block comment without AllowBlockComments", input: []byte("/**/"), expected: []Token{NewToken(ILLEGAL, []byte("/**/"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(false))},
		{msg: "Lex invalid block comment without AllowBlockComments", input: []byte("/*"), expected: []Token{NewToken(ILLEGAL, []byte("/*"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(false))},

		// lex block comment with AllowBlockComments
		{msg: "Lex valid block comment with AllowBlockComments", input: []byte("/**/"), expected: []Token{NewToken(BLOCK_COMMENT, []byte("/**/"), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
		{msg: "Lex invalid block comment with AllowBlockComments", input: []byte("/*"), expected: []Token{NewToken(ILLEGAL, []byte("/*"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
		{msg: "Lex valid block comment with AllowBlockComments", input: []byte("/* This is a block comment*/"), expected: []Token{NewToken(BLOCK_COMMENT, []byte("/* This is a block comment*/"), 1, 0, nil), NewToken(EOF, nil, 1, 28, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunLexerTests(t, tests)
}
