package engine

import "math"

// InternalVars holds mathematical constants as float64.
// The engine converts these to *big.Rat at initialisation time.
// Expressions using these constants (e.g. sin(Pi)) will produce inexact results.
var InternalVars = map[string]float64{
	"Pi": math.Pi,
	"e":  math.E,
}
