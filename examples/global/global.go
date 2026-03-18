package main

import (
	"fmt"

	"github.com/guiyuanju/golisp/evaluator"
)

func main() {
	dynamicDiscountRule := "(var a 0)"

	e := evaluator.WithPrelude()
	_, ok := e.EvalString(dynamicDiscountRule)
	if !ok {
		return
	}

	v, err := e.GetGlobal("a")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)

	e.SetGlobal("a", "a new string value")
	v, _ = e.GetGlobal("a")
	fmt.Println(v)
}
