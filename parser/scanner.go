package parser

import (
	"fmt"
	"slices"
)

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACKET
	RIGHT_BRACKET
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	QUOTE
	INTEGER
	STRING
	TRUE
	FALSE
	IF
	VAR
	SET
	FN
	SYMBOL
)

type Token struct {
	TokenType TokenType
	Line      int
	Column    int
	Length    int
	Value     any
}

func NewToken(tokenType TokenType, line, column, length int, value any) Token {
	return Token{tokenType, line, column, length, value}
}

type Scanner struct {
	s      []byte
	i      int // current parsing index
	line   int // current parsing line
	column int // current parsing column in current line
	length int // length of current parsing token
}

func NewScanner(s string) Scanner {
	return Scanner{
		s:      []byte(s),
		i:      0,
		line:   1,
		column: 1,
		length: 0,
	}
}

func (s *Scanner) newToken(tokenType TokenType, value any) Token {
	res := NewToken(tokenType, s.line, s.column-s.length, s.length, value)
	s.length = 0
	return res
}

func (s *Scanner) cur() byte {
	return s.s[s.i]
}

func (s *Scanner) isEnd() bool {
	return s.i >= len(s.s)
}

func (s *Scanner) advance() {
	s.i++
	s.column++
}

func (s *Scanner) consume(value string) bool {
	start := s.i
	for _, b := range []byte(value) {
		if s.isEnd() || b != s.cur() {
			s.i = start
			s.length = 0
			return false
		}
		s.advance()
		s.length++
	}
	return true
}

func (s *Scanner) errorInfo(info string) string {
	return fmt.Sprintf("repl:%d:%d: %s", s.line, s.column, info)
}

func (s *Scanner) Scan() ([]Token, bool) {
	var res []Token
	for !s.isEnd() {
		switch s.cur() {
		case '\n':
			s.line++
			s.column = 0
			s.advance()
		case ' ':
			s.advance()
		case '(':
			res = append(res, s.newToken(LEFT_PAREN, nil))
			s.advance()
		case ')':
			res = append(res, s.newToken(RIGHT_PAREN, nil))
			s.advance()
		case '[':
			res = append(res, s.newToken(LEFT_BRACKET, nil))
			s.advance()
		case ']':
			res = append(res, s.newToken(RIGHT_BRACKET, nil))
			s.advance()
		case '{':
			res = append(res, s.newToken(LEFT_BRACE, nil))
			s.advance()
		case '}':
			res = append(res, s.newToken(RIGHT_BRACE, nil))
			s.advance()
		case ',':
			res = append(res, s.newToken(COMMA, nil))
			s.advance()
		case '\'':
			res = append(res, s.newToken(QUOTE, nil))
			s.advance()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			res = append(res, s.newToken(INTEGER, s.number()))
		case '"':
			v, ok := s.string()
			if !ok {
				return nil, false
			}
			res = append(res, s.newToken(STRING, v))
		default:
			if s.consume("true") {
				res = append(res, s.newToken(TRUE, nil))
			} else if s.consume("false") {
				res = append(res, s.newToken(FALSE, nil))
			} else if s.consume("if") {
				res = append(res, s.newToken(IF, nil))
			} else if s.consume("var") {
				res = append(res, s.newToken(VAR, nil))
			} else if s.consume("set") {
				res = append(res, s.newToken(SET, nil))
			} else if s.consume("fn") {
				res = append(res, s.newToken(FN, nil))
			} else {
				res = append(res, s.newToken(SYMBOL, s.symbol()))
			}
		}
	}
	return res, true
}

// alphabet, -, number
func (s *Scanner) symbol() string {
	nonSym := []byte{'(', ')', '{', '}', '[', ']', ' ', '\n', '"'}
	var res []byte
	for !s.isEnd() && !slices.Contains(nonSym, s.cur()) {
		res = append(res, s.cur())
		s.advance()
		s.length++
	}
	return string(res)
}

func (s *Scanner) string() (string, bool) {
	var res []byte
	s.advance()
	s.length++
	for !s.isEnd() && s.cur() != '"' && s.cur() != '\n' {
		res = append(res, s.cur())
		s.advance()
		s.length++
	}
	if !s.consume("\"") {
		fmt.Println(s.errorInfo("expect \""))
		return "", false
	}
	s.length++
	return string(res), true
}

func (s *Scanner) number() int {
	isDigit := func(x byte) bool {
		return '0' <= x && x <= '9'
	}
	var res int
	for !s.isEnd() && isDigit(s.cur()) {
		res = res*10 + int(s.cur()-'0')
		s.advance()
		s.length++
	}
	return res
}
