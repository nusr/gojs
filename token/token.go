package token

import (
	"fmt"
	"strconv"
)

type Token struct {
	Type   Type
	Lexeme string
	Line   int
}

func (token Token) String() string {
	return fmt.Sprintf("Type: %s, Lexeme: %s, Line: %d", strconv.Itoa(int(token.Type)), token.Lexeme, token.Line)
}
