package interchaintest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	istypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	// simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
)

const (
	heightDelta = 20
)

// Spin up a quicksilverd chain, push a contract, and get that contract code from chain. Submit a proposal to register zones and query zones.
func TestRegisterZone(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver
	numVals := 3
	numFullNodes := 3
	client, network := interchaintest.DockerSetup(t)

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()

	// Get both chains
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			ChainConfig: ibc.ChainConfig{
				Type:           "cosmos",
				Name:           "quicksilver",
				ChainID:        "quicksilverd",
				Images:         []ibc.DockerImage{QuicksilverImage},
				Bin:            "quicksilverd",
				Bech32Prefix:   "quick",
				Denom:          "uqck",
				GasPrices:      "0.00uqck",
				GasAdjustment:  1.3,
				TrustingPeriod: "504h",
				EncodingConfig: quicksilverEncoding(),
				NoHostMount:    true,
				ModifyGenesis:  ModifyGenesisShortProposals(votingPeriod, maxDepositPeriod),
			},
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "juno",
			Version:       "v14.1.0",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	t.Logf("Calling cf.Chains")
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	quicksilverd, juno := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Get a relayer instance
	r := interchaintest.NewBuiltinRelayerFactory(ibc.Hermes, zaptest.NewLogger(t)).Build(t, client, network)

	// Build the network; spin up the chains and configure the relayer
	t.Logf("NewInterchain")
	ic := interchaintest.NewInterchain().
		AddChain(quicksilverd).
		AddChain(juno).
		AddRelayer(r, "rly").
		AddLink(interchaintest.InterchainLink{
			Chain1:  quicksilverd,
			Chain2:  juno,
			Relayer: r,
			Path:    pathQuicksilverJuno,
		})

	t.Logf("Interchain build options")
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  true, // Skip path creation, so we can have granular control over the process
	}))

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Create and Fund User Wallets
	fundAmount := math.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount.Int64(), quicksilverd, juno)
	quicksilverd1User := users[0]
	juno1User := users[1]

	err = testutil.WaitForBlocks(ctx, 10, quicksilverd, juno)
	require.NoError(t, err)

	quicksilverd1UserBalInitial, err := quicksilverd.GetBalance(ctx, quicksilverd1User.FormattedAddress(), quicksilverd.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount, quicksilverd1UserBalInitial)

	juno1UserBalInitial, err := juno.GetBalance(ctx, juno1User.FormattedAddress(), juno.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount, juno1UserBalInitial)

	// Generate a new IBC path
	err = r.GeneratePath(ctx, eRep, quicksilverd.Config().ChainID, juno.Config().ChainID, pathQuicksilverJuno)
	require.NoError(t, err)

	// Create new clients
	err = r.CreateClients(ctx, eRep, pathQuicksilverJuno, ibc.CreateClientOptions{TrustingPeriod: "330h"})
	require.NoError(t, err)

	// Create a new connection
	err = r.CreateConnections(ctx, eRep, pathQuicksilverJuno)
	require.NoError(t, err)

	connections, err := r.GetConnections(ctx, eRep, quicksilverd.Config().ChainID)
	require.NoError(t, err)
	// require.Equal(t, 1, len(connections))

	// Create a new channel
	err = r.CreateChannel(ctx, eRep, pathQuicksilverJuno, ibc.DefaultChannelOpts())
	require.NoError(t, err)

	// Query for the newly created channel
	_, err = r.GetChannels(ctx, eRep, quicksilverd.Config().ChainID)
	require.NoError(t, err)

	// Start the relayer and set the cleanup function.
	require.NoError(t, r.StartRelayer(ctx, eRep, pathQuicksilverJuno))
	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				panic(fmt.Errorf("an error occurred while stopping the relayer: %s", err))
			}
		},
	)

	proposal := cosmos.TxProposalv1{
		Metadata: "none",
		Deposit:  "500000000" + quicksilverd.Config().Denom, // greater than min deposit
		Title:    "title",
		Summary:  "suma",
	}

	content := istypes.RegisterZoneProposal{
		Title:            "register lstest-1 zone",
		Description:      "register lstest-1 zone with multisend and lsm enabled",
		ConnectionId:     "connection-0",
		BaseDenom:        "ujuno",
		LocalDenom:       "uqjuno",
		AccountPrefix:    "juno",
		DepositsEnabled:  true,
		UnbondingEnabled: true,
		LiquidityModule:  false,
		ReturnToSender:   true,
		Decimals:         6,
	}

	check, err := cdctypes.NewAnyWithValue(&content)
	require.NoError(t, err)

	message := govv1.MsgExecLegacyContent{
		Content:   check,
		Authority: "quick10d07y265gmmuvt4z0w9aw880jnsr700j3xrh0p",
	}
	msg, err := quicksilverd.Config().EncodingConfig.Codec.MarshalInterfaceJSON(&message)
	require.NoError(t, err)
	proposal.Messages = append(proposal.Messages, msg)

	// Submit Proposal
	proposalTx, err := quicksilverd.SubmitProposal(ctx, quicksilverd1User.KeyName(), proposal)
	require.NoError(t, err, "error submitting proposal tx")

	height, err := quicksilverd.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	err = quicksilverd.VoteOnProposalAllValidators(ctx, proposalTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, quicksilverd, height, height+heightDelta, proposalTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = testutil.WaitForBlocks(ctx, 20, quicksilverd)
	require.NoError(t, err)

	stdout, _, err := quicksilverd.Validators[0].ExecQuery(ctx, "interchainstaking", "zones")
	require.NoError(t, err)
	require.NotEmpty(t, stdout)

	var zones istypes.QueryZonesResponse
	err = codec.NewLegacyAmino().UnmarshalJSON(stdout, &zones)
	require.NoError(t, err)

	zone := zones.Zones

	// Deposit Address Check
	depositAddress := zone[0].DepositAddress
	queryICA := []string{
		quicksilverd.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", depositAddress.Address, connections[0].ID,
		"--chain-id", quicksilverd.Config().ChainID,
		"--home", quicksilverd.HomeDir(),
		"--node", quicksilverd.GetRPCAddress(),
	}
	stdout, _, err = quicksilverd.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts := strings.SplitN(string(stdout), ":", 2)
	icaAddr := strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	// Withdrawal Address Check
	withdralAddress := zone[0].WithdrawalAddress
	queryICA = []string{
		quicksilverd.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", withdralAddress.Address, connections[0].ID,
		"--chain-id", quicksilverd.Config().ChainID,
		"--home", quicksilverd.HomeDir(),
		"--node", quicksilverd.GetRPCAddress(),
	}
	stdout, _, err = quicksilverd.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts = strings.SplitN(string(stdout), ":", 2)
	icaAddr = strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	// Delegation Address Check
	delegationAddress := zone[0].DelegationAddress
	queryICA = []string{
		quicksilverd.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", delegationAddress.Address, connections[0].ID,
		"--chain-id", quicksilverd.Config().ChainID,
		"--home", quicksilverd.HomeDir(),
		"--node", quicksilverd.GetRPCAddress(),
	}
	stdout, _, err = quicksilverd.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts = strings.SplitN(string(stdout), ":", 2)
	icaAddr = strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	// Performance Address Check
	performanceAddress := zone[0].DelegationAddress
	queryICA = []string{
		quicksilverd.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", performanceAddress.Address, connections[0].ID,
		"--chain-id", quicksilverd.Config().ChainID,
		"--home", quicksilverd.HomeDir(),
		"--node", quicksilverd.GetRPCAddress(),
	}
	stdout, _, err = quicksilverd.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts = strings.SplitN(string(stdout), ":", 2)
	icaAddr = strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)
}
