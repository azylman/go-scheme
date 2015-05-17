package main

import "fmt"

func eval(ex Expr, env Environment) Expr {
	switch e := ex.(type) {
	case Symbol:
		res, ok := env[e]
		if !ok {
			panic(fmt.Errorf("unknown symbol '%s'", e))
		}
		return res
	case Number:
		return e
	case []Expr:
		switch e[0].(Symbol) {
		case "quote":
			return e[1]
		case "if":
			test, conseq, alt := e[1], e[2], e[3]
			if eval(test, env).(bool) {
				return eval(conseq, env)
			}
			return eval(alt, env)
		case "define":
			variable, exp := e[1], e[2]
			env[variable.(Symbol)] = eval(exp, env)
		case "set!":
			variable, exp := e[1], e[2]
			v := variable.(Symbol)
			if _, ok := env[v]; !ok {
				panic(fmt.Errorf("unknown symbol '%s'", e))
			}
			env[v] = eval(exp, env)
		case "lambda":
			params, body := e[1], e[2]
			return Procedure(func(env Environment, exprs ...Expr) Expr {
				var args []Expr
				switch p := params.(type) {
				case []Expr:
					args = p
				case Symbol:
					args = []Expr{p}
				default:
					panic(fmt.Errorf("must declare functions with symbol arguments"))
				}
				env = env.copy()
				if len(args) != len(exprs) {
					panic(fmt.Errorf("wrong number of arguments"))
				}
				for i, arg := range args {
					env[arg.(Symbol)] = exprs[i]
				}
				return eval(body, env)
			})
		default:
			proc, ok := eval(e[0], env).(Procedure)
			if !ok {
				panic(fmt.Errorf("'%s' isn't a procedure", e[0]))
			}
			args := make([]Expr, len(e)-1)
			for i, arg := range e[1:] {
				args[i] = eval(arg, env)
			}
			return proc(env, args...)
		}
	default:
		panic(fmt.Errorf("unknown expression type %#v", ex))
	}
	return nil
}
