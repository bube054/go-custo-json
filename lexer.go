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
	pos := l.pos
	char := l.char
	nextChar := l.peek()
	// prev := l.prev()

	switch char {

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

		switch nextChar {
		case 47:
			if !l.config.AllowLineComments {
				return NewToken(ILLEGAL, l.input[l.pos:], l.line, pos, nil)
			}

			for !IsNewLine(l.char) && l.char != 0 {
				l.readChar()
			}

			if l.char == 0 {
				return NewToken(LINE_COMMENT, l.input[pos:], l.line, pos, nil)
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

	// Lexing string starts here
	case
		34, // double quotes
		39: // single quote
		if l.char == 39 && !l.config.AllowSingleQuotes {
			return NewToken(ILLEGAL, l.input[l.pos:], l.line, pos, nil)
		}

		l.readChar()

		for {
			next := l.peek()
			prev := l.prev()

			if l.char == 0 {
				break
			}

			if l.char == char && prev != 92 {
				break
			}

			if l.char == 92 && prev != l.char {
				switch next {
				case char:
				case 92: // backward slash
				case 47: // forward slash
				case 98: // b
				case 102: // f
				case 110: // n
				case 114: // r
				case 116: // t
				case 117: // u
					h1 := l.peekBy(2)
					h2 := l.peekBy(3)
					h3 := l.peekBy(4)
					h4 := l.peekBy(5)

					if !Is4HexDigits([4]byte{h1, h2, h3, h4}) {
						return NewToken(ILLEGAL, l.input[pos:], l.line, pos, nil)
					}

				default:
					return NewToken(ILLEGAL, l.input[pos:], l.line, pos, nil)
				}
			}

			l.readChar()
		}

		if l.char == 0 {
			return NewToken(ILLEGAL, l.input[pos:], l.line, pos, nil)
		}

		// fmt.Printf("%+v\n", l)
		return NewToken(STRING, l.input[pos:l.readPos], l.line, pos, l.readChar)
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

func (l *Lexer) peekBy(target int) byte {
	pos := l.pos + target

	if pos > (len(l.input) - 1) {
		return 0
	} else {
		return l.input[pos]
	}
}
