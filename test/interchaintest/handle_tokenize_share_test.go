package interchaintest

import (
	"context"
	"encoding/json"
	"fmt"
	_ "path"
	"strconv"
	"strings"
	"testing"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	paramsutils "github.com/cosmos/cosmos-sdk/x/params/client/utils"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/quicksilver-zone/quicksilver/test/interchaintest/utils"
	istypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/testreporter"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
	"github.com/stretchr/testify/require"
	_ "go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

var GaiaImageVersion = "v14.1.0"

// TestHandleTokenizedShares
func TestHandleTokenizedShares(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver and gaia
	numVals := 3
	numFullNodes := 3

	config, err := createConfig()
	require.NoError(t, err)

	modifyGenesisQuicksilver := []cosmos.GenesisKV{
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
	modifyGenesisGaia := []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.voting_params.voting_period",
			Value: "20s",
		},
		{
			Key:   "app_state.gov.deposit_params.max_deposit_period",
			Value: "20s",
		},
		{
			Key:   "app_state.gov.deposit_params.min_deposit.0.amount",
			Value: "100000",
		},
		{
			Key:   "app_state.gov.deposit_params.min_deposit.0.denom",
			Value: "uatom",
		},

		{
			Key:   "app_state.staking.params.unbonding_time",
			Value: "30s",
		},
	}
	config.ModifyGenesis = cosmos.ModifyGenesis(modifyGenesisQuicksilver)
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "quicksilver",
			ChainConfig:   config,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:    "gaia",
			Version: GaiaImageVersion,
			ChainConfig: ibc.ChainConfig{
				GasPrices:     "0.0uatom",
				ModifyGenesis: cosmos.ModifyGenesis(modifyGenesisGaia),
			},
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
	users2 := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, quicksilver, gaia)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, quicksilver, gaia)
	require.NoError(t, err)
	// Get our Bech32 encoded user addresses
	quickUser, gaiaUser, gaiaUser2 := users[0], users[1], users2[1]

	err = quicksilver.CreateKey(ctx, "default")
	require.NoError(t, err)
	defaultAddr, err := quicksilver.FullNodes[0].AccountKeyBech32(ctx, "default")
	require.NoError(t, err)

	_ = runICQ(ctx, t, quicksilver, gaia)
	fmt.Println("Default Address: ", defaultAddr)
	err = quicksilver.SendFunds(ctx, quickUser.KeyName(), ibc.WalletAmount{
		Denom:   quicksilver.Config().Denom,
		Amount:  math.NewInt(100_000),
		Address: string(defaultAddr),
	})
	require.NoError(t, err)

	quickUserAddr := quickUser.FormattedAddress()
	gaiaUserAddr := gaiaUser.FormattedAddress()

	// Get original account balances
	quicksilverOrigBal, err := quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount.Sub(math.NewInt(100_000)), quicksilverOrigBal)

	gaiaOrigBal, err := gaia.GetBalance(ctx, gaiaUserAddr, gaia.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, gaiaOrigBal)

	// Compose an IBC transfer and send from Quicksilver -> gaia
	transferAmount := math.NewInt(1000)
	transfer := ibc.WalletAmount{
		Address: gaiaUserAddr,
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

	// Get the IBC denom for uqck on gaia
	quicksilverTokenDenom := transfertypes.GetPrefixedDenom(quickChannels[0].Counterparty.PortID, quickChannels[0].Counterparty.ChannelID, quicksilver.Config().Denom)
	quicksilverIBCDenom := transfertypes.ParseDenomTrace(quicksilverTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on gaia and are in the user acc on gaia
	quicksilverUpdateBal, err := quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal.Sub(transferAmount), quicksilverUpdateBal)

	gaiaUpdateBal, err := gaia.GetBalance(ctx, gaiaUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, gaiaUpdateBal)

	// Create new clients
	err = r.CreateClients(ctx, eRep, pathQuicksilverGaia, ibc.CreateClientOptions{TrustingPeriod: "330h"})
	require.NoError(t, err)

	// Create a new connection
	err = r.CreateConnections(ctx, eRep, pathQuicksilverGaia)
	require.NoError(t, err)

	connections, err := r.GetConnections(ctx, eRep, quicksilver.Config().ChainID)
	require.NoError(t, err)

	// Compose an IBC transfer and send from Quicksilver -> gaia
	transfer = ibc.WalletAmount{
		Address: quickUserAddr,
		Denom:   quicksilverIBCDenom,
		Amount:  transferAmount,
	}

	transferTx, err = gaia.SendIBCTransfer(ctx, quickChannels[0].Counterparty.ChannelID, gaiaUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	gaiaHeight, err := gaia.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, gaia, gaiaHeight, gaiaHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Assert that the funds are now back on gaia and not on gaia
	quicksilverUpdateBal, err = quicksilver.GetBalance(ctx, quickUserAddr, quicksilver.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, quicksilverOrigBal, quicksilverUpdateBal)

	gaiaUpdateBal, err = gaia.GetBalance(ctx, gaiaUserAddr, quicksilverIBCDenom)
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), gaiaUpdateBal)

	registerProposal := istypes.RegisterZoneProposal{
		Title:            "Register zone",
		Description:      "Register zone",
		ConnectionId:     "connection-0",
		BaseDenom:        "uatom",
		LocalDenom:       "qatom",
		AccountPrefix:    "quick",
		DepositsEnabled:  true,
		UnbondingEnabled: true,
		LiquidityModule:  false,
		ReturnToSender:   true,
		Decimals:         6,
		TransferChannel:  quickChannels[0].Counterparty.ChannelID,
		DustThreshold:    math.NewInt(10),
	}

	check, err := cdctypes.NewAnyWithValue(&registerProposal)
	require.NoError(t, err)

	message := govv1.MsgExecLegacyContent{
		Content:   check,
		Authority: "quick10d07y265gmmuvt4z0w9aw880jnsr700j3xrh0p",
	}
	msg, err := quicksilver.Config().EncodingConfig.Codec.MarshalInterfaceJSON(&message)
	fmt.Println("Msg: ", string(msg))
	require.NoError(t, err)

	proposal := utils.TxProposalv1{
		Metadata: "none",
		Deposit:  "500000000" + quicksilver.Config().Denom,
		Title:    "title",
		Summary:  "register lstest-1 zone with multisend and lsm enabled",
	}

	// Appending proposal data in messages
	proposal.Messages = append(proposal.Messages, msg)

	require.NoError(t, err)

	// Submitting a proposal on Quicksilver
	proposalID, err := utils.SubmitProposal(ctx, quicksilver, quickUserAddr, proposal)

	require.NoError(t, err)

	// Voting on the proposal
	err = quicksilver.VoteOnProposalAllValidators(ctx, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "Failed to submit votes")

	heightAfterVote, err := quicksilver.Height(ctx)
	require.NoError(t, err, "error fetching height before vote")

	// Checking the proposal with matching ID and status.
	_, err = cosmos.PollForProposalStatus(ctx, quicksilver, heightAfterVote, heightAfterVote+20, proposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "Proposal status did not change to passed in expected number of blocks")
	time.Sleep(15 * time.Second)
	zone, err := utils.QueryZones(ctx, quicksilver)
	require.NoError(t, err)

	// Deposit Address Check
	depositAddress := zone[0].DepositAddress.Address
	icaAddr, err := utils.QueryZoneICAAddress(ctx, quicksilver, depositAddress, connections[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, icaAddr)

	// Withdrawl Address Check
	withdralAddress := zone[0].WithdrawalAddress.Address
	icaAddr, err = utils.QueryZoneICAAddress(ctx, quicksilver, withdralAddress, connections[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, icaAddr)

	// Delegation Address Check
	delegationAddress := zone[0].DelegationAddress.Address
	icaAddr, err = utils.QueryZoneICAAddress(ctx, quicksilver, delegationAddress, connections[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, icaAddr)

	// Performance Address Check
	performanceAddress := zone[0].DelegationAddress.Address
	icaAddr, err = utils.QueryZoneICAAddress(ctx, quicksilver, performanceAddress, connections[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, icaAddr)

	var updateZoneValue []*istypes.UpdateZoneValue
	updateZoneValue = append(updateZoneValue, &istypes.UpdateZoneValue{
		Key:   "",
		Value: "",
	})

	validators, err := utils.QueryStakingValidators(ctx, quicksilver, stakingtypes.Bonded.String())
	fmt.Println(validators)
	require.NoError(t, err)

	delegateTx, err := utils.RequestStakingDelegate(
		ctx,
		quicksilver,
		validators[0].OperatorAddress,
		quickUserAddr,
		"1000"+quicksilver.Config().Denom,
	)
	fmt.Println(delegateTx)
	require.NoError(t, err)

	delegation, err := utils.QueryStakingDelegation(ctx, quicksilver, validators[0].OperatorAddress, quickUserAddr)
	fmt.Println(delegation)
	require.NoError(t, err)

	unbondTx, err := utils.RequestStakingUnbond(
		ctx,
		quicksilver,
		validators[0].OperatorAddress,
		quickUserAddr,
		"1000"+quicksilver.Config().Denom,
	)
	fmt.Println(unbondTx)
	require.NoError(t, err)

	response, err := utils.RequestICSRedeem(
		ctx,
		quicksilver,
		quickUserAddr,
		"1000"+quicksilver.Config().Denom,
	)
	fmt.Println(response)
	require.NoError(t, err)

	// Stake process
	gaiaValidators, err := utils.QueryStakingValidators(ctx, gaia, stakingtypes.BondStatusBonded)
	require.NoError(t, err)
	require.Len(t, gaiaValidators, 3)

	// Staking params gaia
	gaiaStakingParams, err := utils.QueryStakingParamsGRPC(ctx, gaia)

	require.NoError(t, err)
	fmt.Println(gaiaStakingParams)

	paramChangeValue := []string{"0.25", "250.0", "1.0"}
	rawValues := make([]json.RawMessage, len(paramChangeValue))
	for i, v := range paramChangeValue {
		rawValue, err := json.Marshal(v)
		require.NoError(t, err)
		rawValues[i] = rawValue
	}

	newStakingParams := paramsutils.ParamChangeProposalJSON{
		Title:       "Update liquid staking params",
		Description: "Update liquid staking params",
		Changes: paramsutils.ParamChangesJSON{
			paramsutils.ParamChangeJSON{
				Subspace: "staking",
				Key:      "GlobalLiquidStakingCap",
				Value:    rawValues[0],
			},
			paramsutils.ParamChangeJSON{
				Subspace: "staking",
				Key:      "ValidatorBondFactor",
				Value:    rawValues[1],
			},
			paramsutils.ParamChangeJSON{
				Subspace: "staking",
				Key:      "ValidatorLiquidStakingCap",
				Value:    rawValues[2],
			},
		},
		Deposit: "1000000" + gaia.Config().Denom,
	}
	paramTx, err := gaia.ParamChangeProposal(ctx, gaiaUser.KeyName(), &newStakingParams)
	require.NoError(t, err, "error submitting param change proposal tx")

	time.Sleep(10 * time.Second)

	// Voting on the proposal
	err = gaia.VoteOnProposalAllValidators(ctx, paramTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "Failed to submit votes")

	heightAfterVote, err = gaia.Height(ctx)
	require.NoError(t, err, "error fetching height before vote")

	// Checking the proposal with matching ID and status.
	_, err = cosmos.PollForProposalStatus(ctx, gaia, heightAfterVote, heightAfterVote+20, paramTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "Proposal status did not change to passed in expected number of blocks")

	// Verify the staking params have been updated
	param, _ := gaia.QueryParam(ctx, "staking", "ValidatorBondFactor")
	require.NotEmpty(t, param.Value, "MaxValidators value is nil")

	param, _ = gaia.QueryParam(ctx, "staking", "GlobalLiquidStakingCap")
	require.NotEmpty(t, param.Value, "GlobalLiquidStakingCap value is nil")

	param, _ = gaia.QueryParam(ctx, "staking", "ValidatorLiquidStakingCap")
	require.NotEmpty(t, param.Value, "ValidatorLiquidStakingCap value is nil")

	param, _ = gaia.QueryParam(ctx, "staking", "UnbondingTime")
	require.NotEmpty(t, param.Value, "ValidatorLiquidStakingCap value is nil")

	fmt.Println("unbonding time: ", param.Value)
	// do stake on gaia
	// Gaia validator self-bond
	gaiaValidatorAddr := gaiaValidators[0].OperatorAddress

	_, err = utils.RequestStakingDelegate(
		ctx,
		gaia,
		gaiaValidators[0].OperatorAddress,
		gaiaUser2.KeyName(),
		"100000"+gaia.Config().Denom,
	)
	require.NoError(t, err)

	tx, err := utils.ExecuteValidatorBond(ctx, gaia, gaiaUser2.KeyName(), gaiaValidatorAddr)
	require.NoError(t, err)
	require.NotEmpty(t, tx)

	delegateAmount := ibc.WalletAmount{
		Address: "abc",
		Denom:   gaia.Config().Denom,
		Amount:  math.NewInt(100000),
	}

	delegateTx, err = utils.RequestStakingDelegate(
		ctx,
		gaia,
		gaiaValidators[0].OperatorAddress,
		gaiaUserAddr,
		"100000"+gaia.Config().Denom,
	)
	require.NoError(t, err)
	require.NotEmpty(t, delegateTx)

	delegation, err = utils.QueryStakingDelegationGRPC(ctx, gaia, gaiaUserAddr, gaiaValidators[0].OperatorAddress)
	require.NoError(t, err)
	fmt.Println(delegation)
	tokenizeShareAmount := sdk.NewCoin(gaia.Config().Denom, math.NewInt(20000))
	tx, err = utils.ExecuteTokenizeShares(ctx, gaia, gaiaUser.KeyName(), gaiaValidatorAddr, gaiaUserAddr, tokenizeShareAmount)
	require.NoError(t, err)
	require.NotEmpty(t, tx)

	time.Sleep(60 * time.Second)
	// Validate share reduced
	gaiaDelegation, err := utils.QueryStakingDelegationGRPC(ctx, gaia, gaiaUserAddr, gaiaValidators[0].OperatorAddress)
	require.NoError(t, err)
	require.Equal(t, delegateAmount.Amount.Sub(tokenizeShareAmount.Amount).Int64(), gaiaDelegation.Shares.TruncateInt64())

	// Validate tokenized share balances
	recordID := int(1)
	shareDenom := fmt.Sprintf("%s/%s", strings.ToLower(gaiaValidators[0].OperatorAddress), strconv.Itoa(recordID))
	shareBalance, err := gaia.GetBalance(ctx, gaiaUserAddr, shareDenom)
	fmt.Println(shareBalance)
	require.NoError(t, err)
	balances, err := gaia.AllBalances(ctx, gaiaUserAddr)
	require.NoError(t, err)
	fmt.Println("All balances: ", balances)
	require.Equal(t, tokenizeShareAmount.Amount, shareBalance)

	// Attempt to transfer tokenized shares to deposit account
	tokenizeShare := ibc.WalletAmount{
		Address: depositAddress,
		Denom:   shareDenom,
		Amount:  tokenizeShareAmount.Amount,
	}
	err = gaia.SendFunds(ctx, gaiaUser.KeyName(), tokenizeShare)
	require.NoError(t, err)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, quicksilver, gaia)
	require.NoError(t, err)
	// get the balance of the deposit account
	depositBalance, err := gaia.GetBalance(ctx, depositAddress, tokenizeShare.Denom)
	require.NoError(t, err)
	require.Equal(t, tokenizeShare.Amount, depositBalance)

	err = testutil.WaitForBlocks(ctx, 10, quicksilver, gaia)
	require.NoError(t, err)

	delegationBalance, err := gaia.GetBalance(ctx, delegationAddress, tokenizeShare.Denom)
	require.NoError(t, err)
	require.Equal(t, tokenizeShare.Amount, delegationBalance)

}

type IcqKey struct {
	Mnemonic string `json:"mnemonic"`
	Address  string `json:"address"`
}
