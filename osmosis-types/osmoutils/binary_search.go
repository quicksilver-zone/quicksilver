package osmoutils

import (
	"errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrTolerance is used to define a compare function, which checks if two
// ints are within a certain error tolerance of one another.
// ErrTolerance.Compare(a, b) returns true iff:
// |a - b| <= AdditiveTolerance
// |a - b| / min(a, b) <= MultiplicativeTolerance
// Each check is respectively ignored if the entry is nil (sdk.Dec{}, math.Int{})
// Note that if AdditiveTolerance == 0, then this is equivalent to a standard compare.
type ErrTolerance struct {
	AdditiveTolerance       math.Int
	MultiplicativeTolerance sdk.Dec
}

// Compare returns if actual is within errTolerance of expected.
// returns 0 if it is
// returns 1 if not, and expected > actual.
// returns -1 if not, and expected < actual
func (e ErrTolerance) Compare(expected, actual math.Int) int {
	diff := expected.Sub(actual).Abs()

	comparisonSign := 0
	if expected.GT(actual) {
		comparisonSign = 1
	} else {
		comparisonSign = -1
	}

	// Check additive tolerance equations
	if !e.AdditiveTolerance.IsNil() {
		// if no error accepted, do a direct compare.
		if e.AdditiveTolerance.IsZero() {
			if expected.Equal(actual) {
				return 0
			}
		}

		if diff.GT(e.AdditiveTolerance) {
			return comparisonSign
		}
	}
	// Check multiplicative tolerance equations
	if !e.MultiplicativeTolerance.IsNil() && !e.MultiplicativeTolerance.IsZero() {
		errTerm := sdk.NewDecFromInt(diff).Quo(sdk.NewDecFromInt(sdk.MinInt(expected, actual)))
		if errTerm.GT(e.MultiplicativeTolerance) {
			return comparisonSign
		}
	}

	return 0
}

// Binary search inputs between [lowerbound, upperbound] to a monotonic increasing function f.
// We stop once f(found_input) meets the ErrTolerance constraints.
// If we perform more than maxIterations (or equivalently lowerbound = upperbound), we return an error.
func BinarySearch(f func(input math.Int) (math.Int, error),
	lowerbound math.Int,
	upperbound math.Int,
	targetOutput math.Int,
	errTolerance ErrTolerance,
	maxIterations int,
) (math.Int, error) {
	// Setup base case of loop
	curEstimate := lowerbound.Add(upperbound).QuoRaw(2)
	curOutput, err := f(curEstimate)
	if err != nil {
		return math.Int{}, err
	}
	curIteration := 0
	for ; curIteration < maxIterations; curIteration += 1 {
		compRes := errTolerance.Compare(curOutput, targetOutput)
		if compRes > 0 {
			upperbound = curEstimate
		} else if compRes < 0 {
			lowerbound = curEstimate
		} else {
			break
		}
		curEstimate = lowerbound.Add(upperbound).QuoRaw(2)
		curOutput, err = f(curEstimate)
		if err != nil {
			return math.Int{}, err
		}
	}
	if curIteration == maxIterations {
		return math.Int{}, errors.New("hit maximum iterations, did not converge fast enough")
	}
	return curEstimate, nil
}
