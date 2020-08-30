package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/crockeo/lisg/repl"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	context := repl.BaseContext()

	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("encountered error: %s\n", err)
			continue
		}

		symbols := repl.Lex(text)

		ast, err := repl.Parse(symbols)
		if err != nil {
			fmt.Printf("encountered error in parsing: %s\n", err)
			continue
		}

		result, err := repl.Eval(context, context, ast)
		if err != nil {
			fmt.Printf("encountered error in evaluating: %s\n", err)
		}

		fmt.Println(result)
	}
}
