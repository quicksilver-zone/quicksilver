package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// HandleRegisterZoneProposal is a handler for executing a passed community spend proposal
func HandleRegisterZoneProposal(ctx sdk.Context, k Keeper, p *types.RegisterZoneProposal) error {
	// get chain id from connection
	chainID, err := k.GetChainID(ctx, p.ConnectionId)
	if err != nil {
		return fmt.Errorf("unable to obtain chain id: %w", err)
	}

	// get zone
	_, found := k.GetRegisteredZoneInfo(ctx, chainID)
	if found {
		return fmt.Errorf("invalid chain id, zone for \"%s\" already registered", chainID)
	}

	zone := types.RegisteredZone{
		ChainId:            chainID,
		ConnectionId:       p.ConnectionId,
		LocalDenom:         p.LocalDenom,
		BaseDenom:          p.BaseDenom,
		AccountPrefix:      p.AccountPrefix,
		RedemptionRate:     sdk.NewDec(1),
		LastRedemptionRate: sdk.NewDec(1),
		MultiSend:          p.MultiSend,
		LiquidityModule:    p.LiquidityModule,
	}
	k.SetRegisteredZone(ctx, zone)

	// generate deposit account
	portOwner := chainID + ".deposit"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return err
	}

	// generate withdrawal account
	portOwner = chainID + ".withdrawal"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return err
	}

	// generate perf account
	portOwner = chainID + ".performance"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return err
	}

	// generate delegate accounts
	delegateAccountCount := int(k.GetParam(ctx, types.KeyDelegateAccountCount))
	for i := 0; i < delegateAccountCount; i++ {
		portOwner := fmt.Sprintf("%s.delegate.%d", chainID, i)
		if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
			return err
		}
	}
	err = k.EmitValsetRequery(ctx, p.ConnectionId, chainID)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRegisterZone,
			sdk.NewAttribute(types.AttributeKeyConnectionID, p.ConnectionId),
			sdk.NewAttribute(types.AttributeKeyConnectionID, chainID),
		),
	})

	return nil
}

func (k Keeper) registerInterchainAccount(ctx sdk.Context, connectionID string, portOwner string) error {
	if err := k.ICAControllerKeeper.RegisterInterchainAccount(ctx, connectionID, portOwner); err != nil {
		return err
	}
	portID, _ := icatypes.NewControllerPortID(portOwner)
	if err := k.SetConnectionForPort(ctx, connectionID, portID); err != nil {
		return err
	}

	return nil
}

// HandleUpdateZoneProposal is a handler for executing a passed community spend proposal
func HandleUpdateZoneProposal(ctx sdk.Context, k Keeper, p *types.UpdateZoneProposal) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, p.ChainId)
	if !found {
		fmt.Errorf("Unable to get registered zone for chain id: %s", p.ChainId)
	}

	for _, change := range p.Changes {
		switch change.Key {
		case "base_denom":
			if err := sdk.ValidateDenom(change.Value); err != nil {
				return err
			}
			zone.BaseDenom = change.Value
			k.SetRegisteredZone(ctx, zone)
		}
	}

	logger := k.Logger(ctx)
	logger.Info("applied changes to zone", "changes", p.Changes, "zone", zone.ChainId)

	return nil
}
