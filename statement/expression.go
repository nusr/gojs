package statement

import (
	"strings"

	"github.com/nusr/gojs/token"
)

type ExpressionVisitor interface {
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
	VisitTokenExpression(expression TokenExpression) any
	VisitFunctionExpression(expression FunctionExpression) any
	VisitClassExpression(expression ClassExpression) any
	VisitArrayLiteralExpression(expression ArrayLiteralExpression) any
	VisitObjectLiteralExpression(expression ObjectLiteralExpression) any
	VisitNewExpression(expression NewExpression) any
}

type Expression interface {
	Accept(visitor ExpressionVisitor) (result any)
	String() string
}

type AssignExpression struct {
	Name  token.Token
	Value Expression
}

func (expression AssignExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitAssignExpression(expression)
}

func (expression AssignExpression) String() string {
	return expression.Name.String() + "=" + expression.Value.String()
}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (expression BinaryExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitBinaryExpression(expression)
}

func (expression BinaryExpression) String() string {
	return expression.Left.String() + expression.Operator.String() + expression.Right.String()
}

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
}

func (expression CallExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitCallExpression(expression)
}

func (expression CallExpression) String() string {
	var temp []string
	for _, item := range expression.Arguments {
		temp = append(temp, item.String())
	}
	return expression.Callee.String() + "(" + strings.Join(temp, ",") + ")"
}

type GetExpression struct {
	Object   Expression
	Property Expression
	IsSquare bool // []
}

func (expression GetExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitGetExpression(expression)
}

func (expression GetExpression) String() string {
	if expression.IsSquare {
		return expression.Object.String() + "[" + expression.Property.String() + "]"
	}
	return expression.Object.String() + "." + expression.Property.String()
}

type SetExpression struct {
	Object GetExpression
	Value  Expression
}

func (expression SetExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitSetExpression(expression)
}

func (expression SetExpression) String() string {
	return expression.Object.String() + "=" + expression.Value.String()
}

type GroupingExpression struct {
	Expression Expression
}

func (expression GroupingExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitGroupingExpression(expression)
}

func (expression GroupingExpression) String() string {
	return "(" + expression.Expression.String() + ")"
}

type LiteralExpression struct {
	Value string
	Type  token.Type
}

func (expression LiteralExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitLiteralExpression(expression)
}

func (expression LiteralExpression) String() string {
	return expression.Value
}

type LogicalExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (expression LogicalExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitLogicalExpression(expression)
}

func (expression LogicalExpression) String() string {
	return expression.Left.String() + expression.Operator.String() + expression.Right.String()
}

type SuperExpression struct {
	Arguments []Expression
}

func (expression SuperExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitSuperExpression(expression)
}

func (expression SuperExpression) String() string {
	return ""
}

type ThisExpression struct {
	Keyword token.Token
}

func (expression ThisExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitThisExpression(expression)
}

func (expression ThisExpression) String() string {
	return ""
}

type UnaryExpression struct {
	Operator token.Token
	Right    Expression
}

func (expression UnaryExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitUnaryExpression(expression)
}

func (expression UnaryExpression) String() string {
	return expression.Operator.String() + expression.Right.String()
}

type VariableExpression struct {
	Name token.Token
}

func (expression VariableExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitVariableExpression(expression)
}

func (expression VariableExpression) String() string {
	return expression.Name.String()
}

type TokenExpression struct {
	Name token.Token
}

func (expression TokenExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitTokenExpression(expression)
}

func (expression TokenExpression) String() string {
	return expression.Name.String()
}

type FunctionExpression struct {
	Name   *token.Token
	Body   BlockStatement
	Params []token.Token
}

func (expression FunctionExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitFunctionExpression(expression)
}

func (expression FunctionExpression) String() string {
	var temp []string
	for _, item := range expression.Params {
		temp = append(temp, item.String())
	}
	var name string
	if expression.Name != nil {
		name = " " + expression.Name.String()
	}
	return "function" + name + "(" + strings.Join(temp, ",") + ")" + expression.Body.String()
}

type ClassExpression struct {
	Name       *token.Token
	SuperClass *VariableExpression
	Methods    []Statement
}

func (expression ClassExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitClassExpression(expression)
}

func (expression ClassExpression) String() string {
	var temp []string
	for _, item := range expression.Methods {
		t := strings.Split(item.String(), " ")
		temp = append(temp, strings.Join(t[1:], " "))
	}
	var name string
	if expression.Name != nil {
		name = expression.Name.String()
	}
	return "class " + name + "{" + strings.Join(temp, "") + "}"
}

type ArrayLiteralExpression struct {
	Elements []Expression
}

func (expression ArrayLiteralExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitArrayLiteralExpression(expression)
}

func (expression ArrayLiteralExpression) String() string {
	var temp []string
	for _, item := range expression.Elements {
		temp = append(temp, item.String())
	}
	return "[" + strings.Join(temp, ",") + "]"
}

type ObjectLiteralItem struct {
	Key   Expression
	Value Expression
}

type ObjectLiteralExpression struct {
	Properties []ObjectLiteralItem
}

func (expression ObjectLiteralExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitObjectLiteralExpression(expression)
}

func (expression ObjectLiteralExpression) String() string {
	var temp []string
	for _, item := range expression.Properties {
		temp = append(temp, item.Key.String()+":"+item.Value.String())
	}
	return "{" + strings.Join(temp, ",") + "}"
}

type NewExpression struct {
	Expression Expression
}

func (expression NewExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitNewExpression(expression)
}

func (expression NewExpression) String() string {
	return "new " + expression.Expression.String()
}
