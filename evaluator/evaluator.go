package evaluator

import (
	"fmt"
	"golisp/expr"
	"golisp/parser"
	"strconv"
	"strings"
)

type Evaluator struct {
	env       expr.Env
	Positions parser.Positions
	builtins  Builtins
}

func New(positions parser.Positions) Evaluator {
	env := expr.NewEnv()
	builtins := NewBuiltins()
	for name, _ := range builtins {
		env.Add(name, expr.NewBuiltin(name))
	}
	return Evaluator{env, positions, builtins}
}

func (e Evaluator) errorInfo(file string, expr expr.Expr, info ...string) string {
	return fmt.Sprintf("%s:%d:%d: (%s) %s", file, e.Positions[expr.ExprId()].Line, e.Positions[expr.ExprId()].Column, fmt.Sprintf("%T", expr), strings.Join(info, " "))
}

func (evaluator Evaluator) Eval(e expr.Expr) (expr.Expr, bool) {
	switch e := e.(type) {
	case expr.Int, expr.String, expr.Bool, nil:
		return e, true

	case expr.Symbol:
		if v, ok := evaluator.env.Get(e.Value); ok {
			return v, ok
		}
		fmt.Println(evaluator.errorInfo("repl", e, "undefined:", e.Value))
		return nil, false

	case expr.Nil:
		return e, true

	case expr.Quote:
		return e.Value, true

	case expr.Def:
		value, ok := evaluator.Eval(e.Value)
		if !ok {
			fmt.Println(evaluator.errorInfo("repl", e.Value, "failed to evaluate value of def"))
			return nil, false
		}
		if !evaluator.env.Add(e.Name, value) {
			fmt.Println(evaluator.errorInfo("repl", e, "already defined:", e.Name))
			return nil, false
		}
		return nil, true

	case expr.Set:
		value, ok := evaluator.Eval(e.Value)
		if !ok {
			fmt.Println(evaluator.errorInfo("repl", e.Value, "failed to evaluate value of set"))
			return nil, false
		}
		if !evaluator.env.Set(e.Name, value) {
			fmt.Println(evaluator.errorInfo("repl", e, "undefined:", e.Name))
			return nil, false
		}
		return nil, true

	case expr.If:
		pred, ok := evaluator.Eval(e.Pred)
		if !ok {
			fmt.Println(evaluator.errorInfo("repl", e.Pred, "failed to evaluate if pred"))
			return nil, false
		}
		if isTruthy(pred) {
			return evaluator.Eval(e.Then)
		} else {
			return evaluator.Eval(e.Else)
		}

	case expr.Closure:
		e.Env = evaluator.env
		// check param name duplication
		exist := map[string]bool{}
		for _, param := range e.Params {
			if exist[param] {
				fmt.Println(evaluator.errorInfo("repl", e, "parameter name must be unique"))
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
			fmt.Println(evaluator.errorInfo("repl", e.Value[0], "failed to evaluate op in list"))
			return nil, false
		}

		switch operator := operator.(type) {
		case expr.Builtin:
			args := []expr.Expr{operator}
			for _, arg := range e.Value[1:] {
				value, ok := evaluator.Eval(arg)
				if !ok {
					fmt.Println(evaluator.errorInfo("repl", arg, "failed to evaluate argument"))
					return nil, false
				}
				args = append(args, value)
			}
			proc, ok := evaluator.builtins[operator.Name]
			if !ok {
				panic("builtin not found")
			}
			return proc(evaluator, args...)

		// invoke a closure
		case expr.Closure:
			if len(e.Value)-1 != len(operator.Params) {
				fmt.Println(evaluator.errorInfo("repl", operator, "expect", strconv.Itoa(len(operator.Params)), "arguments, got", strconv.Itoa(len(e.Value)-1)))
				return nil, false
			}

			args := []expr.Expr{}
			for _, arg := range e.Value[1:] {
				value, ok := evaluator.Eval(arg)
				if !ok {
					fmt.Println(evaluator.errorInfo("repl", arg, "failed to evaluate argument"))
					return nil, false
				}
				args = append(args, value)
			}

			env := expr.NewEnv()
			for i := range operator.Params {
				env.Add(operator.Params[i], args[i])
			}
			newEnv := operator.Env.AppendEnv(env)
			newEvaluator := New(evaluator.Positions)
			newEvaluator.env = newEnv

			var last expr.Expr
			for _, b := range operator.Body {
				v, ok := newEvaluator.Eval(b)
				if !ok {
					return nil, false
				}
				last = v
			}
			return last, true
		}

		fmt.Println(evaluator.errorInfo("repl", e.Value[0], "expect proc or function"))
		return nil, false
	}

	panic("unexhaustive evaluation")
}

func isTruthy(e expr.Expr) bool {
	switch e := e.(type) {
	case expr.Bool:
		return e.Value
	case expr.Nil:
		return false
	default:
		return true
	}
}
