package gocustojson

import "fmt"

type TokenKind int

const (
	WhiteSpace TokenKind = iota
)

type TokenType struct {
	TokenKind  TokenKind
	TokenValue []byte
	Line       int
	Column     int
}

func (t *TokenType) String() string {
	return fmt.Sprintf("{Kind: %d, Value: %q}", t.TokenKind, t.TokenValue)
}
