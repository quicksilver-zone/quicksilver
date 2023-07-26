package types

import (
	"fmt"

	"github.com/ingenuity-build/quicksilver/third-party-chains/osmosis-types/osmoutils"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EpochHooks interface {
	// the first block whose timestamp is after the duration is counted as the end of the epoch
	AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error
	// new epoch is next block of epoch end block
	BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error
}

var _ EpochHooks = MultiEpochHooks{}

// MultiEpochHooks combine multiple hooks, all hook functions are run in array sequence.
type MultiEpochHooks []EpochHooks

func NewMultiEpochHooks(hooks ...EpochHooks) MultiEpochHooks {
	return hooks
}

// AfterEpochEnd is called when epoch is going to be ended, epochNumber is the number of epoch that is ending.
func (h MultiEpochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	for i := range h {
		panicCatchingEpochHook(ctx, h[i].AfterEpochEnd, epochIdentifier, epochNumber)
	}
	return nil
}

// BeforeEpochStart is called when epoch is going to be started, epochNumber is the number of epoch that is starting.
func (h MultiEpochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	for i := range h {
		panicCatchingEpochHook(ctx, h[i].BeforeEpochStart, epochIdentifier, epochNumber)
	}
	return nil
}

func panicCatchingEpochHook(
	ctx sdk.Context,
	hookFn func(ctx sdk.Context, epochIdentifier string, epochNumber int64) error,
	epochIdentifier string,
	epochNumber int64,
) {
	wrappedHookFn := func(ctx sdk.Context) error {
		return hookFn(ctx, epochIdentifier, epochNumber)
	}

	err := osmoutils.ApplyFuncIfNoError(ctx, wrappedHookFn)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("error in epoch hook %v", err))
	}
}
