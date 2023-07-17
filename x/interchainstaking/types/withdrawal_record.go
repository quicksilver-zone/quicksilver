package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultWithdrawalRequeueDelay = 6 * time.Hour

	// setting WithdrawStatusTokenize as 0 causes the value to be omitted when (un)marshalling :/.
	WithdrawStatusTokenize  int32 = 1
	WithdrawStatusQueued    int32 = 2
	WithdrawStatusUnbond    int32 = 3
	WithdrawStatusSend      int32 = 4
	WithdrawStatusCompleted int32 = 5
)

// DelayCompletion updates a withdrawal record completion date to:
//
//	updatedCompletion = currentTime + delay
func (w *WithdrawalRecord) DelayCompletion(ctx sdk.Context, delay time.Duration) {
	w.CompletionTime = ctx.BlockTime().Add(delay)
}
