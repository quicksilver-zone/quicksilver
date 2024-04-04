package interchaintest

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/testreporter"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/go-bip39"
	qsapp "github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/test/interchaintest/util"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func TestInterchainStaking(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver and Juno
	numVals := 2
	numFullNodes := 1

	config, err := createConfig()
	require.NoError(t, err)

	// config.SidecarConfigs = []ibc.SidecarConfig{
	// 	{
	// 		ProcessName:      "icq",
	// 		Image:            ibc.DockerImage{Repository: "quicksilverzone/interchain-queries", Version: "v1.0.0-beta.2"},
	// 		HomeDir:          "/icq-relayer",
	// 		StartCmd:         []string{"icq-relayer", "start", "--home", "/icq-relayer"},
	// 		PreStart:         true,
	// 		ValidatorProcess: false,
	// 	},
	// }
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "quicksilver",
			ChainConfig:   config,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "juno",
			Version:       "v14.0.0",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				GasPrices:           "0.0ujuno",
				ConfigFileOverrides: map[string]any{"config/config.toml": testutil.Toml{"consensus": testutil.Toml{"timeout_commit": "1s", "timeout_propose": "500ms", "timeout_prevote": "500ms", "timeout_precommit": "500ms"}}},
				ModifyGenesis: cosmos.ModifyGenesis([]cosmos.GenesisKV{
					{
						Key:   "app_state.interchainaccounts.host_genesis_state.params.allow_messages",
						Value: []string{"*"},
					},
					{
						Key:   "app_state.interchainaccounts.host_genesis_state.params.host_enabled",
						Value: true,
					},
				}),
			},
		},
	})
	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	quicksilver, juno := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Create relayer factory to utilize the go-relayer
	client, network := interchaintest.DockerSetup(t)

	r := interchaintest.NewBuiltinRelayerFactory(ibc.Hermes, zaptest.NewLogger(t) /*relayer.CustomDockerImage("informalsystems/hermes", "1.10.3", "1000:1000")*/).Build(t, client, network)

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
	entropy, err := bip39.NewEntropy(256)
	require.NoError(t, err)
	mnemonic, err := bip39.NewMnemonic(entropy)
	require.NoError(t, err)
	quickUser, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, t.Name(), mnemonic, genesisWalletAmount, quicksilver)
	require.NoError(t, err)
	junoUser, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, t.Name(), mnemonic, genesisWalletAmount, juno)
	require.NoError(t, err)
	entropy, err = bip39.NewEntropy(256)
	require.NoError(t, err)
	mnemonic, err = bip39.NewMnemonic(entropy)
	require.NoError(t, err)
	icqUser, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, t.Name(), mnemonic, genesisWalletAmount, quicksilver)
	require.NoError(t, err)

	err = quicksilver.GetNode().RecoverKey(ctx, quickUser.KeyName(), quickUser.Mnemonic())
	require.NoError(t, err)
	err = juno.GetNode().RecoverKey(ctx, junoUser.KeyName(), junoUser.Mnemonic())
	require.NoError(t, err)

	util.RunICQ(t, ctx, quicksilver, juno, icqUser)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, quicksilver, juno)
	require.NoError(t, err)

	// register zone

	propMsg := icstypes.RegisterZoneProposal{
		Title:            "Onboard Juno onto Quicksilver",
		Description:      "Test onboarding",
		ConnectionId:     "connection-0",
		BaseDenom:        "ujuno",
		LocalDenom:       "uqjuno",
		AccountPrefix:    "juno",
		MultiSend:        false,
		LiquidityModule:  false,
		MessagesPerTx:    8,
		ReturnToSender:   false,
		DepositsEnabled:  true,
		UnbondingEnabled: true,
		Decimals:         6,
		Is_118:           true,
		DustThreshold:    math.OneInt(),
		TransferChannel:  "channel-0",
	}

	packedMsg, err := codectypes.NewAnyWithValue(&propMsg)
	require.NoError(t, err)

	govMsg := govtypes.MsgExecLegacyContent{
		Content:   packedMsg,
		Authority: "quick10d07y265gmmuvt4z0w9aw880jnsr700j3xrh0p",
	}

	packedGov, err := codectypes.NewAnyWithValue(&govMsg)
	require.NoError(t, err)

	govBytes := qsapp.MakeEncodingConfig().Marshaler.MustMarshalJSON(packedGov)
	//msgBytes := qsapp.MakeEncodingConfig().Marshaler.MustMarshalJSON(packedMsg)

	//t.Log(string(govBytes))
	require.NoError(t, err)
	prop := cosmos.TxProposalv1{
		Messages: []json.RawMessage{govBytes},
		Title:    "Onboard Juno onto Quicksilver",
		Deposit:  "10000000uqck",
		Summary:  "Onboard Juno onto Quicksilver",
	}

	txid, err := quicksilver.GetNode().SubmitProposal(ctx, quickUser.KeyName(), prop)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, quicksilver, juno)
	require.NoError(t, err)

	var proposalId string
	testutil.WaitForCondition(time.Second*30, time.Second, func() (bool, error) {
		stdOut, _, err := quicksilver.GetNode().ExecQuery(ctx, "tx", txid)
		require.NoError(t, err)

		var res sdk.TxResponse
		json.Unmarshal(stdOut, &res)
		if len(res.Logs) == 0 {
			return false, nil
		}

		for _, event := range res.Logs[0].Events {
			if event.Type == "proposal_deposit" {
				for _, attr := range event.Attributes {
					if attr.Key == "proposal_id" {
						proposalId = attr.Value
						t.Log("proposalId", proposalId)
						break
					}
				}
			}
		}
		if proposalId == "" {
			return false, nil
		}
		return true, nil
	})
	if proposalId == "" {
		t.Fatal("proposalId not found")
	}
	err = quicksilver.VoteOnProposalAllValidators(ctx, proposalId, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")
	err = testutil.WaitForBlocks(ctx, 1, quicksilver, juno)
	require.NoError(t, err)

	height, _ := quicksilver.Height(ctx)

	_, err = cosmos.PollForProposalStatus(ctx, quicksilver, height, height+10, proposalId, "PROPOSAL_STATUS_PASSED")
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = testutil.WaitForBlocks(ctx, 20, quicksilver)
	require.NoError(t, err)
	// query zone here.
	out, _, err := quicksilver.GetNode().ExecQuery(ctx, "ics", "zones")
	require.NoError(t, err)
	zones := icstypes.QueryZonesResponse{}
	err = quicksilver.Config().EncodingConfig.Codec.UnmarshalJSON(out, &zones)
	require.NoError(t, err)
	require.Equal(t, 1, len(zones.Zones))
	require.Equal(t, juno.Config().ChainID, zones.Zones[0].ChainId)
	require.NotNil(t, zones.Zones[0].DelegationAddress)
	require.NotNil(t, zones.Zones[0].DepositAddress)
	require.NotNil(t, zones.Zones[0].PerformanceAddress)
	require.NotNil(t, zones.Zones[0].WithdrawalAddress)

	t.Log("Deposit Address", zones.Zones[0].DepositAddress.Address)
	t.Log("Delegation Address", zones.Zones[0].DelegationAddress.Address)
	t.Log("User Address", junoUser.FormattedAddress())

	err = juno.GetNode().SendFunds(ctx, junoUser.KeyName(), ibc.WalletAmount{
		Address: zones.Zones[0].DepositAddress.Address,
		Amount:  math.NewInt(5000000),
		Denom:   "ujuno",
	})

	require.NoError(t, err)
	err = testutil.WaitForBlocks(ctx, 25, quicksilver, juno)
	require.NoError(t, err)

	qjunoBalance := math.ZeroInt()
	testutil.WaitForCondition(time.Second*30, time.Second*5, func() (bool, error) {
		out, _, err = quicksilver.GetNode().ExecQuery(ctx, "bank", "balances", quickUser.FormattedAddress())
		require.NoError(t, err)
		var res banktypes.QueryAllBalancesResponse
		err = quicksilver.Config().EncodingConfig.Codec.UnmarshalJSON(out, &res)
		require.NoError(t, err)
		for _, balance := range res.Balances {
			if balance.Denom == "uqjuno" {
				qjunoBalance = balance.Amount
				return true, nil
			}
		}
		return false, nil
	})
	require.Equal(t, math.NewInt(5000000), qjunoBalance)

	testutil.WaitForCondition(time.Second*30, time.Second*5, func() (bool, error) {
		out, _, err = juno.GetNode().ExecQuery(ctx, "staking", "delegations", zones.Zones[0].DelegationAddress.Address)
		require.NoError(t, err)
		var res stakingtypes.QueryDelegatorDelegationsResponse
		err = juno.Config().EncodingConfig.Codec.UnmarshalJSON(out, &res)
		require.NoError(t, err)
		if len(res.DelegationResponses) != len(juno.Validators) {
			return false, nil
		}
		t.Log("Delegations", res.DelegationResponses)
		return true, nil
	})

	// find next epoch. wait until then. query after epoch ends.

	out, _, err = quicksilver.GetNode().ExecQuery(ctx, "epochs", "epoch-infos")
	require.NoError(t, err)
	var res epochstypes.QueryEpochsInfoResponse
	err = quicksilver.Config().EncodingConfig.Codec.UnmarshalJSON(out, &res)
	require.NoError(t, err)
	t.Log("Epochs", res.Epochs)
	for _, epoch := range res.Epochs {
		if epoch.Identifier == "epoch" {
			t.Log("Waiting for epoch to end", epoch)
			time.Sleep(time.Until(epoch.StartTime.Add(epoch.Duration)))
			break
		}
	}
	err = testutil.WaitForBlocks(ctx, 25, quicksilver, juno)
	require.NoError(t, err)

	// check RR
	out, _, err = quicksilver.GetNode().ExecQuery(ctx, "ics", "zone", zones.Zones[0].ChainId)
	require.NoError(t, err)
	var zone icstypes.QueryZoneResponse
	err = quicksilver.Config().EncodingConfig.Codec.UnmarshalJSON(out, &zone)
	require.NoError(t, err)
	t.Log("Zone", zone.Zone)
	require.Greater(t, zone.Zone.RedemptionRate.MustFloat64(), 1.0)
	err = testutil.WaitForBlocks(ctx, 25, quicksilver, juno)
	require.NoError(t, err)
	val0Addr, err := quicksilver.Validators[0].AccountKeyBech32(ctx, "validator")
	require.NoError(t, err)
	val0Balance, _, err := quicksilver.GetNode().ExecQuery(ctx, "bank", "balances", val0Addr)
	var balance banktypes.QueryAllBalancesResponse
	err = quicksilver.Config().EncodingConfig.Codec.UnmarshalJSON(val0Balance, &balance)
	require.NoError(t, err)
	t.Log("Validator Balance", balance.Balances)

	// check delegation
	// check rewards

}
