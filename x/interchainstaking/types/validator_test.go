package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestNormalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: []*types.ValidatorIntent{}}
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 0)
}

func TestOrdinalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: []*types.ValidatorIntent{}}
	di = di.Ordinalize(sdk.NewInt(10000))
	require.Equal(t, len(di.Intents), 0)
}

func TestNormalizeIntentWithOneIntent(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(1000)})
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 1)
	require.Equal(t, di.Intents[0].Weight, sdk.OneDec())
}

func TestNormalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(1000)})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.NewDec(1000)})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(1000)})
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 3)
}

func TestNormalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(12108)})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.NewDec(3)})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(4002881)})
	require.NotPanics(t, func() { di.Normalize() })
	require.Equal(t, len(di.Intents), 3)
}

func TestOrdinalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})
	di = di.Ordinalize(sdk.NewInt(3000))
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.Intents[0].Weight.RoundInt(), sdk.NewInt(1000))
}

func TestAddOrdinal(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})

	newIntents := map[string]*types.ValidatorIntent{}
	newIntents["cosmosvaloper12345678"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.NewDec(1000)}
	newIntents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(2000)}

	di = di.AddOrdinal(sdk.NewInt(6000), newIntents)

	require.Equal(t, 3, len(di.Intents))

	// feels risky fetch these by numeric index; can we guarantee ordering? DelegatorIntents should probably be a string indexed map
	require.Equal(t, di.Intents[0].Weight, sdk.NewDec(3).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents[1].Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents[2].Weight, sdk.NewDec(4).QuoTruncate(sdk.NewDec(9)))
}

func TestAddOrdinalWithNewVal(t *testing.T) {
	di := types.DelegatorIntent{Delegator: "cosmos12345667890", Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper12345678", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper23456789", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.OneDec().QuoTruncate(sdk.NewDec(3))})

	newIntents := map[string]*types.ValidatorIntent{}
	// add a validator we haven't seen before here; ensure it is included in output.
	newIntents["cosmosvaloper98765432"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper98765432", Weight: sdk.NewDec(1000)}
	newIntents["cosmosvaloper34567890"] = &types.ValidatorIntent{ValoperAddress: "cosmosvaloper34567890", Weight: sdk.NewDec(2000)}

	di = di.AddOrdinal(sdk.NewInt(6000), newIntents)

	require.Equal(t, 4, len(di.Intents))

	// feels risky fetch these by numeric index; can we guarantee ordering? DelegatorIntents should probably be a string indexed map
	require.Equal(t, di.Intents[0].Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents[1].Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents[2].Weight, sdk.NewDec(4).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.Intents[3].Weight, sdk.NewDec(1).QuoTruncate(sdk.NewDec(9)))
}
