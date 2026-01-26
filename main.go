package main
import (
	"fmt"
	"github.com/EwanClarke/postfixcalc/lexer"
)

func main() {
	var input string = "1 + 1"
	tokens := []lexer.Token{}
	
	// tokenise
	l := lexer.New(input)
	tokens = l.Tokenise()

	fmt.Println(tokens)

	// shunting yard
	// evaluate
}
