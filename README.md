# postfixcalc

A terminal based postfix calculator providing CLI and TUI interfaces. CLI interface provides quick access to a calculator for single expression, where the TUI provides a familiar calculator interface navigable through arrow keys or vim-style bindings (hjkl).

## About

<!-- Brief explanation of what the project does -->

## Features

This project supports basic mathematical operations (+-*/^), brackets and unary functions (sin, cos, tan). However, unary functions have currently only been implemented for the CLI and direct input in edit mode in the TUI, with plans to give them their own buttons.

## Quick Start

### Installation

```
go install github.com/EwanClarke/postfixcalc@latest
```

### Usage

<!-- CLI and TUI usage examples with code blocks -->

## Technical Architecture

<!-- Overview of clean architecture and algorithms -->

### Components

<!-- Brief description of lexer, parser, evaluator, TUI -->

### Algorithms
This project makes use of a self implemented shunting yard algorithm along with stack based evaluation, allowing the program to handle complex expressions with strictly following BODMAS rules.

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

<!-- Screenshot placeholder and description -->

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

<!-- Information about tests and coverage -->

## Contributing

<!-- Contribution guidelines -->

## License

<!-- License information -->
