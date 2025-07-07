(macro and [a b]
    (list 'if a b a))

(macro or [a b]
    (list 'if a a b))

(fn map [f xs]
    (if (= 0 (len xs))
        ()
        (append (map f (: 0 -1 xs)) (f (. -1 xs)))))

(fn pair [seq]
    (if (= 0 (len seq))
        ()
        (append (pair (: 0 -2 seq)) (list (. -2 seq) (. -1 seq)))))

(fn concat [xs ys]
    (if (= 0 (len ys))
        xs
        (concat (append xs (. 0 ys)) (: 1 (len ys) ys))))

(macro let [bindings body]
    (var vars (map (fn [x] (list 'var (. 0 x) (. 1 x)))
                    (pair bindings)))
    (list (concat (concat '(fn []) vars)
            (list body))))

(macro timeit [forms]
    (list 'let '[start (time)]
        (list 'do
            forms
            '(nano->milisec (- (time) start)))))

(fn nano->milisec [x] (/ x 1000000))
