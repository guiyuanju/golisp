package evaluator

import (
	"fmt"
	"strconv"
	"time"

	"github.com/guiyuanju/golisp/expr"
)

type Proc func(e Evaluator, exprs ...expr.Expr) (expr.Expr, bool)

type Builtins map[string]Proc

var RegisteredBuiltins Builtins = Builtins{}

func RegisterBuiltin(name string, proc Proc) {
	RegisteredBuiltins[name] = proc
}

func RegisterDefaultBuiltins() {
	RegisteredBuiltins["+"] = plus
	RegisteredBuiltins["-"] = minus
	RegisteredBuiltins["*"] = multiply
	RegisteredBuiltins["/"] = divide
	RegisteredBuiltins["print"] = print
	RegisteredBuiltins["do"] = do
	RegisteredBuiltins["="] = equal
	RegisteredBuiltins[">"] = greater
	RegisteredBuiltins["<"] = less
	RegisteredBuiltins["<="] = lessEqual
	RegisteredBuiltins[">="] = greaterEqual
	RegisteredBuiltins["append"] = _append
	RegisteredBuiltins[":"] = slice
	RegisteredBuiltins["list"] = list
	RegisteredBuiltins["not"] = not
	RegisteredBuiltins["type"] = _type
	RegisteredBuiltins["macroexpand"] = macroexpand
	RegisteredBuiltins["time"] = _time
	RegisteredBuiltins["."] = dot
	RegisteredBuiltins["len"] = length
	RegisteredBuiltins["eval"] = eval
}

func eval(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	return e.Eval(values[1])
}

func slice(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 4 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need 3 arguments"))
		return nil, false
	}
	start, ok := values[1].(expr.Number)
	if !ok {
		fmt.Println(e.errorInfo("repl", values[0], "expect int"))
		return nil, false
	}
	end, ok := values[2].(expr.Number)
	if !ok {
		fmt.Println(e.errorInfo("repl", values[0], "expect int"))
		return nil, false
	}
	switch seq := values[3].(type) {
	case expr.List:
		startIdx := int(start.Value)
		if startIdx < 0 {
			startIdx += len(seq.Value)
		}
		if startIdx < 0 || startIdx > len(seq.Value) {
			fmt.Println(e.errorInfo("repl", values[1], fmt.Sprintf("index %d out of bound %d", startIdx, len(seq.Value))))
			return nil, false
		}
		endIdx := int(end.Value)
		if endIdx < 0 {
			endIdx += len(seq.Value)
		}
		if endIdx < 0 || endIdx > len(seq.Value) {
			fmt.Println(e.errorInfo("repl", values[1], fmt.Sprintf("index %d out of bound %d", startIdx, len(seq.Value))))
			return nil, false
		}
		if startIdx > endIdx {
			fmt.Println(e.errorInfo("repl", values[1], "start is greater than end"))
			return nil, false
		}
		return expr.NewList(seq.Value[startIdx:endIdx]...), true

	default:
		fmt.Println(e.errorInfo("repl", values[3], "expect vector or list"))
		return nil, false
	}
}

func formalizeIndex(idx int, length int) (int, bool) {
	index := idx
	if idx < 0 {
		index = idx + length
	}
	if index < 0 || index >= length {
		return index, false
	}
	return index, true
}

func length(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 2 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need 1 argument"))
		return nil, false
	}
	switch seq := values[1].(type) {
	case expr.List:
		return expr.NewNum(float64(len(seq.Value))), true
	default:
		fmt.Println(e.errorInfo("repl", values[1], "unsupported type for len"))
		return nil, false
	}
}

func dot(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need 2 arguments"))
		return nil, false
	}
	switch seq := values[2].(type) {
	case expr.List:
		v, ok := values[1].(expr.Number)
		if !ok {
			fmt.Println(e.errorInfo("repl", values[2], "expect int"))
			return nil, false
		}
		idx, ok := formalizeIndex(int(v.Value), len(seq.Value))
		if !ok {
			fmt.Println(e.errorInfo("repl", values[1], fmt.Sprintf("index %d out of bound %d", idx, len(seq.Value))))
			return nil, false
		}
		return seq.Value[idx], true
	default:
		fmt.Println(e.errorInfo("repl", values[2], "unsupported type for dot"))
		return nil, false
	}
}

