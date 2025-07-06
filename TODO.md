- [x] move pos info to separate place, split concern btw pos and ast, can use hash map
- [x] parser
- [ ] test
- [x] String() for Expr, pretty printing
- [ ] move do... to std lib
- [ ] stdlib implement

Features:
- [x] primitives
  - [x] int
  - [x] string
  - [x] symbol
- [x] variable
  - [x] define
  - [x] update
- [x] closure
- [x] do block
- [x] print
- [ ] control flow
  - [x] if
  - [x] bool, nil
  - [ ] comparison
    - [x] =
    - [ ] >
  - [ ] logical
    - [ ] not = if x then false else true
    - [ ] and = if a then b else false
    - [ ] or = if a then true else b
- [ ] quote
- [ ] macro

keword:
- if
- def
- set
- closure
- list
- append
- index
- len

builtin:
- print
- + - * /
- = < ...

stdlib:
- do
- != > >= <=
- not and or


macro
function
syntax sugar
compiler
builtin

(def f [x y xs...] )
(f x y 1 2 3)
(f x y xs...)

builtin: funcitonailyt that cannot be implemented in language itself
function: can be implemented in itself
syntax sugar: can be replaced as language code (may not be able to be implemeted as function), as a language keyword, not redefinable, can be implemented
macro: not precise error info, similar to syntax sugar, but not a language keyword

```
expr = int | string | bool | symbol | nil | quote | var | set | if | fn | | macro | list
var = "(" "var" symbol expr ")"
set = "(" "set" symbol expr ")"
if = "(" "if" expr expr expr? ")"
fn = "(" "fn" symbol? "[" symbol* "]" expr* ")"
quote = "'" expr
macro = "(" "macro" symbol "[" symbol* "]" expr* ")"
list = "(" expr* ")"
```

macro:
(macro and [a b] '(if a b false))
'(if a b false) => (list 'if 'a 'b 'false)
(and true false)
