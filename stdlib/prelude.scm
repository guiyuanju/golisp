(macro and [a b]
    (list 'if a b a))

(macro or [a b]
    (list 'if a a b))

(macro pair [seq]
    (if (= 0 (len seq))
        nil
        [(. 0 seq) (. 1 seq)]))

(macro let [args body]
        )
; (let [a b c d] ())
; (fn []
;     (var ))

; ()
