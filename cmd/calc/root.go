package main

import (
	"fmt"
	"os"

	"github.com/EwanClarke/postfixcalc/internal/engine"
	"github.com/EwanClarke/postfixcalc/internal/tui/calc_ui"
	"github.com/spf13/cobra"
)

var verbose bool

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display output from tokenisation and postfix steps")
}

var rootCmd = &cobra.Command{
	Use:   "calc [expression]",
	Short: "A postfix calculator",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// launch TUI
			if err := calcui.Start(); err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
			return
		}

		input := args[0]
		executeLogic(input)
	},
}

func executeLogic(input string) {
	e := engine.NewEngine()
	res, err := e.Calculate(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if verbose {
		l := engine.NewLexer()
		tokens, err := l.Tokenise(input)
		if err != nil {
			fmt.Printf("Lexical Error: %v\n", err)
			return
		}

		p := engine.NewParser()
		postfix, err := p.Convert(tokens)
		if err != nil {
			fmt.Printf("Parser Error: %v\n", err)
			return
		}

		fmt.Printf("Input: %v\n", engine.CanonicalString(tokens))
		fmt.Printf("Tokens: %v\n", tokens)
		fmt.Printf("Postfix: %v\n", engine.CanonicalString(postfix))
		fmt.Print("Result: ")
	}

	fmt.Println(res)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
