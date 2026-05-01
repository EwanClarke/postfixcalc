package engine

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

// FormatMixed formats a *big.Rat as a mixed fraction string.
//
// Examples:
//
//	7/2   → "3 1/2"
//	-7/2  → "-3 1/2"
//	6/3   → "2"
//	1/3   → "1/3"
//	0     → "0"
var piRat = new(big.Rat).SetFloat64(math.Pi)

func FormatMixed(r *big.Rat) string {
	if r.Sign() == 0 {
		return "0"
	}

	// Check if the result is a simple rational multiple of Pi
	coeff := new(big.Rat).Quo(r, piRat)
	if coeff.Num().IsInt64() && coeff.Denom().IsInt64() {
		n := coeff.Num().Int64()
		d := coeff.Denom().Int64()
		if n > -100000 && n < 100000 && d < 100000 {
			if n == 1 && d == 1 {
				return "π"
			} else if n == -1 && d == 1 {
				return "-π"
			} else if d == 1 {
				return fmt.Sprintf("%dπ", n)
			} else if n == 1 {
				return fmt.Sprintf("π/%d", d)
			} else if n == -1 {
				return fmt.Sprintf("-π/%d", d)
			} else {
				return fmt.Sprintf("%dπ/%d", n, d)
			}
		}
	}

	negative := r.Sign() < 0

	// Work with the absolute value.
	abs := new(big.Rat).Abs(new(big.Rat).Set(r))

	num := abs.Num()   // numerator   (always >= 0 after Abs)
	den := abs.Denom() // denominator (always > 0)

	whole := new(big.Int).Quo(num, den)
	remainder := new(big.Int).Rem(num, den)

	var sb strings.Builder
	if negative {
		sb.WriteByte('-')
	}

	if whole.Sign() > 0 && remainder.Sign() > 0 {
		// Mixed: whole + fraction  e.g. "3 1/2"
		sb.WriteString(whole.String())
		sb.WriteByte(' ')
		sb.WriteString(remainder.String())
		sb.WriteByte('/')
		sb.WriteString(den.String())
	} else if whole.Sign() > 0 {
		// Whole number  e.g. "3"
		sb.WriteString(whole.String())
	} else {
		// Proper fraction (|r| < 1)  e.g. "1/3"
		sb.WriteString(num.String())
		sb.WriteByte('/')
		sb.WriteString(den.String())
	}

	return sb.String()
}

// FormatDecimal formats a *big.Rat as a decimal string using %g formatting,
// matching the previous float64 output behaviour.
func FormatDecimal(r *big.Rat) string {
	f, _ := r.Float64()
	return fmt.Sprintf("%g", f)
}
