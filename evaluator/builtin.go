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
	res["-"] = minus
	res["print"] = print
	res["do"] = do
	res["="] = equal
	res[">"] = greater
	res["<"] = less
	res["<="] = lessEqual
	res[">="] = greaterEqual
	res[":"] = _append
	res["list"] = list
	res["not"] = not
	res["type"] = _type
	return res
}

func _type(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 2 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 1 argumte"))
		return nil, false
	}
	return expr.NewString(values[1].ExprName()), true
}

func list(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	return expr.NewList(values[1:]...), true
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

func not(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 2 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 1 argumtes"))
		return nil, false
	}
	if isTruthy(values[1]) {
		return expr.NewBool(false), true
	}
	return expr.NewBool(true), true
}

func less(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	gt, ok := greater(e, values...)
	if !ok {
		return gt, ok
	}
	eq, ok := equal(e, values...)
	if !ok {
		return eq, ok
	}
	return expr.NewBool(!gt.(expr.Bool).Value && !eq.(expr.Bool).Value), true
}

func greaterEqual(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	eq, ok := equal(e, values...)
	if !ok {
		return eq, ok
	}
	if eq.(expr.Bool).Value {
		return expr.NewBool(true), true
	}
	gt, ok := greater(e, values...)
	if !ok {
		return gt, ok
	}
	if gt.(expr.Bool).Value {
		return expr.NewBool(true), true
	}
	return expr.NewBool(false), true
}

func lessEqual(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	res, ok := greater(e, values...)
	if !ok {
		return res, ok
	}
	return not(e, res)
}

func greater(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 2 argumtes"))
		return nil, false
	}
	switch value := values[1].(type) {
	case expr.Int:
		prev := value.Value
		for i := 2; i < len(values); i++ {
			if v, ok := values[i].(expr.Int); ok {
				if v.Value <= prev {
					return expr.NewBool(false), true
				}
				prev = v.Value
			} else {
				fmt.Println(e.errorInfo("repl", v, "expect int"))
				return nil, false
			}
		}
		return expr.NewBool(true), true
	case expr.String:
		prev := value.Value
		for i := 2; i < len(values); i++ {
			switch v := values[i].(type) {
			case expr.String:
				if v.Value <= prev {
					return expr.NewBool(false), true
				}
				prev = v.Value
			case expr.Int:
				str := strconv.Itoa(v.Value)
				if str <= prev {
					return expr.NewBool(false), true
				}
				prev = str
			default:
				fmt.Println(e.errorInfo("repl", v, "expect string or int"))
				return nil, false
			}
		}
		return expr.NewBool(true), true
	}
	fmt.Println(e.errorInfo("repl", values[1], "expect string or int"))
	return nil, false
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
		fmt.Println(e.errorInfo("repl", op, "need at least two argument"))
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
				return nil, false
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

func minus(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	op := values[0]
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", op, "need at least two argument"))
		return nil, false
	}

	switch value := values[1].(type) {
	case expr.Int:
		res := value.Value
		for i := 2; i < len(values); i++ {
			if v, ok := values[i].(expr.Int); ok {
				res -= v.Value
			} else {
				fmt.Println(e.errorInfo("repl", v, "expect int"))
				return nil, false
			}
		}
		return expr.NewInt(res), true
	}

	fmt.Println(e.errorInfo("repl", values[1], "unsupported operand for -: expect int"))
	return nil, false
}
