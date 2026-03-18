package main

import (
	"fmt"

	"github.com/guiyuanju/golisp/evaluator"
	"github.com/guiyuanju/golisp/expr"
)

func main() {
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
		fmt.Printf("%v -> %v\n", prices[order], res)
	}
}
