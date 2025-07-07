(fn fib [x]
    (if (< x 2)
        x
        (+ (fib (- x 1))
            (fib (- x 2)))))

(var form '(timeit (fib 30)))
(print (macroexpand form))
(print (eval form))
