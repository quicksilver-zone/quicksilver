package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the interchainstaking
// MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RequestRedemption handles MsgRequestRedemption by creating a corresponding withdrawal record queued for unbonding.
func (k msgServer) RequestRedemption(goCtx context.Context, msg *types.MsgRequestRedemption) (*types.MsgRequestRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.Keeper.GetUnbondingEnabled(ctx) {
		return nil, fmt.Errorf("unbonding is currently disabled")
	}

	zone := k.Keeper.GetZoneByLocalDenom(ctx, msg.Value.Denom)

	// does zone exist?
	if zone == nil {
		return nil, fmt.Errorf("unable to find matching zone for denom %s", msg.Value.GetDenom())
	}

	if !zone.UnbondingEnabled {
		return nil, fmt.Errorf("unbonding currently disabled for zone %s", zone.ChainId)
	}

	// does destination address match the prefix registered against the zone?
	if _, err := addressutils.AccAddressFromBech32(msg.DestinationAddress, zone.AccountPrefix); err != nil {
		return nil, fmt.Errorf("destination address %s does not match expected prefix %s [%w]", msg.DestinationAddress, zone.AccountPrefix, err)
	}

	sender, _ := sdk.AccAddressFromBech32(msg.FromAddress) // already validated

	// does the user have sufficient assets to burn
	if !k.BankKeeper.HasBalance(ctx, sender, msg.Value) {
		return nil, errors.New("account has insufficient balance of qasset to burn")
	}

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(ctx.BlockHeight()))
	hash := sha256.Sum256(append(msg.GetSignBytes(), heightBytes...))
	hashString := hex.EncodeToString(hash[:])

	if err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.EscrowModuleAccount, sdk.NewCoins(msg.Value)); err != nil {
		return nil, fmt.Errorf("unable to send coins to escrow account: %w", err)
	}

	if err := k.queueRedemption(ctx, zone, sender, msg.DestinationAddress, msg.Value, hashString); err != nil {
		return nil, fmt.Errorf("unable to queue redemption: %w", err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRedemptionRequest,
			sdk.NewAttribute(types.AttributeKeyBurnAmount, msg.Value.String()),
			sdk.NewAttribute(types.AttributeKeyRecipientAddress, msg.DestinationAddress),
			sdk.NewAttribute(types.AttributeKeyChainID, zone.ChainId),
		),
	})

	return &types.MsgRequestRedemptionResponse{}, nil
}

func (k msgServer) CancelRedemption(goCtx context.Context, msg *types.MsgCancelQueuedRedemption) (*types.MsgCancelQueuedRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	record, found := k.GetWithdrawalRecord(ctx, msg.ChainId, msg.Hash, types.WithdrawStatusQueued)

	if !found {
		return nil, fmt.Errorf("no queued record with hash \"%s\" found", msg.Hash)
	}

	if record.Delegator != msg.FromAddress {
		return nil, fmt.Errorf("incorrect user for record with hash \"%s\"", msg.Hash)
	}

	// all good. delete!
	k.DeleteWithdrawalRecord(ctx, msg.ChainId, msg.Hash, types.WithdrawStatusQueued)

	userAccAddress, err := addressutils.AddressFromBech32(record.Delegator, "")
	if err != nil {
		return nil, err
	}

	// return coins
	if err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.EscrowModuleAccount, userAccAddress, sdk.NewCoins(record.BurnAmount)); err != nil {
		return nil, fmt.Errorf("unable to return coins from escrow account: %w", err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRedemptionCancellation,
			sdk.NewAttribute(types.AttributeKeyReturnedAmount, record.BurnAmount.String()),
			sdk.NewAttribute(types.AttributeKeyUser, msg.FromAddress),
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
		),
	})

	return &types.MsgCancelQueuedRedemptionResponse{Returned: record.BurnAmount}, nil
}

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
		return nil, errors.New("invalid port format; zone not found")
	}

	if err := k.Keeper.registerInterchainAccount(ctx, msg.ConnectionId, portID); err != nil {
		return nil, err
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
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.Keeper.GetGovAuthority(ctx), msg.Authority,
			)
	}

	_, capability, err := k.Keeper.IBCKeeper.ChannelKeeper.LookupModuleByChannel(ctx, msg.PortId, msg.ChannelId)
	if err != nil {
		return nil, err
	}

	if err := k.IBCKeeper.ChannelKeeper.ChanCloseInit(ctx, msg.PortId, msg.ChannelId, capability); err != nil {
		return nil, err
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

// GovSetLsmCaps set the liquid staking caps for a given chain.
func (k msgServer) GovSetLsmCaps(goCtx context.Context, msg *types.MsgGovSetLsmCaps) (*types.MsgGovSetLsmCapsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.Keeper.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.Keeper.GetGovAuthority(ctx), msg.Authority,
			)
	}

	zone, found := k.Keeper.GetZone(ctx, msg.ChainId)
	if !found {
		return nil,
			fmt.Errorf(
				"no zone found for: %s",
				msg.ChainId,
			)
	}
	if !zone.SupportLsm() {
		return nil,
			fmt.Errorf(
				"zone %s does not have LSM support enabled",
				msg.ChainId,
			)
	}

	k.SetLsmCaps(ctx, zone.ChainId, *msg.Caps)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeSetLsmCaps,
			sdk.NewAttribute(types.AttributeLsmValidatorCap, msg.Caps.ValidatorCap.String()),
			sdk.NewAttribute(types.AttributeLsmValidatorBondCap, msg.Caps.ValidatorBondCap.String()),
			sdk.NewAttribute(types.AttributeLsmGlobalCap, msg.Caps.GlobalCap.String()),
		),
	})

	return &types.MsgGovSetLsmCapsResponse{}, nil
}
