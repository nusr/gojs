package interpreter

import (
	"fmt"
	"math"
	"strconv"

	"github.com/nusr/gojs/call"
	"github.com/nusr/gojs/control_flow"
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/expression"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
)

const (
	nanNumber = "NaN"
)

const whileMaxIteration = 100000

func convertBtoF(b bool) float64 {
	if b {
		return float64(1)
	}
	return float64(0)
}

func convertLtoI(left any, right any) (int64, int64, bool) {
	leftInt, leftType := left.(int64)
	rightInt, rightType := right.(int64)
	if leftType && rightType {
		return leftInt, rightInt, true
	}
	return 0, 0, false
}

func convertLtoF(left any, right any) (float64, float64, bool) {
	a := float64(0)
	b := float64(0)
	count := 0
	if left == nil && right == nil {
		count = 2
	}
	if left == nil {
		count++
	}
	if right == nil {
		count++
	}
	if val, ok := left.(float64); ok {
		a = val
		count++
	}
	if val, ok := right.(float64); ok {
		b = val
		count++
	}
	if val, ok := left.(bool); ok {
		a = convertBtoF(val)
		count++
	}
	if val, ok := right.(bool); ok {
		b = convertBtoF(val)
		count++
	}
	if val, ok := left.(int64); ok {
		a = float64(val)
		count++
	}
	if val, ok := right.(int64); ok {
		b = float64(val)
		count++
	}
	return a, b, count >= 2
}

type Interpreter struct {
	environment *environment.Environment
	globals     *environment.Environment
}

func NewInterpreter(environment *environment.Environment) *Interpreter {
	return &Interpreter{
		environment: environment,
		globals:     environment,
	}
}
func (interpreter *Interpreter) GetGlobal() *environment.Environment {
	return interpreter.globals
}
func (interpreter *Interpreter) Interpret(list []statement.Statement) any {
	var result any
	for _, item := range list {
		result = interpreter.execute(item)
		if val, ok := result.(control_flow.ReturnValue); ok {
			return val.Value
		}
	}
	return result
}

func (interpreter *Interpreter) execute(statement statement.Statement) any {
	if statement == nil {
		return nil
	}
	return statement.Accept(interpreter)
}

func (interpreter *Interpreter) evaluate(expression expression.Expression) any {
	if expression == nil {
		return nil
	}
	t := expression.Accept(interpreter)
	if val, ok := t.(control_flow.ReturnValue); ok {
		return val.Value
	}
	return t
}

func (interpreter *Interpreter) isTruth(value any) bool {
	if value == true {
		return true
	}
	list := []interface{}{
		nil,
		"",
		int64(0),
		float64(0),
		false,
	}
	for _, item := range list {
		if item == value {
			return false
		}
	}
	text := token.LiteralTypeToString(value)
	result := text != ""
	return result
}

func (interpreter *Interpreter) ExecuteBlock(statement statement.BlockStatement, environment *environment.Environment) (result any) {
	previous := interpreter.environment
	interpreter.environment = environment
	for _, t := range statement.Statements {
		result = interpreter.execute(t)
		if val, ok := result.(control_flow.ReturnValue); ok {
			interpreter.environment = previous
			return val
		}
	}
	interpreter.environment = previous
	return result
}

func (interpreter *Interpreter) VisitExpressionStatement(statement statement.ExpressionStatement) any {
	return interpreter.evaluate(statement.Expression)
}
func (interpreter *Interpreter) VisitVariableStatement(statement statement.VariableStatement) any {
	if statement.Initializer != nil {
		value := interpreter.evaluate(statement.Initializer)
		interpreter.environment.Define(statement.Name.Lexeme, value)
	} else {
		interpreter.environment.Define(statement.Name.Lexeme, nil)
	}
	return nil
}
func (interpreter *Interpreter) VisitBlockStatement(statement statement.BlockStatement) any {
	return interpreter.ExecuteBlock(statement, environment.NewEnvironment(interpreter.environment))
}
func (interpreter *Interpreter) VisitClassStatement(statement statement.ClassStatement) any {
	// TODO
	return nil
}
func (interpreter *Interpreter) VisitFunctionStatement(statement statement.FunctionStatement) any {
	interpreter.environment.Define(statement.Name.Lexeme, call.NewCallable(statement))
	return nil
}

