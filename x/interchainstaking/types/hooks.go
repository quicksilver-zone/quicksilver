package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple ics hooks, all hook functions are run in array sequence.
var _ IcsHooks = &MultiIcsHooks{}

type MultiIcsHooks []IcsHooks

func NewMultiIcsHooks(hooks ...IcsHooks) MultiIcsHooks {
	return hooks
}

func (h MultiIcsHooks) AfterZoneCreated(ctx sdk.Context, connectionID, chainID, accountPrefix string) error {
	for i := range h {
		if err := h[i].AfterZoneCreated(ctx, connectionID, chainID, accountPrefix); err != nil {
			return err
		}
	}

	return nil
}
