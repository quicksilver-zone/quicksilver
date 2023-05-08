package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
<<<<<<< HEAD
	"sort"
	"strings"
	"time"
=======
	"strings"
>>>>>>> origin/develop

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
<<<<<<< HEAD
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	lsmstakingtypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"
=======
>>>>>>> origin/develop

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

// RequestRedemption handles MsgRequestRedemption by creating a corresponding withdrawal record queued for unbonding.
func (k msgServer) RequestRedemption(goCtx context.Context, msg *types.MsgRequestRedemption) (*types.MsgRequestRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.Keeper.GetUnbondingEnabled(ctx) {
		return nil, fmt.Errorf("unbonding is currently disabled")
	}

	var zone *types.Zone
	k.IterateZones(ctx, func(_ int64, thisZone *types.Zone) bool {
		if thisZone.LocalDenom == msg.Value.GetDenom() {
			zone = thisZone
			return true
		}
		return false
	})

	// does zone exist?
	if zone == nil {
		return nil, fmt.Errorf("unable to find matching zone for denom %s", msg.Value.GetDenom())
	}

	if !zone.UnbondingEnabled {
		return nil, fmt.Errorf("unbonding currently disabled for zone %s", zone.ChainId)
	}

	// does destination address match the prefix registered against the zone?
	if _, err := utils.AccAddressFromBech32(msg.DestinationAddress, zone.AccountPrefix); err != nil {
		return nil, fmt.Errorf("destination address %s does not match expected prefix %s [%w]", msg.DestinationAddress, zone.AccountPrefix, err)
	}

	sender, _ := sdk.AccAddressFromBech32(msg.FromAddress) // already validated

	// does the user have sufficient assets to burn
	if !k.BankKeeper.HasBalance(ctx, sender, msg.Value) {
		return nil, errors.New("account has insufficient balance of qasset to burn")
	}

	// get min of LastRedemptionRate (N-1) and RedemptionRate (N)
	rate := sdk.MinDec(zone.LastRedemptionRate, zone.RedemptionRate)
	nativeTokens := sdk.NewDecFromInt(msg.Value.Amount).Mul(rate).TruncateInt()
	outTokens := sdk.NewCoin(zone.BaseDenom, nativeTokens)
	k.Logger(ctx).Info("tokens to distribute", "amount", outTokens)

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(ctx.BlockHeight()))
	hash := sha256.Sum256(append(msg.GetSignBytes(), heightBytes...))
	hashString := hex.EncodeToString(hash[:])

	if err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.EscrowModuleAccount, sdk.NewCoins(msg.Value)); err != nil {
		return nil, fmt.Errorf("unable to send coins to escrow account: %w", err)
	}

	if zone.LiquidityModule {
		if err := k.processRedemptionForLsm(ctx, zone, sender, msg.DestinationAddress, nativeTokens, msg.Value, hashString); err != nil {
			return nil, fmt.Errorf("unable to process redemption for LSM: %w", err)
		}
	} else {
		if err := k.queueRedemption(ctx, zone, sender, msg.DestinationAddress, nativeTokens, msg.Value, hashString); err != nil {
			return nil, fmt.Errorf("unable to queue redemption: %w", err)
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
			sdk.NewAttribute(types.AttributeKeyChainID, zone.ChainId),
			sdk.NewAttribute(types.AttributeKeyConnectionID, zone.ConnectionId),
		),
	})

	return &types.MsgRequestRedemptionResponse{}, nil
}

