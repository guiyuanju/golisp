package expr

import (
	"fmt"
	"strings"
)

var id int

func getId() int {
	res := id
	id++
	return res
}

type Expr interface {
	ExprId() int
}

type Int struct {
	Id    int
	Value int
}

func (e Int) ExprId() int {
	return e.Id
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
func (e Symbol) String() string {
	return fmt.Sprint(e.Value)
}

type List struct {
	Id    int
	Value []Expr
}

func (e List) ExprId() int {
	return e.Id
}
func (e List) String() string {
	var res []string
	for _, v := range e.Value {
		res = append(res, fmt.Sprint(v))
	}
	return "(" + strings.Join(res, " ") + ")"
}

type Def struct {
	Id    int
	Name  string
	Value Expr
}

func (e Def) ExprId() int {
	return e.Id
}
func (e Def) String() string {
	return fmt.Sprintf("(var %s %s)", e.Name, e.Value)
}

type Set struct {
	Id    int
	Name  string
	Value Expr
}

func (e Set) ExprId() int {
	return e.Id
}
func (e Set) String() string {
	return fmt.Sprintf("(set %s %s)", e.Name, e.Value)
}

type If struct {
	Id   int
	Pred Expr
	Then Expr
	Else Expr
}

func (e If) ExprId() int {
	return e.Id
}
func (e If) String() string {
	return fmt.Sprintf("(if %s %s %s)", e.Pred, e.Then, e.Else)
}

type Builtin struct {
	Id   int
	Name string
}

func (e Builtin) ExprId() int {
	return e.Id
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
func (e Closure) String() string {
	return "<closure>"
}

type Env []map[string]Expr

func NewInt(value int, pos ...int) Int {
	if len(pos) == 2 {
		return Int{getId(), value}
	}
	return Int{getId(), value}
}

func NewString(value string, pos ...int) String {
	if len(pos) == 2 {
		return String{getId(), value}
	}
	return String{getId(), value}
}

func NewBool(value bool, pos ...int) Bool {
	return Bool{getId(), value}
}

func NewSymbol(value string, pos ...int) Symbol {
	if len(pos) == 2 {
		return Symbol{getId(), value}
	}
	return Symbol{getId(), value}
}

func NewList(values ...Expr) List {
	return List{getId(), values}
}

func NewDef(name string, value Expr) Def {
	return Def{getId(), name, value}
}

func NewSet(name string, value Expr) Set {
	return Set{getId(), name, value}
}

func NewIf(pred, then, _else Expr) If {
	return If{getId(), pred, then, _else}
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
