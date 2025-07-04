package evaluator

import (
	"fmt"
	"golisp/builtin"
	"golisp/expr"
	"strconv"
)

type Evaluator struct {
	env expr.Env
}

func New() Evaluator {
	env := expr.NewEnv()
	for _, b := range builtin.New() {
		env.Add(b.Name, b)
	}
	return Evaluator{env}
}

func (evaluator Evaluator) Eval(e expr.Expr) (expr.Expr, bool) {
	switch e := e.(type) {
	case expr.Int, expr.String:
		return e, true

	case expr.Symbol:
		if v, ok := evaluator.env.Get(e.Value); ok {
			return v, ok
		}
		fmt.Println(expr.ErrorInfo("repl", e, "undefined:", e.Value))
		return nil, false

	case expr.Def:
		value, ok := evaluator.Eval(e.Value)
		if !ok {
			fmt.Println(expr.ErrorInfo("repl", e.Value, "failed to evaluate value of def"))
			return nil, false
		}
		if !evaluator.env.Add(e.Name, value) {
			fmt.Println(expr.ErrorInfo("repl", e, "already defined:", e.Name))
			return nil, false
		}
		return nil, true

	case expr.Set:
		value, ok := evaluator.Eval(e.Value)
		if !ok {
			fmt.Println(expr.ErrorInfo("repl", e.Value, "failed to evaluate value of set"))
			return nil, false
		}
		if !evaluator.env.Set(e.Name, value) {
			fmt.Println(expr.ErrorInfo("repl", e, "undefined:", e.Name))
			return nil, false
		}
		return nil, true

	case expr.Closure:
		e.Env = evaluator.env
		// check param name duplication
		exist := map[string]bool{}
		for _, param := range e.Params {
			if exist[param] {
				fmt.Println(expr.ErrorInfo("repl", e, "parameter name must be unique"))
				return nil, false
			}
			exist[param] = true
		}
		return e, true

	case expr.List:
		if len(e.Value) == 0 {
			return e, true
		}

		operator, ok := evaluator.Eval(e.Value[0])
		if !ok {
			fmt.Println(expr.ErrorInfo("repl", e.Value[0], "failed to evaluate op in list"))
			return nil, false
		}

		switch operator := operator.(type) {
		case expr.Builtin:
			args := []expr.Expr{operator}
			for _, arg := range e.Value[1:] {
				value, ok := evaluator.Eval(arg)
				if !ok {
					fmt.Println(expr.ErrorInfo("repl", arg, "failed to evaluate argument"))
					return nil, false
				}
				args = append(args, value)
			}
			return operator.Proc(args...)

		// invoke a closure
		case expr.Closure:
			if len(e.Value)-1 != len(operator.Params) {
				fmt.Println(expr.ErrorInfo("repl", operator, "expect", strconv.Itoa(len(operator.Params)), "arguments, got", strconv.Itoa(len(e.Value)-1)))
				return nil, false
			}

			args := []expr.Expr{}
			for _, arg := range e.Value[1:] {
				value, ok := evaluator.Eval(arg)
				if !ok {
					fmt.Println(expr.ErrorInfo("repl", arg, "failed to evaluate argument"))
					return nil, false
				}
				args = append(args, value)
			}

			env := expr.NewEnv()
			for i := range operator.Params {
				env.Add(operator.Params[i], args[i])
			}
			newEnv := operator.Env.AppendEnv(env)
			newEvaluator := New()
			newEvaluator.env = newEnv
			return newEvaluator.Eval(operator.Body)
		}

		fmt.Println(expr.ErrorInfo("repl", e.Value[0], "expect proc or function"))
		return nil, false
	}

	panic("unexhaustive evaluation")
}
