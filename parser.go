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
// SyntaxError: JSON.parse: unterminated fractional number/
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
	ErrJSONMultipleContent = errors.New("multiple JSON values")
)

type Parser struct {
	tokens Tokens
	config *Config

	curToken  Token
	curPos    int
	peekToken Token
	peekPos   int
	prevToken Token
	prevPos   int
}

func New(config *Config) Parser {
	if config == nil {
		config = NewConfig()
	}

	p := Parser{
		config: config,
	}

	return p
}

func (p *Parser) Parse(input []byte) (JSON, error) {
	l := NewLexer(input, p.config)

	tokens := l.Tokens()
	chunks, err := tokens.Split()
	lastIndex := len(chunks) - 1

	if err != nil {

		if lastIndex > 0 && errors.Is(err, ErrIllegalToken) {
			chunk := chunks[lastIndex]
			illegalToken := tokens[chunk[0]]
			return nil, WrapJSONUnexpectedCharError(illegalToken)
		}

		if lastIndex > 0 && errors.Is(err, ErrUnexpectedToken) {
			chunk := chunks[lastIndex]
			illegalToken := tokens[chunk[0]]
			return nil, WrapJSONUnexpectedCharError(illegalToken)
		}

		return nil, err
	}

	if len(chunks) == 0 {
		return nil, ErrJSONNoContent
	}

	if lastIndex > 0 && len(chunks) != 1 {
		chunk := chunks[lastIndex]
		extraToken := tokens[chunk[0]]
		return nil, WrapJSONMultipleContentError(extraToken)
	}

	chunk := chunks[0]
	p.tokens = tokens[chunk[0] : chunk[1]+1]
	p.curPos = -1
	p.curToken = Token{Kind: EOF}
	p.peekPos = -1
	p.peekToken = Token{Kind: EOF}
	p.prevPos = -1
	p.prevToken = Token{Kind: EOF}

	p.nextToken()
	p.nextToken()

	return p.parse()
}

func (p *Parser) parse() (JSON, error) {
	switch p.curToken.Kind {
	case NULL:
		return p.parseNull()
	case FALSE, TRUE:
		return p.parseBoolean()
	case STRING:
		return p.parseString()
	case NUMBER:
		return p.parseNumber()
	// case LEFT_SQUARE_BRACE:
		// return p.parseArray()
	// case LEFT_CURLY_BRACE:
	// 	return p.parseObject()
	case ILLEGAL:
		return p.parseIllegal()
	default:
		return p.parseDefault()
	}
}

func (p *Parser) parseDefault() (JSON, error) {
	// fmt.Println("prev Token", p.prevToken)
	// fmt.Println("current Token", p.curToken)
	// fmt.Println("peek Token", p.peekToken)
	return nil, WrapJSONUnexpectedCharError(p.curToken)
}

func (p *Parser) parseIllegal() (JSON, error) {
	return nil, WrapJSONUnexpectedCharError(p.curToken)
}

// func (p *Parser) handleComment() (JSON, error) {
// 	return newJSONNull(p.curToken, p.nextToken), nil
// }

func (p *Parser) parseNull() (JSON, error) {
	return newJSONNull(p.curToken, p.nextToken), nil
}

func (p *Parser) parseBoolean() (JSON, error) {
	return newJSONBoolean(p.curToken, p.nextToken), nil
}

func (p *Parser) parseString() (JSON, error) {
	return newJSONString(p.curToken, p.nextToken), nil
}

func (p *Parser) parseNumber() (JSON, error) {
	return newJSONNumber(p.curToken, p.nextToken), nil
}

// func (p *Parser) parseArray() (JSON, error) {
// 	items := []JSON{}
// 	p.nextToken()

// 	for !p.expectCurToken(RIGHT_SQUARE_BRACE) {
// 		item, err := p.Parse()
// 		if err != nil {
// 			return nil, err
// 		}

// 		items = append(items, item)

