package optimization

import (
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
	assert.True(t, math.Abs(result-config.Target) <= .01)
}

func TestNelderMeadPolynomial(t *testing.T) {
	fn := func(vars []float64) float64 {
		// x^2-4x+y^2-y-xy, solution is (3, 2)
		return math.Pow(vars[0], 2) - 4*vars[0] + math.Pow(vars[1], 2) - vars[1] - vars[0]*vars[1]
	}
	config := NelderMeadConfiguration{
		Target: float64(-100),
		Fn:     fn,
		Vars:   []float64{-10, 10},
	}

	result := NelderMead(config)
	calced := fn(result)
	assert.True(t, math.Abs(7-math.Abs(calced)) <= .01)
	assert.True(t, math.Abs(3-result[0]) <= .1)
	assert.True(t, math.Abs(2-result[1]) <= .1)
}

func TestNelderMeadPolynomialMin(t *testing.T) {
	fn := func(vars []float64) float64 {
		// x^2-4x+y^2-y-xy, solution is (3, 2)
		return math.Pow(vars[0], 2) - 4*vars[0] + math.Pow(vars[1], 2) - vars[1] - vars[0]*vars[1]
	}
	config := NelderMeadConfiguration{
		Target: math.Inf(-1),
		Fn:     fn,
		Vars:   []float64{-10, 10},
	}

	result := NelderMead(config)
	calced := fn(result)
	assert.True(t, math.Abs(7-math.Abs(calced)) <= .01)
	assert.True(t, math.Abs(3-result[0]) <= .01)
	assert.True(t, math.Abs(2-result[1]) <= .01)
}

func TestNelderMeadPolynomialMax(t *testing.T) {
	fn := func(vars []float64) float64 {
		// 3+sin(x)+2cos(y)^2, the min on this equation is 2 and the max is 6
		return 3 + math.Sin(vars[0]) + 2*math.Pow(math.Cos(vars[1]), 2)
	}

	config := NelderMeadConfiguration{
		Target: math.Inf(1),
		Fn:     fn,
		Vars:   []float64{-5, 5},
	}

	result := NelderMead(config)
	calced := fn(result)
	assert.True(t, math.Abs(6-math.Abs(calced)) <= .01)
	assert.True(t, math.Abs(4.7-math.Abs(result[0])) <= .1)
	assert.True(t, math.Abs(6.2-math.Abs(result[1])) <= .1)
}
