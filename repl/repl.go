package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/guiyuanju/golisp/evaluator"
	"github.com/guiyuanju/golisp/parser"
)

func Repl(e evaluator.Evaluator) {
	scanner := bufio.NewScanner(os.Stdin)
	hasInput := true
	for hasInput {
		fmt.Print("> ")
		hasInput = scanner.Scan()

		line := scanner.Text()

		s := parser.NewScanner(line)
		tokens, ok := s.Scan()
		if !ok {
			continue
		}

		p := parser.New(tokens)
		exprs, ok := p.Parse()
		if !ok {
			continue
		}

		e.Positions = p.Positions
		for _, expr := range exprs {
			res, ok := e.Eval(expr)
			if !ok {
				continue
			}
			fmt.Println(res)
		}
	}
}
