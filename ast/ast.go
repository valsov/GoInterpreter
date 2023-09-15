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
