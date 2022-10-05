package expression

import (
	"github.com/nusr/gojs/token"
)

type Expression interface {
	Accept(visitor Visitor) (result any)
	String() string
}

type AssignExpression struct {
	Name  *token.Token
	Value Expression
}

func (assignExpression AssignExpression) Accept(visitor Visitor) any {
	return visitor.VisitAssignExpression(assignExpression)
}

func (assignExpression AssignExpression) String() string {
	return ""
}

type BinaryExpression struct {
	Left     Expression
	Operator *token.Token
	Right    Expression
}

func (binaryExpression BinaryExpression) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpression(binaryExpression)
}

func (binaryExpression BinaryExpression) String() string {
	return ""
}

type CallExpression struct {
	Callee       Expression
	Paren        *token.Token
	ArgumentList []Expression
}

func (callExpression CallExpression) Accept(visitor Visitor) any {
	return visitor.VisitCallExpression(callExpression)
}

func (callExpression CallExpression) String() string {
	return ""
}

type GetExpression struct {
	Object Expression
	Name   *token.Token
}

func (getExpression GetExpression) Accept(visitor Visitor) any {
	return visitor.VisitGetExpression(getExpression)
}

func (getExpression GetExpression) String() string {
	return ""
}

type SetExpression struct {
	Object Expression
	Name   *token.Token
	Value  Expression
}

func (setExpression SetExpression) Accept(visitor Visitor) any {
	return visitor.VisitSetExpression(setExpression)
}

func (setExpression SetExpression) String() string {
	return ""
}

type GroupingExpression struct {
	Expression Expression
}

func (groupingExpression GroupingExpression) Accept(visitor Visitor) any {
	return visitor.VisitGroupingExpression(groupingExpression)
}

func (groupingExpression GroupingExpression) String() string {
	return ""
}

type LiteralExpression struct {
	Value     string
	TokenType token.Type
}

func (literalExpression LiteralExpression) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpression(literalExpression)
}

func (literalExpression LiteralExpression) String() string {
	return ""
}

type LogicalExpression struct {
	Left     Expression
	Operator *token.Token
	Right    Expression
}

func (logicalExpression LogicalExpression) Accept(visitor Visitor) any {
	return visitor.VisitLogicalExpression(logicalExpression)
}

func (logicalExpression LogicalExpression) String() string {
	return ""
}

type SuperExpression struct {
	Keyword *token.Token
	Value   Expression
}

func (superExpression SuperExpression) Accept(visitor Visitor) any {
	return visitor.VisitSuperExpression(superExpression)
}

func (superExpression SuperExpression) String() string {
	return ""
}

type ThisExpression struct {
	Keyword *token.Token
}

func (thisExpression ThisExpression) Accept(visitor Visitor) any {
	return visitor.VisitThisExpression(thisExpression)
}

func (thisExpression ThisExpression) String() string {
	return ""
}

type UnaryExpression struct {
	Operator *token.Token
	Right    Expression
}

func (unaryExpression UnaryExpression) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpression(unaryExpression)
}

func (unaryExpression UnaryExpression) String() string {
	return ""
}

type VariableExpression struct {
	Name *token.Token
}

func (variableExpression VariableExpression) Accept(visitor Visitor) any {
	return visitor.VisitVariableExpression(variableExpression)
}

func (variableExpression VariableExpression) String() string {
	return variableExpression.Name.Lexeme
}

type Visitor interface {
	VisitAssignExpression(expression AssignExpression) any
	VisitBinaryExpression(expression BinaryExpression) any
	VisitCallExpression(expression CallExpression) any
	VisitGetExpression(expression GetExpression) any
	VisitSetExpression(expression SetExpression) any
	VisitGroupingExpression(expression GroupingExpression) any
	VisitLiteralExpression(expression LiteralExpression) any
	VisitLogicalExpression(expression LogicalExpression) any
	VisitSuperExpression(expression SuperExpression) any
	VisitThisExpression(expression ThisExpression) any
	VisitUnaryExpression(expression UnaryExpression) any
	VisitVariableExpression(expression VariableExpression) any
}
