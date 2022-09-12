package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	lsmstakingtypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
	"github.com/ingenuity-build/quicksilver/utils"
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

func (k msgServer) RequestRedemption(goCtx context.Context, msg *types.MsgRequestRedemption) (*types.MsgRequestRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate coins are positive
	err := msg.Value.Validate()
	if err != nil {
		return nil, err
	}

	// validate recipient address
	if len(msg.DestinationAddress) == 0 {
		return nil, fmt.Errorf("recipient address not provided")
	}

	if _, _, err = bech32.DecodeAndConvert(msg.DestinationAddress); err != nil {
		return nil, err
	}

	var zone *types.Zone

	k.IterateZones(ctx, func(_ int64, thisZone types.Zone) bool {
		if thisZone.LocalDenom == msg.Value.GetDenom() {
			zone = &thisZone
			return true
		}
		return false
	})

	// does zone exist?
	if nil == zone {
		return nil, fmt.Errorf("unable to find matching zone for denom %s", msg.Value.GetDenom())
	}

	// does destination address match the prefix registered against the zone?
	if _, err := utils.AccAddressFromBech32(msg.DestinationAddress, zone.AccountPrefix); err != nil {
		return nil, fmt.Errorf("destination address %s does not match expected prefix %s", msg.DestinationAddress, zone.AccountPrefix)
	}

	sender, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	// does the user have sufficient assets to burn
	if !k.BankKeeper.HasBalance(ctx, sender, msg.Value) {
		return nil, fmt.Errorf("account has insufficient balance of qasset to burn")
	}

	// get min of LastRedemptionRate (N-1) and RedemptionRate (N)
	var rate sdk.Dec
	rate = zone.LastRedemptionRate
	if zone.RedemptionRate.LT(rate) {
		rate = zone.RedemptionRate
	}

	nativeTokens := sdk.NewDecFromInt(msg.Value.Amount).Mul(rate).TruncateInt()

	outTokens := sdk.NewCoin(zone.BaseDenom, nativeTokens)
	k.Logger(ctx).Error("outtokens", "o", outTokens)

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(ctx.BlockHeight()))
	hash := sha256.Sum256(append(msg.GetSignBytes(), heightBytes...))
	hashString := hex.EncodeToString(hash[:])

	// lock qAssets - how are we tracking this?
	if err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(msg.Value)); err != nil {
		return nil, err
	}

	if zone.LiquidityModule {
		return nil, fmt.Errorf("lsm not currently supported")
		if err = k.processRedemptionForLsm(ctx, *zone, sender, msg.DestinationAddress, nativeTokens, msg.Value, hashString); err != nil {
			return nil, err
		}
	} else {
		if err = k.queueRedemption(ctx, *zone, sender, msg.DestinationAddress, nativeTokens, msg.Value, hashString); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRedemptionRequest,
			sdk.NewAttribute(types.AttributeKeyBurnAmount, msg.Value.String()),
			sdk.NewAttribute(types.AttributeKeyRedeemAmount, nativeTokens.String()),
			sdk.NewAttribute(types.AttributeKeyRecipientAddress, msg.DestinationAddress),
			sdk.NewAttribute(types.AttributeKeyRecipientChain, zone.ChainId),
			sdk.NewAttribute(types.AttributeKeyConnectionID, zone.ConnectionId),
		),
	})

	return &types.MsgRequestRedemptionResponse{}, nil
}

