# postfixcalc

## About

A terminal-based calculator that converts infix expressions to postfix (reverse polish notation) for unambiguous evaluation without precedence rules. It provides both a CLI interface for quick single-expression calculations and a TUI interface with a visual calculator layout navigable via arrow keys or vim-style bindings (hjkl).

## Features

This project supports:
- basic mathematical operations (+-*/^)
- brackets
- unary functions (sin, cos, tan)

Unary functions currently only implemented for the CLI and direct input in edit mode, with plans to give them their own buttons.

## Quick Start

### Installation

```
go install github.com/EwanClarke/postfixcalc@latest
```

## Examples

The calculator can be used from either the command line by following the command with a string containing the expression to evaluate, or through the TUI interface giving a more visual experience either through navigation of the on screen buttons or direct input with immediate evaluation of the current input.

### CLI Examples
#### Direct calculation
```
❯ calc "2(3+2)"
10
```
##### Verbose
The "-v" or "--verbose" flag can be used to view the output of each step
```
❯ calc "4 + 3 * 6" -v
Input: 4+3*6
Tokens: [{Number 4} {Operator +} {Number 3} {Operator *} {Number 6}]
Postfix: 436*+
Result: 22
```
#### Open TUI

```
❯ calc
```

### TUI Demo

<img src="demo.gif" width="600" alt="TUI Demo">

## Technical Architecture

The project follows a clean architecture pattern with clear separation of concerns between interfaces and core logic.

### Components

- **Lexer**: Tokenizes input strings into numbers, operators, and functions
- **Parser**: Converts infix expressions to postfix using the shunting yard algorithm
- **Evaluator**: Evaluates postfix expressions using a stack-based approach
- **TUI**: Interactive terminal interface built with the bubbletea framework

### Algorithms
This project makes use of a self implemented shunting yard algorithm along with stack based evaluation, allowing the program to handle complex expressions with strictly following BODMAS rules.

## Project Structure

The project follows Go's standard project layout with clear architecture principles. The CLI and TUI interfaces are separated from the core calculation logic, which is separated into distinct packages for tokenisation, parsing and evaluation.

```
postfixcalc/
├── cmd/
│   └── calc/
│       ├── main.go          # Entry point
│       └── root.go          # CLI commands and flags
├── internal/
│   ├── lexer/
│   │   ├── lexer.go         # Tokenisation logic
│   │   ├── token.go         # Token definitions
│   │   └── lexer_test.go
│   ├── parser/
│   │   ├── shuntingyard.go  # Infix to postfix conversion
│   │   ├── rules.go         # Parser rules
│   │   └── shuntingyard_test.go
│   ├── evaluator/
│   │   ├── evaluator.go     # Postfix evaluation
│   │   └── evaluator_test.go
│   └── tui/                 # bubbletea TUI
│       ├── model.go
│       ├── view.go
│       ├── update.go
│       ├── styles.go
│       └── tui.go           # TUI entry point
├── go.mod
├── go.sum
└── README.md
```

## Testing

Unit tests have been implemented for low-level functions in each of the tokenisation, conversion and evaluation steps.
These tests can be run with the following command:
```
make test
```
verbose and coverage outputs are also obtainable using:
```
make test-verbose
```
and
```
make test-coverage
```

## Planned Additions

- Number values to be handled internally as bit.Rat instead of float64, allowing for fractional output in addition to decimal
- expandable function tray which hides function buttons and lesser used operations
- Commonly used constants such as e and pi
- standard output for built-in constants

## License

This project is licensed under the MIT License - see the LICENSE file for details.
