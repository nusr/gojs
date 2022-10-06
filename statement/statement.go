package statement

import (
	"strings"

	"github.com/nusr/gojs/token"
)

type Statement interface {
	Accept(visitor StatementVisitor) any
	String() string
}

type BlockStatement struct {
	Statements []Statement
}

func (statement BlockStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitBlockStatement(statement)
}

func (statement BlockStatement) String() string {
	var temp []string
	for _, item := range statement.Statements {
		temp = append(temp, item.String())
	}
	return "{" + strings.Join(temp, ";") + "}"
}

type ClassStatement struct {
	Name       token.Token
	SuperClass VariableExpression
	Methods    []Statement
}

func (statement ClassStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitClassStatement(statement)
}

func (statement ClassStatement) String() string {
	var temp []string
	for _, item := range statement.Methods {
		temp = append(temp, item.String())
	}

	return "class " + statement.Name.String() + "{" + strings.Join(temp, "") + "}"
}

type ExpressionStatement struct {
	Expression Expression
}

func (statement ExpressionStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitExpressionStatement(statement)
}

func (statement ExpressionStatement) String() string {
	return statement.Expression.String() + ";"
}

type FunctionStatement struct {
	Name   token.Token
	Body   BlockStatement
	Params []token.Token
}

func (statement FunctionStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitFunctionStatement(statement)
}

func (statement FunctionStatement) String() string {
	var temp []string
	for _, item := range statement.Params {
		temp = append(temp, item.String())
	}

	return "function " + statement.Name.String() + "(" + strings.Join(temp, ",") + "){" + statement.Body.String() + "}"
}

type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (statement IfStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitIfStatement(statement)
}

func (statement IfStatement) String() string {
	temp := "if(" + statement.Condition.String() + ")" + statement.ThenBranch.String()
	if statement.ElseBranch != nil {
		temp += "else " + statement.ElseBranch.String()
	}
	return temp
}

type PrintStatement struct {
	Expression Expression
}

func (statement PrintStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitPrintStatement(statement)
}

func (statement PrintStatement) String() string {
	return ""
}

type ReturnStatement struct {
	Value Expression
}

func (statement ReturnStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitReturnStatement(statement)
}

func (statement ReturnStatement) String() string {
	if statement.Value == nil {
		return "return;"
	}
	return "return " + statement.Value.String() + ";"
}

type VariableStatement struct {
	Name        token.Token
	Initializer Expression
}

func (statement VariableStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitVariableStatement(statement)
}

func (statement VariableStatement) String() string {
	temp := "var " + statement.Name.String()
	if statement.Initializer != nil {
		temp += "=" + statement.Initializer.String()
	}
	return temp + ";"
}

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func (statement WhileStatement) Accept(visitor StatementVisitor) any {
	return visitor.VisitWhileStatement(statement)
}

func (statement WhileStatement) String() string {
	return "while(" + statement.Condition.String() + ")" + statement.Body.String()
}

type StatementVisitor interface {
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
