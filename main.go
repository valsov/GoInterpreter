package main

import (
	"fmt"
	"os"

	"github.com/valsov/gointerpreter/repl"
)

func main() {
	fmt.Fprintln(os.Stdout, "REPL instance")
	repl.Start(os.Stdin, os.Stdout)
}