func multiply(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 2 argumte"))
		return nil, false
	}
	var res float64 = 1
	for _, v := range values[1:] {
		n, ok := v.(expr.Number)
		if !ok {
			fmt.Println(e.errorInfo("repl", values[0], "expect int"))
			return nil, false
		}
		res *= n.Value
	}
	return expr.NewNum(res), true
}

func divide(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 2 argumte"))
		return nil, false
	}
	v, ok := values[1].(expr.Number)
	if !ok {
		fmt.Println(e.errorInfo("repl", values[0], "expect int"))
		return nil, false
	}
	res := v.Value
	for _, v := range values[2:] {
		n, ok := v.(expr.Number)
		if !ok {
			fmt.Println(e.errorInfo("repl", values[0], "expect int"))
			return nil, false
		}
		res /= n.Value
	}
	return expr.NewNum(res), true
}

func _time(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	nano := time.Now().UnixNano()
	return expr.NewNum(float64(nano)), true
}

func macroexpand(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 2 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need 1 argumte"))
		return nil, false
	}
	arg, ok := values[1].(expr.List)
	if !ok {
		fmt.Println(e.errorInfo("repl", values[1], "expext argument to be a quoted list"))
		return nil, false
	}
	if !e.isMacro(arg.Value[0]) {
		fmt.Println(e.errorInfo("repl", arg.Value[0], "not macro"))
		return nil, false
	}
	return e.macroExpand(arg)
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
	switch target := values[1].(type) {
	case expr.List:
		res := target.Value
		for _, v := range values[2:] {
			res = append(res, v)
		}
		return expr.NewList(res...), true
	default:
		fmt.Println(e.errorInfo("repl", values[1], "type mismatch:", "expect a list or vector"))
		return nil, false
	}
}

func equal(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 2 argumtes"))
		return nil, false
	}
	for i := 2; i < len(values); i++ {
		if !values[i-1].Equal(values[i]) {
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
		return nil, false
	}
	return not(e, []expr.Expr{values[0], res}...)
}

func greater(e Evaluator, values ...expr.Expr) (expr.Expr, bool) {
	if len(values) < 3 {
		fmt.Println(e.errorInfo("repl", values[0], "arity mismatch:", "need at least 2 argumtes"))
		return nil, false
	}
	switch value := values[1].(type) {
	case expr.Number:
		prev := value.Value
		for i := 2; i < len(values); i++ {
			if v, ok := values[i].(expr.Number); ok {
				if v.Value >= prev {
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
				if v.Value >= prev {
					return expr.NewBool(false), true
				}
				prev = v.Value
			case expr.Number:
				str := strconv.FormatFloat(v.Value, 'f', -1, 64)
				if str >= prev {
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
	case expr.Number:
		res := value.Value
		for i := 2; i < len(values); i++ {
			if v, ok := values[i].(expr.Number); ok {
				res += v.Value
			} else {
				fmt.Println(e.errorInfo("repl", v, "expect int"))
				return nil, false
			}
		}
		return expr.NewNum(res), true
	case expr.String:
		res := value.Value
		for i := 2; i < len(values); i++ {
			switch v := values[i].(type) {
			case expr.String:
				res += v.Value
			case expr.Number:
				res += strconv.FormatFloat(v.Value, 'f', -1, 64)
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
	if len(values) < 2 {
		fmt.Println(e.errorInfo("repl", op, "need at least one argument"))
		return nil, false
	}

	switch values[1].(type) {
	case expr.Number:
		if len(values) == 1 {
			return expr.NewNum(-values[1].(expr.Number).Value), true
		}
		res := values[1].(expr.Number).Value
		for i := 2; i < len(values); i++ {
			if v, ok := values[i].(expr.Number); ok {
				res -= v.Value
			} else {
				fmt.Println(e.errorInfo("repl", v, "expect int"))
				return nil, false
			}
		}
		return expr.NewNum(res), true
	}

	fmt.Println(e.errorInfo("repl", values[1], "unsupported operand for -: expect int"))
	return nil, false
}
