# A Lisp in Go

GoLisp is an embedded Lisp in Go, bringing Lisp's powerful macro to Go, provoding a dynamic and extensible execution layer.

Being embedded means you can easily provide existing Go function to GoLisp by defining and register builtin functions.

Typical usage: rule engine, game logic, workflow script, plugin and extension system.

## Usage

```go
import (
    "github.com/guiyuanju/golisp/evaluator"
	"github.com/guiyuanju/golisp/parser"
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

	e := withPrelude()
	res, ok := e.EvalString(program)
	if !ok {
		return
	}
	fmt.Println("result =", res)
}
```

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

## Quick Start

1. `git clone https://github.com/guiyuanju/golisp.git && cd golisp`
2. `go build`
3. `./golisp` for REPL or `./golisp filename` for file

## Features

- [x] primitives
  - [x] bool
  - [x] int
  - [x] string
  - [x] symbol
- [x] variable
- [x] closure / function
- [x] control flow
- [x] quote, eval
- [x] macro
- [ ] module / namespace
- [ ] bytecode virtual machine

## Test

```go
cd test
go test -v ./...
```
