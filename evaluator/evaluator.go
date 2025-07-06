package evaluator

import (
	"fmt"
	"golisp/expr"
	"golisp/parser"
	"strconv"
	"strings"
)

const (
	SF_QUOTE = "quote"
	SF_VAR   = "var"
	SF_SET   = "set"
	SF_IF    = "if"
	SF_FN    = "fn"
	SF_MACRO = "macro"
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

func isSpecialForm(e expr.List) bool {
	if len(e.Value) == 0 {
		return false
	}
	s, ok := e.Value[0].(expr.Symbol)
	if !ok {
		return false
	}
	switch s.Value {
	case SF_QUOTE, SF_VAR, SF_SET, SF_IF, SF_FN, SF_MACRO:
		return true
	default:
		return false
	}
}

func (evaluator Evaluator) evalSpecialForm(e expr.List) (expr.Expr, bool) {
	s := e.Value[0].(expr.Symbol)
	switch s.Value {
	case SF_QUOTE:
		if len(e.Value) != 2 {
			fmt.Println(evaluator.errorInfo("repl", s, "expect 1 argument"))
			return nil, false
		}
		return e.Value[1], true

	case SF_VAR:
		if len(e.Value) < 3 {
			fmt.Println(evaluator.errorInfo("repl", s, "expect 2 arguments"))
			return nil, false
		}
		name, ok := e.Value[1].(expr.Symbol)
		if !ok {
			fmt.Println(evaluator.errorInfo("repl", name, "expect symbol"))
			return nil, false
		}
		value, ok := evaluator.Eval(e.Value[2])
		if !ok {
			return nil, false
		}
		if !evaluator.env.Add(name.Value, value) {
			fmt.Println(evaluator.errorInfo("repl", name, "already defined:", name.Value))
			return nil, false
		}
		return expr.NewNil(), true

	case SF_SET:
		if len(e.Value) < 3 {
			fmt.Println(evaluator.errorInfo("repl", s, "expect 2 arguments"))
			return nil, false
		}
		name, ok := e.Value[1].(expr.Symbol)
		if !ok {
			fmt.Println(evaluator.errorInfo("repl", name, "expect symbol"))
			return nil, false
		}
		value, ok := evaluator.Eval(e.Value[2])
		if !ok {
			return nil, false
		}
		if !evaluator.env.Set(name.Value, value) {
			fmt.Println(evaluator.errorInfo("repl", name, "already defined:", name.Value))
			return nil, false
		}
		return expr.NewNil(), true

	case SF_IF:
		if len(e.Value) < 3 {
			fmt.Println(evaluator.errorInfo("repl", s, "expect at least two arguments"))
			return nil, false
		}
		pred, ok := evaluator.Eval(e.Value[1])
		if !ok {
			return nil, false
		}
		if isTruthy(pred) {
			return evaluator.Eval(e.Value[2])
		}
		if len(e.Value) < 4 {
			return expr.NewNil(), true
		}
		return evaluator.Eval(e.Value[3])

	case SF_FN:
		if len(e.Value) < 2 {
			fmt.Println(evaluator.errorInfo("repl", s, "expect an argument list"))
			return nil, false
		}
		switch first := e.Value[1].(type) {
		case expr.Vector:
			params := []string{}
			body := []expr.Expr{}
			for _, v := range first.Value {
				p, ok := v.(expr.Symbol)
				if !ok {
					fmt.Println(evaluator.errorInfo("repl", p, "expect a symbol"))
					return nil, false
				}
				params = append(params, p.Value)
			}
			body = e.Value[2:]
			// check param name duplication
			exist := map[string]bool{}
			for _, param := range params {
				if exist[param] {
					fmt.Println(evaluator.errorInfo("repl", e, "parameter name must be unique"))
					return nil, false
				}
				exist[param] = true
			}
			closure := expr.NewClosure(evaluator.env, params, body)
			return closure, true
		case expr.Symbol:
			name := first.Value
			if len(e.Value) < 3 {
				fmt.Println(evaluator.errorInfo("repl", s, "expect an argument list"))
				return nil, false
			}
			// redispatch to (var (fn [...] ...))
			newFn := []expr.Expr{e.Value[0]}
			newFn = append(newFn, e.Value[2:]...)
			newVar := expr.NewList(expr.NewSymbol(SF_VAR), expr.NewSymbol(name), expr.NewList(newFn...))
			return evaluator.Eval(newVar)
		default:
			fmt.Println(evaluator.errorInfo("repl", e, "expect a symbol or an argument list"))
			return nil, false
		}

	case SF_MACRO:
		return nil, false

	default:
		panic("unrechable")
	}
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

	// case expr.Quote:
	// 	return e.Value, true

	// case expr.Def:
	// 	value, ok := evaluator.Eval(e.Value)
	// 	if !ok {
	// 		fmt.Println(evaluator.errorInfo("repl", e.Value, "failed to evaluate value of def"))
	// 		return nil, false
	// 	}
	// 	if !evaluator.env.Add(e.Name, value) {
	// 		fmt.Println(evaluator.errorInfo("repl", e, "already defined:", e.Name))
	// 		return nil, false
	// 	}
	// 	return nil, true

	// case expr.Set:
	// 	value, ok := evaluator.Eval(e.Value)
	// 	if !ok {
	// 		fmt.Println(evaluator.errorInfo("repl", e.Value, "failed to evaluate value of set"))
	// 		return nil, false
	// 	}
	// 	if !evaluator.env.Set(e.Name, value) {
	// 		fmt.Println(evaluator.errorInfo("repl", e, "undefined:", e.Name))
	// 		return nil, false
	// 	}
	// 	return nil, true

	// case expr.If:
	// 	pred, ok := evaluator.Eval(e.Pred)
	// 	if !ok {
	// 		fmt.Println(evaluator.errorInfo("repl", e.Pred, "failed to evaluate if pred"))
	// 		return nil, false
	// 	}
	// 	if isTruthy(pred) {
	// 		return evaluator.Eval(e.Then)
	// 	} else {
	// 		return evaluator.Eval(e.Else)
	// 	}

	// case expr.Closure:
	// 	e.Env = evaluator.env
	// 	// check param name duplication
	// 	exist := map[string]bool{}
	// 	for _, param := range e.Params {
	// 		if exist[param] {
	// 			fmt.Println(evaluator.errorInfo("repl", e, "parameter name must be unique"))
	// 			return nil, false
	// 		}
	// 		exist[param] = true
	// 	}
	// 	return e, true

	case expr.Vector:
		var res []expr.Expr
		for _, v := range e.Value {
			value, ok := evaluator.Eval(v)
			if !ok {
				return nil, false
			}
			res = append(res, value)
		}
		return expr.NewVector(res...), true

	case expr.List:
		if len(e.Value) == 0 {
			return e, true
		}

		if isSpecialForm(e) {
			return evaluator.evalSpecialForm(e)
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
