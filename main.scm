(fn fib (x)
    (if (< x 2)
        x
        (+ (fib (- x 1))
            (fib (- x 2)))))

(var form '(timeit (fib 30)))

(fn show (x y & rest)
    (print x y rest))

(apply (show 1 2) '(1 2 3 4 5 6))
