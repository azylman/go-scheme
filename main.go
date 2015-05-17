package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()
	env := standardEnv()
	if len(flag.Args()) == 0 {
		// repl mode
		parser := NewParser(os.Stdin)
		for fmt.Print("> "); parser.Scan(); fmt.Print("> ") {
			fmt.Println(eval(parser.Expression(), env))
		}
		if err := parser.Err(); err != nil {
			log.Fatal("error parsing: %s", err.Error())
		}
	} else {
		// interpret code in file
		f, err := os.Open(flag.Args()[0])
		if err != nil {
			log.Fatalf("error opening file: %s", err.Error())
		}
		var result Expr
		parser := NewParser(f)
		for parser.Scan() {
			result = eval(parser.Expression(), env)
		}
		if err := parser.Err(); err != nil {
			log.Fatalf("error parsing: %s", err.Error())
		}
		fmt.Println(result)
	}
}
