package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndEval(t *testing.T) {
	tests := []struct {
		In   string
		Out  []Expr
		Eval Expr
	}{
		{
			In: `(+ 1 1 1)`,
			Out: []Expr{
				[]Expr{"+", Number(1), Number(1), Number(1)}},
			Eval: Number(3),
		},
		{
			In: `(define r
10)
(* r r)`,
			Out: []Expr{
				[]Expr{"define", "r", Number(10)},
				[]Expr{"*", "r", "r"},
			},
			Eval: Number(100),
		},
		{
			In: `(define circle-area (lambda (r) (* pi (* r r))))
(circle-area 3)`,
			Out: []Expr{
				[]Expr{"define", "circle-area", []Expr{
					"lambda", []Expr{"r"}, []Expr{
						"*", "pi", []Expr{"*", "r", "r"}}}},
				[]Expr{"circle-area", Number(3)}},
			Eval: Number(28.274333882308138),
		},
	}
	for _, test := range tests {
		// Parse the whole thing
		parser := NewParser(bytes.NewBufferString(test.In))
		out := []Expr{}
		for parser.Scan() {
			out = append(out, parser.Expression())
		}
		assert.Equal(t, test.Out, out)
		// Eval each expression
		env := standardEnv()
		var result Expr
		for _, expr := range out {
			result = eval(expr, env)
		}
		assert.Equal(t, test.Eval, result)
	}
}
