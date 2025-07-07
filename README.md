# A Lisp in Go

```scheme
; function definition
(fn fib [x]
    (if (< x 2)
        x
        (+ (fib (- x 1))
           (fib (- x 2)))))

; macro definition
(macro timeit [forms]
    (list 'let '[start (time)]
        (list 'do
            forms
            '(nano->milisec (- (time) start)))))

; variable definition
(var form '(timeit (fib 30)))

; macro expansion
(print (macroexpand form))
; => (let [start (time)] (do (fib 30) (nano->milisec (- (time) start))))

(print (macroexpand (macroexpand form)))
; => ((fn [] (var start (time)) (do (fib 30) (nano->milisec (- (time) start)))))

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
