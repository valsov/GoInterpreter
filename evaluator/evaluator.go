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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBoolean(node.Value)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		parameters := evalExpressions(node.Arguments, env)
		if len(parameters) == 1 && isError(parameters[0]) {
			return parameters[0]
		}

		return applyFunction(function, parameters)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var obj object.Object
	for _, s := range program.Statements {
		obj = Eval(s, env)

		switch obj := obj.(type) {
		case *object.ReturnValue:
			return obj.Value
		case *object.Error:
			return obj
		}
	}
	return obj
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var obj object.Object
	for _, statement := range block.Statements {
		obj = Eval(statement, env)

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
	case left.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
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

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	return &object.String{Value: leftValue + rightValue}
}

func evalIfExpression(ifExp *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ifExp.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTrue(condition) {
		return Eval(ifExp.Consequence, env)
	} else if ifExp.Alternative != nil {
		return Eval(ifExp.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, expression := range expressions {
		eval := Eval(expression, env)
		if isError(eval) {
			return []object.Object{eval}
		}
		result = append(result, eval)
	}
	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrObject := array.(*object.Array)
	indexValue := index.(*object.Integer).Value
	if indexValue < 0 || indexValue > int64(len(arrObject.Elements))-1 {
		// Invalid index
		return NULL
	}
	return arrObject.Elements[indexValue]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	kvPair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		// Not found
		return NULL
	}
	return kvPair.Value
}

func applyFunction(functionObj object.Object, parameters []object.Object) object.Object {
	switch function := functionObj.(type) {
	case *object.Function:
		functionEnvironment := extendFunctionEnv(function, parameters)
		eval := Eval(function.Body, functionEnvironment)
		return unwrapReturnValue(eval)
	case *object.Builtin:
		return function.Fn(parameters...)
	default:
		return newError("not a function: %s", function.Type())
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := map[object.HashKey]object.HashPair{}
	for _, nodesPair := range node.Pairs {
		key := Eval(nodesPair.Key, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(nodesPair.Value, env)
		if isError(value) {
			return value
		}

		hash := hashKey.HashKey()
		pairs[hash] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func extendFunctionEnv(function *object.Function, parameters []object.Object) *object.Environment {
	newEnv := object.NewEnclosedEnvironment(function.Env)
	for i, p := range function.Parameters {
		newEnv.Set(p.Value, parameters[i])
	}
	return newEnv
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Value
	}
	return obj
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
