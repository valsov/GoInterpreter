package parser

import (
	"fmt"

	"github.com/valsov/gointerpreter/ast"
	"github.com/valsov/gointerpreter/lexer"
	"github.com/valsov/gointerpreter/token"
)

type Parser struct {
	l            *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Init token cursors
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken() // Get from lexer
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	statement := ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// [...] todo
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &statement // Need to use a pointer here because LetStatement only implements Statement interface with a pointer receiver
}

func (p *Parser) parseReturnStatement() ast.Statement {
	statement := ast.ReturnStatement{Token: p.currentToken}
	p.nextToken()

	// [...] todo
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &statement
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if !p.peekTokenIs(t) {
		err := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
		p.errors = append(p.errors, err)
		return false
	}
	p.nextToken()
	return true
}
