package interpreter

import (
	"fmt"
	"math"
	"strconv"

	"github.com/nusr/gojs/call"
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/flow"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
)

const (
	nanNumber = "NaN"
)

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
	environment     *environment.Environment
	globals         *environment.Environment
	lastObjectKey   any
	lastObjectValue call.Property
}

func New(environment *environment.Environment) *Interpreter {
	return &Interpreter{
		environment:     environment,
		globals:         environment,
		lastObjectKey:   "",
		lastObjectValue: nil,
	}
}
func (interpreter *Interpreter) GetGlobal() *environment.Environment {
	return interpreter.globals
}
func (interpreter *Interpreter) Interpret(list []statement.Statement) any {
	var result any
	for _, item := range list {
		result = interpreter.Execute(item)
		if val, ok := result.(flow.Return); ok {
			return val.Value
		}
	}
	return result
}

func (interpreter *Interpreter) Execute(statement statement.Statement) any {
	if statement == nil {
		return nil
	}
	return statement.Accept(interpreter)
}

func (interpreter *Interpreter) Evaluate(expression statement.Expression) any {
	if expression == nil {
		return nil
	}
	t := expression.Accept(interpreter)
	if val, ok := t.(flow.Return); ok {
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
	text := token.ConvertAnyToString(value)
	result := text != ""
	return result
}

func (interpreter *Interpreter) ExecuteBlock(statement statement.BlockStatement, environment *environment.Environment) (result any) {
	previous := interpreter.environment
	interpreter.environment = environment
	for _, t := range statement.Statements {
		result = interpreter.Execute(t)
		if val, ok := result.(flow.Return); ok {
			interpreter.environment = previous
			return val
		}
	}
	interpreter.environment = previous
	return result
}

func (interpreter *Interpreter) VisitExpressionStatement(statement statement.ExpressionStatement) any {
	return interpreter.Evaluate(statement.Expression)
}
func (interpreter *Interpreter) VisitVariableStatement(statement statement.VariableStatement) any {
	if statement.Initializer != nil {
		value := interpreter.Evaluate(statement.Initializer)
		interpreter.environment.Define(statement.Name.Lexeme, value)
	} else {
		interpreter.environment.Define(statement.Name.Lexeme, nil)
	}
	return nil
}
func (interpreter *Interpreter) VisitBlockStatement(statement statement.BlockStatement) any {
	return interpreter.ExecuteBlock(statement, environment.New(interpreter.environment))
}
func (interpreter *Interpreter) VisitClassStatement(statement statement.ClassStatement) any {
	class := call.NewClass(statement.Methods)
	interpreter.environment.Define(statement.Name.Lexeme, class)
	return nil
}
func (interpreter *Interpreter) VisitFunctionStatement(statement statement.FunctionStatement) any {
	interpreter.environment.Define(statement.Name.Lexeme, call.NewFunction(statement.Body, statement.Params, interpreter.environment))
	return nil
}

func (interpreter *Interpreter) VisitIfStatement(statement statement.IfStatement) (result any) {
	if interpreter.isTruth(interpreter.Evaluate(statement.Condition)) {
		result = interpreter.Execute(statement.ThenBranch)
	} else if statement.ElseBranch != nil {
		result = interpreter.Execute(statement.ElseBranch)
	}
	if val, ok := result.(flow.Return); ok {
		return val
	}
	return nil
}
func (interpreter *Interpreter) VisitPrintStatement(statement statement.PrintStatement) any {
	result := interpreter.Evaluate(statement.Expression)
	actual := token.ConvertAnyToString(result)
	fmt.Println(actual)
	return result
}

func (interpreter *Interpreter) VisitReturnStatement(statement statement.ReturnStatement) any {
	value := interpreter.Evaluate(statement.Value)
	return flow.NewReturnValue(value)
}
func (interpreter *Interpreter) VisitWhileStatement(statement statement.WhileStatement) any {
	for interpreter.isTruth(interpreter.Evaluate(statement.Condition)) {
		t := interpreter.Execute(statement.Body)
		if val, ok := t.(flow.Return); ok {
			return val
		}
	}
	return nil
}

func (interpreter *Interpreter) VisitVariableExpression(expression statement.VariableExpression) any {
	return interpreter.environment.Get(expression.Name.Lexeme)
}
func (interpreter *Interpreter) VisitLiteralExpression(expr statement.LiteralExpression) any {
	switch expr.Type {
	case token.Null:
		return nil
	case token.String:
		return expr.Value
	case token.Float64:
		{
			result, err := strconv.ParseFloat(expr.Value, 64)
			if err != nil {
				panic(err)
			}
			return result
		}
	case token.Int64:
		{
			result, err := strconv.ParseInt(expr.Value, 10, 64)
			if err != nil {
				panic(err)
			}
			return result
		}
	case token.True:
		return true
	case token.False:
		return false
	}
	return nil
}

func (interpreter *Interpreter) VisitBinaryExpression(expression statement.BinaryExpression) any {
	left := interpreter.Evaluate(expression.Left)
	right := interpreter.Evaluate(expression.Right)
	switch expression.Operator.Type {
	case token.EqualEqual:
		return left == right
	case token.BangEqual:
		return left != right
	case token.Less:
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
	case token.Greater:
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
	case token.Plus:
		{
			_, stringType1 := left.(string)
			_, stringType2 := right.(string)
			if stringType1 || stringType2 {
				return token.ConvertAnyToString(left) + token.ConvertAnyToString(right)
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
	case token.Minus:
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
	case token.Star:
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
	case token.Slash:
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

func (interpreter *Interpreter) VisitCallExpression(expression statement.CallExpression) any {
	callable := interpreter.Evaluate(expression.Callee)
	var params []any
	for _, item := range expression.Arguments {
		params = append(params, interpreter.Evaluate(item))
	}
	val, ok := callable.(call.Callable)
	fmt.Printf("Type of %v is %T, bool: %v", callable, callable, ok)
	if ok {
		return val.Call(interpreter, params)
	}
	panic(any("can only call function and call"))
}
func (interpreter *Interpreter) VisitGetExpression(expression statement.GetExpression) any {
	result := interpreter.Evaluate(expression.Object)
	key := interpreter.Evaluate(expression.Property)
	if val, ok := result.(call.Property); ok {
		interpreter.lastObjectKey = key
		interpreter.lastObjectValue = val
		return val.Get(key)
	}
	return nil
}
func (interpreter *Interpreter) VisitSetExpression(expression statement.SetExpression) any {
	interpreter.lastObjectKey = ""
	interpreter.lastObjectValue = nil
	interpreter.Evaluate(expression.Object)
	key := interpreter.lastObjectKey
	object := interpreter.lastObjectValue
	if object != nil {
		value := interpreter.Evaluate(expression.Value)
		object.Set(key, value)
		return value
	}
	return nil
}
func (interpreter *Interpreter) VisitLogicalExpression(expression statement.LogicalExpression) any {
	left := interpreter.Evaluate(expression.Left)
	check := interpreter.isTruth(left)
	if expression.Operator.Type == token.And {
		if !check {
			return left
		}
	} else {
		if check {
			return left
		}
	}
	return interpreter.Evaluate(expression.Right)
}

func (interpreter *Interpreter) VisitSuperExpression(expression statement.SuperExpression) any {
	// TODO
	return nil
}

func (interpreter *Interpreter) VisitGroupingExpression(expression statement.GroupingExpression) any {
	result := interpreter.Evaluate(expression.Expression)
	return result
}

func (interpreter *Interpreter) VisitThisExpression(expression statement.ThisExpression) any {
	// TODO
	return nil
}
func (interpreter *Interpreter) VisitUnaryExpression(expression statement.UnaryExpression) any {
	result := interpreter.Evaluate(expression.Right)
	switch expression.Operator.Type {
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
	case token.Plus:
		return result
	case token.Minus:
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
	case token.Bang:
		return !interpreter.isTruth(result)
	}
	return nil
}

func (interpreter *Interpreter) VisitAssignExpression(expression statement.AssignExpression) any {
	temp := interpreter.Evaluate(expression.Value)
	interpreter.environment.Assign(expression.Name.Lexeme, temp)
	return temp
}

func (interpreter *Interpreter) VisitTokenExpression(expression statement.TokenExpression) any {
	return expression.Name.Lexeme
}

func (interpreter *Interpreter) VisitFunctionExpression(expression statement.FunctionExpression) any {
	return call.NewFunction(expression.Body, expression.Params, interpreter.environment)
}

func (interpreter *Interpreter) VisitClassExpression(expression statement.ClassExpression) any {
	class := call.NewClass(expression.Methods)
	return class
}

func (interpreter *Interpreter) VisitArrayLiteralExpression(expression statement.ArrayLiteralExpression) any {
	instance := call.NewArray()
	for i, item := range expression.Elements {
		value := interpreter.Evaluate(item)
		instance.Set(i, value)
	}
	return instance
}

func (interpreter *Interpreter) VisitObjectLiteralExpression(expression statement.ObjectLiteralExpression) any {
	instance := call.NewInstance()
	for _, item := range expression.Properties {
		key := interpreter.Evaluate(item.Key)
		value := interpreter.Evaluate(item.Value)
		instance.Set(key, value)
	}
	return instance
}

func (interpreter *Interpreter) VisitNewExpression(expression statement.NewExpression) any {
	if _, ok := expression.Expression.(statement.CallExpression); ok {
		result := interpreter.Evaluate(expression.Expression)
		return result
	}
	panic(`Class constructor cannot be invoked without 'new'`)
}
