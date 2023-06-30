package failsim

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ctxKey uint

const (
	KeyFailSim ctxKey = iota
)

// FailuresFromContext returns a map of maps containing indicated failure
// hooks. The first map indexes the handler, the second map the failure hook.
func FailuresFromContext(ctx context.Context) map[uint8]map[uint8]struct{} {
	failures := make(map[uint8]map[uint8]struct{})
	failuresAny := ctx.Value(KeyFailSim)
	if failuresAny != nil {
		fmt.Println("nothing in context")
		failures = failuresAny.(map[uint8]map[uint8]struct{})
	}
	fmt.Printf("FailuresFromContext: %v\n", failures)
	return failures
}

// SetFailureContext returns a derived context containing simulated failure
// state. It is a map of maps containing indicated failure hooks. The first map
// indexes the handler, the second map the failure hook.
//
// The parameter useHooks is a string that denotes specific failure hooks, as a
// comma separated string, where each value denotes the handler:hook
// relationship, e.g. "0:1,1:9"...
func SetFailureContext(ctx context.Context, useHooks string) (context.Context, error) {
	failures := make(map[uint8]map[uint8]struct{})

	hooksExpr := regexp.MustCompile(`\d+:\d+(,\d+:\d+)*`)
	if !hooksExpr.MatchString(useHooks) {
		return ctx, fmt.Errorf("invalid useHooks format, must match expression %s", hooksExpr.String())
	}

	modstrs := strings.Split(useHooks, ",")
	for i, str := range modstrs {
		hookstrs := strings.Split(str, ":")
		mi, err := strconv.Atoi(hookstrs[0])
		if err != nil {
			return ctx, fmt.Errorf("invalid useHooks format, token %d, %w", i, err)
		}
		hi, err := strconv.Atoi(hookstrs[1])
		if err != nil {
			return ctx, fmt.Errorf("invalid useHooks format, token %d, %w", i, err)
		}
		if _, ok := failures[uint8(mi)]; !ok {
			failures[uint8(mi)] = make(map[uint8]struct{})
		}
		failures[uint8(mi)][uint8(hi)] = struct{}{}
	}

	fmt.Println("setFailureContext", failures)

	return context.WithValue(ctx, KeyFailSim, failures), nil
}

// FailureHook is used to intercept any potential error state and return a
// simulated error if, and only if, no error is presented (i.e. "err == nil").
//
// The given failures map, along with the provided hook index, is used to determine if the simulated failure error should be generated and
// returned.
//
// The simulated failure is deleted from the failures map to ensure single
// failure executions in multi error environments. Thus, failures maps are to
// be used on a per thread basis (i.e. not thread safe).
func FailureHook(failures map[uint8]struct{}, hook uint8, err error, simErrMsg string) error {
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("simulated error [%d]", hook)
	if _, failHere := failures[hook]; failHere {
		delete(failures, hook) // execute these only once (loop)
		return fmt.Errorf("%s: %s", prefix, simErrMsg)
	}

	return nil
}
