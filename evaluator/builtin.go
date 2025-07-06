package evaluator

import (
	"fmt"
	"golisp/expr"
	"strconv"
)

type Proc func(e Evaluator, exprs ...expr.Expr) (expr.Expr, bool)

type Builtins map[string]Proc

func NewBuiltins() Builtins {
	res := map[string]Proc{}
	res["+"] = plus
	res["print"] = print
	res["do"] = do
	res["="] = equal
	res[":"] = _append
	res["list"] = list
	return res
}

func list(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	return expr.NewList(values...), true
}

func _append(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 2 argumtes"))
		return nil, false
	}
	target, ok := values[1].(expr.List)
	if !ok {
		fmt.Println(e.errorInfo("repl", values[1], "type mismatch:", "expect a list"))
		return nil, false
	}
	res := target.Value
	for _, v := range values[2:] {
		res = append(res, v)
	}
	return expr.NewList(res...), true
}

func equal(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 2 argumtes"))
		return nil, false
	}
	var compare func(a, b expr.Expr) bool
	compare = func(a, b expr.Expr) bool {
		// two nils equal
		if a == nil || b == nil {
			return a == b
		}
		la, ok1 := a.(expr.List)
		lb, ok2 := b.(expr.List)
		if ok1 || ok2 {
			if !(ok1 && ok2) {
				return false
			}
			if len(la.Value) != len(lb.Value) {
				return false
			}
			for i := range len(la.Value) {
				if !compare(la.Value[i], lb.Value[i]) {
					return false
				}
			}
			return true
		}
		return a == b
	}
	for i := 1; i < len(values); i++ {
		if !compare(values[i-1], values[i]) {
			return expr.NewBool(false), true
		}
	}
	return expr.NewBool(true), true
}

func do(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) == 0 {
		return nil, true
	}
	return values[len(values)-1], true
}

func print(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) > 1 {
		fmt.Print(values[1])
	}
	for i := 2; i < len(values); i++ {
		fmt.Print(" ", values[i])
	}
	fmt.Println()
	return nil, true
}

func plus(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	op := values[0]
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", op.(expr.Symbol), op.(expr.Symbol).Value, "need at least twp argument"))
		return nil, false
	}

	switch value := values[1].(type) {
	case expr.Int:
		res := value.Value
		for i := 2; i < len(values); i++ {
			if v, ok := values[i].(expr.Int); ok {
				res += v.Value
			} else {
				fmt.Println(e.errorInfo("repl", v, "expect int"))
			}
		}
		return expr.NewInt(res), true
	case expr.String:
		res := value.Value
		for i := 2; i < len(values); i++ {
			switch v := values[i].(type) {
			case expr.String:
				res += v.Value
			case expr.Int:
				res += strconv.Itoa(v.Value)
			default:
				fmt.Println(e.errorInfo("repl", v, "expect string or int"))
				return nil, false
			}
		}
		return expr.NewString(res), true
	}

	fmt.Println(e.errorInfo("repl", values[1], "unsupported operand for +: expect int or string"))
	return nil, false
}
