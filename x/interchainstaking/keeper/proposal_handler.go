package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// HandleRegisterZoneProposal is a handler for executing a passed community spend proposal.
func (k *Keeper) HandleRegisterZoneProposal(ctx sdk.Context, p *types.RegisterZoneProposal) error {
	return errors.New("legacy zone registration proposals are currently disabled")
}

func (k *Keeper) RegisterInterchainAccount(ctx sdk.Context, connectionID, portOwner string) error {
	msg := icacontrollertypes.NewMsgRegisterInterchainAccountWithOrdering(connectionID, portOwner, "", channeltypes.ORDERED)
	ckMsgServer := icacontrollerkeeper.NewMsgServerImpl(&k.ICAControllerKeeper)
	_, err := ckMsgServer.RegisterInterchainAccount(ctx, msg)
	if err != nil {
		return err
	}

	portID, err := icatypes.NewControllerPortID(portOwner)
	if err != nil {
		return err
	}

	k.SetConnectionForPort(ctx, connectionID, portID)
	k.ICAControllerKeeper.SetMiddlewareEnabled(ctx, portID, connectionID)

	return nil
}

// HandleUpdateZoneProposal is a handler for executing a passed community spend proposal.
func (k *Keeper) HandleUpdateZoneProposal(ctx sdk.Context, p *types.UpdateZoneProposal) error {
	return errors.New("legacy zone update proposals are currently disabled")
}
