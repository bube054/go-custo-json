package gocustojson

import (
	"bytes"
	"fmt"
)

// Lexer performs lexical analysis on a JSON-like input stream.
// It breaks the raw input into a sequence of tokens according to the parsing rules
// defined in the provided Config.
//
// The Lexer operates on raw byte slices for performance, and tracks position and line
// information to support error reporting and debugging.
type Lexer struct {
	input  []byte  // input is the raw byte slice being tokenized.
	config *Config // config holds the parsing options that control how the lexer interprets input.

	line    int // line tracks the current line number in the input (starting from 1).
	pos     int // pos is the current character offset in the input (starting from 0).
	readPos int // char is the current character under examination.

	char byte // readPos is the next position to be read.
}

// NewLexer creates a new Lexer instance using the given input and configuration.
//
// It immediately reads the first character to initialize internal state, so the lexer is ready
// for tokenization right after creation.
func NewLexer(input []byte, cfg *Config) *Lexer {
	l := &Lexer{input: input, config: cfg}
	l.readChar()
	return l
}

// String returns a human-readable representation of the current state of the Lexer.
// It includes the input, current line number, character position, read position, and the current character.
func (l *Lexer) String() string {
	return fmt.Sprintf(
		"Lexer{\n  input: %q,\n  line: %d,\n  pos: %d,\n  readPos: %d,\n  char: %q\n}",
		string(l.input), l.line, l.pos, l.readPos, l.char,
	)
}

