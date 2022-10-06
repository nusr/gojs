package parser

import (
	"fmt"

	"github.com/nusr/gojs/expression"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
)

const maxParameterCount = 255

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) *Parser {
	return &Parser{
		current: 0,
		tokens:  tokens,
	}
}

func (parser *Parser) Parse() []statement.Statement {
	var statements []statement.Statement
	for !parser.isAtEnd() {
		statements = append(statements, parser.declaration())
	}
	return statements
}

func (parser *Parser) peek() token.Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) advance() {
	if parser.isAtEnd() {
		return
	}
	parser.current++
}

func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.current-1]
}

func (parser *Parser) consume(tokenType token.Type, message string) token.Token {
	if parser.peek().Type != tokenType {
		panic(any(message))
	}
	parser.advance()
	return parser.previous()
}

func (parser *Parser) check(tokenType token.Type) bool {
	if parser.isAtEnd() {
		return false
	}
	return parser.peek().Type == tokenType
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
	name := parser.consume(token.Identifier, "expect identifier after var")
	var initializer expression.Expression
	if parser.match(token.Equal) {
		initializer = parser.expression()
	}
	parser.match(token.Semicolon)
	return statement.VariableStatement{
		Name:        name,
		Initializer: initializer,
	}
}
func (parser *Parser) primary() expression.Expression {
	if parser.match(token.True) {
		return expression.LiteralExpression{
			TokenType: token.True,
		}
	}
	if parser.match(token.False) {
		return expression.LiteralExpression{
			TokenType: token.False,
		}
	}
	if parser.match(token.Null) {
		return expression.LiteralExpression{
			TokenType: token.Null,
		}
	}
	if parser.match(token.Float64, token.Int64) {
		t := parser.previous()
		return expression.LiteralExpression{
			Value:     t.Lexeme,
			TokenType: t.Type,
		}

	}
	if parser.match(token.String) {
		t := parser.previous()
		return expression.LiteralExpression{
			Value:     t.Lexeme,
			TokenType: t.Type,
		}
	}
	if parser.match(token.Identifier) {
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
		for ok := true; ok; ok = parser.match(token.Comma) {
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
	if parser.match(token.Minus, token.Plus, token.Bang, token.MinusMinus, token.PlusPlus) {
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
	for parser.match(token.Star, token.Slash) {
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
	for parser.match(token.Plus, token.Minus) {
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
	for parser.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
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
	for parser.match(token.And) {
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
	for parser.match(token.Or) {
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
	if parser.match(token.Equal) {
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
	if parser.match(token.Else) {
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

	parser.match(token.Semicolon)

	return statement.PrintStatement{
		Expression: expr,
	}
}

func (parser *Parser) expressionStatement() statement.Statement {
	expr := parser.expression()

	parser.match(token.Semicolon)

	return statement.ExpressionStatement{
		Expression: expr,
	}
}

func (parser *Parser) block() statement.BlockStatement {
	var statements []statement.Statement
	for !parser.isAtEnd() && parser.peek().Type != token.RightBrace {
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
	if parser.match(token.Semicolon) {
		initializer = nil
	} else if parser.match(token.Var) {
		initializer = parser.varDeclaration()
	} else {
		initializer = parser.expressionStatement()
	}

	var condition expression.Expression
	if !parser.check(token.Semicolon) {
		condition = parser.expression()
	}
	parser.consume(token.Semicolon, "expect ;")

	var increment statement.Statement
	if !parser.check(token.RightParen) {
		increment = parser.expressionStatement()
	}
	parser.consume(token.RightParen, "expect )")

	body := parser.statement()

	if condition == nil {
		condition = expression.LiteralExpression{
			TokenType: token.True,
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
	parser.consume(token.While, "expect while")
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
	keyword := parser.previous()
	var expr expression.Expression
	if !parser.check(token.Semicolon) {
		expr = parser.expression()
	}
	parser.match(token.Semicolon)
	return statement.ReturnStatement{
		Keyword: keyword,
		Value:   expr,
	}
}

func (parser *Parser) statement() statement.Statement {
	if parser.match(token.If) {
		return parser.ifStatement()
	}
	if parser.match(token.Return) {
		return parser.returnStatement()
	}
	if parser.match(token.Print) {
		return parser.printStatement()
	}
	if parser.match(token.LeftBrace) {
		return parser.block()
	}
	if parser.match(token.Do) {
		return parser.doWhile()
	}
	if parser.match(token.For) {
		return parser.forStatement()
	}
	if parser.match(token.While) {
		return parser.while()
	}
	return parser.expressionStatement()
}

func (parser *Parser) functionDeclaration() statement.FunctionStatement {
	name := parser.consume(token.Identifier, "expect name")
	parser.consume(token.LeftParen, "expect (")
	var parameters []token.Token
	if !parser.check(token.RightParen) {
		count := 0
		for ok := true; ok; ok = parser.match(token.Comma) {
			parameters = append(parameters, parser.consume(token.Identifier, "expect parameter name"))
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
	name := parser.consume(token.Identifier, "expect call name")
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
	if parser.match(token.Class) {
		return parser.classDeclaration()
	}
	if parser.match(token.Function) {
		return parser.functionDeclaration()
	}
	if parser.match(token.Var) {
		return parser.varDeclaration()
	}

	return parser.statement()
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}
