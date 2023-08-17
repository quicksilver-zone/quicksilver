package interchaintest

import (
	"context"
	"encoding/json"
	"fmt"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/icza/dyno"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"strings"
	"time"

	math "cosmossdk.io/math"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	istypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	heightDelta      = 20
	votingPeriod     = "30s"
	maxDepositPeriod = "10s"
)

func TestQuicksilverE2E(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()
	// Create chain factory with Quicksilver
	//numVals := 3
	//numFullNodes := 3
	//
	//config, err := createConfig()
	//require.NoError(t, err)

	// Create relayer factory to utilize the go-relayer
	client, network := interchaintest.DockerSetup(t)

	ctx := context.Background()

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			ChainConfig: ibc.ChainConfig{
				Type:           "cosmos",
				Name:           "quicksilver",
				ChainID:        "quicksilverd",
				Images:         []ibc.DockerImage{QuicksilverImage},
				Bin:            "quicksilverd",
				Bech32Prefix:   "quick",
				Denom:          "stake",
				GasPrices:      "0.00stake",
				GasAdjustment:  1.3,
				TrustingPeriod: "504h",
				EncodingConfig: quicksilverEncoding(),
				NoHostMount:    true,
				ModifyGenesis:  modifyGenesisShortProposals(votingPeriod, maxDepositPeriod),
			},
		},
		{
			Name:    "juno",
			Version: "v14.1.0",
		},
	})
	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	quicksilverd, juno := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	r := interchaintest.NewBuiltinRelayerFactory(ibc.Hermes, zaptest.NewLogger(t)).Build(t, client, network)

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
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

	// Log location
	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)
	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	// Build interchain
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  true,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})
	// Generate a new IBC path
	err = r.GeneratePath(ctx, eRep, quicksilverd.Config().ChainID, juno.Config().ChainID, pathQuicksilverJuno)
	require.NoError(t, err)

	// Create new clients
	err = r.CreateClients(ctx, eRep, pathQuicksilverJuno, ibc.CreateClientOptions{TrustingPeriod: "330h"})
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, quicksilverd, juno)
	require.NoError(t, err)

	// Create a new connection
	err = r.CreateConnections(ctx, eRep, pathQuicksilverJuno)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, quicksilverd, juno)
	require.NoError(t, err)

	// Query for the newly created connection
	connections, err := r.GetConnections(ctx, eRep, quicksilverd.Config().ChainID)
	require.NoError(t, err)
	require.Equal(t, 1, len(connections))

	// Create a new channel
	err = r.CreateChannel(ctx, eRep, pathQuicksilverJuno, ibc.DefaultChannelOpts())
	require.NoError(t, err)

	// Query for the newly created channel
	quickChannels, err := r.GetChannels(ctx, eRep, quicksilverd.Config().ChainID)
	require.NoError(t, err)
	fmt.Println("Channel", quickChannels)

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
	users := interchaintest.GetAndFundTestUsers(t, ctx, "testing", genesisWalletAmount, quicksilverd, juno)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, quicksilverd, juno)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	quickUser, junoUser := users[0], users[1]

	quickUserAddr := quickUser.FormattedAddress()
	junoUserAddr := junoUser.FormattedAddress()

	// Get original account balances
	quicksilverOrigBal, err := quicksilverd.GetBalance(ctx, quickUserAddr, quicksilverd.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, quicksilverOrigBal)

	junoOrigBal, err := juno.GetBalance(ctx, junoUserAddr, juno.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, junoOrigBal)

	// Compose an IBC transfer and send from Quicksilver -> Juno
	transferAmount := math.NewInt(1_000)
	transfer := ibc.WalletAmount{
		Address: junoUserAddr,
		Denom:   quicksilverd.Config().Denom,
		Amount:  transferAmount,
	}

	transferTx, err := quicksilverd.SendIBCTransfer(ctx, quickChannels[0].ChannelID, quickUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	quicksilverHeight, err := quicksilverd.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, quicksilverd, quicksilverHeight, quicksilverHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Get the IBC denom for uqck on Juno
	quicksilverTokenDenom := transfertypes.GetPrefixedDenom(quickChannels[0].Counterparty.PortID, quickChannels[0].Counterparty.ChannelID, quicksilverd.Config().Denom)
	quicksilverIBCDenom := transfertypes.ParseDenomTrace(quicksilverTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on Juno and are in the user acc on Juno
	quicksilverUpdateBal, err := quicksilverd.GetBalance(ctx, quickUserAddr, quicksilverd.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal.Sub(transferAmount), quicksilverUpdateBal)

	junoUpdateBal, err := juno.GetBalance(ctx, junoUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, junoUpdateBal)

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
	quicksilverUpdateBal, err = quicksilverd.GetBalance(ctx, quickUserAddr, quicksilverd.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal, quicksilverUpdateBal)

	junoUpdateBal, err = juno.GetBalance(ctx, junoUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(0), junoUpdateBal)

	height1, err := quicksilverd.Height(ctx)
	require.NoError(t, err)

	//Creating a proposal on Quicksilver
	messages := istypes.RegisterZoneProposal{

		Title:            "register lstest-1 zone",
		Description:      "register lstest-1 zone ",
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
	check, err := cdctypes.NewAnyWithValue(&messages)

	message := govv1.MsgExecLegacyContent{
		Content:   check,
		Authority: "quick10d07y265gmmuvt4z0w9aw880jnsr700j3xrh0p",
	}
	msg, err := quicksilverd.Config().EncodingConfig.Codec.MarshalInterfaceJSON(&message)
	fmt.Println("Msg: ", string(msg))
	require.NoError(t, err)

	//Appending proposal data in messages
	//messages = append(messages)
	//txProposal, err := quicksilver.BuildProposal(messages, "RegisterZone Proposal For Juno", "Juno <-> Quicksilver", "", "1000_000")
	var propType cosmos.TxProposalv1
	propType.Metadata = ""
	propType.Title = "Adding Juno as a Zone"
	propType.Summary = ""

	propType.Messages = append(propType.Messages, msg)

	require.NoError(t, err)

	//Submitting a proposal on Quicksilver
	tx, err := quicksilverd.SubmitProposal(ctx, users[0].KeyName(), propType)

	//require.NoError(t, tx.Validate())

	require.NoError(t, err)

	//Voting on the proposal
	err = quicksilverd.VoteOnProposalAllValidators(ctx, tx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "Failed to submit votes")

	height2, err := quicksilverd.Height(ctx)
	require.NoError(t, err, "error fetching height before upgrade")

	//Checking the proposal with matching ID and status.
	_, err = cosmos.PollForProposalStatus(ctx, quicksilverd, height1, height2, tx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "Proposal status did not change to passed in expected number of blocks")

	stdout, _, err := quicksilverd.Validators[0].ExecQuery(ctx, "interchainstaking", "zones")

	require.NotEmpty(t, stdout)
	require.NoError(t, err)
	var zones []istypes.Zone
	err = json.Unmarshal([]byte(stdout), &zones)

	//Deposit Address Check
	depositAddress := zones[0].DepositAddress
	queryICA := []string{
		quicksilverd.Config().Bin, "query", "intertx", "interchainaccounts", connections[0].ID, depositAddress.Address,
		"--chain-id", quicksilverd.Config().ChainID,
		"--home", quicksilverd.HomeDir(),
		"--node", quicksilverd.GetRPCAddress(),
	}
	stdout, _, err = quicksilverd.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts := strings.SplitN(string(stdout), ":", 2)
	icaAddr := strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	//Withdrawl Address Check
	withdralAddress := zones[0].WithdrawalAddress
	queryICA = []string{
		quicksilverd.Config().Bin, "query", "intertx", "interchainaccounts", connections[0].ID, withdralAddress.Address,
		"--chain-id", quicksilverd.Config().ChainID,
		"--home", quicksilverd.HomeDir(),
		"--node", quicksilverd.GetRPCAddress(),
	}
	stdout, _, err = quicksilverd.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts = strings.SplitN(string(stdout), ":", 2)
	icaAddr = strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	//Delegation Address Check
	delegationAddress := zones[0].DelegationAddress
	queryICA = []string{
		quicksilverd.Config().Bin, "query", "intertx", "interchainaccounts", connections[0].ID, delegationAddress.Address,
		"--chain-id", quicksilverd.Config().ChainID,
		"--home", quicksilverd.HomeDir(),
		"--node", quicksilverd.GetRPCAddress(),
	}
	stdout, _, err = quicksilverd.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts = strings.SplitN(string(stdout), ":", 2)
	icaAddr = strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	//Performance Address Check
	performanceAddress := zones[0].DelegationAddress
	queryICA = []string{
		quicksilverd.Config().Bin, "query", "intertx", "interchainaccounts", connections[0].ID, performanceAddress.Address,
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

func modifyGenesisShortProposals(votingPeriod, maxDepositPeriod string) func(ibc.ChainConfig, []byte) ([]byte, error) {
	return func(chainConfig ibc.ChainConfig, genbz []byte) ([]byte, error) {
		g := make(map[string]interface{})
		if err := json.Unmarshal(genbz, &g); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
		}
		if err := dyno.Set(g, votingPeriod, "app_state", "gov", "params", "voting_period"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		if err := dyno.Set(g, maxDepositPeriod, "app_state", "gov", "params", "max_deposit_period"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		if err := dyno.Set(g, chainConfig.Denom, "app_state", "gov", "params", "min_deposit", 0, "denom"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		out, err := json.Marshal(g)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
		}
		return out, nil
	}
}