// processRedemptionForLsm will determine based on user intent, the tokens to return to the user, generate Redeem message and send them.
func (k *Keeper) processRedemptionForLsm(ctx sdk.Context, zone types.Zone, sender sdk.AccAddress, destination string, nativeTokens sdk.Int, burnAmount sdk.Coin, hash string) error {
	intent, found := k.GetIntent(ctx, zone, sender.String(), false)
	// msgs is slice of MsgTokenizeShares, so we can handle dust allocation later.
	var msgs []*lsmstakingtypes.MsgTokenizeShares
	intents := intent.Intents
	if !found {
		// if user has no intent set (this can happen if redeeming tokens that were obtained offchain), use global intent.
		// Note: this can be improved; user will receive a bunch of tokens.
		intents = zone.AggregateIntent
	}
	outstanding := nativeTokens
	distribution := make(map[string]sdk.Int, 0)
	for _, valoper := range utils.Keys(intents) {
		intent := intents[valoper]
		thisAmount := intent.Weight.MulInt(nativeTokens).TruncateInt()
		distribution[valoper] = thisAmount
		outstanding = outstanding.Sub(thisAmount)
		msgs = append(msgs, &stakingtypes.MsgTokenizeShares{
			DelegatorAddress:    zone.DelegationAddress.Address,
			ValidatorAddress:    intent.ValoperAddress,
			Amount:              sdk.NewCoin(zone.BaseDenom, thisAmount),
			TokenizedShareOwner: destination,
		})
	}

	// add unallocated dust.
	msgs[0].Amount = msgs[0].Amount.AddAmount(outstanding)
	var sdkMsgs []sdk.Msg
	for _, msg := range msgs {
		sdkMsgs = append(sdkMsgs, sdk.Msg(msg))
	}
	k.AddWithdrawalRecord(ctx, zone, sender.String(), distribution, destination, sdk.Coins{}, burnAmount, hash, WithdrawStatusTokenize, time.Unix(0, 0))

	return k.SubmitTx(ctx, sdkMsgs, zone.DelegationAddress, hash)
}

// queueRedemption will determine based on zone intent, the tokens to unbond, and add a withdrawal record with status QUEUED.
func (k *Keeper) queueRedemption(ctx sdk.Context, zone types.Zone, sender sdk.AccAddress, destination string, nativeTokens sdk.Int, burnAmount sdk.Coin, hash string) error {
	distribution := make(map[string]sdk.Int, 0)
	outstanding := nativeTokens

	for _, valoper := range utils.Keys(zone.AggregateIntent) {
		intent := zone.AggregateIntent[valoper]
		thisAmount := intent.Weight.MulInt(nativeTokens).TruncateInt()
		outstanding = outstanding.Sub(thisAmount)
		distribution[valoper] = thisAmount
	}

	k.AddWithdrawalRecord(ctx, zone, sender.String(), distribution, destination, sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, nativeTokens)), burnAmount, hash, WithdrawStatusQueued, time.Unix(0, 0))

	return nil
}

func IntentSliceToMap(in []*types.ValidatorIntent) (out map[string]*types.ValidatorIntent) {
	out = make(map[string]*types.ValidatorIntent, 0)
	for _, intent := range in {
		out[intent.ValoperAddress] = intent
	}
	return
}

func (k msgServer) SignalIntent(goCtx context.Context, msg *types.MsgSignalIntent) (*types.MsgSignalIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get zone
	zone, ok := k.GetZone(ctx, msg.ChainId)
	if !ok {
		return nil, fmt.Errorf("invalid chain id \"%s\"", msg.ChainId)
	}

	// validate intents (aggregated errors)
	if err := k.validateIntents(zone, msg.Intents); err != nil {
		return nil, err
	}

	intent := types.DelegatorIntent{
		Delegator: msg.FromAddress,
		Intents:   IntentSliceToMap(msg.Intents),
	}

	k.SetIntent(ctx, zone, intent, false)

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

func (k msgServer) validateIntents(zone types.Zone, intents []*types.ValidatorIntent) error {
	errors := make(map[string]error)

	for i, intent := range intents {
		_, found := zone.GetValidatorByValoper(intent.ValoperAddress)
		if !found {
			errors[fmt.Sprintf("intent[%v]", i)] = fmt.Errorf("unable to find valoper %s", intent.ValoperAddress)
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

func (k Keeper) EmitValsetRequery(ctx sdk.Context, connectionID string, chainID string) error {
	query := stakingtypes.QueryValidatorsRequest{}
	bz1, err := k.cdc.Marshal(&query)
	if err != nil {
		return err
	}

	period := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))

	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		chainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz1,
		sdk.NewInt(period),
		types.ModuleName,
		"valset",
		0,
	)
	return nil
}
