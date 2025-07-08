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

func New() Evaluator {
	env := expr.NewEnv()
	builtins := NewBuiltins()
	for name, _ := range builtins {
		env.Add(name, expr.NewBuiltin(name))
	}
	return Evaluator{env, parser.NewPositions(), builtins}
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
	case expr.SF_QUOTE, expr.SF_VAR, expr.SF_SET, expr.SF_IF, expr.SF_FN, expr.SF_MACRO, expr.SF_APPLY:
		return true
	default:
		return false
	}
}

func (evaluator Evaluator) evalSpecialForm(e expr.List) (expr.Expr, bool) {
	s := e.Value[0].(expr.Symbol)
	switch s.Value {
	case expr.SF_QUOTE:
		if len(e.Value) != 2 {
			fmt.Println(evaluator.errorInfo("repl", s, "expect 1 argument"))
			return nil, false
		}
		return e.Value[1], true

	case expr.SF_VAR:
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

	case expr.SF_SET:
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

	case expr.SF_IF:
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

	case expr.SF_FN:
		if len(e.Value) < 3 {
			fmt.Println(evaluator.errorInfo("repl", s, "expect an argument list and a body"))
			return nil, false
		}
		switch first := e.Value[1].(type) {
		case expr.List:
			params := []string{}
			body := []expr.Expr{}
			var i int
			for ; i < len(first.Value); i++ {
				p, ok := first.Value[i].(expr.Symbol)
				if !ok {
					fmt.Println(evaluator.errorInfo("repl", p, "expect a symbol"))
					return nil, false
				}
				if p.Value == "&" {
					break
				}
				params = append(params, p.Value)
			}
			// varparam
			var varparam string
			if i < len(first.Value) {
				if i == len(first.Value)-1 {
					fmt.Println(evaluator.errorInfo("repl", first.Value[i], "expect a symbol after &"))
					return nil, false
				}
				v, ok := first.Value[i+1].(expr.Symbol)
				if !ok {
					fmt.Println(evaluator.errorInfo("repl", first.Value[i+1], "expect a symbol"))
					return nil, false
				}
				varparam = v.Value
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

			closure := expr.NewClosure(evaluator.env, params, varparam, body)
			return closure, true
		case expr.Symbol:
			name := first.Value
			if len(e.Value) < 4 {
				fmt.Println(evaluator.errorInfo("repl", s, "expect an argument list and body"))
				return nil, false
			}
			// redispatch to (var (fn [...] ...))
			newFn := []expr.Expr{e.Value[0]}
			newFn = append(newFn, e.Value[2:]...)
			newVar := expr.NewList(expr.NewSymbol(expr.SF_VAR), expr.NewSymbol(name), expr.NewList(newFn...))
			return evaluator.Eval(newVar)
		default:
			fmt.Println(evaluator.errorInfo("repl", e.Value[1], "expect a symbol or an argument list"))
			return nil, false
		}

	case expr.SF_MACRO:
		if len(e.Value) < 4 {
			fmt.Println(evaluator.errorInfo("repl", e, "expect a symbol, a argument list and body"))
			return nil, false
		}
		name, ok := e.Value[1].(expr.Symbol)
		if !ok {
			fmt.Println(evaluator.errorInfo("repl", e.Value[1], "expect a symbol"))
			return nil, false
		}
		args, ok := e.Value[2].(expr.List)
		if !ok {
			fmt.Println(evaluator.errorInfo("repl", e.Value[2], "expect a argument list"))
			return nil, false
		}
		params := []string{}
		var i int
		for ; i < len(args.Value); i++ {
			p, ok := args.Value[i].(expr.Symbol)
			if !ok {
				fmt.Println(evaluator.errorInfo("repl", args.Value[i], "expect a symbol"))
				return nil, false
			}
			if p.Value == "&" {
				break
			}
			params = append(params, p.Value)
		}
		var varparam string
		if i < len(args.Value) {
			if i == len(args.Value)-1 {
				fmt.Println(evaluator.errorInfo("repl", args.Value[i], "expect a symbol after &"))
				return nil, false
			}
			v, ok := args.Value[i+1].(expr.Symbol)
			if !ok {
				fmt.Println(evaluator.errorInfo("repl", v, "expect a symbol"))
				return nil, false
			}
			varparam = v.Value
		}
		body := e.Value[3:]
		closure := expr.NewClosure(evaluator.env, params, varparam, body)
		macro := expr.NewMacro(name.Value, closure)
		if !evaluator.env.Add(name.Value, macro) {
			fmt.Println(evaluator.errorInfo("repl", name, "already defined"))
			return nil, false
		}
		return expr.NewNil(), true

	case expr.SF_APPLY:
		if len(e.Value)-1 < 2 {
			fmt.Println(evaluator.errorInfo("repl", e.Value[0], "need at least 2 arguments"))
			return nil, false
		}
		rest, ok := evaluator.Eval(e.Value[2])
		if !ok {
			return nil, false
		}
		restList, ok := rest.(expr.List)
		if !ok {
			fmt.Println(evaluator.errorInfo("repl", e.Value[2], "expect a list"))
			return nil, false
		}
		switch f := e.Value[1].(type) {
		case expr.List:
			f.Append(restList)
			f.Value = append(f.Value, restList.Value...)
			return evaluator.Eval(f)
		default:
			v := []expr.Expr{f}
			v = append(v, restList.Value...)
			return evaluator.Eval(expr.NewList(v...))
		}

	default:
		panic("unrechable")
	}
}

func (evaluator Evaluator) isMacro(e expr.Expr) bool {
	symbol, ok := e.(expr.Symbol)
	if !ok {
		return false
	}
	value, ok := evaluator.env.Get(symbol.Value)
	if !ok {
		return false
	}
	_, ok = value.(expr.Macro)
	return ok
}

func (evaluator Evaluator) macroExpand(e expr.List) (expr.Expr, bool) {
	value, _ := evaluator.env.Get(e.Value[0].(expr.Symbol).Value)
	macro := value.(expr.Macro)

	args := e.Value[1:]

	if len(e.Value)-1 < len(macro.Closure.Params) {
		fmt.Println(evaluator.errorInfo("repl", macro, "expect at least", strconv.Itoa(len(macro.Closure.Params)), "arguments, got", strconv.Itoa(len(e.Value)-1)))
		return nil, false
	}

	return apply(evaluator, macro.Closure, args)
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

	case expr.List:
		if len(e.Value) == 0 {
			return e, true
		}

		if evaluator.isMacro(e.Value[0]) {
			expanded, ok := evaluator.macroExpand(e)
			if !ok {
				return nil, false
			}
			return evaluator.Eval(expanded)
		}

		if isSpecialForm(e) {
			return evaluator.evalSpecialForm(e)
		}

		operator, ok := evaluator.Eval(e.Value[0])
		if !ok {
			return nil, false
		}

		switch operator := operator.(type) {
		case expr.Builtin:
			args := []expr.Expr{operator}
			for _, arg := range e.Value[1:] {
				value, ok := evaluator.Eval(arg)
				if !ok {
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
			if len(e.Value)-1 < len(operator.Params) {
				fmt.Println(evaluator.errorInfo("repl", operator, "expect at least", strconv.Itoa(len(operator.Params)), "arguments, got", strconv.Itoa(len(e.Value)-1)))
				return nil, false
			}

			args := []expr.Expr{}
			for _, arg := range e.Value[1:] {
				value, ok := evaluator.Eval(arg)
				if !ok {
					return nil, false
				}
				args = append(args, value)
			}

			return apply(evaluator, operator, args)
		}

		fmt.Println(evaluator.errorInfo("repl", e.Value[0], "expect proc or function"))
		return nil, false
	}

	panic("unexhaustive evaluation")
}

func (e Evaluator) EvalString(code string) (expr.Expr, bool) {
	s := parser.NewScanner(code)
	tokens, ok := s.Scan()
	if !ok {
		return nil, false
	}
	p := parser.New(tokens)
	exprs, ok := p.Parse()
	if !ok {
		return nil, false
	}
	e.Positions = p.Positions
	var last expr.Expr
	for _, expr := range exprs {
		v, ok := e.Eval(expr)
		if !ok {
			return nil, false
		}
		last = v
	}
	return last, true
}

func apply(e Evaluator, closure expr.Closure, args []expr.Expr) (expr.Expr, bool) {
	env := expr.NewEnv()
	var i int
	for ; i < len(closure.Params); i++ {
		env.Add(closure.Params[i], args[i])
	}
	env.Add(closure.VarParam, expr.NewList(args[i:]...))

	newEnv := closure.Env.AppendEnv(env)
	newEvaluator := New()
	newEvaluator.Positions = e.Positions
	newEvaluator.env = newEnv

	var last expr.Expr
	for _, b := range closure.Body {
		v, ok := newEvaluator.Eval(b)
		if !ok {
			return nil, false
		}
		last = v
	}
	return last, true
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
