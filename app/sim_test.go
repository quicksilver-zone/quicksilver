package app_test

import (
	"encoding/json"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	"math/rand"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/ingenuity-build/quicksilver/app"
	appsim "github.com/ingenuity-build/quicksilver/app/simulation"
)

func init() {
	simapp.GetSimulatorFlags()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}
func BenchmarkSimulation(b *testing.B) {
	simapp.FlagVerboseValue = true
	simapp.FlagOnOperationValue = true
	simapp.FlagAllInvariantsValue = true
	simapp.FlagInitialBlockHeightValue = 1

	config, db, dir, _, _, err := simapp.SetupSimulation("goleveldb-app-sim", "Simulation")
	require.NoError(b, err, "simulation setup failed")

	b.Cleanup(func() {
		err := db.Close()
		require.NoError(b, err)
		err = os.RemoveAll(dir)
		require.NoError(b, err)
	})

	quicksilver := app.NewQuicksilver(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		simapp.FlagPeriodValue,
		app.MakeEncodingConfig(),
		wasm.EnableAllProposals,
		simapp.EmptyAppOptions{},
		app.GetWasmOpts(simapp.EmptyAppOptions{}),
		false,
	)

	// Run randomized simulations
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		quicksilver.GetBaseApp(),
		appsim.AppStateFn(quicksilver.AppCodec(), quicksilver.SimulationManager()),
		simulationtypes.RandomAccounts,
		appsim.Operations(quicksilver, quicksilver.AppCodec(), config),
		quicksilver.ModuleAccountAddrs(),
		config,
		quicksilver.AppCodec(),
	)
	require.NoError(b, simErr)

	// export state and simParams before the simulation error is checked
	err = appsim.CheckExportSimulation(quicksilver, config, simParams)
	require.NoError(b, err)

	if config.Commit {
		simapp.PrintStats(db)
	}
}

// TestAppStateDeterminism TODO.
func TestAppStateDeterminism(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = true
	config.AllInvariants = true

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			db := dbm.NewMemDB()

			quicksilver := app.NewQuicksilver(
				log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
				db,
				nil,
				true,
				map[int64]bool{},
				app.DefaultNodeHome,
				simapp.FlagPeriodValue,
				app.MakeEncodingConfig(),
				wasm.EnableAllProposals,
				simapp.EmptyAppOptions{},
				app.GetWasmOpts(simapp.EmptyAppOptions{}),
				false,
				interBlockCacheOpt(),
			)

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				quicksilver.GetBaseApp(),
				appsim.AppStateFn(quicksilver.AppCodec(), quicksilver.SimulationManager()),
				simulationtypes.RandomAccounts,
				appsim.Operations(quicksilver, quicksilver.AppCodec(), config),
				quicksilver.ModuleAccountAddrs(),
				config,
				quicksilver.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simapp.PrintStats(db)
			}

			appHash := quicksilver.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, string(appHashList[0]), string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}
