- [x] move pos info to separate place, split concern btw pos and ast, can use hash map
- [x] parser
- [x] test
- [x] String() for Expr, pretty printing
- [x] fix equal, greater, less...
- [x] support macro definition, change `and` to a custom macro
- [x] recursive macro
- [ ] var args -> remove do in let
- [ ] for loop
- [ ] prepend
- [ ] implement let using macro or builtin?
- [ ] macro simplify support ,
- [ ] tail call optimization
- [ ] variable arity function, implement do with macro
- [ ] index
- [ ] len
- [ ] stdlib implement
- [ ] hygine macro
- [ ] static type
- [ ] reactive
- [ ] go interop
- [ ] compile to go
- [ ] abstract list and vector to seq

Features:
- [x] primitives
  - [x] int
  - [x] string
  - [x] symbol
  - [x] bool, nil
- [x] variable
  - [x] define
  - [x] update
- [x] closure
- [x] control flow
  - [x] if
  - [x] comparison
  - [x] logical
- [x] quote
- [x] macro
- [ ] module / namespace

syntax rules:
```
expr = int | string | bool | symbol | nil | quote | var | set | if | fn | | macro | list
var = "(" "var" symbol expr ")"
set = "(" "set" symbol expr ")"
if = "(" "if" expr expr expr? ")"
fn = "(" "fn" symbol? "[" symbol* "]" expr* ")"
quote = "'" expr
macro = "(" "macro" symbol "[" symbol* "]" expr* ")"
list = "(" expr* ")"

special_form = quote | var | set | if | fn | macro
```
