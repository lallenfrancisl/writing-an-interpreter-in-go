// Package repl contains the REPL implementation for monkey lang
package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/internal/compiler"
	"monkey/internal/lexer"
	"monkey/internal/parser"
	"monkey/internal/vm"
)

const (
	Prompt     = ">>"
	MonkeyFace = `
						__,__
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

	for {
		fmt.Fprintf(out, "%s ", Prompt)

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

		comp := compiler.New()

		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)

			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)

			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MonkeyFace)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, "parser errors:\n")

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
