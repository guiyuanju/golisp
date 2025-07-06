package parser

import (
	"fmt"
	"golisp/expr"
)

type Parser struct {
	i         int
	tokens    []Token
	Positions Positions
}

func New(tokens []Token) Parser {
	return Parser{
		i:         0,
		tokens:    tokens,
		Positions: NewPositions(),
	}
}

func (p *Parser) cur() Token {
	return p.tokens[p.i]
}

func (p *Parser) advance() {
	p.i++
}

func (p *Parser) isEnd() bool {
	return p.i >= len(p.tokens)
}

func (p *Parser) withPosOfToken(expr expr.Expr, token Token) expr.Expr {
	pos := Position{token.Line, token.Column}
	p.Positions[expr.ExprId()] = pos
	return expr
}

func (p *Parser) expr() (expr.Expr, bool) {
	if p.isEnd() {
		fmt.Println(errorInfo(p.previous(), "expect a expr after it"))
		return nil, false
	}
	cur := p.cur()
	switch cur.TokenType {
	case INTEGER:
		p.advance()
		return p.withPosOfToken(expr.NewInt(cur.Value.(int)), cur), true
	case STRING:
		p.advance()
		return p.withPosOfToken(expr.NewString(cur.Value.(string)), cur), true
	case TRUE:
		p.advance()
		return p.withPosOfToken(expr.NewBool(true), cur), true
	case FALSE:
		p.advance()
		return p.withPosOfToken(expr.NewBool(false), cur), true
	case SYMBOL:
		p.advance()
		return p.withPosOfToken(expr.NewSymbol(cur.Value.(string)), cur), true
	case QUOTE:
		p.advance()
		v, ok := p.expr()
		if !ok {
			return nil, false
		}
		return expr.NewList(expr.NewSymbol(expr.SF_QUOTE), v), true
	case LEFT_BRACKET:
		res, ok := p.vector()
		if !ok {
			return nil, false
		}
		_, ok = p.consume(RIGHT_BRACKET)
		if !ok {
			return nil, false
		}
		return res, true
	case LEFT_PAREN:
		res, ok := p.list()
		if !ok {
			return nil, false
		}
		_, ok = p.consume(RIGHT_PAREN)
		if !ok {
			return nil, false
		}
		return res, true
	}
	fmt.Println(errorInfo(cur, "unexpected token"))
	return nil, false
}

func (p *Parser) list() (expr.Expr, bool) {
	p.advance()
	cur := p.cur()
	var res []expr.Expr
	for !p.isEnd() && p.cur().TokenType != RIGHT_PAREN {
		expr, ok := p.expr()
		if !ok {
			return nil, false
		}
		res = append(res, expr)
	}
	return p.withPosOfToken(expr.NewList(res...), cur), true
}

func (p *Parser) vector() (expr.Expr, bool) {
	p.advance()
	cur := p.cur()
	var res []expr.Expr
	for !p.isEnd() && p.cur().TokenType != RIGHT_BRACKET {
		expr, ok := p.expr()
		if !ok {
			return nil, false
		}
		res = append(res, expr)
	}
	return p.withPosOfToken(expr.NewVector(res...), cur), true
}

func (p *Parser) consume(tokenType TokenType) (Token, bool) {
	if p.isEnd() {
		fmt.Println(errorInfo(p.previous(), "unexpected end"))
		return Token{}, false
	}
	if p.cur().TokenType != tokenType {
		fmt.Println(errorInfo(p.cur(), "unexpected token"))
		return Token{}, false
	}
	cur := p.cur()
	p.advance()
	return cur, true
}

func (p *Parser) previous() Token {
	if p.i <= 0 {
		return Token{Line: 0, Column: 0}
	}
	return p.tokens[p.i-1]
}

func errorInfo(token Token, info string) string {
	return fmt.Sprintf("repl:%d:%d: %s", token.Line, token.Column, info)
}

func (p *Parser) Parse() (expr.Expr, bool) {
	var last expr.Expr
	for !p.isEnd() {
		res, ok := p.expr()
		if !ok {
			return last, ok
		}
		last = res
	}
	return last, true
}
