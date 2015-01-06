package optimization

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNelderMead(t *testing.T) {
	fn := func(vars []float64) (float64, bool) {
		return vars[0] * vars[1], true
	}
	config := NelderMeadConfiguration{
		Target: float64(9),
		Fn:     fn,
		Vars:   []float64{2, 4},
	}

	result, _ := fn(NelderMead(config))
	assert.True(t, math.Abs(result-config.Target) <= .01)
}

func TestNelderMeadPolynomial(t *testing.T) {
	fn := func(vars []float64) (float64, bool) {
		// x^2-4x+y^2-y-xy, solution is (3, 2)
		return math.Pow(vars[0], 2) - 4*vars[0] + math.Pow(vars[1], 2) - vars[1] - vars[0]*vars[1], true
	}
	config := NelderMeadConfiguration{
		Target: float64(-100),
		Fn:     fn,
		Vars:   []float64{-10, 10},
	}

	result := NelderMead(config)
	calced, _ := fn(result)
	assert.True(t, math.Abs(7-math.Abs(calced)) <= .01)
	assert.True(t, math.Abs(3-result[0]) <= .1)
	assert.True(t, math.Abs(2-result[1]) <= .1)
}

func TestNelderMeadPolynomialMin(t *testing.T) {
	fn := func(vars []float64) (float64, bool) {
		// x^2-4x+y^2-y-xy, solution is (3, 2)
		return math.Pow(vars[0], 2) - 4*vars[0] + math.Pow(vars[1], 2) - vars[1] - vars[0]*vars[1], true
	}
	config := NelderMeadConfiguration{
		Target: math.Inf(-1),
		Fn:     fn,
		Vars:   []float64{-10, 10},
	}

	result := NelderMead(config)
	calced, _ := fn(result)
	assert.True(t, math.Abs(7-math.Abs(calced)) <= .01)
	assert.True(t, math.Abs(3-result[0]) <= .01)
	assert.True(t, math.Abs(2-result[1]) <= .01)
}

func TestNelderMeadPolynomialMax(t *testing.T) {
	fn := func(vars []float64) (float64, bool) {
		// 3+sin(x)+2cos(y)^2, the min on this equation is 2 and the max is 6
		return 3 + math.Sin(vars[0]) + 2*math.Pow(math.Cos(vars[1]), 2), true
	}

	config := NelderMeadConfiguration{
		Target: math.Inf(1),
		Fn:     fn,
		Vars:   []float64{-5, 5},
	}

	result := NelderMead(config)
	calced, _ := fn(result)
	assert.True(t, math.Abs(6-math.Abs(calced)) <= .01)
}

func TestNelderMeadConstrained(t *testing.T) {
	fn := func(vars []float64) (float64, bool) {
		if vars[0] < 1 || vars[1] < 1 {
			return 0, false
		}
		return math.Pow(vars[0], 2) - 4*vars[0] + math.Pow(vars[1], 2) - vars[1] - vars[0]*vars[1], true
	}
	// by default, converging on this point with the initial
	// guess of (6, 3) will converge to (~.46, ~4.75).  The
	// fn has the added constraint that no guesses may be below
	// 1.  This should now converge to a point (~8.28, ~4.93).
	config := NelderMeadConfiguration{
		Target: float64(14),
		Fn:     fn,
		Vars:   []float64{6, 3},
	}

	result := NelderMead(config)
	calced, _ := fn(result)
	assert.True(t, math.Abs(14-math.Abs(calced)) <= .01)
	assert.True(t, result[0] >= 1)
	assert.True(t, result[1] >= 1)

	fn = func(vars []float64) (float64, bool) {
		if vars[0] < 6 || vars[0] > 8 {
			return 0, false
		}

		if vars[1] < 0 || vars[1] > 2 {
			return 0, false
		}
		return math.Pow(vars[0], 2) - 4*vars[0] + math.Pow(vars[1], 2) - vars[1] - vars[0]*vars[1], true
	}

	config = NelderMeadConfiguration{
		Target: float64(14),
		Fn:     fn,
		Vars:   []float64{6, .5},
	}

	result = NelderMead(config)
	calced, _ = fn(result)
	// there are two local min here
	assert.True(t, math.Abs(14-math.Abs(calced)) <= .01 || math.Abs(8.75-math.Abs(calced)) <= .01)
	assert.True(t, result[0] >= 6 && result[0] <= 8)
	assert.True(t, result[1] >= 0 && result[1] <= 2)
}

func TestNelderMeadConstrainedBadGuess(t *testing.T) {
	fn := func(vars []float64) (float64, bool) {
		if vars[0] < 1 || vars[1] < 1 {
			return 0, false
		}
		return math.Pow(vars[0], 2) - 4*vars[0] + math.Pow(vars[1], 2) - vars[1] - vars[0]*vars[1], true
	}
	// this is a bad guess, as in the initial guess doesn't
	// match the constraints.  In that case, we return the guessed
	// values.
	config := NelderMeadConfiguration{
		Target: float64(14),
		Fn:     fn,
		Vars:   []float64{0, 3},
	}

	result := NelderMead(config)
	assert.Equal(t, 0, result[0])
	assert.Equal(t, 3, result[1])
}
