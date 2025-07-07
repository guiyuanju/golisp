(macro and [a b]
    (list 'if a b a))

(macro or [a b]
    (list 'if a a b))

(macro pair [seq]
    (if (= 0 (len seq))
        nil
        (: (pair ) [(. -2 seq) (. -1 seq)])))

(print (macroexpand '(pair [a 1 b 2])))
