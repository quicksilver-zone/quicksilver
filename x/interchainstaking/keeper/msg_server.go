package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	chainId, err := k.getChainID(ctx, msg.ConnectionId)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain chain id: %w", err)
	}

	// get zone
	_, found := k.GetRegisteredZoneInfo(ctx, chainId)
	if found {
		return nil, fmt.Errorf("invalid chain id, zone for \"%s\" already registered", chainId)
	}

	zone := types.RegisteredZone{Identifier: msg.Identifier, ChainId: chainId, ConnectionId: msg.ConnectionId, LocalDenom: msg.LocalDenom, BaseDenom: msg.BaseDenom, RedemptionRate: sdk.NewDec(1), DelegatorIntent: make(map[string]*types.DelegatorIntent), MultiSend: msg.MultiSend}
	k.SetRegisteredZone(ctx, zone)

	// generate deposit account
	portOwner := chainId + ".deposit"
	if err := k.ICAControllerKeeper.RegisterInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return nil, err
	}
	portId, _ := icatypes.NewControllerPortID(portOwner)
	if err := k.SetConnectionForPort(ctx, msg.ConnectionId, portId); err != nil {
		return nil, err
	}

	// generate delegate addresses
	for i := 0; i < types.DelegationAccountCount; i++ {
		portOwner := fmt.Sprintf("%s.delegate.%d", chainId, i)
		if err := k.ICAControllerKeeper.RegisterInterchainAccount(ctx, zone.ConnectionId, fmt.Sprintf("%s.delegate.%d", chainId, i)); err != nil {
			return nil, err
		}
		portId, _ := icatypes.NewControllerPortID(portOwner)
		if err := k.SetConnectionForPort(ctx, msg.ConnectionId, portId); err != nil {
			return nil, err
		}
	}

	bondedValidatorQuery := k.ICQKeeper.NewPeriodicQuery(ctx, msg.ConnectionId, chainId, "cosmos.staking.v1beta1.Query/Validators", map[string]string{"status": stakingtypes.BondStatusBonded}, sdk.NewInt(types.ValidatorSetInterval))
	k.ICQKeeper.SetPeriodicQuery(ctx, *bondedValidatorQuery)
	unbondedValidatorQuery := k.ICQKeeper.NewPeriodicQuery(ctx, msg.ConnectionId, chainId, "cosmos.staking.v1beta1.Query/Validators", map[string]string{"status": stakingtypes.BondStatusUnbonded}, sdk.NewInt(types.ValidatorSetInterval))
	k.ICQKeeper.SetPeriodicQuery(ctx, *unbondedValidatorQuery)
	unbondingValidatorQuery := k.ICQKeeper.NewPeriodicQuery(ctx, msg.ConnectionId, chainId, "cosmos.staking.v1beta1.Query/Validators", map[string]string{"status": stakingtypes.BondStatusUnbonding}, sdk.NewInt(types.ValidatorSetInterval))
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

func (k msgServer) RequestRedemption(goCtx context.Context, msg *types.MsgRequestRedemption) (*types.MsgRequestRedemptionResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

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
