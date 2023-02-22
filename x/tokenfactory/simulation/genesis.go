package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	tokenfactoryGenesis := types.DefaultGenesis()
	// use bond denom for simulation
	tokenfactoryGenesis.Params.DenomCreationFee = sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10_000_000))
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
