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

	INTEGER     //  positive or negative
	FLOAT       //  positive or negative
	SCI_NOT_INT //  positive or negative sign/expo
	SCI_NOT_FLT //  positive or negative sign/expo
	HEX         //  positive or negative
	INF         //  positive or negative
	NaN         //  positive or negative
)

func (t TokenKindType) String() string {
	m := map[TokenKindType]string{
		0: "STRING_VALUE",
		1: "IDENT",
		2: "INTEGER",
		3: "FLOAT",
		4: "SCI_NOT_INT",
		5: "SCI_NOT_FLT",
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
)

func (t TokenKind) String() string {
	m := map[TokenKind]string{
		0: "EOF",
		1: "ILLEGAL",
		2: "WHITESPACE",
		3: "LINE_COMMENT",
		4: "BLOCK_COMMENT",
		5: "STRING",
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
