package gocustojson

type Lexer struct {
	characters []byte

	currentLine    int
	columnPosition int

	currentPointer    int
	currentCharacter  byte
	nextCharacter     byte
	previousCharacter byte
}

func (l *Lexer) ConsumeCharacter() {
	l.currentCharacter = l.characters[l.currentPointer]
	l.nextCharacter = l.characters[l.currentPointer+1]

	if l.currentPointer > 0 {
		l.previousCharacter = l.characters[l.currentPointer-1]
	}

	if IsNewLine(l.currentCharacter) {
		l.currentLine++
		l.columnPosition = 0
	} else {
		l.columnPosition++
	}

	l.currentPointer++
}

func NewLexer() *Lexer {
	return &Lexer{currentLine: 1}
}
