package statement

import (
	"github.com/nusr/gojs/expression"
	"github.com/nusr/gojs/token"
)

type Statement interface {
	Accept(visitor Visitor) any
}


type BlockStatement struct {
	Statements []Statement
}

func (blockStatement BlockStatement) Accept(visitor Visitor) any {
	return visitor.VisitBlockStatement(blockStatement)
}

type ClassStatement struct {
	Name       *token.Token
	SuperClass expression.VariableExpression
	Methods    []FunctionStatement
}

func (classStatement ClassStatement) Accept(visitor Visitor) any {
	return visitor.VisitClassStatement(classStatement)
}

type ExpressionStatement struct {
	Expression expression.Expression
}

func (expressionStatement ExpressionStatement) Accept(visitor Visitor) any {
	return visitor.VisitExpressionStatement(expressionStatement)
}

type FunctionStatement struct {
	Name   *token.Token
	Body   BlockStatement
	Params []*token.Token
}

func (functionStatement FunctionStatement) Accept(visitor Visitor) any {
	return visitor.VisitFunctionStatement(functionStatement)
}

type IfStatement struct {
	Condition  expression.Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (ifStatement IfStatement) Accept(visitor Visitor) any {
	return visitor.VisitIfStatement(ifStatement)
}

type PrintStatement struct {
	Expression expression.Expression
	Comment    *token.Token
}

func (printStatement PrintStatement) Accept(visitor Visitor) any {
	return visitor.VisitPrintStatement(printStatement)
}

type ReturnStatement struct {
	Keyword *token.Token
	Value   expression.Expression
}

func (returnStatement ReturnStatement) Accept(visitor Visitor) any {
	return visitor.VisitReturnStatement(returnStatement)
}

type VariableStatement struct {
	Name        *token.Token
	Initializer expression.Expression
}

func (variableStatement VariableStatement) Accept(visitor Visitor) any {
	return visitor.VisitVariableStatement(variableStatement)
}

type WhileStatement struct {
	Condition expression.Expression
	Body      Statement
}

func (whileStatement WhileStatement) Accept(visitor Visitor) any {
	return visitor.VisitWhileStatement(whileStatement)
}

type Visitor interface {
	VisitBlockStatement(statement BlockStatement) any
	VisitClassStatement(statement ClassStatement) any
	VisitExpressionStatement(statement ExpressionStatement) any
	VisitFunctionStatement(statement FunctionStatement) any
	VisitIfStatement(statement IfStatement) any
	VisitPrintStatement(statement PrintStatement) any
	VisitReturnStatement(statement ReturnStatement) any
	VisitVariableStatement(statement VariableStatement) any
	VisitWhileStatement(statement WhileStatement) any
}

