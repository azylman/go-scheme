package main

type Expr interface{}
type Number float64
type Symbol string
type Environment map[Symbol]Expr

func (e Environment) copy() Environment {
	out := Environment{}
	for k, v := range e {
		out[k] = v
	}
	return out
}

type Procedure func(Environment, ...Expr) Expr
