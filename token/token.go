package token

type Token struct {
	Type   Type
	Lexeme string
	Line   int
}

func (token Token) String() string {
	return token.Lexeme
}
