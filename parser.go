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

func pad(s string) []byte {
	return []byte(" " + s + " ")
}

func NewParser(r io.Reader) *Parser {
	// To make tokenizing easier, we pad all parens with spaces
	src := r
	r, w := io.Pipe()
	go func() {
		scanner := bufio.NewScanner(src)
		scanner.Split(bufio.ScanRunes)
		var err error
		for scanner.Scan() {
			switch scanner.Text() {
			case "(", ")":
				_, err = w.Write(pad(scanner.Text()))
			default:
				_, err = w.Write(scanner.Bytes())
			}
			if err != nil {
				w.CloseWithError(err)
				return
			}
		}
		if err := scanner.Err(); err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}()
	return &Parser{r: bufio.NewReader(r)}
}

func readUntilWhitespace(in *bufio.Reader) (string, error) {
	var r rune
	var err error
	var buf bytes.Buffer
	for !unicode.IsSpace(r) {
		r, _, err = in.ReadRune()
		if err != nil {
			return buf.String(), err
		}
		if _, err := buf.WriteRune(r); err != nil {
			return buf.String(), err
		}
	}
	return buf.String(), nil
}

func parseToken(in *bufio.Reader) (Expr, error) {
	// Remove all whitespace at the beginning of the token
	r := ' '
	for r == ' ' || r == '\n' {
		var err error
		r, _, err = in.ReadRune()
		if err != nil {
			return nil, err
		}
	}
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
	token, err := readUntilWhitespace(in)
	if err != nil {
		return nil, err
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