func (interpreter *Interpreter) VisitIfStatement(statement statement.IfStatement) (result any) {
	if interpreter.isTruth(interpreter.evaluate(statement.Condition)) {
		result = interpreter.execute(statement.ThenBranch)
	} else if statement.ElseBranch != nil {
		result = interpreter.execute(statement.ElseBranch)
	}
	if val, ok := result.(control_flow.ReturnValue); ok {
		return val
	}
	return nil
}
func (interpreter *Interpreter) VisitPrintStatement(statement statement.PrintStatement) any {
	result := interpreter.evaluate(statement.Expression)
	actual := token.LiteralTypeToString(result)
	fmt.Println(actual)
	return result
}

func (interpreter *Interpreter) VisitReturnStatement(statement statement.ReturnStatement) any {
	value := interpreter.evaluate(statement.Value)
	return control_flow.NewReturnValue(value)
}
func (interpreter *Interpreter) VisitWhileStatement(statement statement.WhileStatement) any {
	for interpreter.isTruth(interpreter.evaluate(statement.Condition)) {
		t := interpreter.execute(statement.Body)
		if val, ok := t.(control_flow.ReturnValue); ok {
			return val
		}
	}
	return nil
}

func (interpreter *Interpreter) VisitVariableExpression(expression expression.VariableExpression) any {
	return interpreter.environment.Get(expression.Name.Lexeme)
}
func (interpreter *Interpreter) VisitLiteralExpression(expression expression.LiteralExpression) any {
	switch expression.TokenType {
	case token.NULL:
		return nil
	case token.STRING:
		return expression.String
	case token.FLOAT64:
		{
			result, _ := strconv.ParseFloat(expression.Value, 64)

			return result
		}
	case token.INT64:
		{
			result, _ := strconv.ParseInt(expression.Value, 10, 64)
			return result
		}
	case token.TRUE:
		return true
	case token.FALSE:
		return false
	}
	return nil
}

func (interpreter *Interpreter) VisitBinaryExpression(expression expression.BinaryExpression) any {
	left := interpreter.evaluate(expression.Left)
	right := interpreter.evaluate(expression.Right)
	switch expression.Operator.TokenType {
	case token.EqualEqual:
		return left == right
	case token.BangEqual:
		return left != right
	case token.LESS:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return false
			}
			a, b, check := convertLtoF(left, right)
			if !check {
				panic(any(fmt.Sprintf("LESS can not handle value left:%v,right:%v", left, right)))
			}
			return a < b
		}
	case token.LessEqual:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return false
			}
			a, b, check := convertLtoF(left, right)
			if !check {
				panic(any(fmt.Sprintf("LESS_EQUAL can not handle value left:%v,right:%v", left, right)))
			}
			return a <= b
		}
	case token.GREATER:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return false
			}
			a, b, check := convertLtoF(left, right)
			if !check {
				panic(any(fmt.Sprintf("GREATER can not handle value left:%v,right:%v", left, right)))
			}
			return a > b
		}
	case token.GreaterEqual:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return false
			}
			a, b, check := convertLtoF(left, right)
			if !check {
				panic(any(fmt.Sprintf("GREATER_EQUAL can not handle value left:%v,right:%v", left, right)))
			}
			return a >= b
		}
	case token.PLUS:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return token.LiteralTypeToString(left) + token.LiteralTypeToString(right)
			}
			if a, b, check := convertLtoI(left, right); check {
				return a + b
			}
			a, b, check := convertLtoF(left, right)
			if !check {
				panic(any(fmt.Sprintf("PLUS can not handle value left:%v,right:%v", left, right)))
			}
			return a + b
		}
	case token.MINUS:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return nanNumber
			}
			if a, b, check := convertLtoI(left, right); check {
				return a - b
			}
			a, b, check := convertLtoF(left, right)
			if !check {
				panic(any(fmt.Sprintf("MINUS can not handle value left:%v,right:%v", left, right)))
			}
			return a - b
		}
	case token.STAR:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return nanNumber
			}
			if a, b, check := convertLtoI(left, right); check {
				return a * b
			}
			a, b, check := convertLtoF(left, right)
			if !check {
				panic(any(fmt.Sprintf("STAR can not handle value left:%v,right:%v", left, right)))
			}
			return a * b
		}
	case token.SLASH:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return nanNumber
			}
			if a, b, check := convertLtoI(left, right); check {
				return a / b
			}
			a, b, check := convertLtoF(left, right)
			if !check {
				panic(any(fmt.Sprintf("STAR can not handle value left:%v,right:%v", left, right)))
			}
			if b == 0 {
				return math.MaxFloat64
			}
			return a / b
		}
	}
	return nil
}

