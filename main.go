package main
import (
	"fmt"
	"github.com/EwanClarke/postfixcalc/lexer"
	// "github.com/EwanClarke/postfixcalc/parser"
	// "github.com/EwanClarke/postfixcalc/evaluator"
)

func main() {
	input := "1+1"
	tokens := []lexer.Token{}
	
	// tokenise
	l := lexer.New(input)
	tokens, err := l.Tokenise()
	if err != nil {
		fmt.Printf("Lexical Error: %v\n", err)
		return
	}
	
	fmt.Println(tokens)
	
	// parse (shunting yard)
	// postfixTokens, err := parser.Convert()
	// if err != nil {
	// 	fmt.Printf("Math Error: %v\n", err)
	// 	return
	// }

	// evaluate
	// e := evaluator.New(postfixTokens)
	// res := e.evaluate()
	// 
	// fmt.Printf("Result: %f\n", res)
}
