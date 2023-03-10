package simulation

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ingenuity-build/quicksilver/app"
	interchainstakingtypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var denoms = []string{"uqatom", "uqosmo", "uqjunox"}

// AppStateFn returns the initial application state using a genesis or the simulation parameters.
// It panics if the user provides files for both of them.
// If a file is not given for the genesis or the sim params, it creates a randomized one.
func AppStateFn(cdc codec.JSONCodec, simManager *module.SimulationManager) simtypes.AppStateFn {
	return func(r *rand.Rand, accs []simtypes.Account, config simtypes.Config,
	) (appState json.RawMessage, simAccs []simtypes.Account, chainID string, genesisTimestamp time.Time) {
		if FlagGenesisTimeValue == 0 {
			genesisTimestamp = simtypes.RandTimestamp(r)
		} else {
			genesisTimestamp = time.Unix(FlagGenesisTimeValue, 0)
		}

		chainID = config.ChainID
		switch {
		case config.ParamsFile != "" && config.GenesisFile != "":
			panic("cannot provide both a genesis file and a params file")

		case config.GenesisFile != "":
			// override the default chain-id from simapp to set it later to the config
			genesisDoc, accounts := AppStateFromGenesisFileFn(r, cdc, config.GenesisFile)

			if FlagGenesisTimeValue == 0 {
				// use genesis timestamp if no custom timestamp is provided (i.e no random timestamp)
				genesisTimestamp = genesisDoc.GenesisTime
			}

			appState = genesisDoc.AppState
			chainID = genesisDoc.ChainID
			simAccs = accounts

		case config.ParamsFile != "":
			appParams := make(simtypes.AppParams)
			bz, err := os.ReadFile(config.ParamsFile)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(bz, &appParams)
			if err != nil {
				panic(err)
			}
			appState, simAccs = AppStateRandomizedFn(simManager, r, cdc, accs, genesisTimestamp, appParams)

		default:
			appParams := make(simtypes.AppParams)
			appState, simAccs = AppStateRandomizedFn(simManager, r, cdc, accs, genesisTimestamp, appParams)
		}

		rawState := make(map[string]json.RawMessage)
		err := json.Unmarshal(appState, &rawState)
		if err != nil {
			panic(err)
		}

		stakingStateBz, ok := rawState[stakingtypes.ModuleName]
		if !ok {
			panic("staking genesis state is missing")
		}

		stakingState := new(stakingtypes.GenesisState)
		err = cdc.UnmarshalJSON(stakingStateBz, stakingState)
		if err != nil {
			panic(err)
		}

		icsStateBz, ok := rawState[interchainstakingtypes.ModuleName]
		if !ok {
			panic("staking genesis state is missing")
		}

		icsState := new(interchainstakingtypes.GenesisState)
		err = cdc.UnmarshalJSON(icsStateBz, icsState)
		if err != nil {
			panic(err)
		}

		// compute not bonded balance
		notBondedTokens := sdk.ZeroInt()
		for _, val := range stakingState.Validators {
			if val.Status != stakingtypes.Unbonded {
				continue
			}
			notBondedTokens = notBondedTokens.Add(val.GetTokens())
		}
		notBondedCoins := sdk.NewCoin(stakingState.Params.BondDenom, notBondedTokens)
		// edit bank state to make it have the not bonded pool tokens
		bankStateBz, ok := rawState[banktypes.ModuleName]
		if !ok {
			panic("bank genesis state is missing")
		}
		bankState := new(banktypes.GenesisState)
		err = cdc.UnmarshalJSON(bankStateBz, bankState)
		if err != nil {
			panic(err)
		}

		govStateBz, ok := rawState[govtypes.ModuleName]
		if !ok {
			panic("gov genesis state is missing")
		}
		govState := new(govv1.GenesisState)
		err = cdc.UnmarshalJSON(govStateBz, govState)
		if err != nil {
			panic(err)
		}

		bankState.Params.DefaultSendEnabled = true
		bankState.Params.SendEnabled = banktypes.SendEnabledParams{
			banktypes.NewSendEnabled(stakingState.Params.BondDenom, true),
		}

		delMap := map[string][]*interchainstakingtypes.Delegation{}
		for _, denom := range denoms {
			delMap[denom] = []*interchainstakingtypes.Delegation{}
		}

		zoneMap := createZoneMap(icsState.Zones)

		stakingAddr := authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String()
		var found bool
		var newSupply sdk.Coins
		for i, balance := range bankState.Balances {
			if balance.Address == stakingAddr {
				found = true
			}

			// increase supply
			var (
				addDel  bool
				newCoin sdk.Coin
			)
			bankState.Balances[i], newCoin, addDel = modifyGenesisBalance(r, balance, sdk.DefaultBondDenom, sdk.DefaultBondDenom)
			newSupply.Add(bankState.Balances[i].Coins...)
			if addDel {
				zone := zoneMap[newCoin.Denom]
				val := zone.Validators[r.Intn(len(zone.Validators))]

				delMap[newCoin.Denom] = append(delMap[newCoin.Denom], &interchainstakingtypes.Delegation{
					DelegationAddress: zone.DelegationAddress.Address,
					ValidatorAddress:  val.ValoperAddress,
					Amount:            newCoin,
					Height:            0,
				})
			}
		}

		var delegations []interchainstakingtypes.DelegationsForZone
		for _, denom := range denoms {
			zone := zoneMap[denom]
			delegations = append(delegations, interchainstakingtypes.DelegationsForZone{
				ChainId:     zone.ChainId,
				Delegations: delMap[denom],
			})
		}
		icsState.Delegations = delegations

		// set up port connections
		for _, zone := range icsState.Zones {
			portID := zone.GetDelegationAddress().GetPortName()
			connectionID := "test" + portID
			icsState.PortConnections = append(icsState.PortConnections, interchainstakingtypes.PortConnectionTuple{
				ConnectionId: connectionID,
				PortId:       portID,
			})
		}

		// set new supply
		bankState.Supply = newSupply

		if !found {
			bankState.Balances = append(bankState.Balances, banktypes.Balance{
				Address: stakingAddr,
				Coins:   sdk.NewCoins(notBondedCoins),
			})
		}

		minDep := sdk.NewCoins(govState.DepositParams.MinDeposit...)
		govState.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(stakingState.Params.BondDenom, minDep.AmountOf(sdk.DefaultBondDenom)))

		// change appState back
		rawState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingState)
		rawState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)
		rawState[govtypes.ModuleName] = cdc.MustMarshalJSON(govState)
		rawState[interchainstakingtypes.ModuleName] = cdc.MustMarshalJSON(icsState)

		// replace appstate
		appState, err = json.Marshal(rawState)
		if err != nil {
			panic(err)
		}
		return appState, simAccs, chainID, genesisTimestamp
	}
}

