package simulation

import (
	"errors"
	"fmt"
	"math/rand"

	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/ingenuity-build/quicksilver/simulation/simtypes"
	"github.com/ingenuity-build/quicksilver/utils"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdksimtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	OpWeightMsgSignalIntent               = "op_weight_msg_signal_intent"      //nolint:gosec // not credentials
	OpWeightMsgRequestRedemption          = "op_weight_msg_request_redemption" //nolint:gosec // not credentials
	DefaultWeightMsgSignalIntent      int = 50
	DefaultWeightMsgRequestRedemption int = 10
)

var (
	TypeMsgSignalIntent      = sdk.MsgTypeURL(&types.MsgSignalIntent{})
	TypeMsgRequestRedemption = sdk.MsgTypeURL(&types.MsgRequestRedemption{})
)

func WeightedOperations(
	appParams sdksimtypes.AppParams,
	cdc codec.JSONCodec,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.ScopedIBCKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgSignalIntent      int
		weightMsgRequestRedemption int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgSignalIntent, &weightMsgSignalIntent, nil,
		func(_ *rand.Rand) {
			weightMsgSignalIntent = DefaultWeightMsgSignalIntent
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRequestRedemption, &weightMsgRequestRedemption, nil,
		func(_ *rand.Rand) {
			weightMsgRequestRedemption = DefaultWeightMsgRequestRedemption
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgSignalIntent,
			SimulateMsgSignalIntent(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgRequestRedemption,
			SimulateMsgRequestRedemption(ak, bk, sk, k),
		),
	}
}

var ibcPortsSetup = make(map[string]bool)

// SimulateMsgSignalIntent generates a MsgSignalIntent with random values.
func SimulateMsgSignalIntent(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, bApp *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgSignalIntent, "TODO"), nil, nil
	}
}

// SimulateMsgRequestRedemption generates a MsgRequestRedemption with random values.
func SimulateMsgRequestRedemption(ak types.AccountKeeper, bk types.BankKeeper, sk types.ScopedIBCKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, bApp *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		from, balance, err := randomSimAccountWithQAsset(ctx, r, accs, bk)
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgRequestRedemption, "could not find acc with q asset"), nil, nil
		}

		amt := sdk.NewInt(r.Int63n(balance.Amount.Int64()))
		value := sdk.NewCoin(balance.Denom, amt)

		var zone *types.Zone
		k.IterateZones(ctx, func(_ int64, thisZone *types.Zone) bool {
			if thisZone.LocalDenom == value.GetDenom() {
				zone = thisZone
				return true
			}
			return false
		})

		dest, err := bech32.ConvertAndEncode(zone.AccountPrefix, utils.GenerateAccAddressForTest(r))
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgRequestRedemption, "unable to generate dest account"), nil, nil
		}

		// ensure that channel exists for the zone
		portID := zone.GetDelegationAddress().GetPortName()
		if found := ibcPortsSetup[portID]; !found {
			channelID := k.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(ctx)
			connectionID, _ := k.GetConnectionForPort(ctx, portID)
			k.IBCKeeper.ChannelKeeper.SetChannel(ctx, portID, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icatypes.PortID, ChannelId: channelID}, ConnectionHops: []string{connectionID}})
			k.ICAControllerKeeper.SetActiveChannelID(ctx, connectionID, portID, channelID)
			path := host.ChannelCapabilityPath(portID, channelID)
			fmt.Printf("\nsetting up mock channel capability, portID: %s, channelID: %s, connectionID: %s, cap: %s\n",
				portID, channelID, connectionID, path)
			chanCap, err := k.ScopedKeeper().NewCapability(
				ctx,
				path,
			)
			if err != nil {
				panic(err)
			}

			err = sk.ClaimCapability(ctx, chanCap, path)
			if err != nil {
				panic(err)
			}
			k.ICAControllerKeeper.SetActiveChannelID(ctx, connectionID, portID, channelID)
			k.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, portID, channelID, 1)

			ibcPortsSetup[portID] = true
		}

		msg := &types.MsgRequestRedemption{
			Value:              value,
			DestinationAddress: dest,
			FromAddress:        from.Address.String(),
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             bApp,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         TypeMsgRequestRedemption,
			CoinsSpentInMsg: sdk.NewCoins(value), // coins burned
			Context:         ctx,
			SimAccount:      from,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func randomSimAccountWithQAsset(ctx sdk.Context, r *rand.Rand, accs []sdksimtypes.Account, bk types.BankKeeper) (sdksimtypes.Account, sdk.Coin, error) {
	coins := sdk.NewCoins(sdk.NewCoin("uqatom", sdk.OneInt()), sdk.NewCoin("uqosmo", sdk.OneInt()), sdk.NewCoin("uqjunox", sdk.OneInt()))
	randomQAsset := coins[r.Intn(len(coins))]

	accHasQAsset := func(acc sdksimtypes.Account) bool {
		spendableCoins := bk.SpendableCoins(ctx, acc.Address)
		if spendableCoins.Empty() {
			return false
		}

		if spendableCoins.AmountOf(randomQAsset.Denom).IsPositive() {
			return true
		}

		return false
	}

	acc, found := simtypes.RandomSimAccountWithConstraint(r, accHasQAsset, accs)
	if !found {
		return sdksimtypes.Account{}, sdk.Coin{}, errors.New("no address with q assets found")
	}

	spendableCoins := bk.SpendableCoins(ctx, acc.Address)
	var asset sdk.Coin
	for _, coin := range coins {
		found, c := spendableCoins.Find(coin.Denom)
		if found {
			asset = c
			break
		}
	}

	return acc, asset, nil
}
