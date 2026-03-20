# GoLisp

<img src="./resources/logo.png" alt="Description of image" width="100">

GoLisp is an embedded Lisp in Go, bringing Lisp's powerful macro to Go, provoding a dynamic and extensible execution layer. Being embedded means you can easily provide existing Go function to GoLisp by defining and register builtin functions. Typical usage: rule engine, game logic, workflow script, plugin and extension system.

## Usage

### As a embedded language (library)

Add dependency:

```sh
go get github.com/guiyuanju/golisp
```

Use GoLisp script to dynamically execute discount strategy:
(Check **examples** folder for more examples)

```go
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
```

### As a standalone language (executable)

Install binary:

```sh
go install github.com/guiyuanju/golisp
```

Run REPL:

```sh
golisp
```

Run with golisp source file:

```sh
golisp main.gl
```

## Feature

Interop:

- Register Go function for GoLisp to use
- Invoke GoLisp function from Go code
- Auto wrap and unwrap value between Go and GoLisp
- Get global value of GoLisp from Go code
- Set global value of GoLisp from Go code

GoLisp:

- primitives
  - bool
  - number (float64)
  - string
  - symbol
- list: `'(1 2 3)` a list
- variable: `(var a 0)`
- control flow: `(if cond true-brach false-branch)`
- closure: `(fn (arg) ...)`
- function: `(fn name (arg) ...)` equals `(var name (fn (arg) ...))`
- quote: `'a`
- eval: `(eval ...)`
- macro: `(macro name [forms] ...)`, `(macroexpand macroname)`
- [ ] module / namespace
- [ ] bytecode virtual machine

GoLisp syntax example:

```scheme
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

; variable definition
(var form '(timeit (fib 30)))

; macro expansion
(print (macroexpand form))
; => (let (start (time)) (do (fib 30) (nano->milisec (- (time) start))))

(print (macroexpand (macroexpand form)))
; => ((fn () (var start (time)) (do (fib 30) (nano->milisec (- (time) start)))))

; eval
(print (eval form) "miliseconds")
; => 8397 miliseconds
```

## Test

```go
go test -v ./...
```
