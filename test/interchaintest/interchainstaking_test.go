package interchaintest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/gogoproto/proto"

	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestInterchainStaking TODO
func TestInterchainStaking(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver and Gaia
	numVals := 3
	numFullNodes := 3

	client, network := interchaintest.DockerSetup(t)

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()

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
			Name:          "gaia",
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

	quicksilver, gaia := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Create relayer factory to utilize the go-relayer
	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t), relayer.CustomDockerImage("ghcr.io/notional-labs/cosmos-relayer", "nguyen-v2.3.1", "1000:1000")).Build(t, client, network)

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain().
		AddChain(quicksilver).
		AddChain(gaia).
		AddRelayer(r, "rly").
		AddLink(interchaintest.InterchainLink{
			Chain1:  quicksilver,
			Chain2:  gaia,
			Relayer: r,
			Path:    "quicksilver-gaia",
		})

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
	require.NoError(t, r.StartRelayer(ctx, eRep, "quicksilver-gaia"))
	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				panic(fmt.Errorf("an error occurred while stopping the relayer: %s", err))
			}
		},
	)

	// Get connections
	connections, err := r.GetConnections(ctx, eRep, quicksilver.Config().ChainID)
	require.NoError(t, err)

	// Get all Validators
	stdout, _, err := gaia.Validators[0].ExecQuery(ctx, "staking", "validators")
	require.NoError(t, err)
	require.NotEmpty(t, stdout)

	var validatorsResp QueryValidatorsResponse
	err = codec.NewLegacyAmino().UnmarshalJSON(stdout, &validatorsResp)
	require.NoError(t, err)

	gaiaValidators := validatorsResp.Validators

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, quicksilver, gaia)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, quicksilver, gaia)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	quickUser, gaiaUser := users[0], users[1]

	quickUserAddr := quickUser.FormattedAddress()
	gaiaUserAddr := gaiaUser.FormattedAddress()
	_ = quickUserAddr
	_ = gaiaUserAddr

	runSidecars(t, ctx, quicksilver, gaia)

	proposal := cosmos.TxProposalv1{
		Metadata: "none",
		Deposit:  "500000000" + quicksilver.Config().Denom, // greater than min deposit
		Title:    "title",
		Summary:  "summary",
	}

	content := icstypes.RegisterZoneProposal{
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

	check, err := cdctypes.NewAnyWithValue(&content)
	require.NoError(t, err)

	message := govv1.MsgExecLegacyContent{
		Content:   check,
		Authority: "quick10d07y265gmmuvt4z0w9aw880jnsr700j3xrh0p",
	}
	msg, err := quicksilver.Config().EncodingConfig.Codec.MarshalInterfaceJSON(&message)
	require.NoError(t, err)
	proposal.Messages = append(proposal.Messages, msg)

	// Submit Proposal
	proposalTx, err := quicksilver.SubmitProposal(ctx, quickUser.KeyName(), proposal)
	require.NoError(t, err, "error submitting proposal tx")

	height, err := quicksilver.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	err = quicksilver.VoteOnProposalAllValidators(ctx, proposalTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, quicksilver, height, height+heightDelta, proposalTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = testutil.WaitForBlocks(ctx, 20, quicksilver)
	require.NoError(t, err)

	stdout, _, err = quicksilver.Validators[0].ExecQuery(ctx, "interchainstaking", "zones")
	require.NoError(t, err)
	require.NotEmpty(t, stdout)

	var zones icstypes.QueryZonesResponse
	err = codec.NewLegacyAmino().UnmarshalJSON(stdout, &zones)
	require.NoError(t, err)

	zone := zones.Zones

	//Deposit Address Check
	depositAddress := zone[0].DepositAddress
	queryICA := []string{
		quicksilver.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", depositAddress.Address, connections[0].ID,
		"--chain-id", quicksilver.Config().ChainID,
		"--home", quicksilver.HomeDir(),
		"--node", quicksilver.GetRPCAddress(),
	}
	stdout, _, err = quicksilver.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts := strings.SplitN(string(stdout), ":", 2)
	icaAddr := strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	//Withdrawl Address Check
	withdralAddress := zone[0].WithdrawalAddress
	queryICA = []string{
		quicksilver.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", withdralAddress.Address, connections[0].ID,
		"--chain-id", quicksilver.Config().ChainID,
		"--home", quicksilver.HomeDir(),
		"--node", quicksilver.GetRPCAddress(),
	}
	stdout, _, err = quicksilver.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts = strings.SplitN(string(stdout), ":", 2)
	icaAddr = strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	//Delegation Address Check
	delegationAddress := zone[0].DelegationAddress
	queryICA = []string{
		quicksilver.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", delegationAddress.Address, connections[0].ID,
		"--chain-id", quicksilver.Config().ChainID,
		"--home", quicksilver.HomeDir(),
		"--node", quicksilver.GetRPCAddress(),
	}
	stdout, _, err = quicksilver.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts = strings.SplitN(string(stdout), ":", 2)
	icaAddr = strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	//Performance Address Check
	performanceAddress := zone[0].DelegationAddress
	queryICA = []string{
		quicksilver.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", performanceAddress.Address, connections[0].ID,
		"--chain-id", quicksilver.Config().ChainID,
		"--home", quicksilver.HomeDir(),
		"--node", quicksilver.GetRPCAddress(),
	}
	stdout, _, err = quicksilver.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	parts = strings.SplitN(string(stdout), ":", 2)
	icaAddr = strings.TrimSpace(parts[1])
	require.NotEmpty(t, icaAddr)

	version := icatypes.NewDefaultMetadataString("connection-0", "connection-0")
	_, err = quicksilver.FullNodes[0].ExecTx(
		ctx, quickUser.KeyName(), "interchain-accounts", "controller", "register", "connection-0",
		"--version", version,
	)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 5, quicksilver, gaia)
	require.NoError(t, err)

	stdout, _, err = quicksilver.Validators[0].ExecQuery(ctx, "interchain-accounts", "controller", "interchain-account", quickUserAddr, "connection-0")
	require.NoError(t, err)
	require.NotEmpty(t, stdout)

	var icaQuickUser icacontrollertypes.QueryInterchainAccountResponse
	err = codec.NewLegacyAmino().UnmarshalJSON(stdout, &icaQuickUser)
	require.NoError(t, err)

	icaQuickUserAddr := icaQuickUser.Address

	// Bank Send for delegation
	msgSend := &banktypes.MsgSend{
		FromAddress: gaiaUserAddr,
		ToAddress:   icaQuickUserAddr,
		Amount:      sdk.NewCoins(sdk.Coin{Denom: "uatom", Amount: sdkmath.NewInt(10_000_000)}),
	}

	cdc := config.EncodingConfig.Codec
	bz, err := icatypes.SerializeCosmosTx(cdc, []proto.Message{msgSend}) //gfddrg
	require.NoError(t, err)

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: bz,
		Memo: EncodeValidators(
			t,
			[]ValidatorDelegation{
				{
					Address: gaiaValidators[0].OperatorAddress,
					Percent: 50,
				},
				{
					Address: gaiaValidators[1].OperatorAddress,
					Percent: 25,
				},
				{
					Address: gaiaValidators[2].OperatorAddress,
					Percent: 25,
				},
			},
		),
	}
	jsonPacketData, err := codec.NewLegacyAmino().MarshalJSON(packetData)
	require.NoError(t, err)

	_, err = quicksilver.FullNodes[0].ExecTx(
		ctx, quickUser.KeyName(), "interchain-accounts", "controller", "send-tx", "connection-0", string(jsonPacketData),
	)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 50, quicksilver, gaia)
	require.NoError(t, err)

	stdout, _, err = quicksilver.Validators[0].ExecQuery(ctx, "interchainstaking", "intent", "gaia-2", quickUserAddr)
	require.NoError(t, err)
	require.NotEmpty(t, stdout)
	t.Logf("Intent: %s", string(stdout))

	stdout, _, err = quicksilver.Validators[0].ExecQuery(ctx, "bank", "balances", quickUserAddr)
	require.NoError(t, err)
	require.NotEmpty(t, stdout)
	t.Logf("User quick bank balances: %s", string(stdout))

	stdout, _, err = gaia.Validators[0].ExecQuery(ctx, "bank", "balances", gaiaUserAddr)
	require.NoError(t, err)
	require.NotEmpty(t, stdout)
	t.Logf("User gaia bank balances: %s", string(stdout))
}

