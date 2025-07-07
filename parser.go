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
	// l         *Lexer
	// curToken  Token
	// peekToken Token
	// prevToken Token
	tokens Tokens
	config *Config

	curToken  Token
	curPos    int
	peekToken Token
	peekPos   int
	prevToken Token
	prevPos   int
}

func New(input []byte, config *Config) Parser {
	if config == nil {
		config = NewConfig()
	}
	// fmt.Println("\033[1mThis is bold text\033[0m")
	l := NewLexer(input, config)
	tokens := l.TokensWithout(WHITESPACE)
	// fmt.Println(tokens)
	p := Parser{
		tokens:  tokens,
		config:  config,
		curPos:  -1,
		peekPos: -1,
		prevPos: -1,
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() (JSONNode, error) {
	switch p.curToken.Kind {
	case NULL:
		return p.parseNull()
	// case FALSE, TRUE:
	// 	return p.parseBoolean()
	// case STRING:
	// 	return p.parseString()
	// case NUMBER:
	// 	return p.parseNumber()
	// case LEFT_SQUARE_BRACE:
	// 	return p.parseArray()
	// case LEFT_CURLY_BRACE:
	// 	return p.parseObject()
	// case ILLEGAL:
	// 	return p.parseIllegal()
	case COMMENT:
		return p.handleComment()
	default:
		return nil, nil
	}
}

// func (p *Parser) parseDefault() (JSONNode, error) {
// 	// fmt.Println("prev Token", p.prevToken)
// 	// fmt.Println("current Token", p.curToken)
// 	// fmt.Println("peek Token", p.peekToken)

// 	// noContent := p.expectPrevToken(EOF) && p.expectCurToken(EOF) && p.expectPeekToken(EOF)

// 	if p.expectCurToken(COMMENT) && p.expectPeekToken(EOF) {
// 		return nil, ErrJSONNoContent
// 	}

// 	// if !p.expectPrevToken(COMMENT) && p.expectCurToken(COMMENT) && !p.expectPeekToken(COMMENT) {
// 	// 	fmt.Println("prev Token", p.prevToken)
// 	// 	fmt.Println("current Token", p.curToken)
// 	// 	fmt.Println("peek Token", p.peekToken)
// 	// 	return nil, fmt.Errorf("%w: %q at line %d, column %d",
// 	// 		ErrJSONUnexpectedChar,
// 	// 		p.peekToken.Literal,
// 	// 		p.peekToken.Line,
// 	// 		p.peekToken.Column,
// 	// 	)
// 	// }

// 	// if p.expectCurToken(COMMENT) {
// 	// if p.expectCurToken(COMMENT) && p.expectPeekToken(COMMENT) {
// 	// 	p.nextToken()
// 	// 	return p.Parse()
// 	// }

// 	if p.expectCurToken(COMMENT) {
// 		p.nextToken()
// 		return p.Parse()
// 	}

// 	return nil, fmt.Errorf("%w: %q at line %d, column %d",
// 		ErrJSONUnexpectedChar,
// 		p.curToken.Literal,
// 		p.curToken.Line,
// 		p.curToken.Column,
// 	)
// }

// func (p *Parser) parseIllegal() (JSONNode, error) {
// 	return nil, fmt.Errorf("%w: %q at line %d, column %d",
// 		ErrJSONSyntax,
// 		p.curToken.Literal,
// 		p.curToken.Line,
// 		p.curToken.Column+1,
// 	)
// }

func (p *Parser) handleComment() (JSONNode, error) {
	// return newJSONNull(p.curToken, p.nextToken), nil
}

func (p *Parser) parseNull() (JSONNode, error) {
	return newJSONNull(p.curToken, p.nextToken), nil
}

// func (p *Parser) parseBoolean() (JSONNode, error) {
// 	if err := p.ensureSingleValidPrimitive(); err != nil {
// 		return nil, err
// 	}
// 	return newJSONBoolean(p.curToken, p.nextToken), nil
// }

// func (p *Parser) parseString() (JSONNode, error) {
// 	if err := p.ensureSingleValidPrimitive(); err != nil {
// 		return nil, err
// 	}

// 	if p.curToken.SubKind == IDENT && p.expectPrevToken(EOF) && p.expectPeekToken(EOF) {
// 		return nil, fmt.Errorf("%w: %q at line %d, column %d",
// 			ErrJSONUnexpectedChar,
// 			p.curToken.Literal,
// 			p.curToken.Line,
// 			p.curToken.Column+1,
// 		)
// 	}

// 	return newJSONString(p.curToken, p.nextToken), nil
// }

// func (p *Parser) parseNumber() (JSONNode, error) {
// 	if err := p.ensureSingleValidPrimitive(); err != nil {
// 		return nil, err
// 	}

// 	return newJSONNumber(p.curToken, p.nextToken), nil
// }

// func (p *Parser) parseArray() (JSONNode, error) {
// 	items := []JSONNode{}
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

// func (p *Parser) parseObject() (JSONNode, error) {
// 	object := JSONObject{Properties: map[string]JSONNode{}}

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

// func (p *Parser) ensureSingleValidPrimitive() error {
// 	prevIsEOFOrComment := p.expectPrevToken(EOF) || p.expectPrevToken(COMMENT)
// 	peekIsNotEOFOrComment := !(p.expectPeekToken(EOF) || p.expectPeekToken(COMMENT))

// 	shouldReportUnexpectedValue := prevIsEOFOrComment && peekIsNotEOFOrComment

// 	if shouldReportUnexpectedValue {
// 		return fmt.Errorf("%w: %q at line %d, column %d",
// 			ErrJSONUnexpectedChar,
// 			p.peekToken.Literal,
// 			p.peekToken.Line,
// 			p.peekToken.Column,
// 		)
// 	}
// 	return nil
// }
