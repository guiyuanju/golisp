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

// expr = int | string | bool | symbol | var | set | if | fn | list
// var = "(" "var" symbol expr ")"
// set = "(" "set" symbol expr ")"
// if = "(" "if" expr expr expr? ")"
// fn = "(" "fn" symbol? "[" symbol* "]" expr* ")"
// list = "(" expr* ")"

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
	case LEFT_PAREN:
		return p.list()
	}
	panic("unreachable")
}

func (p *Parser) list() (expr.Expr, bool) {
	p.advance()
	cur := p.cur()
	switch cur.TokenType {
	case VAR:
		p.advance()
		t, ok := p.consume(SYMBOL)
		if !ok {
			return nil, false
		}
		name := t.Value.(string)

		value, ok := p.expr()
		if !ok {
			return nil, false
		}
		_, ok = p.consume(RIGHT_PAREN)
		if !ok {
			return nil, false
		}
		return p.withPosOfToken(expr.NewDef(name, value), cur), true

	case SET:
		p.advance()
		t, ok := p.consume(SYMBOL)
		if !ok {
			return nil, false
		}
		name := t.Value.(string)

		value, ok := p.expr()
		if !ok {
			return nil, false
		}
		_, ok = p.consume(RIGHT_PAREN)
		if !ok {
			return nil, false
		}
		return p.withPosOfToken(expr.NewSet(name, value), cur), true

	case IF:
		p.advance()
		pred, ok := p.expr()
		if !ok {
			return nil, false
		}
		thenBranch, ok := p.expr()
		if !ok {
			return nil, false
		}
		elseBranch, ok := p.expr()
		if !ok {
			return nil, false
		}
		_, ok = p.consume(RIGHT_PAREN)
		if !ok {
			return nil, false
		}
		return p.withPosOfToken(expr.NewIf(pred, thenBranch, elseBranch), cur), true

	case FN:
		p.advance()
		if p.isEnd() {
			fmt.Println(errorInfo(cur, "expect a symbol or param list after it"))
			return nil, false
		}
		name := ""
		if p.cur().TokenType != LEFT_BRACKET {
			if p.cur().TokenType != SYMBOL {
				fmt.Println(errorInfo(cur, "expect a symbol"))
				return nil, false
			}
			name = p.cur().Value.(string)
			p.advance()
		}

		_, ok := p.consume(LEFT_BRACKET)
		if !ok {
			return nil, false
		}

		params := []string{}
		for p.cur().TokenType != RIGHT_BRACKET {
			t, ok := p.consume(SYMBOL)
			if !ok {
				return nil, false
			}
			params = append(params, t.Value.(string))
		}

		_, ok = p.consume(RIGHT_BRACKET)
		if !ok {
			return nil, false
		}

		body := []expr.Expr{}
		for !p.isEnd() && p.cur().TokenType != RIGHT_PAREN {
			b, ok := p.expr()
			if !ok {
				return nil, false
			}
			body = append(body, b)
		}

		_, ok = p.consume(RIGHT_PAREN)
		if !ok {
			return nil, false
		}

		if name != "" {
			return expr.NewDef(name, expr.NewClosure(nil, params, body)), true
		}
		return p.withPosOfToken(expr.NewClosure(nil, params, body), cur), true

	default:
		var res []expr.Expr
		for !p.isEnd() && p.cur().TokenType != RIGHT_PAREN {
			expr, ok := p.expr()
			if !ok {
				return nil, false
			}
			res = append(res, expr)
		}
		_, ok := p.consume(RIGHT_PAREN)
		if !ok {
			return nil, false
		}
		return p.withPosOfToken(expr.NewList(res...), cur), true
	}
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
	return p.expr()
}
