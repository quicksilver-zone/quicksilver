package interchaintest

import (
	"context"
	"fmt"
	transfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	interchaintest "github.com/strangelove-ventures/ibctest/v5"
	"github.com/strangelove-ventures/ibctest/v5/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v5/ibc"
	"github.com/strangelove-ventures/ibctest/v5/test"
	"github.com/strangelove-ventures/ibctest/v5/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"testing"
)

// TestQuicksilverOsmosisIBCTransfer spins up a Quicksilver and Osmosis network, initializes an IBC connection between them,
// and sends an ICS20 token transfer from Quicksilver->Osmosis and then back from Osmosis->Quicksilver.
func TestQuicksilverOsmosisIBCTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver and osmosis
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
			Name:          "osmosis",
			Version:       "v12.0.0",
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

	quicksilver, osmosis := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Create relayer factory to utilize the go-relayer
	client, network := interchaintest.DockerSetup(t)

	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(t, client, network)

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain().
		AddChain(quicksilver).
		AddChain(osmosis).
		AddRelayer(r, "rly").
		AddLink(interchaintest.InterchainLink{
			Chain1:  quicksilver,
			Chain2:  osmosis,
			Relayer: r,
			Path:    pathQuicksilverOsmosis,
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
	require.NoError(t, r.StartRelayer(ctx, eRep, pathQuicksilverOsmosis))
	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				panic(fmt.Errorf("an error occurred while stopping the relayer: %s", err))
			}
		},
	)

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, quicksilver, osmosis)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = test.WaitForBlocks(ctx, 5, quicksilver, osmosis)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	quickUser, osmosisUser := users[0], users[1]

	quickUserAddr := quickUser.Bech32Address(quicksilver.Config().Bech32Prefix)
	osmosisUserAddr := osmosisUser.Bech32Address(osmosis.Config().Bech32Prefix)

	// Get original account balances
	quicksilverOrigBal, err := quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, quicksilverOrigBal)

	osmosisOrigBal, err := osmosis.GetBalance(ctx, osmosisUserAddr, osmosis.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, osmosisOrigBal)

	// Compose an IBC transfer and send from Quicksilver -> osmosis
	const transferAmount = int64(1_000)
	transfer := ibc.WalletAmount{
		Address: osmosisUserAddr,
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

	// Get the IBC denom for uqck on osmosis
	quicksilverTokenDenom := transfertypes.GetPrefixedDenom(quickChannels[0].Counterparty.PortID, quickChannels[0].Counterparty.ChannelID, quicksilver.Config().Denom)
	quicksilverIBCDenom := transfertypes.ParseDenomTrace(quicksilverTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on Juno and are in the user acc on osmosis
	quicksilverUpdateBal, err := quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal-transferAmount, quicksilverUpdateBal)

	osmosisUpdateBal, err := osmosis.GetBalance(ctx, osmosisUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, osmosisUpdateBal)

	// Compose an IBC transfer and send from osmosis -> Juno
	transfer = ibc.WalletAmount{
		Address: quickUserAddr,
		Denom:   quicksilverIBCDenom,
		Amount:  transferAmount,
	}

	transferTx, err = osmosis.SendIBCTransfer(ctx, quickChannels[0].Counterparty.ChannelID, osmosisUserAddr, transfer, nil)
	require.NoError(t, err)

	osmosisHeight, err := osmosis.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = test.PollForAck(ctx, osmosis, osmosisHeight, osmosisHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Assert that the funds are now back on Juno and not on osmosis
	quicksilverUpdateBal, err = quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal, quicksilverUpdateBal)

	osmosisUpdateBal, err = osmosis.GetBalance(ctx, osmosisUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(0), osmosisUpdateBal)
}
