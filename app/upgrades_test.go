package app

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	airdroptypes "github.com/ingenuity-build/quicksilver/x/airdrop/types"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestReplaceZone(t *testing.T) {
	// set up zone drop record and claims.
	app := Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	denom := app.StakingKeeper.BondDenom(ctx)
	someCoins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000000)))
	// work around airdrop keeper can't mint :)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, someCoins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, app.AirdropKeeper.GetZoneDropAccountAddress("osmotest-4"), someCoins)

	zd := airdroptypes.ZoneDrop{
		ChainId:    "osmotest-4",
		StartTime:  time.Now().AddDate(0, 0, -1),
		Duration:   time.Hour,
		Decay:      time.Hour,
		Allocation: someCoins.AmountOf(denom).Uint64(),
		Actions: []sdk.Dec{
			sdk.OneDec(),
		},
		IsConcluded: false,
	}

	app.AirdropKeeper.SetZoneDrop(ctx, zd)

	claim1 := airdroptypes.ClaimRecord{
		ChainId:          "osmotest-4",
		Address:          "quick1g035r8sl346ttxuj0555yxdwftr52t849t3q39",
		ActionsCompleted: make(map[int32]*airdroptypes.CompletedAction),
		MaxAllocation:    500000,
		BaseValue:        500000,
	}

	claim2 := airdroptypes.ClaimRecord{
		ChainId:          "osmotest-4",
		Address:          "quick1u53f8u6jjdpxquesk8tqxzv9hvqx7qyfzlkdrj",
		ActionsCompleted: make(map[int32]*airdroptypes.CompletedAction),
		MaxAllocation:    500000,
		BaseValue:        500000,
	}

	err := app.AirdropKeeper.SetClaimRecord(ctx, claim1)
	require.NoError(t, err)
	err = app.AirdropKeeper.SetClaimRecord(ctx, claim2)
	require.NoError(t, err)
	claims := app.AirdropKeeper.AllZoneClaimRecords(ctx, "osmotest-4")
	require.Equal(t, 2, len(claims))
	require.True(t, app.AirdropKeeper.GetZoneDropAccountBalance(ctx, "osmotest-4").Amount.Equal(sdk.NewInt(1000000)))
	require.NotPanics(t, func() { ReplaceZoneDropChain(ctx, app, "osmotest-4", "osmo-test-4", ctx.BlockHeader().Time) })
	claimsAfter := app.AirdropKeeper.AllZoneClaimRecords(ctx, "osmotest-4")
	require.Equal(t, 0, len(claimsAfter))
	claimsNew := app.AirdropKeeper.AllZoneClaimRecords(ctx, "osmo-test-4")
	require.Equal(t, 2, len(claimsNew))
	zoneDropsAfter := app.AirdropKeeper.AllZoneDrops(ctx)
	// check we don't suddenly have two airdrops.
	require.Equal(t, 1, len(zoneDropsAfter))
	// check the one aidrop we have has the expected values.
	require.Equal(t, zoneDropsAfter[0].ChainId, "osmo-test-4")
	require.Equal(t, zoneDropsAfter[0].StartTime, ctx.BlockHeader().Time)
	require.False(t, app.AirdropKeeper.GetZoneDropAccountBalance(ctx, "osmotest-4").Amount.Equal(sdk.NewInt(1000000)))
	require.True(t, app.AirdropKeeper.GetZoneDropAccountBalance(ctx, "osmo-test-4").Amount.Equal(sdk.NewInt(1000000)))
}
