package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var (
	acc1 = addressutils.GenerateAddressForTestWithPrefix("cosmos")
	vals = addressutils.GenerateValidatorsDeterministic(4)
)

func TestSharesToTokens(t *testing.T) {
	val := types.Validator{
		ValoperAddress:  vals[0],
		DelegatorShares: sdk.NewDecWithPrec(12345, 3),
		VotingPower:     sdkmath.NewInt(10),
	}
	require.Equal(t, sdkmath.NewInt(81), val.SharesToTokens(sdk.NewDec(100)))

	nilSharesVal := types.Validator{}
	require.Equal(t, sdkmath.ZeroInt(), nilSharesVal.SharesToTokens(sdk.NewDec(100)))

	nolSharesVal := types.Validator{
		DelegatorShares: sdk.ZeroDec(),
	}
	require.Equal(t, sdkmath.ZeroInt(), nolSharesVal.SharesToTokens(sdk.NewDec(100)))
}

func TestNormalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 0)
}

func TestOrdinalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di = di.Ordinalize(sdk.NewDec(10000))
	require.Equal(t, len(di.Intents), 0)
}

func TestNormalizeIntentWithOneIntent(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDec(1000)})
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 1)
	require.Equal(t, di.Intents[0].Weight, sdk.OneDec())
}

func TestNormalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdk.NewDec(1000)},
		{ValoperAddress: vals[1], Weight: sdk.NewDec(1000)},
		{ValoperAddress: vals[2], Weight: sdk.NewDec(1000)},
	}}

	di = di.Normalize()
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.Intents[0].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
	require.Equal(t, di.Intents[1].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
	require.Equal(t, di.Intents[2].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
}

func TestNormalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: addressutils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdk.NewDec(5).Quo(sdk.NewDec(50))},
		{ValoperAddress: vals[1], Weight: sdk.NewDec(10).Quo(sdk.NewDec(50))},
		{ValoperAddress: vals[2], Weight: sdk.NewDec(35).Quo(sdk.NewDec(50))},
	}}

	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.MustIntentForValoper(vals[0]).Weight, sdk.NewDecWithPrec(1, 1))
	require.Equal(t, di.MustIntentForValoper(vals[1]).Weight, sdk.NewDecWithPrec(2, 1))
	require.Equal(t, di.MustIntentForValoper(vals[2]).Weight, sdk.NewDecWithPrec(7, 1))
}

func TestOrdinalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: addressutils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: vals[1], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: vals[2], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
	}}
	di = di.Ordinalize(sdk.NewDec(3000))
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, sdk.NewInt(1000), di.MustIntentForValoper(vals[0]).Weight.RoundInt())
}

func TestOrdinalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: addressutils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdk.NewDec(5).Quo(sdk.NewDec(50))},
		{ValoperAddress: vals[1], Weight: sdk.NewDec(10).Quo(sdk.NewDec(50))},
		{ValoperAddress: vals[2], Weight: sdk.NewDec(35).Quo(sdk.NewDec(50))},
	}}
	di = di.Ordinalize(sdk.NewDec(3000))
	require.Equal(t, di.MustIntentForValoper(vals[0]).Weight.RoundInt(), sdk.NewInt(300))
	require.Equal(t, di.MustIntentForValoper(vals[1]).Weight.RoundInt(), sdk.NewInt(600))
	require.Equal(t, di.MustIntentForValoper(vals[2]).Weight.RoundInt(), sdk.NewInt(2100))
}

func TestAddOrdinal(t *testing.T) {
	di := types.DelegatorIntent{
		Delegator: addressutils.GenerateAccAddressForTest().String(),
		Intents: []*types.ValidatorIntent{
			{ValoperAddress: vals[0], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
			{ValoperAddress: vals[1], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
			{ValoperAddress: vals[2], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		},
	}

	newIntents := types.ValidatorIntents{
		{ValoperAddress: vals[0], Weight: sdk.NewDec(1000)},
		{ValoperAddress: vals[1], Weight: sdk.NewDec(2000)},
	}

	// do nothing for no validator intents
	modified := di.AddOrdinal(sdk.NewDec(600), types.ValidatorIntents{})
	require.Equal(t, di, modified)

	di = di.AddOrdinal(sdk.NewDec(6000), newIntents)

	require.Equal(t, 3, len(di.Intents))

	require.Equal(t, di.MustIntentForValoper(vals[0]).Weight, sdk.NewDec(3).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[1]).Weight, sdk.NewDec(4).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[2]).Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
}

func TestAddOrdinalWithNewVal(t *testing.T) {
	di := types.DelegatorIntent{Delegator: addressutils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: vals[1], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: vals[2], Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
	}}

	newIntents := types.ValidatorIntents{
		{ValoperAddress: vals[3], Weight: sdk.NewDec(1000)},
		{ValoperAddress: vals[2], Weight: sdk.NewDec(2000)},
	}

	di = di.AddOrdinal(sdk.NewDec(6000), newIntents)

	require.Equal(t, 4, len(di.Intents))

	require.Equal(t, di.MustIntentForValoper(vals[0]).Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[1]).Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[2]).Weight, sdk.NewDec(4).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[3]).Weight, sdk.NewDec(1).QuoTruncate(sdk.NewDec(9)))
}
