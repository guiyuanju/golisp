(macro and [a b]
    (list 'if a b a))

(macro or [a b]
    (list 'if a a b))

(fn head [xs] (. 0 xs))
(fn snd [xs] (. 1 xs))
(fn tail [xs] (: 1 (len xs) xs))
(fn init [xs] (: 0 -1 xs))
(fn last [xs] (. -1 xs))

(fn map [f xs]
    (if (= 0 (len xs))
        ()
        (append (map f (init xs)) (f (last xs)))))

(fn pair [seq]
    (if (= 0 (len seq))
        ()
        (append (pair (: 0 -2 seq)) (list (. -2 seq) (last seq)))))

(fn concat [xs ys]
    (if (= 0 (len ys))
        xs
        (concat (append xs (head ys)) (tail ys))))

(macro let [bindings body]
    (var vars (map (fn [x] (list 'var (head x) (snd x)))
                    (pair bindings)))
    (list (concat (concat '(fn []) vars)
            (list body))))

(macro timeit [forms]
    (list 'let '[start (time)]
        (list 'do
            forms
            '(nano->milisec (- (time) start)))))

(fn nano->milisec [x] (/ x 1000000))
