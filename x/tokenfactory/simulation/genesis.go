package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/quicksilver-zone/quicksilver/x/tokenfactory/types"
)

// RandomizedGenState generates a random GenesisState for tokenfactory.
func RandomizedGenState(simState *module.SimulationState) {
	// random fee
	feeAmt := simState.Rand.Int63n(5_000_000) + 7_500_000 // [7_500_000, 12_500_000)

	tokenfactoryGenesis := types.DefaultGenesis()
	// use bond denom for simulation
	tokenfactoryGenesis.Params.DenomCreationFee = sdk.NewCoins(sdkmath.NewInt64Coin(sdk.DefaultBondDenom, feeAmt))
	err := tokenfactoryGenesis.Validate()
	if err != nil {
		panic(err)
	}

	bz, err := json.MarshalIndent(&tokenfactoryGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected deterministically generated tokenfactory parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(tokenfactoryGenesis)
}
