(fn fib [x] (if (< x 2) x (+ (fib (- x 1)) (fib (- x 2)))))
(fn nano->sec [x] (/ x 1000000000))
(macro timeit [form]
    (list 'do
        (list 'var 'start (list 'time))
        form
        (list '- (list 'time) 'start)))

(print (macroexpand '(timeit (+ 1 2))))
