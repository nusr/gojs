package parser

import (
	"fmt"
	"strings"

	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
)

const maxParameterCount = 255

var assignmentMap = map[token.Type]token.Type{
	token.PlusEqual:                  token.Plus,
	token.MinusEqual:                 token.Minus,
	token.StarStarEqual:              token.StarStar,
	token.StarEqual:                  token.Equal,
	token.SlashEqual:                 token.Slash,
	token.PercentEqual:               token.Percent,
	token.BitLeftShiftEqual:          token.BitLeftShift,
	token.BitRightShiftEqual:         token.BitRightShift,
	token.BitUnsignedRightShiftEqual: token.BitUnsignedRightShift,
	token.BitAndEqual:                token.BitAnd,
	token.BitXOrEqual:                token.BitXOr,
	token.BitOrEqual:                 token.BitOr,
	token.AndEqual:                   token.And,
	token.OrEqual:                    token.Or,
}

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

func (parser *Parser) checkNext(tokenType token.Type) bool {
	if parser.current+2 <= len(parser.tokens) {
		return parser.tokens[parser.current+1].Type == tokenType
	}
	return false
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

func (parser *Parser) varDeclaration(isStatic bool) statement.Statement {
	name := parser.consume(token.Identifier, "expect identifier after var")
	var initializer statement.Expression
	if parser.match(token.Equal) {
		initializer = parser.expression()
	}
	parser.match(token.Semicolon)
	return statement.VariableStatement{
		Name:        name,
		Initializer: initializer,
		Static:      isStatic,
	}
}
func (parser *Parser) primary() statement.Expression {
	if parser.match(token.True, token.False, token.Null, token.Float64, token.Int64, token.String) {
		t := parser.previous()
		return statement.LiteralExpression{
			Value: t.Lexeme,
			Type:  t.Type,
		}
	}
	if parser.match(token.Identifier) {
		return statement.VariableExpression{
			Name: parser.previous(),
		}
	}
	if parser.match(token.LeftParen) {
		expr := parser.expression()
		parser.consume(token.RightParen, fmt.Sprintf("parser expected ')', actual:%s", parser.peek()))
		return statement.GroupingExpression{
			Expression: expr,
		}
	}
	if parser.match(token.LeftSquare) {
		list := parser.getExpressionList(token.RightSquare)
		parser.consume(token.RightSquare, "expect ]")
		return statement.ArrayLiteralExpression{
			Elements: list,
		}
	}
	if parser.match(token.LeftBrace) {
		var properties []statement.ObjectLiteralItem
		if !parser.check(token.RightBrace) {
			for ok := true; ok; ok = parser.match(token.Comma) {
				if parser.check(token.RightBrace) {
					break
				}
				key := parser.consume(token.Identifier, "expect object key")
				parser.consume(token.Colon, "expect :")
				value := parser.expression()
				properties = append(properties, statement.ObjectLiteralItem{
					Key: statement.TokenExpression{
						Name: key,
					},
					Value: value,
				})
			}
		}

		parser.consume(token.RightBrace, "expect }")
		return statement.ObjectLiteralExpression{
			Properties: properties,
		}
	}
	if parser.match(token.Function) {
		name := parser.getPartialName()
		parser.consume(token.LeftParen, "expect (")
		parameters := parser.getTokenList()
		parser.consume(token.RightParen, "expect )")
		parser.consume(token.LeftBrace, "expect {")
		body := parser.block()
		return statement.FunctionExpression{
			Name:   name,
			Body:   body,
			Params: parameters,
		}
	}
	if parser.match(token.Class) {
		name := parser.getPartialName()
		methods := parser.getClassBody()
		return statement.ClassExpression{
			Name:    name,
			Methods: methods,
		}
	}
	panic(fmt.Sprintf("parser can not handle token: %s", parser.peek()))
}
func (parser *Parser) getExpressionList(tokenType token.Type) []statement.Expression {
	var params []statement.Expression
	if parser.check(tokenType) {
		return params
	}
	count := 0
	for ok := true; ok; ok = parser.match(token.Comma) {
		if parser.check(tokenType) {
			break
		}
		if parser.match(token.Comma) {
			params = append(params, nil)
		} else {
			params = append(params, parser.expression())
		}
		if tokenType == token.RightParen {
			count++
			if count > maxParameterCount {
				panic("over max parameter count")
			}
		}
	}
	return params
}
func (parser *Parser) finishCall(callee statement.Expression) statement.Expression {
	params := parser.getExpressionList(token.RightParen)
	parser.consume(token.RightParen, "expect )")
	return statement.CallExpression{
		Callee:    callee,
		Arguments: params,
	}
}
func (parser *Parser) call() statement.Expression {
	expr := parser.primary()
	for {
		if parser.match(token.Dot) {
			name := parser.consume(token.Identifier, "expect name")
			expr = statement.GetExpression{
				Object: expr,
				Property: statement.TokenExpression{
					Name: name,
				},
				IsSquare: false,
			}
		} else if parser.match(token.LeftSquare) {
			name := parser.expression()
			parser.consume(token.RightSquare, "expect ]")
			expr = statement.GetExpression{
				Object:   expr,
				Property: name,
				IsSquare: true,
			}
		} else if parser.match(token.LeftParen) {
			expr = parser.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (parser *Parser) new() statement.Expression {
	if parser.match(token.New) {
		call := parser.call()
		return statement.NewExpression{
			Expression: call,
		}
	}
	return parser.call()
}

func (parser *Parser) postUnary() statement.Expression {
	expr := parser.new()
	if parser.match(token.PlusPlus, token.MinusMinus) {
		operator := parser.previous()
		return statement.PostUnaryExpression{
			Operator: operator,
			Left:     expr,
		}
	}
	return expr
}

func (parser *Parser) unary() statement.Expression {
	if parser.match(token.Minus, token.Plus, token.Bang, token.MinusMinus, token.PlusPlus, token.BitNot) {
		operator := parser.previous()
		value := parser.unary()
		return statement.UnaryExpression{
			Operator: operator,
			Right:    value,
		}
	}
	return parser.postUnary()
}

func (parser *Parser) exponentiation() statement.Expression {
	expo := parser.unary()
	for parser.match(token.StarStar) {
		operator := parser.previous()
		right := parser.exponentiation()
		expo = statement.BinaryExpression{
			Left:     expo,
			Operator: operator,
			Right:    right,
		}
	}
	return expo
}

func (parser *Parser) factor() statement.Expression {
	unary := parser.exponentiation()
	for parser.match(token.Star, token.Slash) {
		operator := parser.previous()
		right := parser.exponentiation()
		unary = statement.BinaryExpression{
			Left:     unary,
			Operator: operator,
			Right:    right,
		}
	}
	return unary
}

func (parser *Parser) term() statement.Expression {
	factor := parser.factor()
	for parser.match(token.Plus, token.Minus) {
		operator := parser.previous()
		right := parser.factor()
		factor = statement.BinaryExpression{
			Left:     factor,
			Operator: operator,
			Right:    right,
		}
	}
	return factor
}

func (parser *Parser) bitShift() statement.Expression {
	term := parser.term()
	for parser.match(token.BitLeftShift, token.BitRightShift, token.BitUnsignedRightShift) {
		operator := parser.previous()
		right := parser.term()
		term = statement.BinaryExpression{
			Left:     term,
			Operator: operator,
			Right:    right,
		}
	}
	return term
}

func (parser *Parser) comparison() statement.Expression {
	term := parser.bitShift()
	for parser.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := parser.previous()
		right := parser.bitShift()
		term = statement.BinaryExpression{
			Left:     term,
			Operator: operator,
			Right:    right,
		}
	}
	return term
}

func (parser *Parser) equality() statement.Expression {
	expr := parser.comparison()
	for parser.match(token.BangEqual, token.EqualEqual, token.EqualEqualEqual, token.BangEqualEqual) {
		operator := parser.previous()
		right := parser.comparison()
		expr = statement.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (parser *Parser) bitAnd() statement.Expression {
	expr := parser.equality()
	for parser.match(token.BitAnd) {
		operator := parser.previous()
		right := parser.bitAnd()
		expr = statement.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (parser *Parser) bitXOr() statement.Expression {
	expr := parser.bitAnd()
	for parser.match(token.BitXOr) {
		operator := parser.previous()
		right := parser.bitXOr()
		expr = statement.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (parser *Parser) bitOr() statement.Expression {
	expr := parser.bitXOr()
	for parser.match(token.BitOr) {
		operator := parser.previous()
		right := parser.bitOr()
		expr = statement.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (parser *Parser) and() statement.Expression {
	expr := parser.bitOr()
	for parser.match(token.And) {
		operator := parser.previous()
		right := parser.and()
		expr = statement.LogicalExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}
func (parser *Parser) or() statement.Expression {
	expr := parser.and()
	for parser.match(token.Or) {
		operator := parser.previous()
		right := parser.and()
		expr = statement.LogicalExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}
func (parser *Parser) assignment() statement.Expression {
	expr := parser.or()
	operatorType, check := assignmentMap[parser.peek().Type]
	if parser.match(token.Equal) || check {
		if check {
			parser.advance()
		}
		equal := parser.previous()
		value := parser.assignment()
		if check {
			Operator := token.Token{
				Type:   operatorType,
				Lexeme: strings.Replace(equal.Lexeme, "=", "", 1),
				Line:   equal.Line,
			}
			if equal.Type == token.AndEqual || equal.Type == token.OrEqual {
				value = statement.LogicalExpression{
					Left:     expr,
					Operator: Operator,
					Right:    value,
				}
			} else {
				value = statement.BinaryExpression{
					Left:     expr,
					Operator: Operator,
					Right:    value,
				}
			}
		}
		if val, ok := expr.(statement.VariableExpression); ok {
			return statement.AssignExpression{
				Name:  val.Name,
				Value: value,
			}
		} else if val, ok := expr.(statement.GetExpression); ok {
			return statement.SetExpression{
				Object: val,
				Value:  value,
			}
		}
		panic(fmt.Sprintf("invalid assign target: %s", equal))
	}
	return expr
}

func (parser *Parser) expression() statement.Expression {
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
	name := parser.previous()
	parser.consume(token.LeftParen, "expect (")

	var initializer statement.Statement
	if parser.match(token.Semicolon) {
		initializer = nil
	} else if parser.match(token.Var) {
		initializer = parser.varDeclaration(false)
	} else {
		initializer = parser.expressionStatement()
	}

	var condition statement.Expression
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
		condition = statement.LiteralExpression{
			Type:  token.True,
			Value: "true",
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
		Name:      name,
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
	name := parser.previous()
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
				Name:      name,
			},
		},
	}
}
func (parser *Parser) while() statement.Statement {
	name := parser.previous()
	parser.consume(token.LeftParen, "expect ( after while")
	condition := parser.expression()
	parser.consume(token.RightParen, "expected ) after while")
	body := parser.statement()
	return statement.WhileStatement{
		Condition: condition,
		Body:      body,
		Name:      name,
	}
}

func (parser *Parser) returnStatement() statement.Statement {
	var expr statement.Expression
	if !parser.check(token.Semicolon) {
		expr = parser.expression()
	}
	parser.match(token.Semicolon)
	return statement.ReturnStatement{
		Value: expr,
	}
}

func (parser *Parser) statement() statement.Statement {
	if parser.match(token.If) {
		return parser.ifStatement()
	}
	if parser.match(token.Return) {
		return parser.returnStatement()
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

func (parser *Parser) getTokenList() []token.Token {
	var parameters []token.Token
	if !parser.check(token.RightParen) {
		count := 0
		for ok := true; ok; ok = parser.match(token.Comma) {
			if parser.check(token.RightParen) {
				break
			}
			if parser.match(token.Comma) {
				break
			}
			parameters = append(parameters, parser.consume(token.Identifier, "expect parameter name"))
			count++
			if count > maxParameterCount {
				panic(any("over max parameter count"))
			}
		}
	}
	return parameters
}

func (parser *Parser) functionDeclaration(isStatic bool) statement.FunctionStatement {
	name := parser.consume(token.Identifier, "expect name")
	parser.consume(token.LeftParen, "expect (")
	parameters := parser.getTokenList()
	parser.consume(token.RightParen, "expect )")
	parser.consume(token.LeftBrace, "expect {")
	body := parser.block()
	return statement.FunctionStatement{
		Name:   name,
		Params: parameters,
		Body:   body,
		Static: isStatic,
	}
}

func (parser *Parser) getClassBody() []statement.Statement {
	parser.consume(token.LeftBrace, "expect {")
	var methods []statement.Statement
	for !parser.check(token.RightBrace) && !parser.isAtEnd() {
		isStatic := parser.match(token.Static)
		if parser.checkNext(token.LeftParen) {
			methods = append(methods, parser.functionDeclaration(isStatic))
		} else {
			methods = append(methods, parser.varDeclaration(isStatic))
		}
	}
	parser.consume(token.RightBrace, "expect }")
	return methods
}

func (parser *Parser) getPartialName() *token.Token {
	var name *token.Token
	if parser.check(token.Identifier) {
		parser.advance()
		t := parser.previous()
		name = &t
	}
	return name
}

func (parser *Parser) classDeclaration() statement.ClassStatement {
	name := parser.consume(token.Identifier, "expect call name")
	methods := parser.getClassBody()
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
		return parser.functionDeclaration(false)
	}
	if parser.match(token.Var) {
		return parser.varDeclaration(false)
	}

	return parser.statement()
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}
