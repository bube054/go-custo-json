package gocustojson

type Lexer struct {
	input  []byte
	config *Config

	line    int
	pos     int
	readPos int

	// currentPointer    int
	char byte
	// nextCharacter     byte
	// previousCharacter byte
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
	case 0:
		// fmt.Println("3")
		return NewToken(EOF, nil, l.line, l.pos, nil)
	default:
		// fmt.Println("4")
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
