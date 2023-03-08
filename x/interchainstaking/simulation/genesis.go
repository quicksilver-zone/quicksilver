package simulation

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// DONTCOVER

// RandomizedGenState generates a random GenesisState for interchainstaking
func RandomizedGenState(simState *module.SimulationState) {
	bz, err := os.ReadFile("./testdata/zones.json")
	if err != nil {
		panic(err)
	}

	var zones []types.Zone
	err = json.Unmarshal(bz, &zones)
	if err != nil {
		panic(err)
	}

	bz, err = os.ReadFile("./testdata/delegations.json")
	if err != nil {
		panic(err)
	}

	var delegations []types.DelegationsForZone
	err = json.Unmarshal(bz, &delegations)
	if err != nil {
		panic(err)
	}

	params := types.DefaultParams()
	params.UnbondingEnabled = true

	icsGenesis := &types.GenesisState{
		Params:                 params,
		Zones:                  zones,
		Receipts:               nil,
		Delegations:            delegations,
		PerformanceDelegations: nil,
		DelegatorIntents:       nil,
		PortConnections:        nil,
		WithdrawalRecords:      nil,
	}

	bz, err = json.MarshalIndent(&icsGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	// TODO: Do some randomization later
	fmt.Printf("Selected deterministically generated interchainstaking parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(icsGenesis)
}
