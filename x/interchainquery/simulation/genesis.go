package simulation

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

// RandomizedGenState generates a random GenesisState for interchainquery
func RandomizedGenState(simState *module.SimulationState) {
	icqGenesis := &types.GenesisState{
		Queries: []types.Query{},
	}

	bz, err := json.MarshalIndent(&icqGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	// TODO: Do some randomization later
	fmt.Printf("Selected deterministically generated interchainquery parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(icqGenesis)
}
