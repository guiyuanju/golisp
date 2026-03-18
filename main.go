package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/guiyuanju/golisp/evaluator"
	"github.com/guiyuanju/golisp/expr"
	"github.com/guiyuanju/golisp/parser"
)

func example2() {
	prices := []float64{12, 100, 80, 120, 200, 40, 30}
	getPriceForOrder := func(e evaluator.Evaluator, params ...expr.Expr) (expr.Expr, bool) {
		order, ok := params[1].(expr.Number)
		if !ok {
			return nil, false
		}
		return expr.NewNum(prices[int(order.Value)]), true
	}
	evaluator.RegisterBuiltin("get-price-for-order", getPriceForOrder)

	dynamicDiscountRule := `
	(fn is-discount-applicable (order)
		(>= (get-price-for-order order) 100))
	
	(fn apply-percentage-discount (price)
		(* 0.8 price))

	(fn main (order)
		(let (price (get-price-for-order order))
			(if (is-discount-applicable order)
				(apply-percentage-discount price)
				price)))
	main
	`

	e := evaluator.WithPrelude()
	mainFunc, _ := e.EvalString(dynamicDiscountRule)
	for order := range prices {
		res, ok := e.Eval(expr.NewList(mainFunc, expr.NewNum(float64(order))))
		if !ok {
			return
		}
		fmt.Println(res)
	}
}

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

	e := evaluator.WithPrelude()
	res, ok := e.EvalString(program)
	if !ok {
		return
	}
	fmt.Println("result =", res)
}

func main() {
	example2()
	return

	e := evaluator.WithPrelude()

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
