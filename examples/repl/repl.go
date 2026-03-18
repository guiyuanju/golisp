package main

import (
	"github.com/guiyuanju/golisp/evaluator"
	"github.com/guiyuanju/golisp/repl"
)

func main() {
	e := evaluator.WithPrelude()
	repl.Repl(e)
}