// Token returns the current token being parsed by the lexer.
func (l *Lexer) Token() Token {
	pos := l.pos
	char := l.char
	nextChar := l.peek()

	switch char {

	// Lexing white space starts here
	case
		' ',  // space (U+0020)
		'\n', // line feed (U+000A)
		'\r', // carriage return (U+000D)
		'\t': // horizontal tab (U+0009)
		return NewToken(WHITESPACE, NONE, l.input[l.pos:l.readPos], l.line, l.pos, l.readChar)

	case
		'\v',     // vertical tab (U+000B)
		'\f',     // form feed (U+000C)
		'\u0085', // next line (NEL, U+0085)
		'\u00A0': // no-break space (U+00A0)
		if l.config.AllowExtraWS {
			return NewToken(WHITESPACE, NONE, l.input[l.pos:l.readPos], l.line, l.pos, l.readChar)
		} else {
			return NewToken(ILLEGAL, NONE, l.input[l.pos:], l.line, l.pos, nil)
		}
		// Lexing white space ends here

	// Lexing left curly brace starts here
	case '{':
		return NewToken(LEFT_CURLY_BRACE, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
	// Lexing left curly brace ends here

	// Lexing right brace starts here
	case '}':
		return NewToken(RIGHT_CURLY_BRACE, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
	// Lexing right brace ends here

	// Lexing left square brace starts here
	case '[':
		return NewToken(LEFT_SQUARE_BRACE, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
	// Lexing left square brace ends here

	// Lexing right square brace starts here
	case ']':
		return NewToken(RIGHT_SQUARE_BRACE, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
		// Lexing right square brace ends here

	// Lexing null starts here
	case 'n':
		second := l.peekBy(1) // expect 'u'
		third := l.peekBy(2)  // expect 'l'
		fourth := l.peekBy(3) // expect 'l'

		if second == 'u' && third == 'l' && fourth == 'l' {
			l.readChar()
			l.readChar()
			l.readChar()

			return NewToken(NULL, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
		}

		return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
		// Lexing null ends here

	// Lexing true starts here
	case 't':
		second := l.peekBy(1) // expect 'r'
		third := l.peekBy(2)  // expect 'u'
		fourth := l.peekBy(3) // expect 'e'

		if second == 'r' && third == 'u' && fourth == 'e' {
			l.readChar()
			l.readChar()
			l.readChar()

			return NewToken(TRUE, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
		}

		return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
		// Lexing true ends here

	// Lexing false starts here
	case 'f':
		second := l.peekBy(1) // expect 'a'
		third := l.peekBy(2)  // expect 'l'
		fourth := l.peekBy(3) // expect 's'
		fifth := l.peekBy(4)  // expect 'e'

		if second == 'a' && third == 'l' && fourth == 's' && fifth == 'e' {
			l.readChar()
			l.readChar()
			l.readChar()
			l.readChar()

			return NewToken(FALSE, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
		}

		return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
		// Lexing false ends here

	// Lexing comment starts here
	case '/': // forward slash

		switch nextChar {
		case '/':
			if !l.config.AllowLineComments {
				return NewToken(ILLEGAL, NONE, l.input[l.pos:], l.line, pos, nil)
			}

			for !isNewLine(l.char) && l.char != 0 {
				l.readChar()
			}

			if l.char == 0 {
				return NewToken(LINE_COMMENT, NONE, l.input[pos:], l.line, pos, nil)
			}

			return NewToken(LINE_COMMENT, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)

		case '*': // asterisk
			if !l.config.AllowBlockComments {
				return NewToken(ILLEGAL, NONE, l.input[l.pos:], l.line, pos, nil)
			}

			for !(l.char == 47 && l.prev() == 42) && l.char != 0 {
				l.readChar()
			}

			if l.char == 0 {
				return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
			}

			return NewToken(BLOCK_COMMENT, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
		default:
			return NewToken(ILLEGAL, NONE, l.input[l.pos:l.readPos], l.line, l.pos, nil)
		}
		// Lexing comment ends here

	// Lexing string starts here
	case '"', '\'': // double or single quote
		if l.char == '\'' && !l.config.AllowSingleQuotes {
			return NewToken(ILLEGAL, NONE, l.input[l.pos:], l.line, pos, nil)
		}

		l.readChar()

		for {
			next := l.peek()
			prev := l.prev()
			prevBy2 := l.prevBy(2)

			if l.char == 0 {
				break
			}

			// End string if not escaped, e.g., "\"" or '\''
			if l.char == char && prev != '\\' {
				break
			}

			// End string if properly escaped, e.g., "\\\"" or '\\\''
			if l.char == char && prev == '\\' && prevBy2 == '\\' {
				break
			}

			// Handle escape sequences
			if l.char == '\\' && prev != '\\' {
				switch next {
				case char: // escaped quote
				case '\\': // escaped backslash
				case '/': // escaped forward slash
				case 'b': // backspace
				case 'f': // form feed
				case 'n': // newline
				case 'r': // carriage return
				case 't': // tab
				case 'u': // unicode escape
					second := l.peekBy(2)
					third := l.peekBy(3)
					fourth := l.peekBy(4)
					fifth := l.peekBy(5)

					if !is4HexDigits([4]byte{second, third, fourth, fifth}) {
						return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
					}

				case '\n': // escaped newline
					if !l.config.AllowNewlineInStrings {
						return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
					}

				default:
					if !l.config.AllowOtherEscapeChars {
						return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
					}
				}
			}

			l.readChar()
		}

		if l.char == 0 {
			return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
		}

		return NewToken(STRING, NONE, l.input[pos:l.readPos], l.line, pos, l.readChar)
		// Lexing string ends here

	case 0:
		return NewToken(EOF, NONE, nil, l.line, l.pos, nil)
	default:

		// Lexing number starts here
		if IsPossibleNumber(l.char, l.peek()) {
			l.readChar()

			for IsPossibleNumber(l.char, 0) {
				l.readChar()
			}

			num := l.input[pos:l.pos]

			parts := bytes.Split(num, []byte("."))
			part1 := parts[0]

			hasLeading0 := len(part1) > 1 && part1[0] == '0'
			hasFollowingX := len(part1) > 1 && (part1[1] == 'x' || part1[1] == 'X')

			if hasLeading0 && !hasFollowingX {
				return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
			}

			if startsWithPlus(num) && !l.config.AllowLeadingPlus {
				return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
			}

			isNaNum := isNaN(num)

			if isNaNum && !l.config.AllowNaN {
				return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
			}

			if isNaNum {
				return NewToken(NUMBER, NaN, l.input[pos:], l.line, pos, nil)
			}

			isInfinity := isInf(num)

			if isInfinity && !l.config.AllowInfinity {
				return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
			}

			if isInfinity {
				return NewToken(NUMBER, INF, l.input[pos:], l.line, pos, nil)
			}

			if startsOrEndsWithDot(num) && !l.config.AllowPointEdgeNumbers {
				return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
			}

			if isInteger(num) {
				return NewToken(NUMBER, INTEGER, l.input[pos:], l.line, pos, nil)
			}

			if isFloat(num) {
				return NewToken(NUMBER, FLOAT, l.input[pos:], l.line, pos, nil)
			}

			if isScientificNotation(num) {
				return NewToken(NUMBER, SCI_NOT, l.input[pos:], l.line, pos, nil)
			}

			isHexDec := isHex(num)
			if isHexDec && !l.config.AllowHexNumbers {
				return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
			}

			if isHexDec {
				return NewToken(NUMBER, HEX, l.input[pos:], l.line, pos, nil)
			}

		}
		// Lexing number ends here

		// Lexing ident starts here
		if IsPossibleJSIdentifier(l.char) {
			l.readChar()

			for IsPossibleJSIdentifier(l.char) {
				l.readChar()
			}

			if l.config.AllowUnquoted && IsJSIdentifier(l.input[pos:l.pos]) {
				return NewToken(STRING, IDENT, l.input[pos:l.pos], l.line, pos, nil)
			}
		}
		// Lexing ident ends here

		return NewToken(ILLEGAL, NONE, l.input[pos:], l.line, pos, nil)
	}
}

// Tokens returns a slice of all tokens produced so far by the lexer.
func (l *Lexer) Tokens() Tokens {
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

// readChar advances the lexer to the next character in the input,
// updating the current character, position, and line counters as needed.
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

// peek returns the next character in the input without advancing the lexer.
// If the end of input is reached, it returns 0.
func (l *Lexer) peek() byte {
	if l.readPos > (len(l.input) - 1) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

// peekBy returns the character at a position `target` bytes ahead in the input
// without moving the current position. Returns 0 if out of bounds.
func (l *Lexer) peekBy(target int) byte {
	pos := l.pos + target

	if pos > (len(l.input) - 1) {
		return 0
	} else {
		return l.input[pos]
	}
}

// prev returns the character immediately before the current position.
// Returns 0 if the current position is at the start.
func (l *Lexer) prev() byte {
	if (l.pos - 1) < 0 {
		return 0
	} else {
		return l.input[l.pos-1]
	}
}

// prevBy returns the character `target` bytes before the current position.
// Returns 0 if the target is out of bounds.
func (l *Lexer) prevBy(target int) byte {
	pos := l.pos - target

	if pos < 0 {
		return 0
	} else {
		return l.input[pos]
	}
}
