package simulation

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// RandomizedGenState generates a random GenesisState for interchainstaking
func RandomizedGenState(simState *module.SimulationState) {
	var zones = []types.Zone{
		{
			ConnectionId:    "connection-77001",
			ChainId:         "cosmoshub-4",
			AccountPrefix:   "cosmos",
			LocalDenom:      "uqatom",
			BaseDenom:       "uatom",
			MultiSend:       false,
			LiquidityModule: false,
		},
		{
			ConnectionId:    "connection-77002",
			ChainId:         "osmosis-1",
			AccountPrefix:   "osmo",
			LocalDenom:      "uqosmo",
			BaseDenom:       "uosmo",
			MultiSend:       false,
			LiquidityModule: true,
		},
		{
			ConnectionId:    "connection-77003",
			ChainId:         "uni-5",
			AccountPrefix:   "juno",
			LocalDenom:      "uqjunox",
			BaseDenom:       "ujunox",
			MultiSend:       false,
			LiquidityModule: true,
		},
	}

	for _, z := range zones {
		// set zone validators
		z.Validators = append(z.Validators, &types.Validator{ValoperAddress: z.AccountPrefix + "valoper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded})
		z.Validators = append(z.Validators, &types.Validator{ValoperAddress: z.AccountPrefix + "valoper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded})
		z.Validators = append(z.Validators, &types.Validator{ValoperAddress: z.AccountPrefix + "valoper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded})
		z.Validators = append(z.Validators, &types.Validator{ValoperAddress: z.AccountPrefix + "valoper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded})
		z.Validators = append(z.Validators, &types.Validator{ValoperAddress: z.AccountPrefix + "valoper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded})
		z.Validators = append(z.Validators, &types.Validator{ValoperAddress: z.AccountPrefix + "valoper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded})
	}

	icsGenesis := &types.GenesisState{
		Params:                 types.DefaultParams(),
		Zones:                  zones,
		Receipts:               nil,
		Delegations:            nil,
		PerformanceDelegations: nil,
		DelegatorIntents:       nil,
		PortConnections:        nil,
		WithdrawalRecords:      nil,
	}

	bz, err := json.MarshalIndent(&icsGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	// TODO: Do some randomization later
	fmt.Printf("Selected deterministically generated interchainstaking parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(icsGenesis)
}
