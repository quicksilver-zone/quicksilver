package types_test

import (
	"errors"
	"testing"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func TestErrorsDeterminism(t *testing.T) {
	e := multierror.Errors{
		Errors: map[string]error{
			"a":    errors.New("a"),
			"Z":    errors.New("Z"),
			"🚨":    errors.New("🚨"),
			"a🚨":   errors.New("a🚨"),
			"ABC":  errors.New("ABC"),
			"1one": errors.New("1one"),
			"A":    errors.New("A"),
			"X":    errors.New("X"),
		},
	}

	e0 := e.Error()

	for i := 0; i < 100; i++ {
		ei := e.Error()
		if ei != e0 {
			t.Errorf("Iteration #%d produced non-deterministic data\n\tGot:  %q\n\tWant: %q", i+1, ei, e0)
		}
	}
}
