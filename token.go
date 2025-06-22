package gocustojson

import (
	"fmt"
	"strings"
)

type TokenKind int

const (
	EOF TokenKind = iota
	ILLEGAL
	WHITESPACE
	LINE_COMMENT
	BLOCK_COMMENT
	STRING
	IDENT
)

func (t TokenKind) String() string {
	m := map[TokenKind]string{
		0: "EOF",
		1: "ILLEGAL",
		2: "WHITESPACE",
		3: "LINE_COMMENT",
		4: "BLOCK_COMMENT",
		5: "STRING",
		6: "IDENT",
	}

	str := m[t]
	return str
}

type Token struct {
	Kind   TokenKind
	Value  string
	Line   int
	Column int
}

func (t Token) String() string {
	return fmt.Sprintf(
		"Token{Kind: %s, Value: %s, Line: %d, Column: %d}",
		t.Kind,
		t.Value,
		t.Line,
		t.Column,
	)
}

func NewToken(kind TokenKind, value []byte, line, start int, cb func()) Token {
	if cb != nil {
		cb()
	}

	return Token{
		Kind:   kind,
		Value:  string(value),
		Line:   line,
		Column: start,
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
