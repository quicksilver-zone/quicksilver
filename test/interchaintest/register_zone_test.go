package interchaintest

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/icza/dyno"
	istypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	// simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
)

const (
	heightDelta      = 20
	votingPeriod     = "30s"
	maxDepositPeriod = "10s"
)

// Spin up a quicksilverd chain, push a contract, and get that contract code from chain
func TestRegisterZone(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

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
	})

	t.Logf("Calling cf.Chains")
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	quicksilverd := chains[0]

	t.Logf("NewInterchain")
	ic := interchaintest.NewInterchain().
		AddChain(quicksilverd)

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
	fundAmount := int64(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", int64(fundAmount), quicksilverd)
	quicksilverd1User := users[0]

	err = testutil.WaitForBlocks(ctx, 10, quicksilverd)
	require.NoError(t, err)

	quicksilverd1UserBalInitial, err := quicksilverd.GetBalance(ctx, quicksilverd1User.FormattedAddress(), quicksilverd.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount, quicksilverd1UserBalInitial)

	quicksilverdChain := quicksilverd.(*cosmos.CosmosChain)

	proposal := cosmos.TxProposalv1{
		Metadata: "none",
		Deposit:  "500000000" + quicksilverdChain.Config().Denom, // greater than min deposit
		Title:    "title",
		Summary:  "suma",
	}

	content := istypes.RegisterZoneProposal{
		Title:            "register lstest-1 zone",
		Description:      "register lstest-1 zone with multisend and lsm enabled",
		ConnectionId:     "connection-0",
		BaseDenom:        "uatom",
		LocalDenom:       "uqatom",
		AccountPrefix:    "cosmos",
		DepositsEnabled:  true,
		UnbondingEnabled: true,
		LiquidityModule:  false,
		ReturnToSender:   true,
		Decimals:         6,
	}

	message := govv1.MsgExecLegacyContent{
		Content:   content,
		Authority: "quick10d07y265gmmuvt4z0w9aw880jnsr700j3xrh0p",
	}
	msg, err := quicksilverd.Config().EncodingConfig.Codec.MarshalInterfaceJSON(&message)
	fmt.Println("Msg: ", string(msg))
	require.NoError(t, err)
	proposal.Messages = append(proposal.Messages, msg)
	proposalTx, err := quicksilverdChain.SubmitProposal(ctx, quicksilverd1User.KeyName(), proposal)
	require.NoError(t, err, "error submitting proposal tx")

	height, err := quicksilverd.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	err = quicksilverdChain.VoteOnProposalAllValidators(ctx, proposalTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, quicksilverdChain, height, height+heightDelta, proposalTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = testutil.WaitForBlocks(ctx, 2, quicksilverd)
	require.NoError(t, err)

	address, _, err := quicksilverdChain.Validators[0].ExecQuery(ctx, "interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r")
	require.Equal(t, nil, address)
	require.NoError(t, err)
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
