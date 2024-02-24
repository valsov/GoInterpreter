# Interpreter implemented in Go
Study on interpreters, featuring a complete dynamically-typed language. This work is based on the book [Writing An Interpreter In Go, by Thorsten Ball](https://interpreterbook.com/). The project is fully unit-tested.

## Structure of the project
- `/ast` All available data structures representing the evaluated program
- `/evaluator` Navigates through the AST in order to evaluate its nodes
- `/lexer` Produces tokens from chars, it is responsible of syntax checking
- `/object` Data structures representing the execution results of the AST by the evaluator
- `/parser` The role of the parser is to produce an AST of the program from tokens
- `/repl` Read-Eval-Print Loop, which enables direct development in this language
- `/token` Base representation of the code: a collection of tokens

**Execution flow of the interpreter**:

`string` =[`lexer`]=> `[]token` =[`parser`]=> `ast` =[`evaluator`]=> `stream of objects`

## List of additional features to implement
### Lexer
- [ ] Support UNICODE + UFT8 encoding instead of only ASCII. This requires switchingfrom `byte` to `rune` reading
- [ ] Allow integers as part of a variable or function name (but only is not solely composed of int chars)
- [ ] Support ++, --, +=, -=, *=, /=
- [ ] Support modulo
- [ ] Allow ternary operators

### Parser
- [ ] indicate source line and col when reporting errors (impacts lexer)
- [ ] Support else if(...)