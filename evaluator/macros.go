package evaluator

import (
	"fmt"
	"golisp/expr"
	"strconv"
)

const (
	MACRO_AND string = "and"
)

type Macros map[string]expr.Closure

func NewMacros() Macros {
	macros := map[string]expr.Closure{}
	return macros
}

func (m Macros) isMacro(e expr.List) bool {
	if len(e.Value) == 0 {
		return false
	}
	name, ok := e.Value[0].(expr.Symbol)
	if !ok {
		return false
	}
	_, ok = m[name.Value]
	return ok
}

func (m Macros) AddMacro(name string, closure expr.Closure) {
	m[name] = closure
}

func (m Macros) expand(evaluator Evaluator, e expr.List) (expr.Expr, bool) {
	macro := m[e.Value[0].(expr.Symbol).Value]
	args := e.Value[1:]

	if len(e.Value)-1 != len(macro.Params) {
		fmt.Println(evaluator.errorInfo("repl", macro, "expect", strconv.Itoa(len(macro.Params)), "arguments, got", strconv.Itoa(len(e.Value)-1)))
		return nil, false
	}

	return apply(evaluator, macro, args)
}

// func macroAnd(e Evaluator, exprs ...expr.Expr) (expr.Expr, bool) {
// 	if len(exprs) < 3 {
// 		fmt.Println(e.errorInfo("repl", exprs[0], "expect at least 2 arguments"))
// 		return nil, false
// 	}
// 	var args []expr.Expr
// 	args = append(args, expr.NewSymbol(expr.SF_IF))
// 	args = append(args, exprs[1])
// 	args = append(args, exprs[2])
// 	args = append(args, expr.NewBool(false))
// 	return expr.NewList(args...), true
// }
