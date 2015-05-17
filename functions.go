package main

import (
	"fmt"
	"math"
)

func standardEnv() Environment {
	return map[Symbol]Expr{
		"+": noEnv(func(exprs ...Expr) Expr {
			sum := Number(0)
			for _, ex := range exprs {
				sum += ex.(Number)
			}
			return Number(sum)
		}),
		"*": noEnv(func(exprs ...Expr) Expr {
			out := Number(1)
			for _, ex := range exprs {
				out *= ex.(Number)
			}
			return Number(out)
		}),
		"pi": Number(math.Pi),
		">": noEnv(twoArg(">", func(one, two Expr) Expr {
			return one.(Number) > two.(Number)
		})),
	}
}

func twoArg(name string, fn func(Expr, Expr) Expr) func(...Expr) Expr {
	return func(exprs ...Expr) Expr {
		if len(exprs) != 2 {
			panic(fmt.Errorf("can only call %s with two arguments", name))
		}
		return fn(exprs[0], exprs[1])
	}
}

func noEnv(fn func(...Expr) Expr) Procedure {
	return func(e Environment, exprs ...Expr) Expr {
		return fn(exprs...)
	}
}
