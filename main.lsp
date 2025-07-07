(fn fib [x] (if (< x 2) x (+ (fib (- x 1)) (fib (- x 2)))))

(print (fib 0) (fib 1) (fib 2) (fib 3) (fib 4) (fib 5) (fib 6) (fib 7))
(print (fib 30))
