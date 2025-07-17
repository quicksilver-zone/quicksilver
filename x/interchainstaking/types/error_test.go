package types_test

import (
	"errors"
	"testing"

	"go.uber.org/multierr"
)

func TestErrorsDeterminism(t *testing.T) {
	errs := []error{
		errors.New("a"),
		errors.New("Z"),
		errors.New("🚨"),
		errors.New("a🚨"),
		errors.New("ABC"),
		errors.New("1one"),
		errors.New("A"),
		errors.New("X"),
	}

	e := multierr.Combine(errs...)
	e0 := e.Error()

	for i := 0; i < 100; i++ {
		ei := e.Error()
		if ei != e0 {
			t.Errorf("Iteration #%d produced non-deterministic data\n\tGot:  %q\n\tWant: %q", i+1, ei, e0)
		}
	}
}
