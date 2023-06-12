package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultWithdrawalRequeueDelay = 6 * time.Hour

	// setting WithdrawStatusTokenize as 0 causes the value to be omitted when (un)marshalling :/.
	WithdrawStatusTokenize  int32 = iota + 1
	WithdrawStatusQueued    int32 = iota + 1
	WithdrawStatusUnbond    int32 = iota + 1
	WithdrawStatusSend      int32 = iota + 1
	WithdrawStatusCompleted int32 = iota + 1
)

// DelayCompletion updates a withdrawal record completion date to:
//
//	updatedCompletion = currentTime + delay
func (w *WithdrawalRecord) DelayCompletion(ctx sdk.Context, delay time.Duration) {
	w.CompletionTime = ctx.BlockTime().Add(delay)
}
