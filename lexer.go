package gocustojson

type Lexer struct {
	input  []byte
	config *Config

	line    int
	pos     int
	readPos int

	char byte
}

func NewLexer(input []byte, cfg *Config) *Lexer {
	l := &Lexer{input: input, config: cfg}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.line == 0 {
		l.line = 1
	}

	if l.readPos > len(l.input)-1 {
		l.char = 0
	} else {
		l.char = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) Token() Token {
	switch l.char {

	// Lexing white space starts here
	case
		32, // space
		10, // line feed
		13, // carriage return
		9:  // horizontal return
		return NewToken(WHITESPACE, l.input[l.pos:l.readPos], l.line, l.pos, l.readChar)
	case
		11,  // line tabulation
		12,  // form feed
		133, // next line
		160: // no break space
		if l.config.AllowExtraWS {
			return NewToken(WHITESPACE, l.input[l.pos:l.readPos], l.line, l.pos, l.readChar)
		} else {
			return NewToken(ILLEGAL, l.input[l.pos:l.readPos], l.line, l.pos, nil)
		}
	// Lexing white space ends here

	// Lexing comment starts here
	case 47: // forward slash
		pos := l.pos
		nextChar := l.peek()
		switch nextChar {
		case 47:
			if !l.config.AllowLineComments {
				return NewToken(ILLEGAL, l.input[l.pos:], l.line, pos, nil)
			}

			for !IsNewLine(l.char) && l.char != 0 {
				l.readChar()
			}

			if l.char == 0 {
				return NewToken(ILLEGAL, l.input[pos:], l.line, pos, nil)
			}

			return NewToken(LINE_COMMENT, l.input[pos:l.readPos], l.line, pos, l.readChar)

		case 42: // asterisk
			if !l.config.AllowBlockComments {
				return NewToken(ILLEGAL, l.input[l.pos:], l.line, pos, nil)
			}

			for !(l.char == 47 && l.prev() == 42) && l.char != 0 {
				l.readChar()
			}

			if l.char == 0 {
				return NewToken(ILLEGAL, l.input[pos:], l.line, pos, nil)
			}

			return NewToken(BLOCK_COMMENT, l.input[pos:l.readPos], l.line, pos, l.readChar)
		default:
			return NewToken(ILLEGAL, l.input[l.pos:l.readPos], l.line, l.pos, nil)
		}
	// Lexing comment ends here

	// Lexing string ends here
	// Lexing string ends here

	case 0:
		return NewToken(EOF, nil, l.line, l.pos, nil)
	default:
		return NewToken(ILLEGAL, l.input[l.pos:l.readPos], l.line, l.pos, nil)
	}
}

func (l *Lexer) GenerateTokens() Tokens {
	tokens := []Token{}

	for {
		token := l.Token()

		tokens = append(tokens, token)

		if token.Kind == EOF || token.Kind == ILLEGAL {
			break
		}
	}

	return tokens
}

func (l *Lexer) peek() byte {
	if l.readPos > (len(l.input) - 1) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func (l *Lexer) prev() byte {
	if (l.pos - 1) < 0 {
		return 0
	} else {
		return l.input[l.pos-1]
	}
}

// // Start current line at 1
// if l.line == 0 {
// 	l.line = 1
// }

// // Do not consume nothing
// if len(l.input) == 0 {
// 	return
// }

// // Nothing to consume again
// if l.currentPointer > len(l.input)-1 {
// 	l.currentCharacter = 0
// 	return
// }

// l.currentCharacter = l.input[l.currentPointer]

// if l.currentPointer+1 < len(l.input)-1 {
// 	l.nextCharacter = l.input[l.currentPointer+1]
// } else {
// 	l.nextCharacter = 0
// }

// if l.currentPointer > 0 {
// 	l.previousCharacter = l.input[l.currentPointer-1]
// } else {
// 	l.previousCharacter = 0
// }

// if IsNewLine(l.currentCharacter) {
// 	l.line++
// 	l.pos = 0
// } else {
// 	l.pos++
// }

// l.currentPointer++
