package main

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	r   *bufio.Reader
	err error
	exp Expr
}

func (p *Parser) Scan() bool {
	expr, err := parseToken(p.r)
	if err != nil {
		if err != io.EOF {
			p.err = err
		}
		return false
	}
	p.exp = expr
	return true
}

func (p *Parser) Expression() Expr {
	return p.exp
}

func (p *Parser) Err() error {
	return p.err
}

func NewParser(r io.Reader) *Parser {
	return &Parser{r: bufio.NewReader(r)}
}

func readUntilWhitespace(in *bufio.Reader) (string, error) {
	return readUntil(in, unicode.IsSpace)
}

func readUntil(in *bufio.Reader, test func(rune) bool) (string, error) {
	var r rune
	var err error
	var buf bytes.Buffer
	read := func() error {
		r, _, err = in.ReadRune()
		if err != nil {
			return err
		}
		_, err := buf.WriteRune(r)
		return err
	}
	for read(); !test(r) && err == nil; read() {
	}
	return buf.String(), err
}

func invert(test func(rune) bool) func(rune) bool {
	return func(r rune) bool { return !test(r) }
}

func parseToken(in *bufio.Reader) (Expr, error) {
	// Remove all whitespace at the beginning of the token
	s, err := readUntil(in, invert(unicode.IsSpace))
	if err != nil {
		return nil, err
	}
	r := s[len(s)-1]
	switch r {
	case '(':
		l := []Expr{}
		for {
			expr, err := parseToken(in)
			if err != nil {
				return nil, err
			}
			if expr == nil {
				break
			}
			l = append(l, expr)
		}
		return l, nil
	case ')':
		return nil, nil
	}
	token, err := readUntil(in, func(r rune) bool {
		return unicode.IsSpace(r) || r == '(' || r == ')'
	})
	if err != nil {
		return nil, err
	}
	// If the last character is a token terminator, put it back
	if token[len(token)-1] == ')' {
		in.UnreadRune()
		token = token[:len(token)-1]
	}
	return atom(strings.TrimSpace(string(r) + token)), nil
}

func atom(token string) Expr {
	// First, see if it's a Number
	f, err := strconv.ParseFloat(token, 64)
	if err == nil {
		return Number(f)
	}
	// Assume it's a Symbol
	return Symbol(token)
}
