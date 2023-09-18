package ast

import (
	"fmt"
	"strings"

	"github.com/valsov/gointerpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expresionNode()
}

// Root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	sb := strings.Builder{}
	for _, s := range p.Statements {
		sb.WriteString(s.String())
	}
	return sb.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) <= 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%v %v = ", ls.TokenLiteral(), ls.Name.String()))
	if ls.Value != nil {
		sb.WriteString(ls.Value.String())
	}
	sb.WriteRune(';')
	return sb.String()
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) String() string {
	return id.Value
}
func (id *Identifier) expresionNode() {}
func (id *Identifier) TokenLiteral() string {
	return id.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	sb := strings.Builder{}
	sb.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		sb.WriteString(rs.ReturnValue.String())
	}
	sb.WriteRune(';')
	return sb.String()
}
func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) expresionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // Operator token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}
func (pe *PrefixExpression) expresionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

type InfixExpression struct {
	Token    token.Token // Operator token
	Operator string
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}
func (ie *InfixExpression) expresionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) String() string {
	return b.Token.Literal
}
func (b *Boolean) expresionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       token.Token // if
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("if %s %s", ie.Condition.String(), ie.Consequence.String()))

	if ie.Alternative != nil {
		sb.WriteString(fmt.Sprintf("else %s", ie.Alternative.String()))
	}
	return sb.String()
}
func (ie *IfExpression) expresionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

type BlockStatement struct {
	Token      token.Token // {
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	sb := strings.Builder{}
	for _, s := range bs.Statements {
		sb.WriteString(s.String())
	}
	return sb.String()
}
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

type FunctionLiteral struct {
	Token      token.Token // fn
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) String() string {
	parametersStr := make([]string, len(fl.Parameters))
	for i, s := range fl.Parameters {
		parametersStr[i] = s.String()
	}
	return fmt.Sprintf("%s(%s) %s", fl.TokenLiteral(), strings.Join(parametersStr, ", "), fl.Body.String())
}
func (fl *FunctionLiteral) expresionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

type CallExpression struct {
	Token     token.Token // (
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) String() string {
	args := make([]string, len(ce.Arguments))
	for i, arg := range ce.Arguments {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", ce.Function.String(), strings.Join(args, ", "))
}
func (ce *CallExpression) expresionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