func (interpreter *Interpreter) VisitCallExpression(expression expression.CallExpression) any {
	callable := interpreter.evaluate(expression.Callee)
	var params []any
	for _, item := range expression.ArgumentList {
		params = append(params, interpreter.evaluate(item))
	}
	if val, ok := callable.(call.BaseCallable); ok {
		if val.Size() != len(params) {
			panic(any(fmt.Sprintf("expect %d size arguments, actual get %d\n", val.Size(), len(params))))
		}
		return val.Call(interpreter, params)
	}
	panic(any("can only call function and call"))
}
func (interpreter *Interpreter) VisitGetExpression(expression expression.GetExpression) any {
	// TODO
	return nil
}
func (interpreter *Interpreter) VisitSetExpression(expression expression.SetExpression) any {
	// TODO
	return nil
}
func (interpreter *Interpreter) VisitLogicalExpression(expression expression.LogicalExpression) any {
	left := interpreter.evaluate(expression.Left)
	check := interpreter.isTruth(left)
	if expression.Operator.TokenType == token.AND {
		if !check {
			return left
		}
	} else {
		if check {
			return left
		}
	}
	return interpreter.evaluate(expression.Right)
}

func (interpreter *Interpreter) VisitSuperExpression(expression expression.SuperExpression) any {
	// TODO
	return nil
}

func (interpreter *Interpreter) VisitGroupingExpression(expression expression.GroupingExpression) any {
	result := interpreter.evaluate(expression.Expression)
	return result
}

func (interpreter *Interpreter) VisitThisExpression(expression expression.ThisExpression) any {
	// TODO
	return nil
}
func (interpreter *Interpreter) VisitUnaryExpression(expression expression.UnaryExpression) any {
	result := interpreter.evaluate(expression.Right)
	switch expression.Operator.TokenType {
	case token.PlusPlus:
		{

			var temp any

			if val, check := result.(int64); check {
				temp = val + 1
			} else {
				a, _, check := convertLtoF(result, 0)
				if check {
					temp = a + 1
				} else {
					panic(any("error type"))
				}
			}
			interpreter.environment.Assign(expression.Right.String(), temp)
			return temp

		}

	case token.MinusMinus:
		{
			var temp any
			if val, check := result.(int64); check {
				temp = val - 1
			} else {
				a, _, check := convertLtoF(result, 0)
				if check {
					temp = a - 1
				} else {
					panic(any("error type"))
				}
			}
			interpreter.environment.Assign(expression.Right.String(), temp)
			return temp

		}
	case token.PLUS:
		return result
	case token.MINUS:
		{
			if result == nil {
				return -0
			}
			if val, ok := result.(bool); ok {
				return convertBtoF(val)
			}
			if val, ok := result.(float64); ok {
				return -val
			}
			return nanNumber
		}
	case token.BANG:
		return !interpreter.isTruth(result)
	}
	return nil
}

func (interpreter *Interpreter) VisitAssignExpression(expression expression.AssignExpression) any {
	temp := interpreter.evaluate(expression.Value)
	interpreter.environment.Assign(expression.Name.Lexeme, temp)
	return temp
}