func runSidecars(t *testing.T, ctx context.Context, quicksilver, gaia *cosmos.CosmosChain) {
	t.Helper()

	runICQ(t, ctx, quicksilver, gaia)
	// runXCC(t, ctx, quicksilver, gaia)
}

func runICQ(t *testing.T, ctx context.Context, quicksilver, gaia *cosmos.CosmosChain) {
	t.Helper()

	var icq *cosmos.SidecarProcess
	for _, sidecar := range quicksilver.Sidecars {
		if sidecar.ProcessName == "icq" {
			icq = sidecar
		}
	}
	require.NotNil(t, icq)

	containerCfg := "config.yaml"

	file := fmt.Sprintf(`default_chain: '%s'
chains:
  '%s':
    key: default
    chain-id: '%s'
    rpc-addr: '%s'
    grpc-addr: '%s'
    account-prefix: quick
    keyring-backend: test
    gas-adjustment: 1.2
    gas-prices: 0.01uqck
    min-gas-amount: 0
    key-directory: %s/.icq/keys
    debug: false
    timeout: 20s
    block-timeout: 10s
    output-format: json
    sign-mode: direct
  '%s':
    key: default
    chain-id: '%s'
    rpc-addr: '%s'
    grpc-addr: '%s'
    account-prefix: osmo
    keyring-backend: test
    gas-adjustment: 1.2
    gas-prices: 0.01uosmo
    min-gas-amount: 0
    key-directory: %s/.icq/keys
    debug: false
    timeout: 20s
    block-timeout: 10s
    output-format: json
    sign-mode: direct
`,
		quicksilver.Config().ChainID,
		quicksilver.Config().ChainID,
		quicksilver.Config().ChainID,
		quicksilver.GetRPCAddress(),
		quicksilver.GetGRPCAddress(),
		icq.HomeDir(),
		gaia.Config().ChainID,
		gaia.Config().ChainID,
		gaia.GetRPCAddress(),
		gaia.GetGRPCAddress(),
		icq.HomeDir(),
	)

	err := icq.WriteFile(ctx, []byte(file), containerCfg)
	require.NoError(t, err)
	_, err = icq.ReadFile(ctx, containerCfg)
	require.NoError(t, err)

	err = icq.StartContainer(ctx)
	require.NoError(t, err)
}

func runXCC(t *testing.T, ctx context.Context, quicksilver, gaia *cosmos.CosmosChain) {
	t.Helper()

	var xcc *cosmos.SidecarProcess
	for _, sidecar := range quicksilver.Sidecars {
		if sidecar.ProcessName == "xcc" {
			xcc = sidecar
		}
	}
	require.NotNil(t, xcc)

	containerCfg := "config.yaml"

	file := fmt.Sprintf(`source_chain: '%s'
chains:
  quick-1: '%s'
  gaia-1: '%s'
`,
		quicksilver.Config().ChainID,
		quicksilver.GetRPCAddress(),
		gaia.GetRPCAddress(),
	)

	err := xcc.WriteFile(ctx, []byte(file), containerCfg)
	require.NoError(t, err)
	_, err = xcc.ReadFile(ctx, containerCfg)
	require.NoError(t, err)

	err = xcc.StartContainer(ctx)
	require.NoError(t, err)
}
