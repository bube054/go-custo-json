package gocustojson

import (
	"fmt"
	"strings"
)

type TokenKindType int

const (
	NONE TokenKindType = iota

	STRING_VALUE
	IDENT

	INTEGER //  positive or negative
	FLOAT   //  positive or negative
	SCI_NOT //  positive or negative sign/expo
	HEX     //  positive or negative
	INF     //  positive or negative
	NaN     //  positive or negative
)

func (t TokenKindType) String() string {
	m := map[TokenKindType]string{
		0: "NONE",
		1: "STRING_VALUE",
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

type TokenKind int

const (
	EOF TokenKind = iota
	ILLEGAL
	WHITESPACE
	LINE_COMMENT
	BLOCK_COMMENT
	STRING
	NUMBER
	NULL
	TRUE
	FALSE
	COMMA
	LEFT_SQUARE_BRACE
	RIGHT_SQUARE_BRACE
	LEFT_CURLY_BRACE
	RIGHT_CURLY_BRACE
)

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
		11: "LEFT_SQUARE_BRACE",
		12: "RIGHT_SQUARE_BRACE",
		13: "LEFT_CURLY_BRACE",
		14: "RIGHT_CURLY_BRACE",
	}

	str := m[t]
	return str
}

type Token struct {
	Kind     TokenKind
	KindType TokenKindType
	Value    string
	Line     int
	Column   int
}

func (t Token) String() string {
	return fmt.Sprintf(
		"Token{Kind: %s, KindType: %s, Value: %s, Line: %d, Column: %d}",
		t.Kind,
		t.KindType,
		t.Value,
		t.Line,
		t.Column,
	)
}

func NewToken(kind TokenKind, kindType TokenKindType, value []byte, line, start int, cb func()) Token {
	if cb != nil {
		cb()
	}

	return Token{
		Kind:     kind,
		KindType: kindType,
		Value:    string(value),
		Line:     line,
		Column:   start,
	}
}

type Tokens []Token

func (t Tokens) String() string {
	parts := make([]string, len(t))

	for i, token := range t {
		parts[i] = token.String()
	}

	return fmt.Sprintf("\n[\n%s\n]\n", strings.Join(parts, ",\n"))
}
