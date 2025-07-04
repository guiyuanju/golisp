package main

import (
	"fmt"
	"golisp/evaluator"
	"golisp/expr"
)

func main() {
	repl()
}

func repl() {
	exprs := []expr.Expr{
		expr.NewDef("count", expr.NewInt(0)),
		expr.NewDef("plus-two", expr.NewClosure(nil, []string{"x"},
			expr.NewList(
				expr.NewSymbol("do"),
				expr.NewSet("count", expr.NewList(expr.NewSymbol("+"), expr.NewSymbol("count"), expr.NewInt(1))),
				expr.NewList(expr.NewSymbol("+"), expr.NewSymbol("x"), expr.NewInt(2)),
			),
		)),
		expr.NewList(expr.NewSymbol("plus-two"), expr.NewInt(200)),
		expr.NewList(expr.NewSymbol("print"), expr.NewSymbol("count")),
	}

	e := evaluator.New()
	for _, expr := range exprs {
		v, ok := e.Eval(expr)
		if !ok {
			break
		}
		fmt.Println(v)
	}

	// for {
	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	for scanner.Scan() {
	// 		line := scanner.Text()
	// 		_ = parser.New(line)
	// 		res, ok := e.Eval(l)
	// 		fmt.Println(res, ok)
	// 	}
	// }
}
