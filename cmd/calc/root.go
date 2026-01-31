package main
import (
	"fmt"
	"os"

	"github.com/EwanClarke/postfixcalc/internal/lexer"
	"github.com/EwanClarke/postfixcalc/internal/parser"
	"github.com/EwanClarke/postfixcalc/internal/evaluator"
	"github.com/EwanClarke/postfixcalc/internal/tui"
	"github.com/spf13/cobra"
)

var verbose bool
func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display output from tokenisation and postfix steps")
}

var rootCmd = &cobra.Command {
	Use: "calc [expression]",
	Short: "A postfix calculator",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// launch TUI
			if err := tui.Start(); err != nil {
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
	// tokenise
	l := lexer.New(input)
	tokens, err := l.Tokenise()
	if err != nil {
		fmt.Printf("Lexical Error: %v\n", err)
		return
	}

	// parse (shunting yard)
	p := parser.New()
	postfix, err := p.Convert(tokens)
	if err != nil {
		fmt.Printf("Parser Error: %v\n", err)
		return
	}

	// evaluate
	e := evaluator.New()
	res, err := e.Evaluate(postfix)
	if err != nil {
		fmt.Printf("Evaluation Error: %v\n", err)
		return
	}
	
	if verbose {
		fmt.Printf("Input: %v\n", lexer.CanonicalString(tokens))
		fmt.Printf("Tokens: %v\n", tokens)
		fmt.Printf("Postfix: %v\n", lexer.CanonicalString(postfix))
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

