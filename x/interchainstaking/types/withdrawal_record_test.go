package types_test

import (
	"testing"
	"time"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestWithdrawalRecord_DelayCompletion(t *testing.T) {
	// test context
	ctx := sdk.Context{}.WithBlockTime(time.Now())

	wdr := types.WithdrawalRecord{
		ChainId:        "test",
		Delegator:      "test",
		Recipient:      "test",
		BurnAmount:     sdk.NewCoin("test", sdk.NewInt(10000)),
		Txhash:         "test",
		Status:         types.WithdrawStatusSend,
		CompletionTime: ctx.BlockTime(),
	}

	wdr.DelayCompletion(ctx, time.Hour)
	require.Equal(t, ctx.BlockTime().Add(time.Hour), wdr.CompletionTime)
}