func randQAssetBalaces(r *rand.Rand, balance banktypes.Balance) (banktypes.Balance, sdk.Coin, bool) {
	// do not add qassets for some accounts
	if r.Intn(100)%5 == 0 {
		return balance, sdk.Coin{}, false
	}

	denom := denoms[r.Intn(len(denoms))]
	amount := sdk.NewInt(1_000_000_000 + r.Int63n(1_000_000_000_000))

	newCoin := sdk.NewCoin(denom, amount)
	balance.Coins = balance.Coins.Add(sdk.NewCoins(newCoin)...)

	return balance, newCoin, true
}

func modifyGenesisBalance(r *rand.Rand, balance banktypes.Balance, oldBond, newBond string) (banktypes.Balance, sdk.Coin, bool) {
	amt := balance.Coins.AmountOf(oldBond)
	if amt.IsPositive() {
		balance.Coins = sdk.NewCoins(sdk.NewCoin(newBond, amt))
	}

	return randQAssetBalaces(r, balance)
}

// AppStateRandomizedFn creates calls each module's GenesisState generator function
// and creates the simulation params
func AppStateRandomizedFn(
	simManager *module.SimulationManager, r *rand.Rand, cdc codec.JSONCodec,
	accs []simtypes.Account, genesisTimestamp time.Time, appParams simtypes.AppParams,
) (json.RawMessage, []simtypes.Account) {
	numAccs := int64(len(accs))
	genesisState := app.NewDefaultGenesisState()

	// generate a random amount of initial stake coins and a random initial
	// number of bonded accounts
	var initialStake, numInitiallyBonded int64
	appParams.GetOrGenerate(
		cdc, simappparams.StakePerAccount, &initialStake, r,
		func(r *rand.Rand) { initialStake = r.Int63n(1e12) },
	)
	appParams.GetOrGenerate(
		cdc, simappparams.InitiallyBondedValidators, &numInitiallyBonded, r,
		func(r *rand.Rand) { numInitiallyBonded = int64(r.Intn(300)) },
	)

	if numInitiallyBonded > numAccs {
		numInitiallyBonded = numAccs
	}

	fmt.Printf(
		`Selected randomly generated parameters for simulated genesis:
{
  stake_per_account: "%d",
  initially_bonded_validators: "%d"
}
`, initialStake, numInitiallyBonded,
	)

	simState := &module.SimulationState{
		AppParams:    appParams,
		Cdc:          cdc,
		Rand:         r,
		GenState:     genesisState,
		Accounts:     accs,
		InitialStake: sdkmath.NewInt(initialStake),
		NumBonded:    numInitiallyBonded,
		GenTimestamp: genesisTimestamp,
	}

	simManager.GenerateGenesisStates(simState)

	appState, err := json.Marshal(genesisState)
	if err != nil {
		panic(err)
	}

	return appState, accs
}