// 		hasComma := p.expectCurToken(COMMA)
// 		isClosingBracket := p.expectCurToken(RIGHT_SQUARE_BRACE)
// 		isNextClosingBracket := p.expectPeekToken(RIGHT_SQUARE_BRACE)

// 		isTrailingComma := hasComma && isNextClosingBracket
// 		isValidArrayEnd := isClosingBracket || isTrailingComma

// 		if !p.l.config.AllowTrailingCommaArray && isTrailingComma {
// 			return nil, fmt.Errorf("%w: %q at line %d, column %d",
// 				ErrJSONSyntax,
// 				p.curToken.Literal,
// 				p.curToken.Line,
// 				p.curToken.Column+1,
// 			)
// 		}

// 		if !isValidArrayEnd && !hasComma {
// 			return nil, fmt.Errorf("%w: %q at line %d, column %d",
// 				ErrJSONSyntax,
// 				p.curToken.Literal,
// 				p.curToken.Line,
// 				p.curToken.Column+1,
// 			)
// 		}

// 		if isValidArrayEnd {
// 			break
// 		}

// 		p.nextToken()
// 	}

// 	return newJSONArray(items, p.nextToken), nil
// }

// func (p *Parser) parseObject() (JSON, error) {
// 	object := JSONObject{Properties: map[string]JSON{}}

// 	p.nextToken()

// 	for !p.expectCurToken(RIGHT_CURLY_BRACE) {
// 		if p.expectCurToken(EOF) {
// 			return object, errors.New("object syntax error")
// 		}

// 		if p.expectCurToken(ILLEGAL) {
// 			return object, errors.New("syntax error")
// 		}

// 		key, err := p.Parse()

// 		if err != nil {
// 			return object, errors.New("illegal object key")
// 		}

// 		stringNode, ok := key.(JSONString)
// 		// _, ok := key.(JSONString)

// 		if !ok {
// 			return object, errors.New("only object key string required")
// 		}

// 		p.nextToken()

// 		if !p.expectCurToken(COLON) {
// 			return object, errors.New("colon object required")
// 		}

// 		p.nextToken()

// 		value, err := p.Parse()

// 		if err != nil {
// 			return object, errors.New("illegal object value")
// 		}

// 		object.Properties[stringNode.Literal()] = value

// 		p.nextToken()

// 		if !p.expectCurToken(COMMA) {
// 			return object, errors.New("object key pair should end with comma")
// 		}

// 		p.nextToken()
// 	}

// 	for p.expectCurToken(RIGHT_CURLY_BRACE) {
// 		p.nextToken()
// 	}

// 	return object, nil
// }

func (p *Parser) nextToken() {
	p.prevToken = p.curToken
	p.prevPos = p.curPos
	// fmt.Println("prevToken", p.prevToken)

	p.curToken = p.peekToken
	p.curPos = p.peekPos
	// fmt.Println("curToken", p.curToken)

	if p.peekPos+1 >= len(p.tokens) {
		p.peekToken = Token{Kind: EOF, SubKind: NONE}
		p.peekPos = len(p.tokens)
	} else {
		p.peekPos++
		p.peekToken = p.tokens[p.peekPos]
		// fmt.Println("peekToken", p.peekToken)
	}
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

func WrapUnexpectedCharError(baseErr error, token Token) error {
	return fmt.Errorf("%w: %q at line %d, column %d", baseErr, token.Literal, token.Line, token.Column)
}
func WrapJSONUnexpectedCharError(token Token) error {
	return WrapUnexpectedCharError(ErrJSONUnexpectedChar, token)
}

func WrapJSONSyntaxError(token Token) error {
	return WrapUnexpectedCharError(ErrJSONSyntax, token)
}

func WrapJSONMultipleContentError(token Token) error {
	return fmt.Errorf("%w: extra value %q at line %d, column %d",
		ErrJSONMultipleContent,
		token.Literal,
		token.Line,
		token.Column,
	)
}
