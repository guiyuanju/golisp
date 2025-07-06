package parser

import (
	"testing"
)

func TestScanInt(t *testing.T) {
	s := NewScanner("123")
	ts, ok := s.Scan()
	if !ok || ts[0].Value.(int) != 123 {
		t.Error("expect 123, got", ts[0].Value)
	}
}

func TestScanString(t *testing.T) {
	str := "\"this is a string\""
	s := NewScanner(str)
	ts, _ := s.Scan()
	if ts[0].Value.(string) != str[1:len(str)-1] {
		t.Error("expect string, got:", ts[0].Value)
	}
}

func TestScanSymbol(t *testing.T) {
	str := "function-name"
	s := NewScanner(str)
	ts, _ := s.Scan()
	if ts[0].Value.(string) != str {
		t.Error("expect symbol, got:", ts[0].Value)
	}
}

func TestScanComplex(t *testing.T) {
	type testCase struct {
		name   string
		input  string
		expect []TokenType
	}
	cases := []testCase{
		{"list", "(+ 1 2 3)", []TokenType{LEFT_PAREN, SYMBOL, INTEGER, INTEGER, INTEGER, RIGHT_PAREN}},
		{"number", "123", []TokenType{INTEGER}},
		{"string", "\"a string\"", []TokenType{STRING}},
		{"true", "true", []TokenType{TRUE}},
		{"complex", "(if true (set a (+ a 1)) b)",
			[]TokenType{LEFT_PAREN, IF, TRUE, LEFT_PAREN, SET, SYMBOL, LEFT_PAREN,
				SYMBOL, SYMBOL, INTEGER, RIGHT_PAREN, RIGHT_PAREN, SYMBOL, RIGHT_PAREN}},
	}
	for _, tc := range cases {
		s := NewScanner(tc.input)
		res, ok := s.Scan()
		if !ok {
			t.Fatal("scan failed at", tc.name)
		}
		if len(res) != len(tc.expect) {
			t.Fatal("scan failed at", tc.name, "lenght not equal")
		}
		for i, r := range res {
			if r.TokenType != tc.expect[i] {
				t.Fatal("scan failed at", tc.name, "tokentype not equal")
			}
		}
	}
}

func TestScanPosition(t *testing.T) {
	str := "(+ 1\n   23)"
	s := NewScanner(str)
	res, _ := s.Scan()
	one := res[2]
	if !(one.Value.(int) == 1 && one.Line == 1 && one.Column == 4) {
		t.Error("postion info of 1 incorrect", one)
	}
	twoThree := res[3]
	if !(twoThree.Value.(int) == 23 && twoThree.Line == 2 && twoThree.Column == 4 && twoThree.Length == 2) {
		t.Error("postion info of 2 incorrect", twoThree)
	}
}
