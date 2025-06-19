package gocustojson

import "fmt"

type TokenKind int

const (
	EOF TokenKind = iota
	ILLEGAL
)

func (t TokenKind) String() string {
	m := map[TokenKind]string{
		0: "EOF",
		1: "ILLEGAL",
	}

	str := m[t]
	return str
}

type Token struct {
	Kind        TokenKind
	Value       string
	Line        int
	StartColumn int
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"Token{Kind: %s, Value: %q, Line: %d, StartColumn: %d}",
		t.Kind,
		t.Value,
		t.Line,
		t.StartColumn,
	)
}

func NewToken(kind TokenKind, value []byte, line, start int) Token {
	return Token{
		Kind:        kind,
		Value:       string(value),
		Line:        line,
		StartColumn: start,
	}
}
