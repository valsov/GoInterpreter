package evaluator

import (
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
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
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
		return &object.ReturnValue{Value: val}
	default:
		return nil
	}
}

func evalProgram(program *ast.Program) object.Object {
	var obj object.Object
	for _, s := range program.Statements {
		obj = Eval(s)

		if returnValue, ok := obj.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return obj
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var obj object.Object
	for _, statement := range block.Statements {
		obj = Eval(statement)

		if obj != nil && obj.Type() == object.RETURN_VALUE_OBJ {
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
		return NULL
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		// Compare by memory address, this works for booleans since we use the same instances TRUE & FALSE
		return nativeBoolToBoolean(left == right)
	case operator == "!=":
		return nativeBoolToBoolean(left != right)
	default:
		return NULL
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
		return NULL
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
		return NULL
	}
}

func evalIfExpression(ifExp *ast.IfExpression) object.Object {
	condition := Eval(ifExp.Condition)

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
