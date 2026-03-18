package main

import (
	"fmt"

	"github.com/guiyuanju/golisp/evaluator"
)

func main() {
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

	e := evaluator.WithPrelude()
	res, ok := e.EvalString(program)
	if !ok {
		return
	}
	fmt.Println("result =", res)
}
