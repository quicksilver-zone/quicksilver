package interchaintest

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/strangelove-ventures/interchaintest/v5"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	_ "go.uber.org/zap"
	_ "path"
	"path/filepath"
	"testing"
	"time"

	transfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	istypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/strangelove-ventures/interchaintest/v5/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v5/ibc"
	"github.com/strangelove-ventures/interchaintest/v5/testreporter"
	"github.com/strangelove-ventures/interchaintest/v5/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TxProposalv1 contains chain proposal transaction detail for gov module v1 (sdk v0.46.0+)
type TxProposalv1 struct {
	Messages []json.RawMessage `json:"messages"`
	Metadata string            `json:"metadata"`
	Deposit  string            `json:"deposit"`
	Title    string            `json:"title"`
	Summary  string            `json:"summary"`

	// SDK v50 only
	Proposer  string `json:"proposer,omitempty"`
	Expedited bool   `json:"expedited,omitempty"`
}

// TestHandleTokenizedShares spins up a Quicksilver and Juno network, initializes an IBC connection between them,
// and sends an ICS20 token transfer from Quicksilver->Juno and then back from Juno->Quicksilver.
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

	content := istypes.RegisterZoneProposal{
		Title:            "register lstest-1 zone",
		Description:      "register lstest-1 zone with multisend and lsm enabled",
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

	check, err := cdctypes.NewAnyWithValue(&content)

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

	heightBeforeVote, err := quicksilver.Height(ctx)
	require.NoError(t, err, "error fetching height before vote")

	//Voting on the proposal
	err = quicksilver.VoteOnProposalAllValidators(ctx, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "Failed to submit votes")

	//Checking the proposal with matching ID and status.
	proposalStatusResponse, err := cosmos.PollForProposalStatus(ctx, quicksilver, heightBeforeVote, heightBeforeVote+20, proposalID, cosmos.ProposalStatusPassed)
	fmt.Println("Proposal status response", proposalStatusResponse)
	require.NoError(t, err, "Proposal status did not change to passed in expected number of blocks")

	stdout, _, err := quicksilver.Validators[0].ExecQuery(ctx, "interchainstaking", "zones")

	require.NotEmpty(t, stdout)
	require.NoError(t, err)
}

func SubmitProposal(ctx context.Context, c *cosmos.CosmosChain, keyName string, prop TxProposalv1) (string, error) {
	tn := c.Validators[0]
	if len(c.FullNodes) > 0 {
		tn = c.FullNodes[0]
	}

	propJson, err := json.MarshalIndent(prop, "", " ")
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(propJson)
	proposalFilename := fmt.Sprintf("%x.json", hash)

	err = tn.WriteFile(ctx, propJson, proposalFilename)
	if err != nil {
		return "", fmt.Errorf("writing param change proposal: %w", err)
	}

	proposalPath := filepath.Join(tn.HomeDir(), proposalFilename)

	command := []string{
		"gov", "submit-proposal",
		proposalPath,
		"--gas", "auto",
	}
	txHash, err := tn.ExecTx(ctx, keyName, command...)
	if err != nil {
		return txHash, fmt.Errorf("failed to submit gov v1 proposal: %w", err)
	}

	return TxProposal(tn, txHash)
}

func TxProposal(tn *cosmos.ChainNode, txHash string) (string, error) {
	var txResp *types.TxResponse
	err := retry.Do(func() error {
		var err error
		txResp, err = authTx.QueryTx(tn.CliContext(), txHash)
		fmt.Println("Tx proposal response: ", txResp)
		return err
	},
		// retry for total of 3 seconds
		retry.Attempts(15),
		retry.Delay(200*time.Millisecond),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction %s: %w", txHash, err)
	}
	events := txResp.Events
	evtSubmitProp := "submit_proposal"
	proposalID, _ := AttributeValue(events, evtSubmitProp, "proposal_id")

	return proposalID, nil
}

func AttributeValue(events []abcitypes.Event, eventType, attrKey string) (string, bool) {
	for _, event := range events {
		if event.Type != eventType {
			continue
		}
		for _, attr := range event.Attributes {
			if string(attr.Key) == attrKey {
				return string(attr.Value), true
			}
		}
	}
	return "", false
}
