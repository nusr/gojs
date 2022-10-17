package scanner

import (
	"fmt"
	"strings"

	"github.com/nusr/gojs/token"
)

const (
	EmptyData = 0
)

var KeywordMap = map[string]token.Type{
	"class":    token.Class,
	"else":     token.Else,
	"false":    token.False,
	"for":      token.For,
	"function": token.Function,
	"if":       token.If,
	"null":     token.Null,
	"return":   token.Return,
	"super":    token.Super,
	// "this":     token.This,
	"true":   token.True,
	"var":    token.Var,
	"while":  token.While,
	"do":     token.Do,
	"new":    token.New,
	"static": token.Static,
}

type Scanner struct {
	source  []rune
	tokens  []token.Token
	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source:  []rune(source), // unicode support
		tokens:  []token.Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}
func (scanner *Scanner) getChar(index int) rune {
	return scanner.source[index]
}
func (scanner *Scanner) peek() rune {
	if scanner.isAtEnd() {
		return EmptyData
	}
	return scanner.getChar(scanner.current)
}
func (scanner *Scanner) peekNext() rune {
	if scanner.current+1 < len(scanner.source) {
		return scanner.getChar(scanner.current + 1)
	}
	return EmptyData
}
func (scanner *Scanner) advance() rune {
	c := scanner.getChar(scanner.current)
	scanner.current++
	return c
}

func (scanner *Scanner) getSubString(start int, end int) string {
	text := string(scanner.source[start:end])
	return text
}

func (scanner *Scanner) appendToken(tokenType token.Type, text string) {
	scanner.tokens = append(scanner.tokens, token.Token{
		Type:   tokenType,
		Lexeme: text,
		Line:   scanner.line,
	})
}

func (scanner *Scanner) addToken(tokenType token.Type) {
	text := scanner.getSubString(scanner.start, scanner.current)
	scanner.appendToken(tokenType, text)
}
func (scanner *Scanner) match(char rune) bool {
	if scanner.isAtEnd() {
		return false
	}
	if scanner.getChar(scanner.current) != char {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (scanner *Scanner) isWhiteSpace(c rune) bool {
	return c == ' ' || c == '\r' || c == '\n' || c == '\t'
}

func (scanner *Scanner) isIdentifierChar(c rune) bool {
	return !scanner.isAtEnd() && !scanner.isWhiteSpace(c) && !strings.ContainsRune("[]{}(),.+-*/%;:?&|!=><\"'", c)
}

func (scanner *Scanner) number() {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}
	if scanner.peek() == '.' && scanner.isDigit(scanner.peekNext()) {
		scanner.advance()
		for scanner.isDigit(scanner.peek()) {
			scanner.advance()
		}
		scanner.addToken(token.Float64)
	} else {
		scanner.addToken(token.Int64)
	}

}

func (scanner *Scanner) string(end rune) {
	for scanner.peek() != end && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}
	if scanner.isAtEnd() {
		panic("unterminated string")
	}
	scanner.advance() // skip "
	text := scanner.getSubString(scanner.start+1, scanner.current-1)
	scanner.appendToken(token.String, text)
}

func (scanner *Scanner) identifier() {
	for scanner.isIdentifierChar(scanner.peek()) {
		scanner.advance()
	}
	text := scanner.getSubString(scanner.start, scanner.current)
	tokenType := token.Identifier
	if val, ok := KeywordMap[text]; ok {
		tokenType = val
	}
	scanner.addToken(tokenType)
}

