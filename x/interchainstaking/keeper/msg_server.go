package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the interchainstaking
// MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: &keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RegisterZone(goCtx context.Context, msg *types.MsgRegisterZone) (*types.MsgRegisterZoneResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get chain id from connection
	chainId, err := k.GetChainID(ctx, msg.ConnectionId)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain chain id: %w", err)
	}

	// get zone
	_, found := k.GetRegisteredZoneInfo(ctx, chainId)
	if found {
		return nil, fmt.Errorf("invalid chain id, zone for \"%s\" already registered", chainId)
	}

	zone := types.RegisteredZone{
		Identifier:         msg.Identifier,
		ChainId:            chainId,
		ConnectionId:       msg.ConnectionId,
		LocalDenom:         msg.LocalDenom,
		BaseDenom:          msg.BaseDenom,
		RedemptionRate:     sdk.NewDec(1),
		LastRedemptionRate: sdk.NewDec(1),
		DelegatorIntent:    make(map[string]*types.DelegatorIntent),
		MultiSend:          msg.MultiSend,
	}
	k.SetRegisteredZone(ctx, zone)

	// generate deposit account
	portOwner := chainId + ".deposit"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return nil, err
	}

	// generate fee account
	portOwner = chainId + ".fee"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return nil, err
	}

	// generate withdrawal account
	portOwner = chainId + ".withdrawal"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return nil, err
	}

	// generate delegate accounts
	delegateAccountCount := int(k.GetParam(ctx, types.KeyDelegateAccountCount))
	for i := 0; i < delegateAccountCount; i++ {
		portOwner := fmt.Sprintf("%s.delegate.%d", chainId, i)
		if err := k.ICAControllerKeeper.RegisterInterchainAccount(
			ctx,
			zone.ConnectionId,
			portOwner,
		); err != nil {
			return nil, err
		}
		portId, _ := icatypes.NewControllerPortID(portOwner)
		if err := k.SetConnectionForPort(ctx, msg.ConnectionId, portId); err != nil {
			return nil, err
		}
	}

	valsetInterval := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))
	bondedValidatorQuery := k.ICQKeeper.NewPeriodicQuery(
		ctx,
		msg.ConnectionId,
		chainId,
		"cosmos.staking.v1beta1.Query/Validators",
		map[string]string{"status": stakingtypes.BondStatusBonded},
		sdk.NewInt(valsetInterval),
	)
	k.ICQKeeper.SetPeriodicQuery(ctx, *bondedValidatorQuery)
	unbondedValidatorQuery := k.ICQKeeper.NewPeriodicQuery(
		ctx,
		msg.ConnectionId,
		chainId,
		"cosmos.staking.v1beta1.Query/Validators",
		map[string]string{"status": stakingtypes.BondStatusUnbonded},
		sdk.NewInt(valsetInterval),
	)
	k.ICQKeeper.SetPeriodicQuery(ctx, *unbondedValidatorQuery)
	unbondingValidatorQuery := k.ICQKeeper.NewPeriodicQuery(
		ctx,
		msg.ConnectionId,
		chainId,
		"cosmos.staking.v1beta1.Query/Validators",
		map[string]string{"status": stakingtypes.BondStatusUnbonding},
		sdk.NewInt(valsetInterval),
	)
	k.ICQKeeper.SetPeriodicQuery(ctx, *unbondingValidatorQuery)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRegisterZone,
			sdk.NewAttribute(types.AttributeKeyConnectionId, msg.ConnectionId),
			sdk.NewAttribute(types.AttributeKeyConnectionId, chainId),
		),
	})

	return &types.MsgRegisterZoneResponse{}, nil
}

func (k msgServer) registerInterchainAccount(ctx sdk.Context, connectionId string, portOwner string) error {
	if err := k.ICAControllerKeeper.RegisterInterchainAccount(ctx, connectionId, portOwner); err != nil {
		return err
	}
	portId, _ := icatypes.NewControllerPortID(portOwner)
	if err := k.SetConnectionForPort(ctx, connectionId, portId); err != nil {
		return err
	}

	return nil
}

