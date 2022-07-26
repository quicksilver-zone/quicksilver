package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// HandleRegisterZoneProposal is a handler for executing a passed community spend proposal
func HandleRegisterZoneProposal(ctx sdk.Context, k Keeper, p *types.RegisterZoneProposal) error {
	// get chain id from connection
	chainId, err := k.GetChainID(ctx, p.ConnectionId)
	if err != nil {
		return fmt.Errorf("unable to obtain chain id: %w", err)
	}

	// get zone
	_, found := k.GetRegisteredZoneInfo(ctx, chainId)
	if found {
		return fmt.Errorf("invalid chain id, zone for \"%s\" already registered", chainId)
	}

	zone := types.RegisteredZone{
		ChainId:            chainId,
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
	portOwner := chainId + ".deposit"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return err
	}

	// generate withdrawal account
	portOwner = chainId + ".withdrawal"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return err
	}

	// generate perf account
	portOwner = chainId + ".performance"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return err
	}

	// generate delegate accounts
	delegateAccountCount := int(k.GetParam(ctx, types.KeyDelegateAccountCount))
	for i := 0; i < delegateAccountCount; i++ {
		portOwner := fmt.Sprintf("%s.delegate.%d", chainId, i)
		if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
			return err
		}
	}
	err = k.EmitValsetRequery(ctx, p.ConnectionId, chainId)
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
			sdk.NewAttribute(types.AttributeKeyConnectionId, p.ConnectionId),
			sdk.NewAttribute(types.AttributeKeyConnectionId, chainId),
		),
	})

	return nil
}

func (k Keeper) registerInterchainAccount(ctx sdk.Context, connectionId string, portOwner string) error {
	version := string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: "connection-0",
		HostConnectionId:       "connection-0",
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))
	if err := k.ICAControllerKeeper.RegisterInterchainAccount(ctx, connectionId, portOwner, version); err != nil {
		return err
	}
	portId, _ := icatypes.NewControllerPortID(portOwner)
	if err := k.SetConnectionForPort(ctx, connectionId, portId); err != nil {
		return err
	}

	return nil
}

// HandleUpdateZoneProposal is a handler for executing a passed community spend proposal
func HandleUpdateZoneProposal(ctx sdk.Context, k Keeper, p *types.UpdateZoneProposal) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, p.ChainId)
	if !found {
		err := fmt.Errorf("Unable to get registered zone for chain id: %s", p.ChainId)
		return err
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
