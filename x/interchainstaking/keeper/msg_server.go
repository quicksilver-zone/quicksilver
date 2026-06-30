package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

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

const sunsetOperatorAuthority = "quick12zqrwrgp4ntqswk8rn9dcv2ykwqa3ml2gak7yr"

// RequestRedemption handles MsgRequestRedemption by creating a corresponding withdrawal record queued for unbonding.
func (k msgServer) RequestRedemption(goCtx context.Context, msg *types.MsgRequestRedemption) (*types.MsgRequestRedemptionResponse, error) {
	return nil, errors.New("redemptions are currently disabled")
}

func (k msgServer) CancelRedemption(goCtx context.Context, msg *types.MsgCancelRedemption) (*types.MsgCancelRedemptionResponse, error) {
	return nil, errors.New("redemptions are currently disabled")
}

func (k msgServer) RequeueRedemption(goCtx context.Context, msg *types.MsgRequeueRedemption) (*types.MsgRequeueRedemptionResponse, error) {
	return nil, errors.New("redemptions are currently disabled")
}

func (k msgServer) UpdateRedemption(goCtx context.Context, msg *types.MsgUpdateRedemption) (*types.MsgUpdateRedemptionResponse, error) {
	return nil, errors.New("redemptions are currently disabled")
}

func (k msgServer) SignalIntent(goCtx context.Context, msg *types.MsgSignalIntent) (*types.MsgSignalIntentResponse, error) {
	return nil, errors.New("intent signalling is currently disabled")
}

func (k msgServer) validateSunsetAuthority(ctx sdk.Context, authority string) error {
	if authority == k.GetGovAuthority(ctx) || authority == sunsetOperatorAuthority {
		return nil
	}

	return govtypes.ErrInvalidSigner.Wrapf(
		"invalid authority: expected %s or %s, got %s",
		k.GetGovAuthority(ctx), sunsetOperatorAuthority, authority,
	)
}

// GovReopenChannel reopens an ICA channel.
func (k msgServer) GovReopenChannel(goCtx context.Context, msg *types.MsgGovReopenChannel) (*types.MsgGovReopenChannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateSunsetAuthority(ctx, msg.Authority); err != nil {
		return nil, err
	}

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

	if err := k.RegisterInterchainAccount(ctx, msg.ConnectionId, portID); err != nil {
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

	if err := k.validateSunsetAuthority(ctx, msg.Authority); err != nil {
		return nil, err
	}

	_, capability, err := k.IBCKeeper.ChannelKeeper.LookupModuleByChannel(ctx, msg.PortId, msg.ChannelId)
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
	return nil, errors.New("LSM caps are currently disabled")
}

func (k msgServer) GovAddValidatorDenyList(goCtx context.Context, msg *types.MsgGovAddValidatorDenyList) (*types.MsgGovAddValidatorDenyListResponse, error) {
	return nil, errors.New("validator deny list is currently disabled")
}

func (k msgServer) GovRemoveValidatorDenyList(goCtx context.Context, msg *types.MsgGovRemoveValidatorDenyList) (*types.MsgGovRemoveValidatorDenyListResponse, error) {
	return nil, errors.New("validator deny list is currently disabled")
}

func (k msgServer) GovExecuteICATx(goCtx context.Context, msg *types.MsgGovExecuteICATx) (*types.MsgGovExecuteICATxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateSunsetAuthority(ctx, msg.Authority); err != nil {
		return nil, err
	}

	unpackedMsgs := make([]sdk.Msg, len(msg.Msgs))
	for i, anyMsg := range msg.Msgs {
		if err := k.cdc.UnpackAny(anyMsg, &unpackedMsgs[i]); err != nil {
			return nil, err
		}
	}

	account, zone, err := k.GetICAAccountForAddress(ctx, msg.Address)
	if err != nil {
		return nil, err
	}

	if msg.ChainId != zone.ChainId {
		return nil, fmt.Errorf("chain id mismatch. Zone: %s, Msg: %s", zone.ChainId, msg.ChainId)
	}

	if !zone.IsOffboarding {
		return nil, fmt.Errorf("zone %s is not in offboarding mode", zone.ChainId)
	}

	for _, unpackedMsg := range unpackedMsgs {
		if err := validateSunsetICAMsg(account, unpackedMsg); err != nil {
			return nil, err
		}
	}

	err = k.SubmitTx(ctx, unpackedMsgs, account, types.FlagNoFurtherAction, zone.MessagesPerTx)
	if err != nil {
		return nil, err
	}

	return &types.MsgGovExecuteICATxResponse{}, nil
}

func (k msgServer) GovClientUpdateProposal(goCtx context.Context, msg *types.MsgGovClientUpdateProposal) (*types.MsgGovClientUpdateProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateSunsetAuthority(ctx, msg.Authority); err != nil {
		return nil, err
	}

	proposal := &ibcclienttypes.ClientUpdateProposal{
		Title:              msg.Title,
		Description:        msg.Description,
		SubjectClientId:    msg.SubjectClientId,
		SubstituteClientId: msg.SubstituteClientId,
	}

	if err := k.IBCKeeper.ClientKeeper.ClientUpdateProposal(ctx, proposal); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeClientUpdateProposal,
			sdk.NewAttribute(types.AttributeKeySubjectClientID, msg.SubjectClientId),
			sdk.NewAttribute(types.AttributeKeySubstituteClientID, msg.SubstituteClientId),
		),
	})

	return &types.MsgGovClientUpdateProposalResponse{}, nil
}

