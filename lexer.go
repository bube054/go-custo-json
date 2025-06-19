package gocustojson

type Lexer struct {
	input  []byte
	config *Config

	currentLine    int
	columnPosition int

	currentPointer    int
	currentCharacter  byte
	nextCharacter     byte
	previousCharacter byte
}

func NewLexer(input string, cfg *Config) *Lexer {
	l := &Lexer{input: []byte(input), config: cfg}
	l.ConsumeCharacter()
	return l
}

func (l *Lexer) ConsumeCharacter() {
	if l.currentLine == 0 {
		l.currentLine = 1
	}

	if len(l.input) == 0 {
		return
	}

	l.currentCharacter = l.input[l.currentPointer]
	l.nextCharacter = l.input[l.currentPointer+1]

	if l.currentPointer > 0 {
		l.previousCharacter = l.input[l.currentPointer-1]
	}

	if IsNewLine(l.currentCharacter) {
		l.currentLine++
		l.columnPosition = 0
	} else {
		l.columnPosition++
	}

	l.currentPointer++
}

func (l *Lexer) Token() Token {
	switch l.currentCharacter {
	case '0', '\x00':
		return NewToken(EOF, nil, l.currentLine, l.columnPosition)
	default:
		return NewToken(ILLEGAL, l.input[l.currentPointer:], l.currentLine, l.columnPosition)
	}
}

func (l *Lexer) GenerateTokens() []Token {
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
