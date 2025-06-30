// Package gocustojson provides configurable options for parsing of JSON.
//
// It supports parsing a broad range of JSON syntax variants — from strict [ECMA-404-compliant JSON]
// to more permissive formats like [JSON5]. The parser behavior can be customized via the Config struct,
// which exposes fine-grained toggles for non-standard features such as comments, trailing commas,
// unquoted keys, single-quoted strings, and more.
//
// [ECMA-404-compliant JSON]: https://datatracker.ietf.org/doc/html/rfc7159
// [JSON5]: https://json5.org/
package gocustojson

import (
	"errors"
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func New(input []byte, cfg *Config) Parser {
	p := Parser{l: NewLexer(input, cfg)}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() (JSONNode, error) {
	var v JSONNode
	switch p.curToken.Kind {
	case NULL:
		return p.parseNull()
	case FALSE, TRUE:
		return p.parseBoolean()
	case STRING:
		return p.parseString()
	case NUMBER:
		return p.parseNumber()
	case LEFT_SQUARE_BRACE:
		return p.parseArray()
	case LEFT_CURLY_BRACE:
		return p.parseObject()
	case ILLEGAL:
		return v, errors.New("illegal character")
	default:
		return v, errors.New("unexpected input")
	}
}

func (p *Parser) parseNull() (JSONNode, error) {
	return JSONNull{
		token: p.curToken,
	}, nil
}

func (p *Parser) parseBoolean() (JSONNode, error) {
	return JSONBoolean{
		token: p.curToken,
	}, nil
}

func (p *Parser) parseString() (JSONNode, error) {
	return JSONString{
		token: p.curToken,
	}, nil
}

func (p *Parser) parseNumber() (JSONNode, error) {
	return JSONNumber{
		token: p.curToken,
	}, nil
}

func (p *Parser) parseArray() (JSONNode, error) {
	array := JSONArray{items: []JSONNode{}}

	p.nextToken()

	for !p.expectCurToken(RIGHT_SQUARE_BRACE) {
		if p.expectCurToken(EOF) {
			return array, errors.New("array syntax error")
		}

		if p.expectCurToken(ILLEGAL) {
			return array, errors.New("syntax error")
		}

		if p.expectCurToken(COMMA) {
			p.nextToken()
			continue
		}

		// lastItem := p.expectPeekToken(RIGHT_SQUARE_BRACE)
		// nextIsComma := p.expectPeekToken(COMMA)

		// if !nextIsComma && !lastItem {
		// 	return array, errors.New("array comma syntax error")
		// }

		item, err := p.Parse()

		if err != nil {
			return array, err
		}

		array.items = append(array.items, item)

		p.nextToken()
	}

	for p.expectCurToken(RIGHT_CURLY_BRACE) {
		p.nextToken()
	}

	return array, nil
}

func (p *Parser) parseObject() (JSONNode, error) {
	object := JSONObject{properties: map[string]JSONNode{}}

	p.nextToken()

	for !p.expectCurToken(RIGHT_CURLY_BRACE) {
		if p.expectCurToken(EOF) {
			return object, errors.New("object syntax error")
		}

		if p.expectCurToken(ILLEGAL) {
			return object, errors.New("syntax error")
		}

		key, err := p.Parse()

		if err != nil {
			return object, errors.New("illegal object key")
		}

		stringNode, ok := key.(JSONString)
		// _, ok := key.(JSONString)

		if !ok {
			return object, errors.New("only object key string required")
		}

		p.nextToken()

		if !p.expectCurToken(COLON) {
			return object, errors.New("colon object required")
		}

		p.nextToken()

		value, err := p.Parse()

		if err != nil {
			return object, errors.New("illegal object value")
		}

		object.properties[stringNode.Literal()] = value

		p.nextToken()

		if !p.expectCurToken(COMMA) {
			return object, errors.New("object key pair should end with comma")
		}

		p.nextToken()
	}

	for p.expectCurToken(RIGHT_CURLY_BRACE) {
		p.nextToken()
	}

	return object, nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextUsefulToken()
}

func (p *Parser) expectCurToken(kind TokenKind) bool {
	if p.curToken.Kind == kind {
		return true
	} else {
		return false
	}
}

func (p *Parser) expectPeekToken(kind TokenKind) bool {
	if p.peekToken.Kind == kind {
		return true
	} else {
		return false
	}
}

func (p *Parser) curTokenIsOneOf(kinds ...TokenKind) bool {
	for _, kind := range kinds {
		if p.expectCurToken(kind) {
			return true
		}
	}

	return false
}
