package main

import (
	"fmt"

	"github.com/guiyuanju/golisp/evaluator"
)

func main() {
	// data is Go side
	prices := []float64{12, 100, 80, 120, 200, 40, 30}

	// define a Go function to retrieve data, register it for GoLisp to use
	evaluator.RegisterBuiltin("get-price-for-order", func(params ...any) (any, error) {
		// ignoring all error handling
		return prices[int(params[0].(float64))], nil
	})

	// define script, can be provided dynamically
	dynamicDiscountRule :=
		`
		;; use Go function to get data
		(fn is-discount-applicable (order)
			(>= (get-price-for-order order) 100))
		
		(fn apply-percentage-discount (price)
			(* 0.8 price))

		;; function defined in script can be invoked from Go side
		(fn get-discounted-price (order)
			(let (price (get-price-for-order order))
				(if (is-discount-applicable order)
					(apply-percentage-discount price)
					price)))
		`

	// initialize evaluator with standard library
	e := evaluator.WithPrelude()

	// evaluate script
	e.EvalString(dynamicDiscountRule)

	// invoke function in script for each order, get discounted price
	for order := range prices {
		res, _ := e.InvokeFunc("get-discounted-price", order)
		fmt.Printf("%v -> %v\n", prices[order], res)
	}
}
