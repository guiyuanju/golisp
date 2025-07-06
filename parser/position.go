package parser

type Position struct {
	Line, Column int
}

type Positions map[int]Position

func NewPositions() Positions {
	return map[int]Position{}
}
