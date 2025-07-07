(fn fib [x] (if (< x 2) x (+ (fib (- x 1)) (fib (- x 2)))))
(print (macroexpand '(timeit (fib 30))))
(print (timeit (fib 30)))
