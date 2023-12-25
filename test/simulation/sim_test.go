package simulation_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"cosmossdk.io/store"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	sdksimulation "github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/test/simulation"
)

func init() {
	simulation.SimulatorFlags()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func BenchmarkSimulation(b *testing.B) {
	simulation.FlagVerboseValue = true
	simulation.FlagOnOperationValue = true
	simulation.FlagAllInvariantsValue = true
	simulation.FlagInitialBlockHeightValue = 1

	config, db, dir, _, _, err := simulation.SetupSimulation("goleveldb-app-sim", "Simulation")
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
		simulation.FlagPeriodValue,
		app.MakeEncodingConfig(),
		wasm.EnableAllProposals,
		app.EmptyAppOptions{},
		app.GetWasmOpts(app.EmptyAppOptions{}),
		false,
	)

	// Run randomized simulations
	_, simParams, simErr := sdksimulation.SimulateFromSeed(
		b,
		os.Stdout,
		quicksilver.GetBaseApp(),
		simulation.AppStateFn(quicksilver.AppCodec(), quicksilver.SimulationManager()),
		simulationtypes.RandomAccounts,
		simulation.Operations(quicksilver, quicksilver.AppCodec(), config),
		quicksilver.ModuleAccountAddrs(),
		config,
		quicksilver.AppCodec(),
	)
	require.NoError(b, simErr)

	// export state and simParams before the simulation error is checked
	err = simulation.CheckExportSimulation(quicksilver, config, simParams)
	require.NoError(b, err)

	if config.Commit {
		simulation.PrintStats(db)
	}
}

// TestAppStateDeterminism TODO.
func TestAppStateDeterminism(t *testing.T) {
	if !simulation.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simulation.NewConfigFromFlags()
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
				simulation.FlagPeriodValue,
				app.MakeEncodingConfig(),
				wasm.EnableAllProposals,
				app.EmptyAppOptions{},
				app.GetWasmOpts(app.EmptyAppOptions{}),
				false,
				interBlockCacheOpt(),
			)

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := sdksimulation.SimulateFromSeed(
				t,
				os.Stdout,
				quicksilver.GetBaseApp(),
				simulation.AppStateFn(quicksilver.AppCodec(), quicksilver.SimulationManager()),
				simulationtypes.RandomAccounts,
				simulation.Operations(quicksilver, quicksilver.AppCodec(), config),
				quicksilver.ModuleAccountAddrs(),
				config,
				quicksilver.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simulation.PrintStats(db)
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
