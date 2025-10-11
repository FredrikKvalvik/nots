package parser

import (
	"errors"
	"fmt"

	"github.com/fredrikkvalvik/nots/pkg/template/ast"
	"github.com/fredrikkvalvik/nots/pkg/template/lexer"
	"github.com/fredrikkvalvik/nots/pkg/template/token"
)

type prefixFn = func() ast.Expr
type infixFn = func(left ast.Expr) ast.Expr

const (
	_ int = iota
	LOWEST
	PIPE   // |
	CONCAT // .

	// OR          // or
	// AND         // and
	// EQUALS      // ==
	// LESSGREATER // > or <
	// SUM         //+ -
	// PRODUCT     //* /
	// PREFIX      //-X or !X
	CALL // myFunction(X)
)

var stickinessMap = map[token.TokenType]int{
	token.TokenTypePipe:   PIPE,
	token.TokenTypeDot:    CONCAT,
	token.TokenTypeLParen: CALL,
}

type stateFn func(p *Parser) stateFn

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	infixParselets  map[token.TokenType]infixFn
	prefixParselets map[token.TokenType]prefixFn

	template ast.Template

	errors []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:               l,
		infixParselets:  map[token.TokenType]infixFn{},
		prefixParselets: map[token.TokenType]prefixFn{},
	}
	// prepare curToken and peekToken
	p.advance()
	p.advance()

	// PREFIX PARSLETS
	p.registerPrefix(token.TokenTypeIdentifier, p.parseIdentifier)
	p.registerPrefix(token.TokenTypeNumber, p.parseNumberLiteral)
	p.registerPrefix(token.TokenTypeString, p.parseStringLiteral)

	// INFIX PARSLETS
	p.registerInfix(token.TokenTypePipe, p.parsePipeExpression)
	p.registerInfix(token.TokenTypeDot, p.parseBinaryExpression)
	p.registerInfix(token.TokenTypeLParen, p.parseFunctionExpression)

	return p
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixFn) {
	p.prefixParselets[t] = fn
}
func (p *Parser) registerInfix(t token.TokenType, fn infixFn) {
	p.infixParselets[t] = fn
}

func (p *Parser) Parse() (*ast.Template, error) {
	for state := parseText; state != nil; {
		state = state(p)
	}

	if len(p.errors) > 0 {
		return nil, errors.Join(p.errors...)
	}

	return &p.template, nil
}

// add a new block to the template
func (p *Parser) emitBlock(n ast.Block) {
	p.template.Blocks = append(p.template.Blocks, n)
}

// consume current token
func (p *Parser) advance() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekStickiness() int {
	if s, ok := stickinessMap[p.peekToken.Type]; ok {
		return s
	}
	return LOWEST
}

func (p *Parser) curStickiness() int {
	if s, ok := stickinessMap[p.curToken.Type]; ok {
		return s
	}
	return LOWEST
}

func (p *Parser) curTokenIs(typ token.TokenType) bool {
	return p.curToken.Type == typ
}

func (p *Parser) peekTokenIs(typ token.TokenType) bool {
	return p.peekToken.Type == typ
}

// when true, advance, making curToken.Type == typ
// stay if false
func (p *Parser) expectPeek(typ token.TokenType) bool {
	if p.peekTokenIs(typ) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) expectCur(typ token.TokenType) bool {
	if p.curTokenIs(typ) {
		p.advance()
		return true
	}
	return false
}
func (p *Parser) errorf(msg string, v ...any) {
	p.errors = append(p.errors, fmt.Errorf(msg, v...))
}

// reports an error the the curToken is not the expected token
func (p *Parser) expectError(t token.TokenType) {
	p.errorf("expected %s, got %s", t, p.curToken)
}

func parseText(p *Parser) stateFn {
	for {
		switch t := p.curToken.Type; t {
		case token.TokenTypeEOF:
			return nil
		case token.TokenTypeError:
			p.errorf("lexing error: %s", p.curToken.Val)
			// TODO: report the error somewhere
			return nil

		case token.TokenTypeLMeta:
			// expects the current
			return parseExpression

		case token.TokenTypeText:
			p.emitBlock(&ast.BlockText{Text: p.curToken.Val})
			p.advance()

		default:
			p.errorf("unexpected token: %s", p.curToken)
			return nil
		}

	}
}

// expects the current token to be LeftMeta
func parseExpression(p *Parser) stateFn {
	if !p.expectCur(token.TokenTypeLMeta) {
		p.expectError(token.TokenTypeLMeta)
		return nil
	}
	// parse expression block
	expr := p.parseExpression(LOWEST)

	p.emitBlock(&ast.BlockExpression{
		Expression: expr,
	})

	// move over to what is expected to be RightMeta
	p.advance()

	// TODO: allow the for recovery by running to end of the expression and continue parsing to collect more errors

	// // expect RightMeta, anything else is an error
	if !p.expectCur(token.TokenTypeRMeta) {
		p.expectError(token.TokenTypeRMeta)
		return nil
	}

	return parseText
}
