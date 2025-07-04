package parser

import "golisp/expr"

type Parser struct{}

func New(s string) Parser {
	return Parser{}
}

func (p *Parser) Parse() expr.Expr {
	return nil
}
