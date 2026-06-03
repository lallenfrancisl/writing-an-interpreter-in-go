// Package repl contains the REPL implementation for monkey lang
package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/internal/evaluator"
	"monkey/internal/lexer"
	"monkey/internal/object"
	"monkey/internal/parser"
)

const (
	PROMPT      = ">>"
	MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()

	for {
		fmt.Fprintf(out, "%s ", PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(out, "scanner error: %v\n", err)
			}

			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())

			return
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, "parser errors:\n")

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
