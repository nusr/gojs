package parser

import (
	"fmt"

	"github.com/nusr/gojs/expression"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
)

const maxParameterCount = 255

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		current: 0,
		tokens:  tokens,
	}
}
func (parser *Parser) peek() *token.Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) advance() {
	if parser.isAtEnd() {
		return
	}
	parser.current++
}

func (parser *Parser) previous() *token.Token {
	return parser.tokens[parser.current-1]
}

func (parser *Parser) consume(tokenType token.Type, message string) *token.Token {
	if parser.peek().TokenType != tokenType {
		panic(any(message))
	}
	parser.advance()
	return parser.previous()
}

func (parser *Parser) check(tokenType token.Type) bool {
	if parser.isAtEnd() {
		return false
	}
	return parser.peek().TokenType == tokenType
}

func (parser *Parser) match(tokenTypes ...token.Type) bool {
	for _, tokenType := range tokenTypes {
		if parser.check(tokenType) {
			parser.advance()
			return true
		}
	}
	return false
}

func (parser *Parser) varDeclaration() statement.Statement {
	name := parser.consume(token.IDENTIFIER, "expect identifier after var")
	if parser.match(token.EQUAL) {
		initializer := parser.expression()
		parser.match(token.SEMICOLON)
		return statement.VariableStatement{
			Name:        name,
			Initializer: initializer,
		}
	} else {
		parser.match(token.SEMICOLON)
		return statement.VariableStatement{
			Name: name,
		}
	}

}
func (parser *Parser) primary() expression.Expression {
	if parser.match(token.TRUE) {
		return expression.LiteralExpression{
			TokenType: token.TRUE,
		}
	}
	if parser.match(token.FALSE) {
		return expression.LiteralExpression{
			TokenType: token.FALSE,
		}
	}
	if parser.match(token.NULL) {
		return expression.LiteralExpression{
			TokenType: token.NULL,
		}
	}
	if parser.match(token.FLOAT64, token.INT64) {
		t := parser.previous()
		return expression.LiteralExpression{
			Value:     t.Lexeme,
			TokenType: t.TokenType,
		}

	}
	if parser.match(token.STRING) {
		return expression.LiteralExpression{
			Value:     parser.previous().Lexeme,
			TokenType: token.STRING,
		}
	}
	if parser.match(token.IDENTIFIER) {
		return expression.VariableExpression{
			Name: parser.previous(),
		}
	}
	if parser.match(token.LeftParen) {
		expr := parser.expression()
		parser.consume(token.RightParen, fmt.Sprintf("parser expected ')', actual:%s", parser.peek()))
		return expression.GroupingExpression{
			Expression: expr,
		}
	}
	panic(any(fmt.Sprintf("parser can not handle token: %s", parser.peek())))
}
func (parser *Parser) finishCall(callee expression.Expression) expression.Expression {
	var params []expression.Expression
	if !parser.check(token.RightParen) {
		count := 0
		for ok := true; ok; ok = parser.match(token.COMMA) {
			params = append(params, parser.expression())
			count++
			if count > maxParameterCount {
				panic(any("over max parameter count"))
			}
		}
	}
	paren := parser.consume(token.RightParen, "expect )")
	return expression.CallExpression{
		Callee:       callee,
		ArgumentList: params,
		Paren:        paren,
	}
}
func (parser *Parser) call() expression.Expression {
	expr := parser.primary()
	for {
		if parser.match(token.LeftParen) {
			expr = parser.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (parser *Parser) unary() expression.Expression {
	if parser.match(token.MINUS, token.PLUS, token.BANG, token.MinusMinus, token.PlusPlus) {
		operator := parser.previous()
		value := parser.unary()
		return expression.UnaryExpression{
			Operator: operator,
			Right:    value,
		}
	}
	return parser.call()

}
func (parser *Parser) factor() expression.Expression {
	unary := parser.unary()
	for parser.match(token.STAR, token.SLASH) {
		operator := parser.previous()
		right := parser.unary()
		unary = expression.BinaryExpression{
			Left:     unary,
			Operator: operator,
			Right:    right,
		}
	}
	return unary
}

func (parser *Parser) term() expression.Expression {
	factor := parser.factor()
	for parser.match(token.PLUS, token.MINUS) {
		operator := parser.previous()
		right := parser.factor()
		factor = expression.BinaryExpression{
			Left:     factor,
			Operator: operator,
			Right:    right,
		}
	}
	return factor
}

func (parser *Parser) comparison() expression.Expression {
	term := parser.term()
	for parser.match(token.GREATER, token.GreaterEqual, token.LESS, token.LessEqual) {
		operator := parser.previous()
		right := parser.term()
		term = expression.BinaryExpression{
			Left:     term,
			Operator: operator,
			Right:    right,
		}
	}
	return term
}

func (parser *Parser) equality() expression.Expression {
	expr := parser.comparison()
	for parser.match(token.BangEqual, token.EqualEqual) {
		operator := parser.previous()
		right := parser.comparison()
		expr = expression.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (parser *Parser) and() expression.Expression {
	expr := parser.equality()
	for parser.match(token.AND) {
		operator := parser.previous()
		right := parser.and()
		expr = expression.LogicalExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}
func (parser *Parser) or() expression.Expression {
	expr := parser.and()
	for parser.match(token.OR) {
		operator := parser.previous()
		right := parser.and()
		expr = expression.LogicalExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}
func (parser *Parser) assignment() expression.Expression {
	expr := parser.or()
	if parser.match(token.EQUAL) {
		equal := parser.previous()
		value := parser.assignment()
		if val, ok := expr.(expression.VariableExpression); ok {
			return expression.AssignExpression{
				Name:  val.Name,
				Value: value,
			}
		}
		panic(any(fmt.Sprintf("invalid assign target: %s", equal)))
	}
	return expr
}

func (parser *Parser) expression() expression.Expression {
	return parser.assignment()
}

func (parser *Parser) ifStatement() statement.Statement {
	parser.consume(token.LeftParen, "expect ( after if")
	expr := parser.expression()
	parser.consume(token.RightParen, "expected ) after if")
	thenBranch := parser.statement()
	if parser.match(token.ELSE) {
		elseBranch := parser.statement()
		return statement.IfStatement{
			Condition:  expr,
			ThenBranch: thenBranch,
			ElseBranch: elseBranch,
		}
	} else {
		return statement.IfStatement{
			Condition:  expr,
			ThenBranch: thenBranch,
		}
	}
}

func (parser *Parser) printStatement() statement.Statement {
	expr := parser.expression()
	if !parser.isAtEnd() {
		parser.match(token.SEMICOLON)
	}
	if parser.match(token.LineComment) {
		comment := parser.previous()
		return statement.PrintStatement{
			Expression: expr,
			Comment:    comment,
		}
	}
	return statement.PrintStatement{
		Expression: expr,
		Comment:    nil,
	}
}

func (parser *Parser) expressionStatement() statement.Statement {
	expr := parser.expression()
	if !parser.isAtEnd() {
		parser.match(token.SEMICOLON)
	}
	return statement.ExpressionStatement{
		Expression: expr,
	}
}

func (parser *Parser) block() statement.BlockStatement {
	var statements []statement.Statement
	for !parser.isAtEnd() && parser.peek().TokenType != token.RightBrace {
		statements = append(statements, parser.declaration())
	}
	parser.consume(token.RightBrace, "expected } after block")
	return statement.BlockStatement{
		Statements: statements,
	}
}

func (parser *Parser) forStatement() statement.Statement {
	parser.consume(token.LeftParen, "expect (")

	var initializer statement.Statement
	if parser.match(token.SEMICOLON) {
		initializer = nil
	} else if parser.match(token.VAR) {
		initializer = parser.varDeclaration()
	} else {
		initializer = parser.expressionStatement()
	}

	var condition expression.Expression
	if !parser.check(token.SEMICOLON) {
		condition = parser.expression()
	}
	parser.consume(token.SEMICOLON, "expect ;")

	var increment statement.Statement
	if !parser.check(token.RightParen) {
		increment = parser.expressionStatement()
	}
	parser.consume(token.RightParen, "expect )")

	body := parser.statement()

	if condition == nil {
		condition = expression.LiteralExpression{
			TokenType: token.TRUE,
		}
	}

	if increment != nil {
		body = statement.BlockStatement{
			Statements: []statement.Statement{
				body,
				increment,
			},
		}
	}

	body = statement.WhileStatement{
		Body:      body,
		Condition: condition,
	}

	if initializer != nil {
		body = statement.BlockStatement{
			Statements: []statement.Statement{
				initializer,
				body,
			},
		}
	}
	return body
}
func (parser *Parser) doWhile() statement.Statement {
	parser.consume(token.LeftBrace, "expect {")
	body := parser.block()
	parser.consume(token.WHILE, "expect while")
	parser.consume(token.LeftParen, "expect (")
	condition := parser.expression()
	parser.consume(token.RightParen, "expect )")
	return statement.BlockStatement{
		Statements: []statement.Statement{
			body,
			statement.WhileStatement{
				Body:      body,
				Condition: condition,
			},
		},
	}
}
func (parser *Parser) while() statement.Statement {
	parser.consume(token.LeftParen, "expect ( after while")
	condition := parser.expression()
	parser.consume(token.RightParen, "expected ) after while")
	body := parser.statement()
	return statement.WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func (parser *Parser) returnStatement() statement.Statement {
	t := parser.previous()
	expr := parser.expression()
	parser.match(token.SEMICOLON)
	return statement.ReturnStatement{
		Keyword: t,
		Value:   expr,
	}
}

func (parser *Parser) statement() statement.Statement {
	if parser.match(token.IF) {
		return parser.ifStatement()
	}
	if parser.match(token.RETURN) {
		return parser.returnStatement()
	}
	if parser.match(token.PRINT) {
		return parser.printStatement()
	}
	if parser.match(token.LeftBrace) {
		return parser.block()
	}
	if parser.match(token.DO) {
		return parser.doWhile()
	}
	if parser.match(token.FOR) {
		return parser.forStatement()
	}
	if parser.match(token.WHILE) {
		return parser.while()
	}
	return parser.expressionStatement()
}

func (parser *Parser) functionDeclaration() statement.FunctionStatement {
	name := parser.consume(token.IDENTIFIER, "expect name")
	parser.consume(token.LeftParen, "expect (")
	var parameters []*token.Token
	if !parser.check(token.RightParen) {
		count := 0
		for ok := true; ok; ok = parser.match(token.COMMA) {
			parameters = append(parameters, parser.consume(token.IDENTIFIER, "expect parameter name"))
			count++
			if count > maxParameterCount {
				panic(any("over max parameter count"))
			}
		}
	}
	parser.consume(token.RightParen, "expect )")
	parser.consume(token.LeftBrace, "expect {")
	body := parser.block()
	return statement.FunctionStatement{
		Name:   name,
		Params: parameters,
		Body:   body,
	}
}
func (parser *Parser) classDeclaration() statement.ClassStatement {
	name := parser.consume(token.IDENTIFIER, "expect call name")
	parser.consume(token.LeftBrace, "expect {")
	var methods []statement.FunctionStatement
	for !parser.check(token.RightBrace) && !parser.isAtEnd() {
		methods = append(methods, parser.functionDeclaration())
	}

	parser.consume(token.RightBrace, "expect }")
	return statement.ClassStatement{
		Methods: methods,
		Name:    name,
	}
}
func (parser *Parser) declaration() statement.Statement {
	if parser.match(token.CLASS) {
		return parser.classDeclaration()
	}
	if parser.match(token.FUNCTION) {
		return parser.functionDeclaration()
	}
	if parser.match(token.VAR) {
		return parser.varDeclaration()
	}

	return parser.statement()
}

func (parser *Parser) Parse() []statement.Statement {
	var statements []statement.Statement
	for !parser.isAtEnd() {
		statements = append(statements, parser.declaration())
	}
	return statements
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek().TokenType == token.EOF
}
