// Package jsonvx provides configurable options for parsing of JSON.
//
// It supports parsing a broad range of JSON syntax variants — from strict [ECMA-404-compliant JSON]
// to more permissive formats like [JSON5]. The parser behavior can be customized via the Config struct,
// which exposes fine-grained toggles for non-standard features such as comments, trailing commas,
// unquoted keys, single-quoted strings, and more.
//
// [ECMA-404-compliant JSON]: https://datatracker.ietf.org/doc/html/rfc7159
// [JSON5]: https://json5.org/
package jsonvx

import (
	"errors"
	"fmt"
)

// SyntaxError: JSON.parse: unterminated string literal
// SyntaxError: JSON.parse: bad control character in string literal
// SyntaxError: JSON.parse: bad character in string literal
// SyntaxError: JSON.parse: bad Unicode escape
// SyntaxError: JSON.parse: bad escape character
// SyntaxError: JSON.parse: unterminated string
// SyntaxError: JSON.parse: no number after minus sign
// SyntaxError: JSON.parse: unexpected non-digit
// SyntaxError: JSON.parse: missing digits after decimal point
// SyntaxError: JSON.parse: unterminated fractional number
// SyntaxError: JSON.parse: missing digits after exponent indicator
// SyntaxError: JSON.parse: missing digits after exponent sign
// SyntaxError: JSON.parse: exponent part is missing a number
// SyntaxError: JSON.parse: unexpected end of data
// SyntaxError: JSON.parse: unexpected keyword
// SyntaxError: JSON.parse: unexpected character
// SyntaxError: JSON.parse: end of data while reading object contents
// SyntaxError: JSON.parse: expected property name or '}'
// SyntaxError: JSON.parse: end of data when ',' or ']' was expected
// SyntaxError: JSON.parse: expected ',' or ']' after array element
// SyntaxError: JSON.parse: end of data when property name was expected
// SyntaxError: JSON.parse: expected double-quoted property name
// SyntaxError: JSON.parse: end of data after property name when ':' was expected
// SyntaxError: JSON.parse: expected ':' after property name in object
// SyntaxError: JSON.parse: end of data after property value in object
// SyntaxError: JSON.parse: expected ',' or '}' after property value in object
// SyntaxError: JSON.parse: expected ',' or '}' after property-value pair in object literal
// SyntaxError: JSON.parse: property names must be double-quoted strings
// SyntaxError: JSON.parse: expected property name or '}'
// SyntaxError: JSON.parse: unexpected character
// SyntaxError: JSON.parse: unexpected non-whitespace character after JSON data

var (
	ErrJSONSyntax         = errors.New("JSON syntax error")
	ErrJSONUnexpectedChar = errors.New("unexpected character in JSON input")
	ErrJSONNoContent      = errors.New("no meaningful content to parse")
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	prevToken Token
}

func New(input []byte, cfg *Config) Parser {
	p := Parser{l: NewLexer(input, cfg)}
	// fmt.Println(p.l.Tokens())
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() (JSONNode, error) {
	switch p.curToken.Kind {
	case NULL:
		return p.parseNull()
	case FALSE, TRUE:
		return p.parseBoolean()
	case STRING:
		return p.parseString()
	// case NUMBER:
	// 	return p.parseNumber()
	// case LEFT_SQUARE_BRACE:
	// 	return p.parseArray()
	// case LEFT_CURLY_BRACE:
	// 	return p.parseObject()
	case ILLEGAL:
		return p.parseIllegal()
	default:
		return p.parseDefault()
	}
}

func (p *Parser) parseDefault() (JSONNode, error) {
	// fmt.Println("prev token", p.prevToken)
	// fmt.Println("current token", p.curToken)
	// fmt.Println("peek token", p.peekToken)

	if p.expectPrevToken(EOF) && p.expectCurToken(EOF) && p.expectPeekToken(EOF) {
		return nil, ErrJSONNoContent
	}

	return nil, fmt.Errorf("%w: %q at line %d, column %d",
		ErrJSONUnexpectedChar,
		p.curToken.Literal,
		p.curToken.Line,
		p.curToken.Column+1,
	)
}

func (p *Parser) parseIllegal() (JSONNode, error) {
	return nil, fmt.Errorf("%w: %q at line %d, column %d",
		ErrJSONSyntax,
		p.curToken.Literal,
		p.curToken.Line,
		p.curToken.Column+1,
	)
}

func (p *Parser) ensureValidPrimitive() error {
	if p.expectPrevToken(EOF) && !p.expectPeekToken(EOF) {
		return fmt.Errorf("%w: %q at line %d, column %d",
			ErrJSONUnexpectedChar,
			p.peekToken.Literal,
			p.peekToken.Line,
			p.peekToken.Column+1,
		)
	}
	return nil
}

func (p *Parser) parseNull() (JSONNode, error) {
	if err := p.ensureValidPrimitive(); err != nil {
		return nil, err
	}
	return JSONNull{token: p.curToken}, nil
}

func (p *Parser) parseBoolean() (JSONNode, error) {
	if err := p.ensureValidPrimitive(); err != nil {
		return nil, err
	}
	return JSONBoolean{token: p.curToken}, nil
}

func (p *Parser) parseString() (JSONNode, error) {
	if err := p.ensureValidPrimitive(); err != nil {
		return nil, err
	}
	return JSONString{token: p.curToken}, nil
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
	p.prevToken = p.curToken
	// fmt.Println("prevToken", p.prevToken)
	p.curToken = p.peekToken
	// fmt.Println("curToken", p.curToken)
	p.peekToken = p.l.NextUsefulToken()
	// fmt.Println("peekToken", p.peekToken)
	// fmt.Println("///////////////////")
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

func (p *Parser) expectPrevToken(kind TokenKind) bool {
	if p.prevToken.Kind == kind {
		return true
	} else {
		return false
	}
}

// func (p *Parser) curTokenIsOneOf(kinds ...TokenKind) bool {
// 	for _, kind := range kinds {
// 		if p.expectCurToken(kind) {
// 			return true
// 		}
// 	}

// 	return false
// }
