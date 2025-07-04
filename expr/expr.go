package expr

import (
	"fmt"
	"strings"
)

type Expr interface {
	Line() int
	Column() int
}

type Int struct {
	Pos
	Value int
}

type String struct {
	Pos
	Value string
}

type Symbol struct {
	Pos
	Value string
}

type List struct {
	Pos
	Value []Expr
}

type Def struct {
	Pos
	Name  string
	Value Expr
}

type Set struct {
	Pos
	Name  string
	Value Expr
}

type Proc func(...Expr) (Expr, bool)

type Builtin struct {
	Pos
	Name string
	Proc Proc
}

type Closure struct {
	Pos
	Env    Env
	Params []string
	Body   Expr
}

type Pos struct {
	line   int
	column int
}

type Env []map[string]Expr

func ErrorInfo(file string, expr Expr, info ...string) string {
	return fmt.Sprintf("%s:%d:%d: (%s) %s", file, expr.Line(), expr.Column(), fmt.Sprintf("%T", expr), strings.Join(info, " "))
}

func NewPos(line, column int) Pos {
	return Pos{line, column}
}

func (p Pos) Line() int {
	return p.line
}

func (p Pos) Column() int {
	return p.column
}

func NewInt(value int, pos ...int) Int {
	if len(pos) == 2 {
		return Int{Pos{pos[0], pos[1]}, value}
	}
	return Int{Pos{0, 0}, value}
}

func NewString(value string, pos ...int) String {
	if len(pos) == 2 {
		return String{Pos{pos[0], pos[1]}, value}
	}
	return String{Pos{0, 0}, value}
}

func NewSymbol(value string, pos ...int) Symbol {
	if len(pos) == 2 {
		return Symbol{Pos{pos[0], pos[1]}, value}
	}
	return Symbol{Pos{0, 0}, value}
}

func NewList(values ...Expr) List {
	return List{Pos{0, 0}, values}
}

func NewDef(name string, value Expr) Def {
	return Def{Pos{0, 0}, name, value}
}

func NewSet(name string, value Expr) Set {
	return Set{NewPos(0, 0), name, value}
}

func NewBuiltin(name string, proc Proc) Builtin {
	return Builtin{NewPos(0, 0), name, proc}
}

func NewClosure(env Env, params []string, body Expr) Closure {
	return Closure{NewPos(0, 0), env, params, body}
}

func NewEnv() Env {
	m := map[string]Expr{}
	return Env([]map[string]Expr{m})
}

func (e Env) Get(name string) (Expr, bool) {
	if len(e) == 0 {
		return nil, false
	}
	for i := len(e) - 1; i >= 0; i-- {
		if v, ok := e[i][name]; ok {
			return v, ok
		}
	}
	return nil, false
}

// if exists, return false
func (e Env) Add(name string, value Expr) bool {
	if len(e) == 0 {
		panic("env len is zero")
	}
	cur := e[len(e)-1]
	if _, ok := cur[name]; ok {
		return false
	}
	cur[name] = value
	return true
}

func (e Env) Set(name string, value Expr) bool {
	if len(e) == 0 {
		panic("env len is zero")
	}
	for i := len(e) - 1; i >= 0; i-- {
		if _, ok := e[i][name]; ok {
			e[i][name] = value
			return true
		}
	}
	return false
}

func (e Env) AppendEnv(env Env) Env {
	if len(env) == 0 {
		return e
	}
	if len(env) > 1 {
		panic("AppendEnv only allow append env with length 1")
	}
	return append(e, env[0])
}
