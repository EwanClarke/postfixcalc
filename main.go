package main
import (
	"fmt"
	"github.com/EwanClarke/postfixcalc/lexer"
	"github.com/EwanClarke/postfixcalc/parser"
	// "github.com/EwanClarke/postfixcalc/evaluator"
)

func main() {
	input := "2(1+1)"
	tokens := []lexer.Token{}
	
	// tokenise
	l := lexer.New(input)
	tokens, err := l.Tokenise()
	if err != nil {
		fmt.Printf("Lexical Error: %v\n", err)
		return
	}
	
	// parse (shunting yard)
	p := parser.New(tokens)
	postfixTokens, err := p.Convert(tokens)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// evaluate
	// e := evaluator.New(postfixTokens)
	// res := e.evaluate()
	// 
	// fmt.Printf("Result: %f\n", res)
	
	fmt.Printf("Raw input: %s\n", input)
	fmt.Printf("Tokenised: %v\n",tokens)
	fmt.Printf("Postfix: %v\n",postfixTokens)
}
