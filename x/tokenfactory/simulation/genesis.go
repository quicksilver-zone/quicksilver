package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

// RandomizedGenState generates a random GenesisState for tokenfactory.
func RandomizedGenState(simState *module.SimulationState) {
	feeAmt := simState.Rand.Int63n(5_000_000) + 5_000_000 // [5_000_000, 10_000_000)

	tokenfactoryGenesis := types.DefaultGenesis()
	// use bond denom for simulation
	tokenfactoryGenesis.Params.DenomCreationFee = sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, feeAmt))
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
