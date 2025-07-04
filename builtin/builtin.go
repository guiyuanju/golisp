package builtin

import (
	"fmt"
	"golisp/expr"
	"strconv"
)

func New() []expr.Builtin {
	return []expr.Builtin{
		expr.NewBuiltin("+", plus),
		expr.NewBuiltin("print", print),
		expr.NewBuiltin("do", do),
	}
}

func do(values ...expr.Expr) (expr.Expr, bool) {
	if len(values) == 0 {
		return nil, true
	}
	return values[len(values)-1], true
}

func print(values ...expr.Expr) (expr.Expr, bool) {
	if len(values) > 1 {
		fmt.Print(values[1])
	}
	for i := 2; i < len(values); i++ {
		fmt.Print(" ", values[i])
	}
	fmt.Println()
	return nil, true
}

func plus(values ...expr.Expr) (expr.Expr, bool) {
	op := values[0]
	if len(values) == 1 {
		fmt.Println(expr.ErrorInfo("repl", op.(expr.Symbol), op.(expr.Symbol).Value, "need at least one argument"))
		return nil, false
	}

	switch value := values[1].(type) {
	case expr.Int:
		res := value.Value
		for i := 2; i < len(values); i++ {
			if v, ok := values[i].(expr.Int); ok {
				res += v.Value
			} else {
				fmt.Println(expr.ErrorInfo("repl", v, "expect int"))
			}
		}
		return expr.NewInt(res, op.Line(), op.Column()), true
	case expr.String:
		res := value.Value
		for i := 2; i < len(values); i++ {
			switch v := values[i].(type) {
			case expr.String:
				res += v.Value
			case expr.Int:
				res += strconv.Itoa(v.Value)
			default:
				fmt.Println(expr.ErrorInfo("repl", v, "expect string or int"))
				return nil, false
			}
		}
		return expr.NewString(res, op.Line(), op.Column()), true
	}

	fmt.Println(expr.ErrorInfo("repl", values[1], "unsupported operand for +: expect int or string"))
	return nil, false
}
