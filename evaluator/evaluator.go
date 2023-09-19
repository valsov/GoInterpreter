package evaluator

import (
	"fmt"

	"github.com/valsov/gointerpreter/ast"
	"github.com/valsov/gointerpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBoolean(node.Value)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	default:
		return nil
	}
}

func evalProgram(program *ast.Program) object.Object {
	var obj object.Object
	for _, s := range program.Statements {
		obj = Eval(s)

		switch obj := obj.(type) {
		case *object.ReturnValue:
			return obj.Value
		case *object.Error:
			return obj
		}
	}
	return obj
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var obj object.Object
	for _, statement := range block.Statements {
		obj = Eval(statement)

		if obj != nil && (obj.Type() == object.RETURN_VALUE_OBJ || obj.Type() == object.ERROR_OBJ) {
			return obj
		}
	}
	return obj
}

func nativeBoolToBoolean(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		// Expect same type from left and right
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		// Compare by memory address, this works for booleans since we use the same instances TRUE & FALSE
		return nativeBoolToBoolean(left == right)
	case operator == "!=":
		return nativeBoolToBoolean(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value} // Apply minus here
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBoolean(leftValue < rightValue)
	case ">":
		return nativeBoolToBoolean(leftValue > rightValue)
	case "==":
		return nativeBoolToBoolean(leftValue == rightValue)
	case "!=":
		return nativeBoolToBoolean(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ifExp *ast.IfExpression) object.Object {
	condition := Eval(ifExp.Condition)
	if isError(condition) {
		return condition
	}

	if isTrue(condition) {
		return Eval(ifExp.Consequence)
	} else if ifExp.Alternative != nil {
		return Eval(ifExp.Alternative)
	} else {
		return NULL
	}
}

func isTrue(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func newError(format string, values ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, values...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Inspect() == object.ERROR_OBJ
}