<<<<<<< HEAD
// processRedemptionForLsm will determine based on user intent, the tokens to return to the user, generate Redeem message and send them.
func (k *Keeper) processRedemptionForLsm(ctx sdk.Context, zone types.Zone, sender sdk.AccAddress, destination string, nativeTokens math.Int, burnAmount sdk.Coin, hash string) error {
	intent, found := k.GetIntent(ctx, zone, sender.String(), false)
	// msgs is slice of MsgTokenizeShares, so we can handle dust allocation later.
	msgs := make([]*lsmstakingtypes.MsgTokenizeShares, 0)
	intents := intent.Intents
	if !found || len(intents) == 0 {
		// if user has no intent set (this can happen if redeeming tokens that were obtained offchain), use global intent.
		// Note: this can be improved; user will receive a bunch of tokens.
		intents = zone.GetAggregateIntentOrDefault()
	}
	outstanding := nativeTokens
	distribution := make(map[string]uint64, 0)

	availablePerValidator := k.GetUnlockedTokensForZone(ctx, &zone)

	for _, intent := range intents.Sort() {
		thisAmount := intent.Weight.MulInt(nativeTokens).TruncateInt()
		if thisAmount.Int64() > availablePerValidator[intent.ValoperAddress] {
			return errors.New("unable to satisfy unbond request; delegations may be locked")
		}
		distribution[intent.ValoperAddress] = thisAmount.Uint64()
		outstanding = outstanding.Sub(thisAmount)
	}

	distribution[intents[0].ValoperAddress] += outstanding.Uint64()

	for _, valoper := range utils.Keys(distribution) {
		msgs = append(msgs, &lsmstakingtypes.MsgTokenizeShares{
			DelegatorAddress:    zone.DelegationAddress.Address,
			ValidatorAddress:    valoper,
			Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewIntFromUint64(distribution[valoper])),
			TokenizedShareOwner: destination,
		})
	}
	// add unallocated dust.
	msgs[0].Amount = msgs[0].Amount.AddAmount(outstanding)
	sdkMsgs := make([]sdk.Msg, 0)
	for _, msg := range msgs {
		sdkMsgs = append(sdkMsgs, sdk.Msg(msg))
	}
	k.AddWithdrawalRecord(ctx, zone.ChainId, sender.String(), []*types.Distribution{}, destination, sdk.Coins{}, burnAmount, hash, WithdrawStatusTokenize, time.Unix(0, 0))

	return k.SubmitTx(ctx, sdkMsgs, zone.DelegationAddress, hash, zone.MessagesPerTx)
}

// queueRedemption will determine based on zone intent, the tokens to unbond, and add a withdrawal record with status QUEUED.
func (k *Keeper) queueRedemption(
	ctx sdk.Context,
	zone types.Zone,
	sender sdk.AccAddress,
	destination string,
	nativeTokens math.Int,
	burnAmount sdk.Coin,
	hash string,
) error { //nolint:unparam // we know that the error is always nil
	distribution := make([]*types.Distribution, 0)
	outstanding := nativeTokens

	aggregateIntent := zone.GetAggregateIntentOrDefault()
	for _, intent := range aggregateIntent {
		thisAmount := intent.Weight.MulInt(nativeTokens).TruncateInt()
		outstanding = outstanding.Sub(thisAmount)
		dist := types.Distribution{
			Valoper: intent.ValoperAddress,
			Amount:  thisAmount.Uint64(),
		}

		distribution = append(distribution, &dist)
	}
	// handle dust ? ok to do uint64 calc here or do we use math.Int (just more verbose) ?
	distribution[0].Amount += outstanding.Uint64()

	amount := sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, nativeTokens))
	k.AddWithdrawalRecord(
		ctx,
		zone.ChainId,
		sender.String(),
		distribution,
		destination,
		amount,
		burnAmount,
		hash,
		WithdrawStatusQueued,
		time.Time{},
	)

	return nil
}

// this can be removed: unused
func IntentSliceToMap(in []*types.ValidatorIntent) (out map[string]*types.ValidatorIntent) {
	out = make(map[string]*types.ValidatorIntent, 0)

	sort.SliceStable(in, func(i, j int) bool {
		return in[i].ValoperAddress > in[j].ValoperAddress
	})

	for _, intent := range in {
		out[intent.ValoperAddress] = intent
	}
	return
}

=======
>>>>>>> origin/develop
func (k msgServer) SignalIntent(goCtx context.Context, msg *types.MsgSignalIntent) (*types.MsgSignalIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get zone
	zone, ok := k.GetZone(ctx, msg.ChainId)
	if !ok {
		return nil, fmt.Errorf("invalid chain id \"%s\"", msg.ChainId)
	}

	// validate intents (aggregated errors)
	intents, err := types.IntentsFromString(msg.Intents)
	if err != nil {
		return nil, err
	}

	if err := k.validateValidatorIntents(ctx, zone, intents); err != nil {
		return nil, err
	}

	intent := types.DelegatorIntent{
		Delegator: msg.FromAddress,
		Intents:   intents,
	}

	k.SetDelegatorIntent(ctx, &zone, intent, false)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeSetIntent,
			sdk.NewAttribute(types.AttributeKeyUser, msg.FromAddress),
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
		),
	})

	return &types.MsgSignalIntentResponse{}, nil
}

