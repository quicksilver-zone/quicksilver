package simulation

import (
	"encoding/json"
	"fmt"
	
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// RandomizedGenState generates a random GenesisState for interchainstaking
func RandomizedGenState(simState *module.SimulationState) {
	icsGenesis := types.DefaultGenesis()

	bz, err := json.MarshalIndent(&icsGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	// TODO: Do some randomization later
	fmt.Printf("Selected deterministically generated interchainstaking parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(icsGenesis)
}
