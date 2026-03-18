package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/guiyuanju/golisp/evaluator"
	"github.com/guiyuanju/golisp/parser"
)

func example() {
	program := `
; function definition
(fn fib (x)
    (if (< x 2)
        x
        (+ (fib (- x 1))
           (fib (- x 2)))))

; macro definition
(macro timeit (forms)
    (list 'let '(start (time))
        (list 'do
            forms
            '(nano->milisec (- (time) start)))))

(var form '(timeit (fib 30)))

(print (macroexpand form))

(print (macroexpand (macroexpand form)))

(print (eval form) "miliseconds")
	`

	e := withPrelude()
	res, ok := e.EvalString(program)
	if !ok {
		return
	}
	fmt.Println("result =", res)
}

func main() {
	example()
	return

	e := withPrelude()

	args := os.Args[1:]
	if len(args) == 0 {
		repl(e)
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

func repl(e evaluator.Evaluator) {
	scanner := bufio.NewScanner(os.Stdin)
	hasInput := true
	for hasInput {
		fmt.Print("> ")
		hasInput = scanner.Scan()

		line := scanner.Text()

		s := parser.NewScanner(line)
		tokens, ok := s.Scan()
		if !ok {
			continue
		}

		p := parser.New(tokens)
		exprs, ok := p.Parse()
		if !ok {
			continue
		}

		e.Positions = p.Positions
		for _, expr := range exprs {
			res, ok := e.Eval(expr)
			if !ok {
				continue
			}
			fmt.Println(res)
		}
	}
}

func withPrelude() evaluator.Evaluator {
	prelude, err := os.ReadFile("./stdlib/prelude.scm")
	if err != nil {
		log.Fatal(err)
	}
	e := evaluator.New()
	_, ok := e.EvalString(string(prelude))
	if !ok {
		os.Exit(1)
	}
	return e
}