// GovReopenChannel reopens an ICA channel.
func (k msgServer) GovReopenChannel(goCtx context.Context, msg *types.MsgGovReopenChannel) (*types.MsgGovReopenChannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// remove leading prefix icacontroller- if passed in msg
	portID := strings.ReplaceAll(msg.PortId, "icacontroller-", "")

	// validate the zone exists, and the format is valid (e.g. quickgaia-1.delegate)
	parts := strings.Split(portID, ".")

	// portId and connectionId format validated in validateBasic, so not duplicated here.

	// assert chainId matches connectionId
	chainID, err := k.GetChainID(ctx, msg.ConnectionId)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain chain id: %w", err)
	}

	if chainID != parts[0] {
		return nil, fmt.Errorf("chainID / connectionID mismatch. Connection: %s, Port: %s", chainID, parts[0])
	}

	if _, found := k.GetZone(ctx, chainID); !found {
		return &types.MsgGovReopenChannelResponse{}, errors.New("invalid port format; zone not found")
	}

	if err := k.Keeper.registerInterchainAccount(ctx, msg.ConnectionId, portID); err != nil {
		return &types.MsgGovReopenChannelResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeReopenICA,
			sdk.NewAttribute(types.AttributeKeyPortID, portID),
			sdk.NewAttribute(types.AttributeKeyConnectionID, msg.ConnectionId),
		),
	})

	return &types.MsgGovReopenChannelResponse{}, nil
}

// GovCloseChannel closes an ICA channel.
func (k msgServer) GovCloseChannel(goCtx context.Context, msg *types.MsgGovCloseChannel) (*types.MsgGovCloseChannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.Keeper.GetGovAuthority(ctx) != msg.Authority {
		return &types.MsgGovCloseChannelResponse{},
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.Keeper.GetGovAuthority(ctx), msg.Authority,
			)
	}

	_, capability, err := k.Keeper.IBCKeeper.ChannelKeeper.LookupModuleByChannel(ctx, msg.PortId, msg.ChannelId)
	if err != nil {
		return &types.MsgGovCloseChannelResponse{}, err
	}

	if err := k.IBCKeeper.ChannelKeeper.ChanCloseInit(ctx, msg.PortId, msg.ChannelId, capability); err != nil {
		return &types.MsgGovCloseChannelResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeReopenICA,
			sdk.NewAttribute(types.AttributeKeyPortID, msg.PortId),
			sdk.NewAttribute(types.AttributeKeyChannelID, msg.ChannelId),
		),
	})

	return &types.MsgGovCloseChannelResponse{}, nil
}

// GovReopenChannel reopens an ICA channel.
func (k msgServer) GovReopenChannel(goCtx context.Context, msg *types.MsgGovReopenChannel) (*types.MsgGovReopenChannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// remove leanding prefix icacontroller- if passed in msg
	portID := strings.ReplaceAll(msg.PortId, "icacontroller-", "")

	// validate the zone exists, and the format is valid (e.g. quickgaia-1.delegate)
	parts := strings.Split(portID, ".")

	if len(parts) != 2 {
		return &types.MsgGovReopenChannelResponse{}, errors.New("invalid port format")
	}

	if _, found := k.GetZone(ctx, parts[0]); !found {
		return &types.MsgGovReopenChannelResponse{}, errors.New("invalid port format; zone not found")
	}

	if parts[1] != "delegate" && parts[1] != "deposit" && parts[1] != "performance" && parts[1] != "withdrawal" {
		return &types.MsgGovReopenChannelResponse{}, errors.New("invalid port format; unexpected account")
	}

	if err := k.Keeper.registerInterchainAccount(ctx, msg.ConnectionId, portID); err != nil {
		return &types.MsgGovReopenChannelResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeReopenICA,
			sdk.NewAttribute(types.AttributeKeyPortID, portID),
			sdk.NewAttribute(types.AttributeKeyConnectionID, msg.ConnectionId),
		),
	})

	return &types.MsgGovReopenChannelResponse{}, nil
}

// GovCloseChannel closes an ICA channel.
func (k msgServer) GovCloseChannel(goCtx context.Context, msg *types.MsgGovCloseChannel) (*types.MsgGovCloseChannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.Keeper.GetGovAuthority(ctx) != msg.Authority {
		return &types.MsgGovCloseChannelResponse{},
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.Keeper.GetGovAuthority(ctx), msg.Authority,
			)
	}

	_, capability, err := k.Keeper.IBCKeeper.ChannelKeeper.LookupModuleByChannel(ctx, msg.PortId, msg.ChannelId)
	if err != nil {
		return &types.MsgGovCloseChannelResponse{}, err
	}

	if err := k.IBCKeeper.ChannelKeeper.ChanCloseInit(ctx, msg.PortId, msg.ChannelId, capability); err != nil {
		return &types.MsgGovCloseChannelResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeReopenICA,
			sdk.NewAttribute(types.AttributeKeyPortID, msg.PortId),
			sdk.NewAttribute(types.AttributeKeyChannelID, msg.ChannelId),
		),
	})

	return &types.MsgGovCloseChannelResponse{}, nil
}