func (k msgServer) RequestRedemption(goCtx context.Context, msg *types.MsgRequestRedemption) (*types.MsgRequestRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate coins are positive
	inCoin, err := sdk.ParseCoinNormalized(msg.Coin)
	if err != nil {
		return nil, err
	}

	if !inCoin.IsPositive() {
		return nil, fmt.Errorf("invalid input coin value")
	}

	// validate recipient address
	if len(msg.DestinationAddress) == 0 {
		return nil, fmt.Errorf("recipient address not provided")
	}

	_, _, err = bech32.DecodeAndConvert(msg.DestinationAddress)
	if err != nil {
		return nil, err
	}

	// TODO: store HRP of RegisteredZone to validate destination address (don't let users do stupid things!)

	var zone types.RegisteredZone

	k.IterateRegisteredZones(ctx, func(_ int64, thisZone types.RegisteredZone) bool {
		if thisZone.LocalDenom == inCoin.GetDenom() {
			zone = thisZone
			return true
		}
		return false
	})
	k.Logger(ctx).Error("DEBUG 6")

	sender, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	if !k.BankKeeper.HasBalance(ctx, sender, inCoin) {
		return nil, fmt.Errorf("account has insufficient balance of qasset to burn")
	}
	k.Logger(ctx).Error("DEBUG 7")

	// get min of LastRedemptionRate (N-1) and RedemptionRate (N)
	var rate sdk.Dec
	rate = zone.LastRedemptionRate
	if zone.RedemptionRate.LT(rate) {
		rate = zone.RedemptionRate
	}
	native_tokens := inCoin.Amount.ToDec().Mul(rate).TruncateInt()
	k.Logger(ctx).Error("DEBUG 8")

	outTokens := sdk.NewCoin(zone.BaseDenom, native_tokens)
	k.Logger(ctx).Error("DEBUG 9", "outTokens", outTokens, "nativeTokens", native_tokens)

	// lock qAssets - how are we tracking this?
	k.BankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msg.FromAddress), types.ModuleName, sdk.NewCoins(inCoin))
	k.Logger(ctx).Error("DEBUG 10")

	// send message
	userIntent, found := k.GetIntent(ctx, zone, msg.FromAddress)
	k.Logger(ctx).Error("DEBUG 11", "intent", userIntent)

	if !found {
		// here we should use the tokens the zone WANTS to get rid of!
		// fetch cachedIntent vs  currentState from zone.
		// for now:
		userIntent = types.DelegatorIntent{Delegator: msg.FromAddress, Intents: []*types.ValidatorIntent{}}
		k.Logger(ctx).Error("DEBUG 11a")

	}

	intentMap := userIntent.ToMap(native_tokens)
	k.Logger(ctx).Error("DEBUG 11", "intent", intentMap)

	targets := zone.GetRedemptionTargets(intentMap, zone.BaseDenom) // map[string][string]sdk.Coin
	k.Logger(ctx).Error("DEBUG 12")

	if len(targets) == 0 {
		return nil, fmt.Errorf("targets can never be zero length")
	}

	sumAmount := sdk.NewCoins()

	for delegator, validators := range targets {
		k.Logger(ctx).Error("DEBUG 13", "delegator", delegator)

		msgs := make([]sdk.Msg, 0)
		for validator, amount := range validators {
			k.Logger(ctx).Error("DEBUG 13a", "delegator", delegator, "validator", validator)
			msgs = append(msgs, &stakingtypes.MsgTokenizeShares{
				DelegatorAddress:    delegator,
				ValidatorAddress:    validator,
				Amount:              amount,
				TokenizedShareOwner: msg.DestinationAddress,
			})
			sumAmount = sumAmount.Add(amount)
			// what happens here if we revert below? this is stored in the KV store, so it should be rolled back. Check me.
			k.AddWithdrawalRecord(ctx, delegator, validator, msg.DestinationAddress, amount)
		}
		icaAccount, err := zone.GetDelegationAccountByAddress(delegator)
		k.Logger(ctx).Error("DEBUG 13b", "delegator", delegator, "ica", icaAccount)
		if err != nil {
			panic(err) // panic here because something is terribly wrong if we cann't find the delegation bucket here!!!
		}
		k.Logger(ctx).Error("DEBUG 13c", "msgs", msgs)
		k.SubmitTx(ctx, msgs, icaAccount)
	}
	k.Logger(ctx).Error("DEBUG 14")

	if !sumAmount.IsAllLTE(sdk.NewCoins(outTokens)) {
		k.Logger(ctx).Error("Output coins > than expected!", "sum", sumAmount, "expected", outTokens)
		panic("argh")
	}
	k.Logger(ctx).Error("DEBUG 15")

	//msgs = append(msgs, &banktypes.MsgSend{t.DelegatorsAddress(), msg.DestinationAddress, t.Amount()})
	// on confirmation of asset dispersal, burn qAssets?

	// burn qAssets

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRedemptionRequest,
			sdk.NewAttribute(types.AttributeKeyBurnAmount, msg.Coin),
			sdk.NewAttribute(types.AttributeKeyRedeemAmount, sumAmount.String()),
			sdk.NewAttribute(types.AttributeKeyRecipientAddress, msg.DestinationAddress),
			sdk.NewAttribute(types.AttributeKeyRecipientChain, zone.ChainId),
			sdk.NewAttribute(types.AttributeKeyConnectionId, zone.ConnectionId),
		),
	})

	return &types.MsgRequestRedemptionResponse{}, nil
}

func (k msgServer) SignalIntent(goCtx context.Context, msg *types.MsgSignalIntent) (*types.MsgSignalIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get zone
	zone, ok := k.GetRegisteredZoneInfo(ctx, msg.ChainId)
	if !ok {
		return nil, fmt.Errorf("invalid chain id \"%s\"", msg.ChainId)
	}

	// validate intents (aggregated errors)
	if err := k.validateIntents(zone, msg.Intents); err != nil {
		return nil, err
	}

	intent := types.DelegatorIntent{
		Delegator: msg.FromAddress,
		Intents:   msg.Intents,
	}

	k.SetIntent(ctx, zone, intent)

	// ctx.EventManager().EmitEvents(sdk.Events{
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 	),
	// 	sdk.NewEvent(
	// 		types.EventTypeRegisterZone,
	// 		sdk.NewAttribute(types.AttributeKeyConnectionId, msg.ConnectionId),
	// 		sdk.NewAttribute(types.AttributeKeyConnectionId, msg.ChainId),
	// 	),
	// })

	return &types.MsgSignalIntentResponse{}, nil
}

func (k msgServer) validateIntents(zone types.RegisteredZone, intents []*types.ValidatorIntent) error {
	errors := make(map[string]error)

	for i, intent := range intents {
		_, err := zone.GetValidatorByValoper(intent.ValoperAddress)
		if err != nil {
			errors[fmt.Sprintf("intent[%v]", i)] = err
		}
	}

	if len(errors) > 0 {
		return types.NewMultiError(errors)
	}

	return nil
}
