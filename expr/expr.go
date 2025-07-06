package expr

import (
	"fmt"
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

var id int

func getId() int {
	res := id
	id++
	return res
}

type Expr interface {
	ExprId() int
	ExprName() string
}

type Int struct {
	Id    int
	Value int
}

func (e Int) ExprId() int {
	return e.Id
}
func (e Int) ExprName() string {
	return "int"
}
func (e Int) String() string {
	return fmt.Sprint(e.Value)
}

type String struct {
	Id    int
	Value string
}

func (e String) ExprId() int {
	return e.Id
}
func (e String) ExprName() string {
	return "string"
}
func (e String) String() string {
	return fmt.Sprint(e.Value)
}

type Bool struct {
	Id    int
	Value bool
}

func (e Bool) ExprId() int {
	return e.Id
}
func (e Bool) ExprName() string {
	return "bool"
}
func (e Bool) String() string {
	return fmt.Sprint(e.Value)
}

type Symbol struct {
	Id    int
	Value string
}

func (e Symbol) ExprId() int {
	return e.Id
}
func (e Symbol) ExprName() string {
	return "symbol"
}
func (e Symbol) String() string {
	return fmt.Sprint(e.Value)
}

type Nil struct {
	Id int
}

func (e Nil) ExprId() int {
	return e.Id
}
func (e Nil) ExprName() string {
	return "nil"
}
func (e Nil) String() string {
	return "nil"
}

type Macro struct {
	Id     int
	Name   string
	Params []string
	Body   []Expr
}

func (e Macro) ExprId() int {
	return e.Id
}
func (e Macro) ExprName() string {
	return "macro"
}
func (e Macro) String() string {
	return fmt.Sprintf("<macro %s>", e.Name)
}

type Vector struct {
	Id    int
	Value []Expr
}

func (e Vector) ExprId() int {
	return e.Id
}
func (e Vector) ExprName() string {
	return "vector"
}
func (e Vector) String() string {
	var res []string
	for _, v := range e.Value {
		res = append(res, fmt.Sprint(v))
	}
	return "[" + strings.Join(res, " ") + "]"
}

type List struct {
	Id    int
	Value []Expr
}

func (e List) ExprId() int {
	return e.Id
}
func (e List) ExprName() string {
	return "list"
}
func (e List) String() string {
	var res []string
	for _, v := range e.Value {
		res = append(res, fmt.Sprint(v))
	}
	return "(" + strings.Join(res, " ") + ")"
}

type Builtin struct {
	Id   int
	Name string
}

func (e Builtin) ExprId() int {
	return e.Id
}
func (e Builtin) ExprName() string {
	return "builtin"
}
func (e Builtin) String() string {
	return fmt.Sprintf("<builtin %s>", e.Name)
}

type Closure struct {
	Id     int
	Env    Env
	Params []string
	Body   []Expr
}

func (e Closure) ExprId() int {
	return e.Id
}
func (e Closure) ExprName() string {
	return "closure"
}
func (e Closure) String() string {
	return "<closure>"
}

type Env []map[string]Expr

func NewInt(value int) Int {
	return Int{getId(), value}
}

func NewString(value string) String {
	return String{getId(), value}
}

func NewBool(value bool) Bool {
	return Bool{getId(), value}
}

func NewSymbol(value string) Symbol {
	return Symbol{getId(), value}
}

func NewNil() Nil {
	return Nil{getId()}
}

func NewMacro(name string, params []string, body []Expr) Macro {
	return Macro{getId(), name, params, body}
}

func NewList(values ...Expr) List {
	return List{getId(), values}
}

func NewVector(values ...Expr) Vector {
	return Vector{getId(), values}
}

func NewBuiltin(name string) Builtin {
	return Builtin{getId(), name}
}

func NewClosure(env Env, params []string, body []Expr) Closure {
	return Closure{getId(), env, params, body}
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
