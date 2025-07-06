package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tokens := []Token{
		NewToken(LEFT_PAREN, 1, 1, 0, nil),
		NewToken(VAR, 0, 2, 0, nil),
		NewToken(SYMBOL, 0, 3, 0, "a"),
		NewToken(INTEGER, 0, 4, 0, 10),
		NewToken(RIGHT_PAREN, 0, 14, 0, nil),

		NewToken(LEFT_PAREN, 0, 5, 0, nil),
		NewToken(SET, 0, 6, 0, nil),
		NewToken(SYMBOL, 0, 7, 0, "a"),
		NewToken(LEFT_PAREN, 0, 8, 0, nil),
		NewToken(SYMBOL, 0, 9, 0, "+"),
		NewToken(SYMBOL, 0, 10, 0, "a"),
		NewToken(INTEGER, 0, 11, 0, 1),
		NewToken(RIGHT_PAREN, 0, 12, 0, nil),
		NewToken(RIGHT_PAREN, 0, 13, 0, nil),
	}

	p := New(tokens)
	res, ok := p.Parse()
	fmt.Println(res, ok)
}
