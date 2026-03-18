package main

import (
	"log"
	"os"

	"github.com/guiyuanju/golisp/evaluator"
	"github.com/guiyuanju/golisp/repl"
)

func main() {
	e := evaluator.WithPrelude()

	args := os.Args[1:]
	if len(args) == 0 {
		repl.Repl(e)
		return
	}

	filename := args[0]
	code, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, ok := e.EvalString(string(code))
	if !ok {
		os.Exit(1)
	}
}
