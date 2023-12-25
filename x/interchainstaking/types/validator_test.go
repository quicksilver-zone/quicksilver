package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

var (
	acc1 = addressutils.GenerateAddressForTestWithPrefix("cosmos")
	vals = addressutils.GenerateValidatorsDeterministic(4)
)

func TestSharesToTokens(t *testing.T) {
	val := types.Validator{
		ValoperAddress:  vals[0],
		DelegatorShares: sdkmath.LegacyNewDecWithPrec(12345, 3),
		VotingPower:     sdkmath.NewInt(10),
	}
	require.Equal(t, sdkmath.NewInt(81), val.SharesToTokens(sdkmath.LegacyNewDec(100)))

	nilSharesVal := types.Validator{}
	require.Equal(t, sdkmath.ZeroInt(), nilSharesVal.SharesToTokens(sdkmath.LegacyNewDec(100)))

	nolSharesVal := types.Validator{
		DelegatorShares: sdkmath.LegacyZeroDec(),
	}
	require.Equal(t, sdkmath.ZeroInt(), nolSharesVal.SharesToTokens(sdkmath.LegacyNewDec(100)))
}

func TestNormalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 0)
}

func TestOrdinalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di = di.Ordinalize(sdkmath.LegacyNewDec(10000))
	require.Equal(t, len(di.Intents), 0)
}

func TestNormalizeIntentWithOneIntent(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdkmath.LegacyNewDec(1000)})
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 1)
	require.Equal(t, di.Intents[0].Weight, sdkmath.LegacyOneDec())
}

func TestNormalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdkmath.LegacyNewDec(1000)},
		{ValoperAddress: vals[1], Weight: sdkmath.LegacyNewDec(1000)},
		{ValoperAddress: vals[2], Weight: sdkmath.LegacyNewDec(1000)},
	}}

	di = di.Normalize()
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.Intents[0].Weight, sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3)))
	require.Equal(t, di.Intents[1].Weight, sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3)))
	require.Equal(t, di.Intents[2].Weight, sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3)))
}

func TestNormalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: addressutils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdkmath.LegacyNewDec(5).Quo(sdkmath.LegacyNewDec(50))},
		{ValoperAddress: vals[1], Weight: sdkmath.LegacyNewDec(10).Quo(sdkmath.LegacyNewDec(50))},
		{ValoperAddress: vals[2], Weight: sdkmath.LegacyNewDec(35).Quo(sdkmath.LegacyNewDec(50))},
	}}

	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.MustIntentForValoper(vals[0]).Weight, sdkmath.LegacyNewDecWithPrec(1, 1))
	require.Equal(t, di.MustIntentForValoper(vals[1]).Weight, sdkmath.LegacyNewDecWithPrec(2, 1))
	require.Equal(t, di.MustIntentForValoper(vals[2]).Weight, sdkmath.LegacyNewDecWithPrec(7, 1))
}

func TestOrdinalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: addressutils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
		{ValoperAddress: vals[1], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
		{ValoperAddress: vals[2], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
	}}
	di = di.Ordinalize(sdkmath.LegacyNewDec(3000))
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, sdkmath.NewInt(1000), di.MustIntentForValoper(vals[0]).Weight.RoundInt())
}

func TestOrdinalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: addressutils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdkmath.LegacyNewDec(5).Quo(sdkmath.LegacyNewDec(50))},
		{ValoperAddress: vals[1], Weight: sdkmath.LegacyNewDec(10).Quo(sdkmath.LegacyNewDec(50))},
		{ValoperAddress: vals[2], Weight: sdkmath.LegacyNewDec(35).Quo(sdkmath.LegacyNewDec(50))},
	}}
	di = di.Ordinalize(sdkmath.LegacyNewDec(3000))
	require.Equal(t, di.MustIntentForValoper(vals[0]).Weight.RoundInt(), sdkmath.NewInt(300))
	require.Equal(t, di.MustIntentForValoper(vals[1]).Weight.RoundInt(), sdkmath.NewInt(600))
	require.Equal(t, di.MustIntentForValoper(vals[2]).Weight.RoundInt(), sdkmath.NewInt(2100))
}

func TestAddOrdinal(t *testing.T) {
	di := types.DelegatorIntent{
		Delegator: addressutils.GenerateAccAddressForTest().String(),
		Intents: []*types.ValidatorIntent{
			{ValoperAddress: vals[0], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
			{ValoperAddress: vals[1], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
			{ValoperAddress: vals[2], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
		},
	}

	newIntents := types.ValidatorIntents{
		{ValoperAddress: vals[0], Weight: sdkmath.LegacyNewDec(1000)},
		{ValoperAddress: vals[1], Weight: sdkmath.LegacyNewDec(2000)},
	}

	// do nothing for no validator intents
	modified := di.AddOrdinal(sdkmath.LegacyNewDec(600), types.ValidatorIntents{})
	require.Equal(t, di, modified)

	di = di.AddOrdinal(sdkmath.LegacyNewDec(6000), newIntents)

	require.Equal(t, 3, len(di.Intents))

	require.Equal(t, di.MustIntentForValoper(vals[0]).Weight, sdkmath.LegacyNewDec(3).QuoTruncate(sdkmath.LegacyNewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[1]).Weight, sdkmath.LegacyNewDec(4).QuoTruncate(sdkmath.LegacyNewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[2]).Weight, sdkmath.LegacyNewDec(2).QuoTruncate(sdkmath.LegacyNewDec(9)))
}

func TestAddOrdinalWithNewVal(t *testing.T) {
	di := types.DelegatorIntent{Delegator: addressutils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: vals[0], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
		{ValoperAddress: vals[1], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
		{ValoperAddress: vals[2], Weight: sdkmath.LegacyOneDec().Quo(sdkmath.LegacyNewDec(3))},
	}}

	newIntents := types.ValidatorIntents{
		{ValoperAddress: vals[3], Weight: sdkmath.LegacyNewDec(1000)},
		{ValoperAddress: vals[2], Weight: sdkmath.LegacyNewDec(2000)},
	}

	di = di.AddOrdinal(sdkmath.LegacyNewDec(6000), newIntents)

	require.Equal(t, 4, len(di.Intents))

	require.Equal(t, di.MustIntentForValoper(vals[0]).Weight, sdkmath.LegacyNewDec(2).QuoTruncate(sdkmath.LegacyNewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[1]).Weight, sdkmath.LegacyNewDec(2).QuoTruncate(sdkmath.LegacyNewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[2]).Weight, sdkmath.LegacyNewDec(4).QuoTruncate(sdkmath.LegacyNewDec(9)))
	require.Equal(t, di.MustIntentForValoper(vals[3]).Weight, sdkmath.LegacyNewDec(1).QuoTruncate(sdkmath.LegacyNewDec(9)))
}