func (scanner *Scanner) scanToken() {
	c := scanner.advance()
	switch c {
	case '(':
		scanner.addToken(token.LeftParen)
	case ')':
		scanner.addToken(token.RightParen)
	case '{':
		scanner.addToken(token.LeftBrace)
	case '}':
		scanner.addToken(token.RightBrace)
	case '[':
		scanner.addToken(token.LeftSquare)
	case ']':
		scanner.addToken(token.RightSquare)
	case ',':
		scanner.addToken(token.Comma)
	case '.':
		scanner.addToken(token.Dot)
	case '-':
		if scanner.match('-') {
			scanner.addToken(token.MinusMinus)
		} else if scanner.match('=') {
			scanner.addToken(token.MinusEqual)
		} else {
			scanner.addToken(token.Minus)
		}
	case '+':
		if scanner.match('+') {
			scanner.addToken(token.PlusPlus)
		} else if scanner.match('=') {
			scanner.addToken(token.PlusEqual)
		} else {
			scanner.addToken(token.Plus)
		}
	case ';':
		scanner.addToken(token.Semicolon)
	case ':':
		scanner.addToken(token.Colon)
	case '^':
		if scanner.match('=') {
			scanner.addToken(token.BitXOrEqual)
		} else {
			scanner.addToken(token.BitXOr)
		}
	case '%':
		if scanner.match('=') {
			scanner.addToken(token.PercentEqual)
		} else {
			scanner.addToken(token.Percent)
		}
	case '?':
		scanner.addToken(token.Mark)
	case '~':
		scanner.addToken(token.BitNot)
	case '&':
		if scanner.match('&') {
			if scanner.match('=') {
				scanner.addToken(token.AndEqual)
			} else {
				scanner.addToken(token.And)
			}
		} else if scanner.match('=') {
			scanner.addToken(token.BitAndEqual)
		} else {
			scanner.addToken(token.BitAnd)
		}
	case '|':
		if scanner.match('|') {
			if scanner.match('=') {
				scanner.addToken(token.OrEqual)
			} else {
				scanner.addToken(token.Or)
			}
		} else if scanner.match('=') {
			scanner.addToken(token.BitOrEqual)
		} else {
			scanner.addToken(token.BitOr)
		}
	case '*':
		if scanner.match('=') {
			scanner.addToken(token.StarEqual)
		} else if scanner.match('*') {
			if scanner.match('=') {
				scanner.addToken(token.StarStarEqual)
			} else {
				scanner.addToken(token.StarStar)
			}
		} else {
			scanner.addToken(token.Star)
		}
	case '!':
		if scanner.match('=') {
			if scanner.match('=') {
				scanner.addToken(token.BangEqualEqual)
			} else {
				scanner.addToken(token.BangEqual)
			}
		} else {
			scanner.addToken(token.Bang)
		}
	case '=':
		if scanner.match('=') {
			if scanner.match('=') {
				scanner.addToken(token.EqualEqualEqual)
			} else {
				scanner.addToken(token.EqualEqual)
			}
		} else {
			scanner.addToken(token.Equal)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(token.GreaterEqual)
		} else if scanner.match('>') {
			if scanner.match('>') {
				if scanner.match('=') {
					scanner.addToken(token.BitUnsignedRightShiftEqual)
				} else {
					scanner.addToken(token.BitUnsignedRightShift)
				}
			} else if scanner.match('=') {
				scanner.addToken(token.BitRightShiftEqual)
			} else {
				scanner.addToken(token.BitRightShift)
			}
		} else {
			scanner.addToken(token.Greater)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(token.LessEqual)
		} else if scanner.match('<') {
			if scanner.match('=') {
				scanner.addToken(token.BitLeftShiftEqual)
			} else {
				scanner.addToken(token.BitLeftShift)
			}
		} else {
			scanner.addToken(token.Less)
		}
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else if scanner.match('=') {
			scanner.addToken(token.SlashEqual)
		} else if scanner.match('*') {
			for !((scanner.peek() == '*' && scanner.peekNext() == '/') || scanner.isAtEnd()) {
				scanner.advance()
			}
			scanner.advance() // skip *
			scanner.advance() // skip /
		} else {
			scanner.addToken(token.Slash)
		}
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		scanner.line++
	case '\'':
		scanner.string(c)
	case '"':
		scanner.string(c)
	default:
		if scanner.isDigit(c) {
			scanner.number()
		} else if scanner.isIdentifierChar(c) {
			scanner.identifier()
		} else {
			fmt.Printf("Unexpected character:%c\n", c)
		}
	}
}

func (scanner *Scanner) Scan() []token.Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}
	scanner.appendToken(token.EOF, "")
	return scanner.tokens
}
