package parser

import (
	"strconv"
	"strings"

	"github.com/fredrikkvalvik/nots/pkg/template/ast"
	"github.com/fredrikkvalvik/nots/pkg/template/token"
)

// advances one token and tries to parse an expression based on curToken
func (p *Parser) parseExpression(stickiness int) ast.Expr {
	// {{ ident|ident }}
	// ^
	prefix, ok := p.prefixParselets[p.curToken.Type]

	if !ok {
		p.errorf("no existing parslest for token=%s", p.curToken)
		return nil
	}

	// {{ ident|ident }}
	//    ^
	left := prefix()

	// {{ ident|ident }}
	//         ^

	for !p.peekTokenIs(token.TokenTypeRMeta) && stickiness < p.peekStickiness() {
		//   2      +     2
		//   left   op    right
		//   ^      ^peeking op
		// standing at end of prefix
		// peek next token to see if we can continue
		infix, ok := p.infixParselets[p.peekToken.Type]

		if !ok {
			return left
		}

		p.advance()
		//   2      +     2
		//   left   op    right
		//          ^
		// parse op as infix
		left = infix(left)
	}

	return left
}

func (p *Parser) parseIdentifier() ast.Expr {
	// {{ ident ... }}
	//    ^
	if !p.curTokenIs(token.TokenTypeIdentifier) {
		p.expectError(token.TokenTypeIdentifier)
		return nil
	}

	// {{ ident ... }}
	//    ^
	i := &ast.IdentifierExpr{
		Value: p.curToken.Val,
	}

	return i
}

func (p *Parser) parseNumberLiteral() ast.Expr {
	// {{ number ... }}
	//    ^
	if !p.curTokenIs(token.TokenTypeNumber) {
		p.expectError(token.TokenTypeNumber)
		return nil
	}

	numStr := p.curToken.Val
	numStr = strings.TrimSuffix(numStr, ".") // remove possible trailing '.'
	num, _ := strconv.ParseFloat(numStr, 64)

	// {{ number ... }}
	//         ^
	i := &ast.NumberLiteralExpr{
		Value: num,
	}

	return i
}

func (p *Parser) parseStringLiteral() ast.Expr {
	// {{ string ... }}
	//    ^
	if !p.curTokenIs(token.TokenTypeString) {
		p.expectError(token.TokenTypeString)
		return nil
	}

	// remove the leading and trailing '"'
	runes := []rune(p.curToken.Val)
	str := string(runes[1 : len(runes)-1])

	// {{ string ... }}
	//         ^
	i := &ast.StringLiteralExpr{
		Value: str,
	}

	return i
}

func (p *Parser) parsePipeExpression(left ast.Expr) ast.Expr {
	// {{ left | right }}
	//         ^
	pipe := &ast.PipeExpr{Left: left}
	pipeStick := p.curStickiness()

	if !p.expectCur(token.TokenTypePipe) {
		p.expectError(token.TokenTypePipe)
		return nil
	}
	// {{ left | right }}
	//           ^
	pipe.Right = p.parseExpression(pipeStick)
	return pipe
}

func (p *Parser) parseFunctionExpression(left ast.Expr) ast.Expr {
	fncall := &ast.FunctionCallExpr{Callee: left}

	fncall.Arguments = p.parseExpressionList(token.TokenTypeRParen)

	return fncall
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expr {
	list := []ast.Expr{}

	// ( expr1, expr2 )
	// ^
	if p.peekTokenIs(end) {
		p.advance()
		// ( expr1, expr2 )
		//                ^
		return list
	}

	p.advance()
	// ( expr1, expr2 )
	//   ^

	list = append(list, p.parseExpression(LOWEST))
	// ( expr1, expr2 )
	//       ^

	for p.peekTokenIs(token.TokenTypeComma) {
		p.advance()
		// ( expr1, expr2 )
		//        ^
		p.advance()
		// ( expr1, expr2 )
		//          ^
		list = append(list, p.parseExpression(LOWEST))
	}

	// ( expr1, expr2 )
	//              ^
	// error if end of list without seeing `end` token
	if !p.expectPeek(end) {
		return nil
	}
	// ( expr1, expr2 )
	//                ^

	return list
}
