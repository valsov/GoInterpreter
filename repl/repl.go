package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/valsov/gointerpreter/evaluator"
	"github.com/valsov/gointerpreter/lexer"
	"github.com/valsov/gointerpreter/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, " --- Parser errors:\n")
	for _, err := range errors {
		io.WriteString(out, fmt.Sprintf("    %s\n", err))
	}
}
