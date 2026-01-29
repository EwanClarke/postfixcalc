package evaluator
import (
	"github.com/EwanClarke/postfixcalc/lexer"
	"strconv"
	"errors"
	"fmt"
	"math"
)

type Evaluator struct {
	resultStack []float64
}

func New() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) pushValue(value float64) {
	e.resultStack = append(e.resultStack, value)
}

func (e *Evaluator) popValue() (float64, error) {
	if len(e.resultStack) == 0 {
		return 0, errors.New("Stack Underflow: attempted removal from empty stack")
	}

	n := len(e.resultStack)-1
	value := e.resultStack[n]
	e.resultStack = e.resultStack[:n]
	return value, nil
}

func (e *Evaluator) Evaluate(postfixTokens []lexer.Token) (float64, error) {
	for _, token := range postfixTokens {
		switch token.Type {
		case lexer.Number:
			val, _ := strconv.ParseFloat(token.Value, 64)
			e.pushValue(val)

		case lexer.Operator:
			if err := e.evaluateBinary(token); err != nil {
				return 0, err
			}

		case lexer.Negation, lexer.Function:
			if err := e.evaluateUnary(token); err != nil {
				return 0, err
			}
		}
	}
	err := e.validateResult()
	if err != nil {
		return 0, err
	}
	result, _ := e.popValue()
	return result, nil
}

func (e *Evaluator) evaluateBinary(token lexer.Token) error {
	if len(e.resultStack) < 2 {
		return fmt.Errorf("Stack Underflow: missing operands for '%s'", token.Value)
	}
	a, err := e.popValue()
	if err != nil {return err}
	b, err := e.popValue()
	if err != nil {return err}

	result, err := e.applyBinaryOp(token.Value, a, b)
	if err != nil {return err}

	e.pushValue(result)
	return nil
}

func (e *Evaluator) evaluateUnary(token lexer.Token) error {
	if len(e.resultStack) < 1 {
		return fmt.Errorf("Stack Underflow: missing operands for '%s'", token.Value)
	}

	a, err := e.popValue()
	if err != nil {return err}

	result, err := e.applyUnaryOp(token.Value, a)
	if err != nil {return err}

	e.pushValue(result)
	return nil
}

func (e *Evaluator) applyBinaryOp(operation string, a, b float64) (float64, error) {
	switch operation {
	case "+": return a + b, nil
	case "-": return a - b, nil
	case "*": return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("Math Error: division by zero")
		}
		return a / b, nil
	case "^":
		if a == 0 && b < 0 {
			return 0, errors.New("Math Error: zero cannot be raised to a negative power")
		}
		return math.Pow(a, b), nil
		default: return 0, fmt.Errorf("Unknown Operator: operator '%s' usage not defined", operation)
	}
}

func (e *Evaluator) applyUnaryOp(operation string, a float64) (float64, error) {
	switch operation {
	case "-u": return -a, nil
	case "sin": return math.Sin(a), nil
	case "cos": return math.Cos(a), nil
	case "tan": return math.Tan(a), nil
	default: return 0, fmt.Errorf("Unknown Operation: usage of Operation '%s' not defined", operation)
	}
}

func (e *Evaluator) validateResult() error {
	if len(e.resultStack) > 1 {
		return errors.New("Missing Operator(s)")
	} else if len(e.resultStack) == 0 {
		return errors.New("Empty expression")
	}
	return nil
}
