package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestNormalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 0)
}

func TestOrdinalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di = di.Ordinalize(sdk.NewDec(10000))
	require.Equal(t, len(di.Intents), 0)
}

func TestNormalizeIntentWithOneIntent(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di.Intents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(1000)}
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 1)
	require.Equal(t, di.Intents["cosmosvaloper12345678"].Weight, sdk.OneDec())
}

func TestNormalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di.Intents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(1000)}
	di.Intents["cosmosvaloper23456789"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.NewDec(1000)}
	di.Intents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(1000)}

	di = di.Normalize()
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.Intents["cosmosvaloper12345678"].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
	require.Equal(t, di.Intents["cosmosvaloper23456789"].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
	require.Equal(t, di.Intents["cosmosvaloper34567890"].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
}

func TestNormalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di.Intents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(5)}
	di.Intents["cosmosvaloper23456789"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.NewDec(10)}
	di.Intents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(35)}
	require.NotPanics(t, func() { di.Normalize() })
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.Intents["cosmosvaloper12345678"].Weight, sdk.NewDecWithPrec(1, 1))
	require.Equal(t, di.Intents["cosmosvaloper23456789"].Weight, sdk.NewDecWithPrec(2, 1))
	require.Equal(t, di.Intents["cosmosvaloper34567890"].Weight, sdk.NewDecWithPrec(7, 1))
}

func TestOrdinalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di.Intents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(3)}
	di.Intents["cosmosvaloper23456789"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.NewDec(3)}
	di.Intents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(3)}
	require.NotPanics(t, func() { di.Normalize() }) // normalise here because are already ordinal
	di = di.Ordinalize(sdk.NewDec(3000))
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, sdk.NewInt(1000), di.Intents["cosmosvaloper12345678"].Weight.RoundInt())
}

func TestOrdinalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di.Intents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(5)}
	di.Intents["cosmosvaloper23456789"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.NewDec(10)}
	di.Intents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(35)}
	require.NotPanics(t, func() { di.Normalize() }) // normalise here because are already ordinal
	di = di.Ordinalize(sdk.NewDec(3000))
	require.Equal(t, di.Intents["cosmosvaloper12345678"].Weight.RoundInt(), sdk.NewInt(300))
	require.Equal(t, di.Intents["cosmosvaloper23456789"].Weight.RoundInt(), sdk.NewInt(600))
	require.Equal(t, di.Intents["cosmosvaloper34567890"].Weight.RoundInt(), sdk.NewInt(2100))
}

func TestAddOrdinal(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di.Intents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(3)}
	di.Intents["cosmosvaloper23456789"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.NewDec(3)}
	di.Intents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(3)}
	require.NotPanics(t, func() { di.Normalize() }) // normalise here because are already ordinal

	newIntents := map[string]*types.ValidatorIntent{}
	newIntents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(1000)}
	newIntents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(2000)}

	di = di.AddOrdinal(sdk.NewDec(6000), newIntents)

	require.Equal(t, 3, len(di.Intents))

	// feels risky fetch these by numeric index; can we guarantee ordering? DelegatorIntents should probably be a string indexed map
	require.Equal(t, di.Intents["cosmosvaloper12345678"].Weight, sdk.NewDec(3).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents["cosmosvaloper23456789"].Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents["cosmosvaloper34567890"].Weight, sdk.NewDec(4).QuoTruncate(sdk.NewDec(9)))
}

func TestAddOrdinalWithNewVal(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: map[string]*types.ValidatorIntent{}}
	di.Intents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))}
	di.Intents["cosmosvaloper23456789"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))}
	di.Intents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))}

	newIntents := map[string]*types.ValidatorIntent{}
	// add a validator we haven't seen before here; ensure it is included in output.
	newIntents["cosmosvaloper98765432"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper98765432", Weight: sdk.NewDec(1000)}
	newIntents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(2000)}

	di = di.AddOrdinal(sdk.NewDec(6000), newIntents)

	require.Equal(t, 4, len(di.Intents))

	// feels risky fetch these by numeric index; can we guarantee ordering? DelegatorIntents should probably be a string indexed map
	require.Equal(t, di.Intents["cosmosvaloper12345678"].Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents["cosmosvaloper23456789"].Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents["cosmosvaloper34567890"].Weight, sdk.NewDec(4).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents["cosmosvaloper98765432"].Weight, sdk.NewDec(1).QuoTruncate(sdk.NewDec(9)))
}
