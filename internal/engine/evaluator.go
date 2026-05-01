package engine

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

type Evaluator struct {
	resultStack []*big.Rat
	exact       bool
	angleMode   AngleMode
}

func NewEvaluator(angleMode AngleMode) *Evaluator {
	return &Evaluator{exact: true, angleMode: angleMode}
}

func (e *Evaluator) pushValue(value *big.Rat) {
	e.resultStack = append(e.resultStack, value)
}

func (e *Evaluator) popValue() (*big.Rat, error) {
	if len(e.resultStack) == 0 {
		return nil, errors.New("Stack Underflow: attempted removal from empty stack")
	}

	n := len(e.resultStack) - 1
	value := e.resultStack[n]
	e.resultStack = e.resultStack[:n]
	return value, nil
}

// Evaluate evaluates a postfix token stream.
// Returns the result as a *big.Rat, a bool indicating whether the result is
// exact (true) or required a float64 approximation (false), and any error.
func (e *Evaluator) Evaluate(postfixTokens []Token, vars map[string]*big.Rat) (*big.Rat, bool, error) {
	for _, token := range postfixTokens {
		switch token.Type {
		case Number:
			r := new(big.Rat)
			if _, ok := r.SetString(token.Value); !ok {
				// Fallback: parse via float64 for values like "3.14"
				f, err := strconv.ParseFloat(token.Value, 64)
				if err != nil {
					return nil, false, fmt.Errorf("invalid number token: %s", token.Value)
				}
				r.SetFloat64(f)
				e.exact = false
			}
			e.pushValue(r)

		case Variable:
			val, ok := vars[token.Value]
			if !ok {
				return nil, false, fmt.Errorf("variable %s not defined", token.Value)
			}
			e.pushValue(new(big.Rat).Set(val))

		case Operator:
			if err := e.evaluateBinary(token); err != nil {
				return nil, false, err
			}

		case Negation, Function:
			if err := e.evaluateUnary(token); err != nil {
				return nil, false, err
			}
		}
	}
	if err := e.validateResult(); err != nil {
		return nil, false, err
	}
	result, _ := e.popValue()
	return result, e.exact, nil
}

func (e *Evaluator) evaluateBinary(token Token) error {
	if len(e.resultStack) < 2 {
		return fmt.Errorf("Stack Underflow: missing operands for '%s'", token.Value)
	}
	b, err := e.popValue()
	if err != nil {
		return err
	}
	a, err := e.popValue()
	if err != nil {
		return err
	}

	result, err := e.applyBinaryOp(token.Value, a, b)
	if err != nil {
		return err
	}

	e.pushValue(result)
	return nil
}

func (e *Evaluator) evaluateUnary(token Token) error {
	if len(e.resultStack) < 1 {
		return fmt.Errorf("Stack Underflow: missing operands for '%s'", token.Value)
	}

	a, err := e.popValue()
	if err != nil {
		return err
	}

	result, err := e.applyUnaryOp(token.Value, a)
	if err != nil {
		return err
	}

	e.pushValue(result)
	return nil
}

func (e *Evaluator) applyBinaryOp(operation string, a, b *big.Rat) (*big.Rat, error) {
	result := new(big.Rat)
	switch operation {
	case "+":
		result.Add(a, b)
	case "-":
		result.Sub(a, b)
	case "*":
		result.Mul(a, b)
	case "/":
		if b.Sign() == 0 {
			return nil, errors.New("Math Error: division by zero")
		}
		result.Quo(a, b)
	case "^":
		if b.IsInt() && b.Num().IsInt64() {
			e64 := b.Num().Int64()
			// To prevent massive allocations, limit exact exponentiation to a reasonable range
			if e64 >= -10000 && e64 <= 10000 {
				if e64 == 0 {
					result.SetInt64(1)
					return result, nil
				}
				isNeg := e64 < 0
				if isNeg {
					e64 = -e64
				}
				num := new(big.Int).Exp(a.Num(), big.NewInt(e64), nil)
				den := new(big.Int).Exp(a.Denom(), big.NewInt(e64), nil)
				if isNeg {
					if num.Sign() == 0 {
						return nil, errors.New("Math Error: zero cannot be raised to a negative power")
					}
					result.SetFrac(den, num)
				} else {
					result.SetFrac(num, den)
				}
				return result, nil
			}
		}

		// Power is generally irrational; round-trip through float64.
		af, _ := a.Float64()
		bf, _ := b.Float64()
		if af == 0 && bf < 0 {
			return nil, errors.New("Math Error: zero cannot be raised to a negative power")
		}
		result.SetFloat64(math.Pow(af, bf))
		e.exact = false
	default:
		return nil, fmt.Errorf("Unknown Operator: operator '%s' usage not defined", operation)
	}
	return result, nil
}

func (e *Evaluator) applyUnaryOp(operation string, a *big.Rat) (*big.Rat, error) {
	switch operation {
	case "-u":
		// Negation is exact.
		return new(big.Rat).Neg(a), nil
	case "sin", "cos", "tan":
		// Transcendental — float64 round-trip.
		af, _ := a.Float64()
		if e.angleMode == Degrees {
			af = af * math.Pi / 180.0
		}
		var f float64
		switch operation {
		case "sin":
			f = math.Sin(af)
		case "cos":
			f = math.Cos(af)
		case "tan":
			f = math.Tan(af)
		}
		e.exact = false
		return new(big.Rat).SetFloat64(f), nil
	default:
		return nil, fmt.Errorf("Unknown Operation: usage of Operation '%s' not defined", operation)
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
