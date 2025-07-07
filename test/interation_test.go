package test

import (
	"golisp/evaluator"
	"testing"
)

type testCase struct {
	name   string
	code   string
	expect string
}

type testSuite struct {
	name      string
	testcases []testCase
}

var TSS []testSuite = []testSuite{
	{
		"comparison",
		[]testCase{
			{">", "(> 2 1)", "true"},
			{">", "(> 1 2)", "false"},
			{"<", "(< 1 2)", "true"},
			{"<", "(< 2 1)", "false"},
			{"=", "(= 1 1)", "true"},
			{"=", "(= nil nil)", "true"},
			{"=", "(= 1 2)", "false"},
			{"=", "(= 1 nil)", "false"},
			{">=", "(>= 2 1)", "true"},
			{">=", "(>= 1 2)", "false"},
			{"<=", "(<= 1 2)", "true"},
			{"<=", "(<= 2 1)", "false"},
		},
	},
	{
		"function",
		[]testCase{
			{"fib", "(fn fib [x] (if (< x 2) x (+ (fib (- x 1)) (fib (- x 2))))) (fib 10)", "55"},
		},
	},
}

func TestSuites(t *testing.T) {
	for _, ts := range TSS {
		for _, tc := range ts.testcases {
			e := evaluator.New()
			expr, ok := e.EvalString(tc.code)
			if !ok {
				t.Fatalf("%s: %s failed, EvalString not ok", ts.name, tc.name)
			}
			if expr.String() != tc.expect {
				t.Fatalf("%s: %s failed, expect %s, got %s", ts.name, tc.name, tc.expect, expr.String())
			}
		}
	}
}
