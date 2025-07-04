package jsonvx

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
			got := lexer.Tokens()

			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("got %v, expected %v", got, test.expected)
			}
		})
	}
}

func TestLexNothing(t *testing.T) {
	var tests = []LexerTest{
		// lex nothing
		{msg: "Lex nothing", input: []byte(``), expected: []Token{NewToken(EOF, NONE, nil, 1, 0, nil)}},
	}

	RunLexerTests(t, tests)
}

func TestLexWhiteSpace(t *testing.T) {
	var tests = []LexerTest{
		// lex standard whitespace without AllowExtraWS
		{msg: "Lex ws, space without AllowExtraWS", input: []byte(" "), expected: []Token{NewToken(WHITESPACE, NONE, []byte(" "), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, line feed without AllowExtraWS", input: []byte("\n"), expected: []Token{NewToken(WHITESPACE, NONE, []byte("\n"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, carriage return without AllowExtraWS", input: []byte("\r"), expected: []Token{NewToken(WHITESPACE, NONE, []byte("\r"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, horizontal tab without AllowExtraWS", input: []byte("\t"), expected: []Token{NewToken(WHITESPACE, NONE, []byte("\t"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},

		// lex standard whitespace with AllowExtraWS
		{msg: "Lex ws, space with AllowExtraWS", input: []byte(" "), expected: []Token{NewToken(WHITESPACE, NONE, []byte(" "), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, line feed with AllowExtraWS", input: []byte("\n"), expected: []Token{NewToken(WHITESPACE, NONE, []byte("\n"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, carriage return with AllowExtraWS", input: []byte("\r"), expected: []Token{NewToken(WHITESPACE, NONE, []byte("\r"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, horizontal tab with AllowExtraWS", input: []byte("\t"), expected: []Token{NewToken(WHITESPACE, NONE, []byte("\t"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},

		// lex extra whitespace without AllowExtraWS
		{msg: "Lex ws, line tabulation without AllowExtraWS", input: []byte("\v"), expected: []Token{NewToken(ILLEGAL, INVALID_WHITESPACE, []byte("\v"), 1, 0, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, form feed without AllowExtraWS", input: []byte("\f"), expected: []Token{NewToken(ILLEGAL, INVALID_WHITESPACE, []byte("\f"), 1, 0, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, next line without AllowExtraWS", input: []byte{0x85}, expected: []Token{NewToken(ILLEGAL, INVALID_WHITESPACE, []byte{0x85}, 1, 0, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},
		{msg: "Lex ws, no-break space without AllowExtraWS", input: []byte{0xA0}, expected: []Token{NewToken(ILLEGAL, INVALID_WHITESPACE, []byte{0xA0}, 1, 0, nil)}, cfg: NewConfig(WithAllowExtraWS(false))},

		// lex extra whitespace with AllowExtraWS
		{msg: "Lex ws, line tabulation with AllowExtraWS", input: []byte("\v"), expected: []Token{NewToken(WHITESPACE, NONE, []byte("\v"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, form feed with AllowExtraWS", input: []byte("\f"), expected: []Token{NewToken(WHITESPACE, NONE, []byte("\f"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, next line with AllowExtraWS", input: []byte{0x85}, expected: []Token{NewToken(WHITESPACE, NONE, []byte{0x85}, 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
		{msg: "Lex ws, no-break space with AllowExtraWS", input: []byte{0xA0}, expected: []Token{NewToken(WHITESPACE, NONE, []byte{0xA0}, 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowExtraWS(true))},
	}

	RunLexerTests(t, tests)
}

func TestLexDelimiters(t *testing.T) {
	var tests = []LexerTest{
		{msg: "Lex [", input: []byte("["), expected: []Token{NewToken(LEFT_SQUARE_BRACE, NONE, []byte("["), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig()},
		{msg: "Lex ]", input: []byte("]"), expected: []Token{NewToken(RIGHT_SQUARE_BRACE, NONE, []byte("]"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig()},
		{msg: "Lex {", input: []byte("{"), expected: []Token{NewToken(LEFT_CURLY_BRACE, NONE, []byte("{"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig()},
		{msg: "Lex }", input: []byte("}"), expected: []Token{NewToken(RIGHT_CURLY_BRACE, NONE, []byte("}"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig()},
		{msg: "Lex ,", input: []byte(","), expected: []Token{NewToken(COMMA, NONE, []byte(","), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig()},
		{msg: "Lex :", input: []byte(":"), expected: []Token{NewToken(COLON, NONE, []byte(":"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig()},
	}

	RunLexerTests(t, tests)
}

func TestLexComment(t *testing.T) {
	var tests = []LexerTest{
		// lex line comment without AllowLineComments
		{msg: "Lex valid line comment without AllowLineComments with newline", input: []byte("//\n"), expected: []Token{NewToken(ILLEGAL, INVALID_LINE_COMMENT, []byte("//\n"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(false))},
		{msg: "Lex valid line comment without AllowLineComments without newline", input: []byte("//"), expected: []Token{NewToken(ILLEGAL, INVALID_LINE_COMMENT, []byte("//"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(false))},
		{msg: "Lex invalid line comment without AllowLineComments", input: []byte("/"), expected: []Token{NewToken(ILLEGAL, INVALID_COMMENT, []byte("/"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(false))},

		// lex line comment with AllowLineComments
		{msg: "Lex valid line comment with AllowLineComments", input: []byte("//\n"), expected: []Token{NewToken(COMMENT, LINE_COMMENT, []byte("//\n"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 3, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex valid line comment with AllowLineComments without newline", input: []byte("//"), expected: []Token{NewToken(COMMENT, LINE_COMMENT, []byte("//"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 2, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex invalid line comment with AllowLineComments", input: []byte("/"), expected: []Token{NewToken(ILLEGAL, INVALID_COMMENT, []byte("/"), 1, 0, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex valid line comment with AllowLineComments and text", input: []byte("// This is a line comment\n"), expected: []Token{NewToken(COMMENT, LINE_COMMENT, []byte("// This is a line comment\n"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 26, nil)}, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Lex valid line comment with AllowLineComments and text without newline", input: []byte("// This is a line comment"), expected: []Token{NewToken(COMMENT, LINE_COMMENT, []byte("// This is a line comment"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 25, nil)}, cfg: NewConfig(WithAllowLineComments(true))},

		// lex block comment without AllowBlockComments
		{msg: "Lex valid block comment without AllowBlockComments", input: []byte("/**/"), expected: []Token{NewToken(ILLEGAL, INVALID_BLOCK_COMMENT, []byte("/**/"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(false))},
		{msg: "Lex invalid block comment without AllowBlockComments", input: []byte("/*"), expected: []Token{NewToken(ILLEGAL, INVALID_BLOCK_COMMENT, []byte("/*"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(false))},

		// lex block comment with AllowBlockComments
		{msg: "Lex valid block comment with AllowBlockComments", input: []byte("/**/"), expected: []Token{NewToken(COMMENT, BLOCK_COMMENT, []byte("/**/"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
		{msg: "Lex invalid block comment with AllowBlockComments", input: []byte("/*"), expected: []Token{NewToken(ILLEGAL, INVALID_BLOCK_COMMENT, []byte("/*"), 1, 0, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
		{msg: "Lex valid block comment with AllowBlockComments", input: []byte("/* This is a block comment*/"), expected: []Token{NewToken(COMMENT, BLOCK_COMMENT, []byte("/* This is a block comment*/"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 28, nil)}, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunLexerTests(t, tests)
}

func TestLexString(t *testing.T) {
	var tests = []LexerTest{
		// lex string with single quotes without AllowSingleQuotes
		{msg: "Lex valid single quote string without AllowLineComments", input: []byte(`''`), expected: []Token{NewToken(ILLEGAL, INVALID_STRING, []byte(`''`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(false))},
		{msg: "Lex invalid single quote string without AllowLineComments", input: []byte(`'`), expected: []Token{NewToken(ILLEGAL, INVALID_STRING, []byte(`'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(false))},

		// lex string with single quotes with AllowSingleQuotes
		{msg: "Lex invalid single quote string with AllowLineComments", input: []byte(`'`), expected: []Token{NewToken(ILLEGAL, INVALID_STRING, []byte(`'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex valid simple single quote string with AllowLineComments", input: []byte(`'This is a single quote string'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'This is a single quote string'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 31, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex empty single string", input: []byte(`''`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`''`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 2, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex non-empty single string", input: []byte(`'text'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'text'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 6, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc single string", input: []byte(`'\''`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\''`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc backward slash single string", input: []byte(`'\\'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\\'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex invalid esc backward slash single string", input: []byte(`'\'`), expected: []Token{NewToken(ILLEGAL, INVALID_STRING, []byte(`'\'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc forward slash single string", input: []byte(`'\/'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\/'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc backspace single string", input: []byte(`'\b'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\b'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc form feed single string", input: []byte(`'\f'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\f'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc new line single string", input: []byte(`'\n'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\n'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc carriage return single string", input: []byte(`'\r'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\r'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc tab single string", input: []byte(`'\t'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\t'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc 4 hex digits single string", input: []byte(`'\u597D'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\u597D'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex invalid esc 4 hex digits single string", input: []byte(`'\u00Z1'`), expected: []Token{NewToken(ILLEGAL, INVALID_HEX_STRING, []byte(`'\u00Z1'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex invalid escape single string without AllowEscapeChars", input: []byte(`'\Users'`), expected: []Token{NewToken(ILLEGAL, INVALID_ESCAPED_STRING, []byte(`'\Users'`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowOtherEscapeChars(false))},
		{msg: "Lex invalid escape single string with AllowEscapeChars", input: []byte(`'\Users'`), expected: []Token{NewToken(STRING, SINGLE_QUOTED, []byte(`'\Users'`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowOtherEscapeChars(true))},
		{msg: "Lex invalid escape newline double string with AllowNewlineInStrings", input: []byte(`"\
		"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\
		"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 6, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowNewlineInStrings(true))},

		// lex string with double quotes with AllowSingleQuotes
		{msg: "Lex invalid double string with AllowLineComments", input: []byte(`"`), expected: []Token{NewToken(ILLEGAL, INVALID_STRING, []byte(`"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex valid double single quote string with AllowLineComments", input: []byte(`"This is a double quote string"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"This is a double quote string"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 31, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex empty double string", input: []byte(`""`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`""`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 2, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex non-empty double string", input: []byte(`"text"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"text"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 6, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc double string", input: []byte(`"\""`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\""`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc backward slash double string", input: []byte(`"\\"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\\"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex invalid esc backward slash double string", input: []byte(`"\"`), expected: []Token{NewToken(ILLEGAL, INVALID_STRING, []byte(`"\"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc forward slash double string", input: []byte(`"\/"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\/"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc backspace double string", input: []byte(`"\b"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\b"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc form feed double string", input: []byte(`"\f"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\f"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc new line double string", input: []byte(`"\n"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\n"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc carriage return double string", input: []byte(`"\r"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\r"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc tab double string", input: []byte(`"\t"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\t"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex esc 4 hex digits double string", input: []byte(`"\u4F60"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\u4F60"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex invalid esc 4 hex digits double string", input: []byte(`"\u00G1"`), expected: []Token{NewToken(ILLEGAL, INVALID_HEX_STRING, []byte(`"\u00G1"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Lex invalid escape double string without AllowEscapeChars", input: []byte(`"\Users"`), expected: []Token{NewToken(ILLEGAL, INVALID_ESCAPED_STRING, []byte(`"\Users"`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowOtherEscapeChars(false))},
		{msg: "Lex invalid escape double string with AllowEscapeChars", input: []byte(`"\Users"`), expected: []Token{NewToken(STRING, DOUBLE_QUOTED, []byte(`"\Users"`), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowOtherEscapeChars(true))},
		{msg: "Lex invalid escape newline double string without AllowNewlineInStrings", input: []byte(`"\
		// "`), expected: []Token{NewToken(ILLEGAL, INVALID_NEWLINE_STRING, []byte(`"\
		// "`), 1, 0, nil)}, cfg: NewConfig(WithAllowSingleQuotes(true), WithAllowNewlineInStrings(false))},
	}

	RunLexerTests(t, tests)
}

func TestLexIdent(t *testing.T) {
	var tests = []LexerTest{
		// lex ident without AllowUnquoted
		{msg: "Lex valid ident without AllowUnquoted", input: []byte("ident"), expected: []Token{NewToken(ILLEGAL, INVALID_CHARACTER, []byte("ident"), 1, 0, nil)}, cfg: NewConfig(WithAllowUnquoted(false))},
		{msg: "Lex invalid ident without AllowUnquoted", input: []byte("1ident"), expected: []Token{NewToken(ILLEGAL, INVALID_CHARACTER, []byte("1ident"), 1, 0, nil)}, cfg: NewConfig(WithAllowUnquoted(false))},

		// lex ident with AllowUnquoted
		{msg: "Lex valid ident with AllowUnquoted", input: []byte("ident"), expected: []Token{NewToken(STRING, IDENT, []byte("ident"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 5, nil)}, cfg: NewConfig(WithAllowUnquoted(true))},
		{msg: "Lex invalid ident with AllowUnquoted", input: []byte("1ident"), expected: []Token{NewToken(ILLEGAL, INVALID_CHARACTER, []byte("1ident"), 1, 0, nil)}, cfg: NewConfig(WithAllowUnquoted(true))},
		{msg: "Lex valid ident with AllowUnquoted with space", input: []byte("ident "), expected: []Token{NewToken(STRING, IDENT, []byte("ident"), 1, 0, nil), NewToken(WHITESPACE, NONE, []byte(" "), 1, 5, nil), NewToken(EOF, NONE, nil, 1, 6, nil)}, cfg: NewConfig(WithAllowUnquoted(true))},
	}

	RunLexerTests(t, tests)
}

func TestLexNull(t *testing.T) {
	var tests = []LexerTest{
		{msg: "Lex valid null", input: []byte("null"), expected: []Token{NewToken(NULL, NONE, []byte("null"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig()},
		{msg: "Lex invalid NULL case-sensitive", input: []byte("NULL"), expected: []Token{NewToken(ILLEGAL, INVALID_CHARACTER, []byte("NULL"), 1, 0, nil)}, cfg: NewConfig()},
		{msg: "Lex invalid null nil", input: []byte("nil"), expected: []Token{NewToken(ILLEGAL, INVALID_NULL, []byte("nil"), 1, 0, nil)}, cfg: NewConfig()},
	}

	RunLexerTests(t, tests)
}

func TestLexTrue(t *testing.T) {
	var tests = []LexerTest{
		{msg: "Lex valid true", input: []byte("true"), expected: []Token{NewToken(TRUE, NONE, []byte("true"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig()},
		{msg: "Lex invalid tRuE case-sensitive", input: []byte("tRuE"), expected: []Token{NewToken(ILLEGAL, INVALID_TRUE, []byte("tRuE"), 1, 0, nil)}, cfg: NewConfig()},
		{msg: "Lex invalid t", input: []byte("t"), expected: []Token{NewToken(ILLEGAL, INVALID_TRUE, []byte("t"), 1, 0, nil)}, cfg: NewConfig()},
	}

	RunLexerTests(t, tests)
}

func TestLexFalse(t *testing.T) {
	var tests = []LexerTest{
		{msg: "Lex valid false", input: []byte("false"), expected: []Token{NewToken(FALSE, NONE, []byte("false"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 5, nil)}, cfg: NewConfig()},
		{msg: "Lex invalid fAlSe case-sensitive", input: []byte("fAlSe"), expected: []Token{NewToken(ILLEGAL, INVALID_FALSE, []byte("fAlSe"), 1, 0, nil)}, cfg: NewConfig()},
		{msg: "Lex invalid f", input: []byte("f"), expected: []Token{NewToken(ILLEGAL, INVALID_FALSE, []byte("f"), 1, 0, nil)}, cfg: NewConfig()},
	}

	RunLexerTests(t, tests)
}

func TestLexNumber(t *testing.T) {
	var tests = []LexerTest{
		// Lex NaN
		{msg: "Lex NaN without AllowNaN", input: []byte("NaN"), expected: []Token{NewToken(ILLEGAL, INVALID_NaN, []byte("NaN"), 1, 0, nil)}, cfg: NewConfig()},
		{msg: "Lex NaN with AllowNaN", input: []byte("NaN"), expected: []Token{NewToken(NUMBER, NaN, []byte("NaN"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 3, nil)}, cfg: NewConfig(WithAllowNaN(true))},
		{msg: "Lex pos NaN with AllowNaN", input: []byte("+NaN"), expected: []Token{NewToken(NUMBER, NaN, []byte("+NaN"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowNaN(true), WithAllowLeadingPlus(true))},
		{msg: "Lex neg NaN with AllowNaN", input: []byte("-NaN"), expected: []Token{NewToken(NUMBER, NaN, []byte("-NaN"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowNaN(true))},

		// Lex Infinity
		{msg: "Lex Infinity without AllowInfinity", input: []byte("Infinity"), expected: []Token{NewToken(ILLEGAL, INVALID_INF, []byte("Infinity"), 1, 0, nil)}, cfg: NewConfig(WithAllowInfinity(false))},
		{msg: "Lex Infinity with AllowInfinity", input: []byte("Infinity"), expected: []Token{NewToken(NUMBER, INF, []byte("Infinity"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowInfinity(true))},
		{msg: "Lex pos Infinity without AllowInfinity and AllowLeadingPlus", input: []byte("+Infinity"), expected: []Token{NewToken(ILLEGAL, INVALID_LEADING_PLUS, []byte("+Infinity"), 1, 0, nil)}, cfg: NewConfig(WithAllowInfinity(false), WithAllowLeadingPlus(false))},
		{msg: "Lex pos Infinity with AllowInfinity and AllowLeadingPlus", input: []byte("+Infinity"), expected: []Token{NewToken(NUMBER, INF, []byte("+Infinity"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 9, nil)}, cfg: NewConfig(WithAllowInfinity(true), WithAllowLeadingPlus(true))},
		{msg: "Lex neg Infinity without AllowInfinity and AllowLeadingPlus", input: []byte("-Infinity"), expected: []Token{NewToken(ILLEGAL, INVALID_INF, []byte("-Infinity"), 1, 0, nil)}, cfg: NewConfig(WithAllowInfinity(false), WithAllowLeadingPlus(false))},
		{msg: "Lex neg Infinity with AllowInfinity and AllowLeadingPlus", input: []byte("-Infinity"), expected: []Token{NewToken(NUMBER, INF, []byte("-Infinity"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 9, nil)}, cfg: NewConfig(WithAllowInfinity(true), WithAllowLeadingPlus(true))},

		// Lex integer
		{msg: "Lex valid integer", input: []byte("0"), expected: []Token{NewToken(NUMBER, INTEGER, []byte("0"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 1, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex invalid integer leading zero", input: []byte("00"), expected: []Token{NewToken(ILLEGAL, INVALID_LEADING_ZERO, []byte("00"), 1, 0, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex valid long integer", input: []byte("1234567890"), expected: []Token{NewToken(NUMBER, INTEGER, []byte("1234567890"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 10, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex valid pos integer", input: []byte("+1234567890"), expected: []Token{NewToken(NUMBER, INTEGER, []byte("+1234567890"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 11, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},

		// Lex float
		{msg: "Lex valid float", input: []byte("0.0"), expected: []Token{NewToken(NUMBER, FLOAT, []byte("0.0"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 3, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex invalid float leading zero", input: []byte("03.4"), expected: []Token{NewToken(ILLEGAL, INVALID_LEADING_ZERO, []byte("03.4"), 1, 0, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex valid float", input: []byte("123.456"), expected: []Token{NewToken(NUMBER, FLOAT, []byte("123.456"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 7, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex valid pos float", input: []byte("+123.456"), expected: []Token{NewToken(NUMBER, FLOAT, []byte("+123.456"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex valid neg float", input: []byte("-123.456"), expected: []Token{NewToken(NUMBER, FLOAT, []byte("-123.456"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 8, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},

		// lex scientific notation
		{msg: "Lex valid e sci-not with", input: []byte("123e456"), expected: []Token{NewToken(NUMBER, SCI_NOT, []byte("123e456"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 7, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex invalid e sci-not with", input: []byte("e456"), expected: []Token{NewToken(ILLEGAL, INVALID_CHARACTER, []byte("e456"), 1, 0, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex valid E sci-not with", input: []byte("123E456"), expected: []Token{NewToken(NUMBER, SCI_NOT, []byte("123E456"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 7, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},
		{msg: "Lex invalid E sci-not with", input: []byte("E456"), expected: []Token{NewToken(ILLEGAL, INVALID_CHARACTER, []byte("E456"), 1, 0, nil)}, cfg: NewConfig(WithAllowLeadingPlus(true))},

		// Lex hex
		{msg: "Lex hex without AllowHexNumbers", input: []byte("0x1A"), expected: []Token{NewToken(ILLEGAL, INVALID_HEX_NUMBER, []byte("0x1A"), 1, 0, nil)}, cfg: NewConfig(WithAllowHexNumbers(false))},
		{msg: "Lex hex with AllowHexNumbers", input: []byte("0x1A"), expected: []Token{NewToken(NUMBER, HEX, []byte("0x1A"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 4, nil)}, cfg: NewConfig(WithAllowHexNumbers(true))},
		{msg: "Lex pos hex without AllowHexNumbers and AllowLeadingPlus", input: []byte("+0x1A"), expected: []Token{NewToken(ILLEGAL, INVALID_LEADING_PLUS, []byte("+0x1A"), 1, 0, nil)}, cfg: NewConfig(WithAllowHexNumbers(false), WithAllowLeadingPlus(false))},
		{msg: "Lex pos hex with AllowHexNumbers and AllowLeadingPlus", input: []byte("+0x1A"), expected: []Token{NewToken(NUMBER, HEX, []byte("+0x1A"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 5, nil)}, cfg: NewConfig(WithAllowHexNumbers(true), WithAllowLeadingPlus(true))},
		{msg: "Lex neg hex without AllowHexNumbers and AllowLeadingPlus", input: []byte("-0x1A"), expected: []Token{NewToken(ILLEGAL, INVALID_HEX_NUMBER, []byte("-0x1A"), 1, 0, nil)}, cfg: NewConfig(WithAllowHexNumbers(false), WithAllowLeadingPlus(false))},
		{msg: "Lex neg hex with AllowHexNumbers and AllowLeadingPlus", input: []byte("-0x1A"), expected: []Token{NewToken(NUMBER, HEX, []byte("-0x1A"), 1, 0, nil), NewToken(EOF, NONE, nil, 1, 5, nil)}, cfg: NewConfig(WithAllowHexNumbers(true), WithAllowLeadingPlus(true))},
	}

	RunLexerTests(t, tests)
}

// func TestXYZ(t *testing.T) {
// 	b := []byte(`// comment
// null // comment
// // comment`)
// 	lexer := NewLexer(b, NewConfig(WithAllowBlockComments(true), WithAllowLineComments(true)))

// 	fmt.Println(lexer.Tokens())
// }

func BenchmarkLexerReadChar(b *testing.B) {
	l := NewLexer([]byte("lexer"), NewConfig())
	for i := 0; i < b.N; i++ {
		l.readChar()
	}
}

func BenchmarkLexerPeek(b *testing.B) {
	l := NewLexer([]byte("lexer"), NewConfig())
	for i := 0; i < b.N; i++ {
		l.peek()
	}
}

func BenchmarkLexerPeekBy(b *testing.B) {
	l := NewLexer([]byte("lexer"), NewConfig())
	for i := 0; i < b.N; i++ {
		l.peekBy(1)
	}
}

func BenchmarkLexerPrev(b *testing.B) {
	l := NewLexer([]byte("lexer"), NewConfig())
	for i := 0; i < b.N; i++ {
		l.prev()
	}
}

func BenchmarkLexerPrevBy(b *testing.B) {
	l := NewLexer([]byte("lexer"), NewConfig())
	for i := 0; i < b.N; i++ {
		l.prevBy(1)
	}
}
