package optimization

import (
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNelderMead(t *testing.T) {
	fn := func(vars []float64) float64 {
		return vars[0] * vars[1]
	}
	config := NelderMeadConfiguration{
		Target: float64(9),
		Fn:     fn,
		Vars:   []float64{2, 4},
	}

	result := fn(NelderMead(config))
	assert.True(t, math.Abs(result-config.Target) <= delta)
}

func TestNelderMeadPolynomial(t *testing.T) {
	fn := func(vars []float64) float64 {
		// x^2-4x+y^2-y-xy, solution is (3, 2)
		return math.Pow(vars[0], 2) - 4*vars[0] + math.Pow(vars[1], 2) - vars[1] - vars[0]*vars[1]
	}
	config := NelderMeadConfiguration{
		Target: float64(-1000000),
		Fn:     fn,
		Vars:   []float64{10, 10},
	}

	result := NelderMead(config)
	log.Printf(`CALED: %+v`, fn([]float64{2.5, 2.5}))
	log.Printf(`RESULT: %+v`, result)
}
