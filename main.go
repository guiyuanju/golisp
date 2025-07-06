package main

import (
	"bufio"
	"fmt"
	"golisp/evaluator"
	"golisp/parser"
	"os"
)

func main() {
	repl()
}

// func repl() {
// 	exprs := []expr.Expr{
// 		expr.NewDef("count", expr.NewInt(0)),
// 		expr.NewDef("plus-two", expr.NewClosure(nil, []string{"x"},
// 			[]expr.Expr{
// 				expr.NewSet("count", expr.NewList(expr.NewSymbol("+"), expr.NewSymbol("count"), expr.NewInt(1))),
// 				expr.NewList(expr.NewSymbol("+"), expr.NewSymbol("x"), expr.NewInt(2)),
// 			},
// 		)),
// 		expr.NewIf(nil,
// 			expr.NewList(expr.NewSymbol("plus-two"), expr.NewInt(200)),
// 			expr.NewList(expr.NewSymbol("plus-two"), expr.NewInt(400)),
// 		),
// 		expr.NewList(expr.NewSymbol("print"), expr.NewSymbol("count")),
// 	}

// 	e := evaluator.New(parser.NewPositions())
// 	for _, expr := range exprs {
// 		v, ok := e.Eval(expr)
// 		if !ok {
// 			break
// 		}
// 		fmt.Println(v)
// 	}

// }

func repl() {
	e := evaluator.New(nil)
	scanner := bufio.NewScanner(os.Stdin)
	hasInput := true
	for hasInput {
		fmt.Print("> ")
		hasInput = scanner.Scan()

		line := scanner.Text()

		s := parser.NewScanner(line)
		tokens, ok := s.Scan()
		if !ok {
			continue
		}

		p := parser.New(tokens)
		exprs, ok := p.Parse()
		if !ok {
			continue
		}

		e.Positions = p.Positions
		res, ok := e.Eval(exprs)
		if !ok {
			continue
		}
		fmt.Println(res)
	}
}
