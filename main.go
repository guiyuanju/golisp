package main

import (
	"bufio"
	"fmt"
	"golisp/evaluator"
	"golisp/parser"
	"log"
	"os"
)

func main() {
	prelude, err := os.ReadFile("./stdlib/prelude.lip")
	if err != nil {
		log.Fatal(err)
	}

	e := evaluator.New(nil)
	e.EvalString(string(prelude))
	repl(e)
}

func repl(e evaluator.Evaluator) {
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
