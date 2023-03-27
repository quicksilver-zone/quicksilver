package interchaintest

import (
	"context"
	"fmt"
	"testing"
	"time"

	transfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	interchaintest "github.com/strangelove-ventures/ibctest/v5"
	"github.com/strangelove-ventures/ibctest/v5/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v5/ibc"
	"github.com/strangelove-ventures/ibctest/v5/test"
	"github.com/strangelove-ventures/ibctest/v5/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestJunoGaiaIBCTransfer spins up a Quicksilver and Gaia network, initializes an IBC connection between them,
// and sends an ICS20 token transfer from Quicksilver->Gaia and then back from Gaia->Quicksilver.
func TestQuicksilverGaiaIBCTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver and Gaia
	numVals := 1
	numFullNodes := 1

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "quicksilver",
			ChainConfig:   config,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "gaia",
			Version:       "v9.0.1",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	quicksilver, gaia := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Create relayer factory to utilize the go-relayer
	client, network := interchaintest.DockerSetup(t)

	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(t, client, network)

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain().
		AddChain(quicksilver).
		AddChain(gaia).
		AddRelayer(r, "rly").
		AddLink(interchaintest.InterchainLink{
			Chain1:  quicksilver,
			Chain2:  gaia,
			Relayer: r,
			Path:    pathQuicksilverGaia,
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
	require.NoError(t, r.StartRelayer(ctx, eRep, pathQuicksilverGaia))
	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				panic(fmt.Errorf("an error occurred while stopping the relayer: %s", err))
			}
		},
	)

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, quicksilver, gaia)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = test.WaitForBlocks(ctx, 5, quicksilver, gaia)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	quickUser, gaiaUser := users[0], users[1]

	quickUserAddr := quickUser.Bech32Address(quicksilver.Config().Bech32Prefix)
	gaiaUserAddr := gaiaUser.Bech32Address(gaia.Config().Bech32Prefix)

	// Get original account balances
	quicksilverOrigBal, err := quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, quicksilverOrigBal)

	gaiaOrigBal, err := gaia.GetBalance(ctx, gaiaUserAddr, gaia.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, gaiaOrigBal)

	// Compose an IBC transfer and send from Quicksilver -> Gaia
	const transferAmount = int64(1_000)
	transfer := ibc.WalletAmount{
		Address: gaiaUserAddr,
		Denom:   quicksilver.Config().Denom,
		Amount:  transferAmount,
	}

	quickChannels, err := r.GetChannels(ctx, eRep, quicksilver.Config().ChainID)
	require.NoError(t, err)

	transferTx, err := quicksilver.SendIBCTransfer(ctx, quickChannels[0].ChannelID, quickUserAddr, transfer, nil)
	require.NoError(t, err)

	quicksilverHeight, err := quicksilver.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = test.PollForAck(ctx, quicksilver, quicksilverHeight, quicksilverHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Get the IBC denom for uqck on Gaia
	quicksilverTokenDenom := transfertypes.GetPrefixedDenom(quickChannels[0].Counterparty.PortID, quickChannels[0].Counterparty.ChannelID, quicksilver.Config().Denom)
	quicksilverIBCDenom := transfertypes.ParseDenomTrace(quicksilverTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on Juno and are in the user acc on Gaia
	quicksilverUpdateBal, err := quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverUpdateBal-transferAmount, quicksilverUpdateBal)

	gaiaUpdateBal, err := gaia.GetBalance(ctx, gaiaUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, gaiaUpdateBal)

	// Compose an IBC transfer and send from Gaia -> Juno
	transfer = ibc.WalletAmount{
		Address: quickUserAddr,
		Denom:   quicksilverIBCDenom,
		Amount:  transferAmount,
	}

	transferTx, err = gaia.SendIBCTransfer(ctx, quickChannels[0].Counterparty.ChannelID, gaiaUserAddr, transfer, nil)
	require.NoError(t, err)

	gaiaHeight, err := gaia.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = test.PollForAck(ctx, gaia, gaiaHeight, gaiaHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Assert that the funds are now back on Juno and not on Gaia
	quicksilverUpdateBal, err = quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal, quicksilverUpdateBal)

	gaiaUpdateBal, err = gaia.GetBalance(ctx, gaiaUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(0), gaiaUpdateBal)
}

// This test is meant to be used as a basic ibctest tutorial.
// Code snippets are broken down in ./docs/upAndRunning.md
func TestLearn(t *testing.T) {
	ctx := context.Background()

	// Chain Factory
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{Name: "gaia", Version: "v9.0.0", ChainConfig: ibc.ChainConfig{
			GasPrices: "0.0uatom",
		}},
		{Name: "osmosis", Version: "v11.0.0"},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	gaia, osmosis := chains[0], chains[1]

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)
	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(
		t, client, network)

	// Prep Interchain
	const ibcPath = "gaia-osmo-demo"
	ic := interchaintest.NewInterchain().
		AddChain(gaia).
		AddChain(osmosis).
		AddRelayer(r, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  gaia,
			Chain2:  osmosis,
			Relayer: r,
			Path:    ibcPath,
		})

	// Log location
	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)
	// Reporter/logs
	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	// Build interchain
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

		SkipPathCreation: false},
	),
	)

	// Create and Fund User Wallets
	fundAmount := int64(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", int64(fundAmount), gaia, osmosis)
	gaiaUser := users[0]
	osmosisUser := users[1]

	gaiaUserBalInitial, err := gaia.GetBalance(ctx, gaiaUser.Bech32Address(gaia.Config().Bech32Prefix), gaia.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount, gaiaUserBalInitial)

	// Get Channel ID
	gaiaChannelInfo, err := r.GetChannels(ctx, eRep, gaia.Config().ChainID)
	require.NoError(t, err)
	gaiaChannelID := gaiaChannelInfo[0].ChannelID

	osmoChannelInfo, err := r.GetChannels(ctx, eRep, osmosis.Config().ChainID)
	require.NoError(t, err)
	osmoChannelID := osmoChannelInfo[0].ChannelID

	// Send Transaction
	amountToSend := int64(1_000_000)
	dstAddress := osmosisUser.Bech32Address(osmosis.Config().Bech32Prefix)
	tx, err := gaia.SendIBCTransfer(ctx, gaiaChannelID, gaiaUser.KeyName, ibc.WalletAmount{
		Address: dstAddress,
		Denom:   gaia.Config().Denom,
		Amount:  amountToSend,
	},
		nil,
	)
	require.NoError(t, err)
	require.NoError(t, tx.Validate())

	// relay packets and acknoledgments
	require.NoError(t, r.FlushPackets(ctx, eRep, ibcPath, osmoChannelID))
	require.NoError(t, r.FlushAcknowledgements(ctx, eRep, ibcPath, gaiaChannelID))

	// test source wallet has decreased funds
	expectedBal := gaiaUserBalInitial - amountToSend
	gaiaUserBalNew, err := gaia.GetBalance(ctx, gaiaUser.Bech32Address(gaia.Config().Bech32Prefix), gaia.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, expectedBal, gaiaUserBalNew)

	// Trace IBC Denom
	srcDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom("transfer", gaiaChannelID, gaia.Config().Denom))
	dstIbcDenom := srcDenomTrace.IBCDenom()

	// Test destination wallet has increased funds
	osmosUserBalNew, err := osmosis.GetBalance(ctx, osmosisUser.Bech32Address(osmosis.Config().Bech32Prefix), dstIbcDenom)
	require.NoError(t, err)
	require.Equal(t, amountToSend, osmosUserBalNew)

}
