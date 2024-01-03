package interchaintest

import (
	"context"
	"fmt"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	transfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	istypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/strangelove-ventures/interchaintest/v5"
	"github.com/strangelove-ventures/interchaintest/v5/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v5/ibc"
	"github.com/strangelove-ventures/interchaintest/v5/testreporter"
	"github.com/strangelove-ventures/interchaintest/v5/testutil"
	"github.com/stretchr/testify/require"
	_ "go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	_ "path"
	"testing"
	"time"
)

// TestHandleTokenizedShares
func TestHandleTokenizedShares(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver and Juno
	numVals := 3
	numFullNodes := 3

	config, err := createConfig()
	require.NoError(t, err)

	modifyGenesis := []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.voting_params.voting_period",
			Value: "20s",
		},
		{
			Key:   "app_state.staking.params.unbonding_time",
			Value: "60s",
		},
		{
			Key:   "app_state.interchainstaking.params.unbonding_enabled",
			Value: true,
		},
	}
	config.ModifyGenesis = cosmos.ModifyGenesis(modifyGenesis)

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
			ChainConfig: ibc.ChainConfig{
				ModifyGenesis: cosmos.ModifyGenesis([]cosmos.GenesisKV{
					{
						Key:   "app_state.staking.params.unbonding_time",
						Value: "60s",
					},
				}),
			},
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

	// Get original account balances
	quicksilverOrigBal, err := quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, quicksilverOrigBal)

	junoOrigBal, err := juno.GetBalance(ctx, junoUserAddr, juno.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, junoOrigBal)

	// Compose an IBC transfer and send from Quicksilver -> Juno
	const transferAmount = int64(1_000)
	transfer := ibc.WalletAmount{
		Address: junoUserAddr,
		Denom:   quicksilver.Config().Denom,
		Amount:  transferAmount,
	}

	quickChannels, err := r.GetChannels(ctx, eRep, quicksilver.Config().ChainID)
	require.NoError(t, err)

	transferTx, err := quicksilver.SendIBCTransfer(ctx, quickChannels[0].ChannelID, quickUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	quicksilverHeight, err := quicksilver.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, quicksilver, quicksilverHeight, quicksilverHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Get the IBC denom for uqck on Juno
	quicksilverTokenDenom := transfertypes.GetPrefixedDenom(quickChannels[0].Counterparty.PortID, quickChannels[0].Counterparty.ChannelID, quicksilver.Config().Denom)
	quicksilverIBCDenom := transfertypes.ParseDenomTrace(quicksilverTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on Juno and are in the user acc on Juno
	quicksilverUpdateBal, err := quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal-transferAmount, quicksilverUpdateBal)

	junoUpdateBal, err := juno.GetBalance(ctx, junoUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, junoUpdateBal)

	// Create new clients
	err = r.CreateClients(ctx, eRep, pathQuicksilverJuno, ibc.CreateClientOptions{TrustingPeriod: "330h"})
	require.NoError(t, err)

	// Create a new connection
	err = r.CreateConnections(ctx, eRep, pathQuicksilverJuno)
	require.NoError(t, err)

	connections, err := r.GetConnections(ctx, eRep, quicksilver.Config().ChainID)
	require.NoError(t, err)

	// Compose an IBC transfer and send from Quicksilver -> Juno
	transfer = ibc.WalletAmount{
		Address: quickUserAddr,
		Denom:   quicksilverIBCDenom,
		Amount:  transferAmount,
	}

	transferTx, err = juno.SendIBCTransfer(ctx, quickChannels[0].Counterparty.ChannelID, junoUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	junoHeight, err := juno.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, juno, junoHeight, junoHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Assert that the funds are now back on Juno and not on Juno
	quicksilverUpdateBal, err = quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal, quicksilverUpdateBal)

	junoUpdateBal, err = juno.GetBalance(ctx, junoUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(0), junoUpdateBal)

	registerProposal := istypes.RegisterZoneProposal{
		Title:            "Register zone",
		Description:      "Register zone",
		ConnectionId:     "connection-0",
		BaseDenom:        quicksilver.Config().Denom,
		LocalDenom:       quicksilver.Config().Denom,
		AccountPrefix:    "quick",
		DepositsEnabled:  true,
		UnbondingEnabled: true,
		LiquidityModule:  false,
		ReturnToSender:   true,
		Decimals:         6,
	}

	check, err := cdctypes.NewAnyWithValue(&registerProposal)

	message := govv1.MsgExecLegacyContent{
		Content:   check,
		Authority: "quick10d07y265gmmuvt4z0w9aw880jnsr700j3xrh0p",
	}
	msg, err := quicksilver.Config().EncodingConfig.Codec.MarshalInterfaceJSON(&message)
	fmt.Println("Msg: ", string(msg))
	require.NoError(t, err)

	proposal := TxProposalv1{
		Metadata: "none",
		Deposit:  "500000000" + quicksilver.Config().Denom,
		Title:    "title",
		Summary:  "register lstest-1 zone with multisend and lsm enabled",
	}

	//Appending proposal data in messages
	proposal.Messages = append(proposal.Messages, msg)

	require.NoError(t, err)

	//Submitting a proposal on Quicksilver
	proposalID, err := SubmitProposal(ctx, quicksilver, quickUserAddr, proposal)

	require.NoError(t, err)

	//Voting on the proposal
	err = quicksilver.VoteOnProposalAllValidators(ctx, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "Failed to submit votes")

	heightAfterVote, err := quicksilver.Height(ctx)
	require.NoError(t, err, "error fetching height before vote")

	//Checking the proposal with matching ID and status.
	_, err = cosmos.PollForProposalStatus(ctx, quicksilver, heightAfterVote, heightAfterVote+20, proposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "Proposal status did not change to passed in expected number of blocks")
	time.Sleep(10 * time.Second)
	zone, err := QueryZones(ctx, quicksilver)
	require.NoError(t, err)

	//Deposit Address Check
	depositAddress := zone[0].DepositAddress.Address
	icaAddr, err := QueryZoneICAAddress(ctx, quicksilver, depositAddress, connections[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, icaAddr)

	//Withdrawl Address Check
	withdralAddress := zone[0].WithdrawalAddress.Address
	icaAddr, err = QueryZoneICAAddress(ctx, quicksilver, withdralAddress, connections[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, icaAddr)

	//Delegation Address Check
	delegationAddress := zone[0].DelegationAddress.Address
	icaAddr, err = QueryZoneICAAddress(ctx, quicksilver, delegationAddress, connections[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, icaAddr)

	//Performance Address Check
	performanceAddress := zone[0].DelegationAddress.Address
	icaAddr, err = QueryZoneICAAddress(ctx, quicksilver, performanceAddress, connections[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, icaAddr)

	var updateZoneValue []*istypes.UpdateZoneValue
	updateZoneValue = append(updateZoneValue, &istypes.UpdateZoneValue{
		Key:   "",
		Value: "",
	})

	validators, err := QueryStakingValidators(ctx, quicksilver)
	fmt.Println(validators)
	require.NoError(t, err)

	delegateTx, err := RequestStakingDelegate(
		ctx,
		quicksilver,
		validators[0].OperatorAddress,
		quickUserAddr,
		"1000"+quicksilver.Config().Denom,
	)
	fmt.Println(delegateTx)
	require.NoError(t, err)

	delegation, err := QueryStakingDelegation(ctx, quicksilver, validators[0].OperatorAddress, quickUserAddr)
	fmt.Println(delegation)
	require.NoError(t, err)

	unbondTx, err := RequestStakingUnbond(
		ctx,
		quicksilver,
		validators[0].OperatorAddress,
		quickUserAddr,
		"1000"+quicksilver.Config().Denom,
	)
	fmt.Println(unbondTx)
	require.NoError(t, err)

	response, err := RequestICSRedeem(
		ctx,
		quicksilver,
		quickUserAddr,
		"1000"+quicksilver.Config().Denom,
	)
	fmt.Println(response)
	require.NoError(t, err)
}