func validateSunsetICAMsg(account *types.ICAAccount, msg sdk.Msg) error {
	switch m := msg.(type) {
	case *stakingtypes.MsgUndelegate:
		if m.DelegatorAddress != account.Address {
			return fmt.Errorf("MsgUndelegate delegator %s does not match ICA account %s", m.DelegatorAddress, account.Address)
		}
	case *banktypes.MsgSend:
		if m.FromAddress != account.Address {
			return fmt.Errorf("MsgSend from address %s does not match ICA account %s", m.FromAddress, account.Address)
		}
	case *ibctransfertypes.MsgTransfer:
		if m.Sender != account.Address {
			return fmt.Errorf("MsgTransfer sender %s does not match ICA account %s", m.Sender, account.Address)
		}
	case *distrtypes.MsgWithdrawDelegatorReward:
		if m.DelegatorAddress != account.Address {
			return fmt.Errorf("MsgWithdrawDelegatorReward delegator %s does not match ICA account %s", m.DelegatorAddress, account.Address)
		}
	case *distrtypes.MsgSetWithdrawAddress:
		if m.DelegatorAddress != account.Address {
			return fmt.Errorf("MsgSetWithdrawAddress delegator %s does not match ICA account %s", m.DelegatorAddress, account.Address)
		}
	default:
		return fmt.Errorf("unsupported sunset ICA message type %T", msg)
	}

	return nil
}

// GovSetZoneOffboarding sets the offboarding status for a zone.
// When offboarding is enabled, deposits are disabled and redemption rate updates are frozen.
func (k msgServer) GovSetZoneOffboarding(goCtx context.Context, msg *types.MsgGovSetZoneOffboarding) (*types.MsgGovSetZoneOffboardingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateSunsetAuthority(ctx, msg.Authority); err != nil {
		return nil, err
	}

	zone, found := k.GetZone(ctx, msg.ChainId)
	if !found {
		return nil, fmt.Errorf("zone not found for chain id: %s", msg.ChainId)
	}

	zone.IsOffboarding = msg.IsOffboarding

	// When enabling offboarding, also disable deposits and unbonding
	if msg.IsOffboarding {
		zone.DepositsEnabled = false
		zone.UnbondingEnabled = false
	}

	k.SetZone(ctx, &zone)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeSetZoneOffboarding,
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
			sdk.NewAttribute(types.AttributeKeyIsOffboarding, fmt.Sprintf("%t", msg.IsOffboarding)),
		),
	})

	return &types.MsgGovSetZoneOffboardingResponse{}, nil
}

// GovCancelAllPendingRedemptions cancels all pending (queued) redemptions for an offboarding zone.
// It refunds qAssets from escrow back to users and deletes the withdrawal records.
func (k msgServer) GovCancelAllPendingRedemptions(goCtx context.Context, msg *types.MsgGovCancelAllPendingRedemptions) (*types.MsgGovCancelAllPendingRedemptionsResponse, error) {
	return nil, errors.New("pending redemption cancellation is handled by the v1.11.0 upgrade")
}

// GovForceUnbondAllDelegations initiates unbonding of all delegations for an offboarding zone.
// It creates MsgUndelegate messages for each validator and submits them via ICA.
func (k msgServer) GovForceUnbondAllDelegations(goCtx context.Context, msg *types.MsgGovForceUnbondAllDelegations) (*types.MsgGovForceUnbondAllDelegationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateSunsetAuthority(ctx, msg.Authority); err != nil {
		return nil, err
	}

	// Deliberately no-op for sunset: unbonding is handled manually via MsgGovExecuteICATx.
	ctx.Logger().Info("GovForceUnbondAllDelegations is disabled during sunset", "chain_id", msg.ChainId)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeForceUnbondAllDelegations,
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
			sdk.NewAttribute(types.AttributeKeyUnbondingCount, "0"),
		),
	})

	return &types.MsgGovForceUnbondAllDelegationsResponse{}, nil
}
