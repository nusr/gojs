package scanner

import (
	"fmt"
	"github.com/nusr/gojs/token"
	"strings"
)

const (
	EmptyData = 0
)

var KeywordMap = map[string]token.Type{
	"call":         token.CLASS,
	"else":         token.ELSE,
	"false":        token.FALSE,
	"for":          token.FOR,
	"fun":          token.FUNCTION,
	"if":           token.IF,
	"null":         token.NULL,
	"print":        token.PRINT,
	"control_flow": token.RETURN,
	"super":        token.SUPER,
	"this":         token.THIS,
	"true":         token.TRUE,
	"var":          token.VAR,
	"while":        token.WHILE,
	"do":           token.DO,
}

type Scanner struct {
	source  []rune
	tokens  []*token.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	var tokens []*token.Token
	return &Scanner{
		source:  []rune(source),
		tokens:  tokens,
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

func (scanner *Scanner) addOneToken(tokenType token.Type) {
	text := scanner.getSubString(scanner.start, scanner.current)
	scanner.appendToken(tokenType, text)
}

func (scanner *Scanner) appendToken(tokenType token.Type, text string) {
	scanner.tokens = append(scanner.tokens, &token.Token{
		TokenType: tokenType,
		Lexeme:    text,
		Line:      scanner.line,
	})
}

func (scanner *Scanner) addToken(tokenType token.Type) {
	scanner.addOneToken(tokenType)
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

func (scanner *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '\\' || c == '_' || c == '$' || c == '#' || (c >= '\u4e00' && c <= '\u9fa5')
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
		scanner.addToken(token.FLOAT64)
	} else {
		scanner.addToken(token.INT64)
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
		fmt.Println("unterminated string")
		return
	}
	scanner.advance() // skip "
	text := scanner.getSubString(scanner.start+1, scanner.current-1)
	scanner.appendToken(token.STRING, text)
}

func (scanner *Scanner) identifier() {
	for scanner.isAlpha(scanner.peek()) {
		scanner.advance()
	}
	text := scanner.getSubString(scanner.start, scanner.current)
	tokenType := token.IDENTIFIER
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
		scanner.addToken(token.COMMA)
	case '.':
		scanner.addToken(token.DOT)
	case '-':
		if scanner.match('-') {
			scanner.addToken(token.MinusMinus)
		} else {
			scanner.addToken(token.MINUS)
		}
	case '+':
		if scanner.match('+') {
			scanner.addToken(token.PlusPlus)
		} else {
			scanner.addToken(token.PLUS)
		}
	case ';':
		scanner.addToken(token.SEMICOLON)
	case ':':
		scanner.addToken(token.COLON)
	case '%':
		scanner.addToken(token.PERCENT)
	case '?':
		scanner.addToken(token.MARK)
	case '&':
		if scanner.match('&') {
			scanner.addToken(token.AND)
		} else {
			scanner.addToken(token.BitAnd)
		}
	case '|':
		if scanner.match('|') {
			scanner.addToken(token.OR)
		} else {
			scanner.addToken(token.BitOr)
		}
	case '*':
		scanner.addToken(token.STAR)
	case '!':
		if scanner.match('=') {
			scanner.addToken(token.BangEqual)
		} else {
			scanner.addToken(token.BANG)
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(token.EqualEqual)
		} else {
			scanner.addToken(token.EQUAL)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(token.GreaterEqual)
		} else {
			scanner.addToken(token.GREATER)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(token.LessEqual)
		} else {
			scanner.addToken(token.LESS)
		}
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
			text := scanner.getSubString(scanner.start, scanner.current)
			if strings.Contains(text, "expect:") {
				scanner.appendToken(token.LineComment, text)
			}
		} else if scanner.match('*') {
			for !((scanner.peek() == '*' && scanner.peekNext() == '/') || scanner.isAtEnd()) {
				scanner.advance()
			}
			scanner.advance() // skip *
			scanner.advance() // skip /
		} else {
			scanner.addToken(token.SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		scanner.line++
	case '\'':
	case '"':
		scanner.string(c)
	default:
		if scanner.isDigit(c) {
			scanner.number()
		} else if scanner.isAlpha(c) {
			scanner.identifier()
		} else {
			fmt.Printf("Unexpected character:%c\n", c)
		}
	}
}

func (scanner *Scanner) ScanTokens() []*token.Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}
	scanner.appendToken(token.EOF, "")
	return scanner.tokens
}
