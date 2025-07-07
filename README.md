# A Lisp in Go

```scheme
(fn fib [x]
    (if (< x 2)
        x
        (+ (fib (- x 1))
           (fib (- x 2)))))

(var form '(timeit (fib 30)))
(print (macroexpand form))
; => (let [start (time)] (do (fib 30) (nano->milisec (- (time) start))))
(print (eval form))
; => 8397 miliseconds
```

Features:
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
