package interchaintest

import (
	"context"
	"fmt"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v5"
	"github.com/strangelove-ventures/interchaintest/v5/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v5/ibc"
	"github.com/strangelove-ventures/interchaintest/v5/testreporter"
	"github.com/strangelove-ventures/interchaintest/v5/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestInterchainStaking TODO
func TestInterchainStaking(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver and Juno
	numVals := 3
	numFullNodes := 3

	config, err := createConfig()
	require.NoError(t, err)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "quicksilver",
			ChainConfig:   config,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "juno",
			Version:       "v14.1.0",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
			//ChainConfig: ibc.ChainConfig{
			//	GasPrices: "0.0uatom",
			//},
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	quicksilver, juno := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Create relayer factory to utilize the go-relayer
	client, network := interchaintest.DockerSetup(t)

	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(t, client, network)

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain().
		AddChain(quicksilver).
		AddChain(juno).
		AddRelayer(r, "rly").
		AddLink(interchaintest.InterchainLink{
			Chain1:  quicksilver,
			Chain2:  juno,
			Relayer: r,
			Path:    pathQuicksilverJuno,
		})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()

	err = ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,

		// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Start the relayer
	require.NoError(t, r.StartRelayer(ctx, eRep, pathQuicksilverJuno))
	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				panic(fmt.Errorf("an error occurred while stopping the relayer: %s", err))
			}
		},
	)

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, quicksilver, juno)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, quicksilver, juno)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	quickUser, junoUser := users[0], users[1]

	quickUserAddr := quickUser.FormattedAddress()
	junoUserAddr := junoUser.FormattedAddress()
	_ = quickUserAddr
	_ = junoUserAddr

	runSidecars(t, ctx, quicksilver, juno)
}

func runSidecars(t *testing.T, ctx context.Context, quicksilver, juno *cosmos.CosmosChain) {
	t.Helper()

	runICQ(t, ctx, quicksilver, juno)
	// runXCC(t, ctx, quicksilver, juno)
}

func runICQ(t *testing.T, ctx context.Context, quicksilver, juno *cosmos.CosmosChain) {
	//	t.Helper()
	//
	//	var icq *cosmos.SidecarProcess
	//	for _, sidecar := range quicksilver.Sidecars {
	//		if sidecar.ProcessName == "icq" {
	//			icq = sidecar
	//		}
	//	}
	//	require.NotNil(t, icq)
	//
	//	containerCfg := "config.yaml"
	//
	//	file := fmt.Sprintf(`default_chain: '%s'
	//chains:
	//  '%s':
	//    key: default
	//    chain-id: '%s'
	//    rpc-addr: '%s'
	//    grpc-addr: '%s'
	//    account-prefix: quick
	//    keyring-backend: test
	//    gas-adjustment: 1.2
	//    gas-prices: 0.01uqck
	//    min-gas-amount: 0
	//    key-directory: %s/.icq/keys
	//    debug: false
	//    timeout: 20s
	//    block-timeout: 10s
	//    output-format: json
	//    sign-mode: direct
	//  '%s':
	//    key: default
	//    chain-id: '%s'
	//    rpc-addr: '%s'
	//    grpc-addr: '%s'
	//    account-prefix: osmo
	//    keyring-backend: test
	//    gas-adjustment: 1.2
	//    gas-prices: 0.01uosmo
	//    min-gas-amount: 0
	//    key-directory: %s/.icq/keys
	//    debug: false
	//    timeout: 20s
	//    block-timeout: 10s
	//    output-format: json
	//    sign-mode: direct
	//`,
	//		quicksilver.Config().ChainID,
	//		quicksilver.Config().ChainID,
	//		quicksilver.Config().ChainID,
	//		quicksilver.GetRPCAddress(),
	//		quicksilver.GetGRPCAddress(),
	//		icq.HomeDir(),
	//		juno.Config().ChainID,
	//		juno.Config().ChainID,
	//		juno.GetRPCAddress(),
	//		juno.GetGRPCAddress(),
	//		icq.HomeDir(),
	//	)
	//
	//	err := icq.WriteFile(ctx, []byte(file), containerCfg)
	//	require.NoError(t, err)
	//	_, err = icq.ReadFile(ctx, containerCfg)
	//	require.NoError(t, err)
	//
	//	err = icq.StartContainer(ctx)
	//	require.NoError(t, err)
	//
	//	err = icq.Running(ctx)
	//	require.NoError(t, err)
}

func runXCC(t *testing.T, ctx context.Context, quicksilver, juno *cosmos.CosmosChain) {
	//	t.Helper()
	//
	//	var xcc *cosmos.SidecarProcess
	//	for _, sidecar := range quicksilver.Sidecars {
	//		if sidecar.ProcessName == "xcc" {
	//			xcc = sidecar
	//		}
	//	}
	//	require.NotNil(t, xcc)
	//
	//	containerCfg := "config.yaml"
	//
	//	file := fmt.Sprintf(`source_chain: '%s'
	//chains:
	//  quick-1: '%s'
	//  juno-1: '%s'
	//`,
	//		quicksilver.Config().ChainID,
	//		quicksilver.GetRPCAddress(),
	//		juno.GetRPCAddress(),
	//	)
	//
	//	err := xcc.WriteFile(ctx, []byte(file), containerCfg)
	//	require.NoError(t, err)
	//	_, err = xcc.ReadFile(ctx, containerCfg)
	//	require.NoError(t, err)
	//
	//	err = xcc.StartContainer(ctx)
	//	require.NoError(t, err)
	//
	//	err = xcc.Running(ctx)
	//	require.NoError(t, err)
}
