package gocustojson

import (
	"fmt"
	"strings"
)

// TokenSubKind represents the subtype of a token, such as numeric format or identifier type.
type TokenSubKind int

const (
	NONE TokenSubKind = iota // NONE represents the absence of a sub kind.

	QUOTED // Quoted represents a quoted string value.
	IDENT  // IDENT represents an unquoted identifier.

	INTEGER // INTEGER represents an integer (positive or negative).
	FLOAT   // FLOAT represents a floating-point number (positive or negative).
	SCI_NOT // SCI_NOT represents scientific notation (e.g., 1e10).
	HEX     // HEX represents a hexadecimal number (e.g., 0xFF).
	INF     // INF represents an Infinity literal (positive or negative).
	NaN     // NaN represents a Not-a-Number literal.
)

// String returns a string representation of the TokenSubKind.
func (t TokenSubKind) String() string {
	m := map[TokenSubKind]string{
		0: "NONE",
		1: "QUOTED",
		2: "IDENT",
		3: "INTEGER",
		4: "FLOAT",
		5: "SCI_NOT",
		6: "HEX",
		7: "INF",
		8: "NaN",
	}

	str := m[t]
	return str
}

// TokenKind represents the primary kind of a token, used in lexical analysis.
type TokenKind int

const (
	EOF                TokenKind = iota // EOF indicates the end of input.
	ILLEGAL                             // ILLEGAL indicates an unrecognized or invalid token.
	WHITESPACE                          // WHITESPACE represents any space character.
	LINE_COMMENT                        // LINE_COMMENT represents a single-line comment (// ...).
	BLOCK_COMMENT                       // BLOCK_COMMENT represents a block comment (/* ... */).
	STRING                              // STRING represents a string literal.
	NUMBER                              // NUMBER represents any numeric literal.
	NULL                                // NULL represents a null value.
	TRUE                                // TRUE represents a boolean true.
	FALSE                               // FALSE represents a boolean false.
	COMMA                               // COMMA represents a ',' separator.
	COLON                               // COLON represents a ':' separator.
	LEFT_SQUARE_BRACE                   // LEFT_SQUARE_BRACE represents '['.
	RIGHT_SQUARE_BRACE                  // RIGHT_SQUARE_BRACE represents ']'.
	LEFT_CURLY_BRACE                    // LEFT_CURLY_BRACE represents '{'.
	RIGHT_CURLY_BRACE                   // RIGHT_CURLY_BRACE represents '}'.
)

// String returns a string representation of the TokenKind.
func (t TokenKind) String() string {
	m := map[TokenKind]string{
		0:  "EOF",
		1:  "ILLEGAL",
		2:  "WHITESPACE",
		3:  "LINE_COMMENT",
		4:  "BLOCK_COMMENT",
		5:  "STRING",
		6:  "NUMBER",
		7:  "NULL",
		8:  "TRUE",
		9:  "FALSE",
		10: "COMMA",
		11: "COLON",
		12: "LEFT_SQUARE_BRACE",
		13: "RIGHT_SQUARE_BRACE",
		14: "LEFT_CURLY_BRACE",
		15: "RIGHT_CURLY_BRACE",
	}

	str := m[t]
	return str
}

// Token represents a single lexical token, including its kind, value, and position.
type Token struct {
	Kind    TokenKind    // The general kind of the token (e.g., STRING, NUMBER).
	SubKind TokenSubKind // The specific sub kind within a kind (e.g., INTEGER vs FLOAT).
	Literal []byte       // The literal value of the token.
	Line    int          // The line number where the token appears.
	Column  int          // The column number (character position) in the line.
}

// NewToken creates and returns a new Token.
// If cb is non-nil, it is called during token creation.
func NewToken(kind TokenKind, subKind TokenSubKind, literal []byte, line, start int, cb func()) Token {
	if cb != nil {
		cb()
	}

	return Token{
		Kind:    kind,
		SubKind: subKind,
		Literal: literal,
		Line:    line,
		Column:  start,
	}
}

// String returns a human-readable representation of the token.
func (t Token) String() string {
	return fmt.Sprintf(
		"Token{Kind: %s, SubKind: %s, Literal: %s, Line: %d, Column: %d}",
		t.Kind,
		t.SubKind,
		string(t.Literal),
		t.Line,
		t.Column,
	)
}

func (t Token) Value() any {
	switch t.Kind {
	case NULL:
		return nil
	case FALSE:
		return false
	case TRUE:
		return true
	case STRING:
		switch t.SubKind {
		case QUOTED:
			return quoteValue(t.Literal)
		case IDENT:
			return string(t.Literal)
		default:
			return nil
		}
	case NUMBER:
		switch t.SubKind {
		case INTEGER:
			return intValue(t.Literal)
		case FLOAT:
			return floatValue(t.Literal)
		default:
			return nil
		}
	default:
		return nil
	}
}

// Tokens is a slice of Token.
type Tokens []Token

// String returns a formatted string listing all tokens in the slice.
func (t Tokens) String() string {
	parts := make([]string, len(t))

	for i, token := range t {
		parts[i] = token.String()
	}

	return fmt.Sprintf("\n[\n%s\n]\n", strings.Join(parts, ",\n"))
}
