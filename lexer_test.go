package gocustojson

import (
	"fmt"
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
		{msg: "Lex valid line comment without AllowLineComments with newline", input: []byte("//\n"), expected: []Token{NewToken(ILLEGAL, []byte("//\n"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(false))},
		{msg: "Lex valid line comment without AllowLineComments without newline", input: []byte("//"), expected: []Token{NewToken(ILLEGAL, []byte("//"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(false))},
		{msg: "Lex invalid line comment without AllowLineComments", input: []byte("/"), expected: []Token{NewToken(ILLEGAL, []byte("/"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(false))},

		// lex line comment with AllowLineComments
		{msg: "Lex valid line comment with AllowLineComments", input: []byte("//\n"), expected: []Token{NewToken(LINE_COMMENT, []byte("//\n"), 1, 0, nil), NewToken(EOF, nil, 1, 3, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex valid line comment with AllowLineComments without newline", input: []byte("//"), expected: []Token{NewToken(LINE_COMMENT, []byte("//"), 1, 0, nil), NewToken(EOF, nil, 1, 2, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex invalid line comment with AllowLineComments", input: []byte("/"), expected: []Token{NewToken(ILLEGAL, []byte("/"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex valid line comment with AllowLineComments and text", input: []byte("// This is a line comment\n"), expected: []Token{NewToken(LINE_COMMENT, []byte("// This is a line comment\n"), 1, 0, nil), NewToken(EOF, nil, 1, 26, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex valid line comment with AllowLineComments and text without newline", input: []byte("// This is a line comment"), expected: []Token{NewToken(LINE_COMMENT, []byte("// This is a line comment"), 1, 0, nil), NewToken(EOF, nil, 1, 25, nil)}, cfg: NewConfig(WithAllowLineComments(true))},

		// lex block comment without AllowBlockComments
		{msg: "Lex valid block comment without AllowBlockComments", input: []byte("/**/"), expected: []Token{NewToken(ILLEGAL, []byte("/**/"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(false))},
		{msg: "Lex invalid block comment without AllowBlockComments", input: []byte("/*"), expected: []Token{NewToken(ILLEGAL, []byte("/*"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(false))},

		// lex block comment with AllowBlockComments
		{msg: "Lex valid block comment with AllowBlockComments", input: []byte("/**/"), expected: []Token{NewToken(BLOCK_COMMENT, []byte("/**/"), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
		// {msg: "Lex invalid block comment with AllowBlockComments", input: []byte("/*"), expected: []Token{NewToken(ILLEGAL, []byte("/*"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
		{msg: "Lex valid block comment with AllowBlockComments", input: []byte("/* This is a block comment*/"), expected: []Token{NewToken(BLOCK_COMMENT, []byte("/* This is a block comment*/"), 1, 0, nil), NewToken(EOF, nil, 1, 28, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunLexerTests(t, tests)
}

func TestLexString(t *testing.T) {
	var tests = []LexerTest{
		// lex string with single quotes without AllowSingleQuotes
		// {msg: "Lex valid single quote string without AllowLineComments", input: []byte(`''`), expected: []Token{NewToken(ILLEGAL, []byte(`''`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(false))},
		// {msg: "Lex invalid single quote string without AllowLineComments", input: []byte(`'`), expected: []Token{NewToken(ILLEGAL, []byte(`'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(false))},

		// lex string with single quotes with AllowSingleQuotes
		// {msg: "Lex invalid single quote string with AllowLineComments", input: []byte(`'`), expected: []Token{NewToken(ILLEGAL, []byte(`'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex valid simple single quote string with AllowLineComments", input: []byte(`'This is a single quote string'`), expected: []Token{NewToken(STRING, []byte(`'This is a single quote string'`), 1, 0, nil), NewToken(EOF, nil, 1, 31, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// lex string with double quotes with AllowSingleQuotes
		// {msg: "Lex invalid double string with AllowLineComments", input: []byte(`"`), expected: []Token{NewToken(ILLEGAL, []byte(`"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex valid double single quote string with AllowLineComments", input: []byte(`"This is a double quote string"`), expected: []Token{NewToken(STRING, []byte(`"This is a double quote string"`), 1, 0, nil), NewToken(EOF, nil, 1, 31, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		//////////////

		// {msg: "Lex empty double string", input: []byte(`""`), expected: []Token{NewToken(STRING, []byte(`""`), 1, 0, nil), NewToken(EOF, nil, 1, 2, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex empty single string", input: []byte(`''`), expected: []Token{NewToken(STRING, []byte(`''`), 1, 0, nil), NewToken(EOF, nil, 1, 2, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex non-empty double string", input: []byte(`"text"`), expected: []Token{NewToken(STRING, []byte(`"text"`), 1, 0, nil), NewToken(EOF, nil, 1, 6, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex non-empty single string", input: []byte(`'text'`), expected: []Token{NewToken(STRING, []byte(`'text'`), 1, 0, nil), NewToken(EOF, nil, 1, 6, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc double string", input: []byte(`"\""`), expected: []Token{NewToken(STRING, []byte(`"\""`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc single string", input: []byte(`'\''`), expected: []Token{NewToken(STRING, []byte(`'\''`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc backward slash double string", input: []byte(`"\\"`), expected: []Token{NewToken(STRING, []byte(`"\\"`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc backward slash single string", input: []byte(`'\\'`), expected: []Token{NewToken(STRING, []byte(`'\\'`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex invalid esc backward slash double string", input: []byte(`"\"`), expected: []Token{NewToken(ILLEGAL, []byte(`"\"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex invalid esc backward slash single string", input: []byte(`'\'`), expected: []Token{NewToken(ILLEGAL, []byte(`'\'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc forward slash double string", input: []byte(`"\/"`), expected: []Token{NewToken(STRING, []byte(`"\/"`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc forward slash single string", input: []byte(`'\/'`), expected: []Token{NewToken(STRING, []byte(`'\/'`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc backspace double string", input: []byte(`"\b"`), expected: []Token{NewToken(STRING, []byte(`"\b"`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc backspace single string", input: []byte(`'\b'`), expected: []Token{NewToken(STRING, []byte(`'\b'`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc form feed double string", input: []byte(`"\f"`), expected: []Token{NewToken(STRING, []byte(`"\f"`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc form feed single string", input: []byte(`'\f'`), expected: []Token{NewToken(STRING, []byte(`'\f'`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc new line double string", input: []byte(`"\n"`), expected: []Token{NewToken(STRING, []byte(`"\n"`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc new line single string", input: []byte(`'\n'`), expected: []Token{NewToken(STRING, []byte(`'\n'`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc carriage return double string", input: []byte(`"\r"`), expected: []Token{NewToken(STRING, []byte(`"\r"`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc carriage return single string", input: []byte(`'\r'`), expected: []Token{NewToken(STRING, []byte(`'\r'`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc tab double string", input: []byte(`"\t"`), expected: []Token{NewToken(STRING, []byte(`"\t"`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc tab single string", input: []byte(`'\t'`), expected: []Token{NewToken(STRING, []byte(`'\t'`), 1, 0, nil), NewToken(EOF, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex esc 4 hex digits double string", input: []byte(`"\u4F60"`), expected: []Token{NewToken(STRING, []byte(`"\u4F60"`), 1, 0, nil), NewToken(EOF, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex esc 4 hex digits single string", input: []byte(`'\u597D'`), expected: []Token{NewToken(STRING, []byte(`'\u597D'`), 1, 0, nil), NewToken(EOF, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex invalid esc 4 hex digits double string", input: []byte(`"\u00G1"`), expected: []Token{NewToken(ILLEGAL, []byte(`"\u00G1"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		// {msg: "Lex invalid esc 4 hex digits single string", input: []byte(`'\u00Z1'`), expected: []Token{NewToken(ILLEGAL, []byte(`'\u00Z1'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},

		// {msg: "Lex invalid escape double string without AllowEscapeChars", input: []byte(`"\Users"`), expected: []Token{NewToken(ILLEGAL, []byte(`"\Users"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowEscapeChars(false))},
		// {msg: "Lex invalid escape double string with AllowEscapeChars", input: []byte(`"\Users"`), expected: []Token{NewToken(STRING, []byte(`"\Users"`), 1, 0, nil), NewToken(EOF, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowEscapeChars(true))},
		// {msg: "Lex invalid escape single string without AllowEscapeChars", input: []byte(`'\Users'`), expected: []Token{NewToken(ILLEGAL, []byte(`'\Users'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowEscapeChars(false))},
		// {msg: "Lex invalid escape single string with AllowEscapeChars", input: []byte(`'\Users'`), expected: []Token{NewToken(STRING, []byte(`'\Users'`), 1, 0, nil), NewToken(EOF, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowEscapeChars(true))},

		{msg: "Lex invalid escape newline double string without AllowNewlineInStrings", input: []byte(`"\
		"`), expected: []Token{NewToken(ILLEGAL, []byte(`"\
		"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowNewlineInStrings(false))},
		{msg: "Lex invalid escape newline double string with AllowNewlineInStrings", input: []byte(`"\
		"`), expected: []Token{NewToken(STRING, []byte(`"\
		"`), 1, 0, nil), NewToken(EOF, nil, 1, 6, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowNewlineInStrings(true))},
	}

	for i, b := range []byte(`'This is a long string \
that continues on the next line \
without breaking.'`) {
		fmt.Printf("i: %d, p: %d, b: %q\n", i, b, b)
	}

	RunLexerTests(t, tests)
}
