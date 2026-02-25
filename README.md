# postfixcalc

## About

A terminal-based calculator that converts infix expressions to postfix (reverse polish notation) for unambiguous evaluation without precedence rules. It provides both a CLI interface for quick single-expression calculations and TUI interfaces with a visual calculator layout navigable via arrow keys or vim-style bindings (hjkl) in addition to a graph plotting TUI for visualizing mathematical functions.

## Features

This project supports:
- basic mathematical operations (+-*/^)
- brackets
- unary functions (sin, cos, tan)
- constants (e, Pi)

## Quick Start

### Installation

```
go install github.com/EwanClarke/postfixcalc/cmd/calc@latest
go install github.com/EwanClarke/postfixcalc/cmd/plot@latest
```

### Docker

Build the Docker image:
```
docker build -t postfixcalc .
```

Run the calculator:
```
docker run -it postfixcalc ./calc
docker run -it postfixcalc ./calc "2(3+2)"
```

Run the graph plotter:
```
docker run -it postfixcalc ./plot
docker run -it postfixcalc ./plot "x^2"
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

<img src="demo.gif" width="400" alt="TUI Demo">

## Technical Architecture

The project follows a clean architecture pattern with clear separation of concerns between interfaces and core logic.

### Components

- **Lexer**: Tokenizes input strings into numbers, operators, and functions
- **Parser**: Converts infix expressions to postfix using the shunting yard algorithm
- **Evaluator**: Evaluates postfix expressions using a stack-based approach
- **TUI**: Interactive terminal interfaces built with the bubbletea framework (calculator and graph plotter)

### Algorithms
This project makes use of a self implemented shunting yard algorithm along with stack based evaluation, allowing the program to handle complex expressions with strictly following BODMAS rules.

## Project Structure

The project follows Go's standard project layout with clear architecture principles. The CLI and TUI interfaces are separated from the core calculation logic, which is contained in the engine package.

```
postfixcalc/
├── cmd/
│   ├── calc/
│   │   ├── main.go          # Entry point
│   │   └── root.go          # CLI commands and flags
│   └── plot/
│       ├── main.go          # Entry point
│       └── root.go          # CLI commands and flags
├── internal/
│   ├── engine/              # Core calculation logic
│   │   ├── lexer.go         # Tokenisation logic
│   │   ├── token.go         # Token definitions
│   │   ├── parser.go        # Infix to postfix conversion
│   │   ├── evaluator.go     # Postfix evaluation
│   │   ├── constants.go    # Mathematical constants
│   │   └── engine_test.go
│   └── tui/
│       ├── calc_ui/         # Calculator TUI
│       │   ├── model.go
│       │   ├── view.go
│       │   ├── update.go
│       │   ├── styles.go
│       │   └── tui.go
│       └── graph_ui/        # Graph plotter TUI
│           ├── model.go
│           ├── view.go
│           ├── update.go
│           ├── styles.go
│           ├── canvas.go
│           └── tui.go
├── go.mod
├── go.sum
├── Dockerfile
├── entrypoint.sh
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

## License

This project is licensed under the MIT License - see the LICENSE file for details.
