# List of features to implement

## Lexer
- Support UNICODE + UFT8 encoding instead of only ASCII. This requires switchingfrom `byte` to `rune` reading
- Allow integers as part of a variable or function name (but only is not solely composed of int chars)
- Support ++, --, +=, -=, *=, /=

## Parser
- indicate source line and col when reporting errors (impacts lexer)