// AppStateFromGenesisFileFn util function to generate the genesis AppState
// from a genesis.json file.
func AppStateFromGenesisFileFn(r io.Reader, cdc codec.JSONCodec, genesisFile string) (tmtypes.GenesisDoc, []simtypes.Account) {
	bytes, err := os.ReadFile(genesisFile)
	if err != nil {
		panic(err)
	}

	var genesis tmtypes.GenesisDoc
	// NOTE: Tendermint uses a custom JSON decoder for GenesisDoc
	err = tmjson.Unmarshal(bytes, &genesis)
	if err != nil {
		panic(err)
	}

	var appState app.GenesisState
	err = json.Unmarshal(genesis.AppState, &appState)
	if err != nil {
		panic(err)
	}

	var authGenesis authtypes.GenesisState
	if appState[authtypes.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[authtypes.ModuleName], &authGenesis)
	}

	newAccs := make([]simtypes.Account, len(authGenesis.Accounts))
	for i, acc := range authGenesis.Accounts {
		// Pick a random private key, since we don't know the actual key
		// This should be fine as it's only used for mock Tendermint validators
		// and these keys are never actually used to sign by mock Tendermint.
		privkeySeed := make([]byte, 15)
		if _, err := r.Read(privkeySeed); err != nil {
			panic(err)
		}

		privKey := secp256k1.GenPrivKeyFromSecret(privkeySeed)

		a, ok := acc.GetCachedValue().(authtypes.AccountI)
		if !ok {
			panic("expected account")
		}

		// create simulator accounts
		simAcc := simtypes.Account{PrivKey: privKey, PubKey: privKey.PubKey(), Address: a.GetAddress()}
		newAccs[i] = simAcc
	}

	return genesis, newAccs
}

func createZoneMap(zones []interchainstakingtypes.Zone) map[string]interchainstakingtypes.Zone {
	m := make(map[string]interchainstakingtypes.Zone)
	for _, zone := range zones {
		m[zone.LocalDenom] = zone
	}

	return m
}
