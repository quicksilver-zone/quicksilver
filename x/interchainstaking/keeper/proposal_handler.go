package keeper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	ibcexported "github.com/cosmos/ibc-go/v5/modules/core/exported"
	tmclienttypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// HandleRegisterZoneProposal is a handler for executing a passed community spend proposal.
func (k *Keeper) HandleRegisterZoneProposal(ctx sdk.Context, p *types.RegisterZoneProposal) error {
	// get chain id from connection
	chainID, err := k.GetChainID(ctx, p.ConnectionId)
	if err != nil {
		return fmt.Errorf("unable to obtain chain id: %w", err)
	}

	// get zone
	_, found := k.GetZone(ctx, chainID)
	if found {
		return fmt.Errorf("invalid chain id, zone for \"%s\" already registered", chainID)
	}

	connection, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, p.ConnectionId)
	if !found {
		return errors.New("unable to fetch connection")
	}

	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return errors.New("unable to fetch client state")
	}

	tmClientState, ok := clientState.(*tmclienttypes.ClientState)
	if !ok {
		return errors.New("error unmarshaling client state")
	}

	if tmClientState.Status(ctx, k.IBCKeeper.ClientKeeper.ClientStore(ctx, connection.ClientId), k.IBCKeeper.Codec()) != ibcexported.Active {
		return errors.New("client state is not active")
	}

	zone := &types.Zone{
		ChainId:            chainID,
		ConnectionId:       p.ConnectionId,
		LocalDenom:         p.LocalDenom,
		BaseDenom:          p.BaseDenom,
		AccountPrefix:      p.AccountPrefix,
		RedemptionRate:     sdk.NewDec(1),
		LastRedemptionRate: sdk.NewDec(1),
		UnbondingEnabled:   p.UnbondingEnabled,
		ReturnToSender:     p.ReturnToSender,
		LiquidityModule:    p.LiquidityModule,
		DepositsEnabled:    p.DepositsEnabled,
		Decimals:           p.Decimals,
		UnbondingPeriod:    int64(tmClientState.UnbondingPeriod),
		MessagesPerTx:      p.MessagesPerTx,
		Is_118:             p.Is_118,
	}
	k.SetZone(ctx, zone)

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
	portOwner = chainID + ".delegate"
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
		return err
	}

	period := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))
	query := stakingTypes.QueryValidatorsRequest{}
	err = k.EmitValSetQuery(ctx, zone.ConnectionId, zone.ChainId, query, sdkmath.NewInt(period))
	if err != nil {
		return err
	}

	err = k.hooks.AfterZoneCreated(ctx, zone.ConnectionId, zone.ChainId, zone.AccountPrefix)
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
			sdk.NewAttribute(types.AttributeKeyChainID, chainID),
		),
	})

	return nil
}

func (k *Keeper) registerInterchainAccount(ctx sdk.Context, connectionID, portOwner string) error {
	if err := k.ICAControllerKeeper.RegisterInterchainAccount(ctx, connectionID, portOwner, ""); err != nil { // todo: add version
		return err
	}
	portID, err := icatypes.NewControllerPortID(portOwner)
	if err != nil {
		return err
	}

	k.SetConnectionForPort(ctx, connectionID, portID)

	return nil
}

// HandleUpdateZoneProposal is a handler for executing a passed community spend proposal.
func (k *Keeper) HandleUpdateZoneProposal(ctx sdk.Context, p *types.UpdateZoneProposal) error {
	zone, found := k.GetZone(ctx, p.ChainId)
	if !found {
		err := fmt.Errorf("unable to get registered zone for chain id: %s", p.ChainId)
		return err
	}

	for _, change := range p.Changes {
		switch change.Key {
		case "base_denom":
			if err := sdk.ValidateDenom(change.Value); err != nil {
				return err
			}
			if k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.IsPositive() {
				return errors.New("zone has assets minted, cannot update base_denom without potentially losing assets")
			}
			zone.BaseDenom = change.Value

		case "local_denom":
			if err := sdk.ValidateDenom(change.Value); err != nil {
				return err
			}
			if k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.IsPositive() {
				return errors.New("zone has assets minted, cannot update local_denom without potentially losing assets")
			}
			zone.LocalDenom = change.Value

		case "liquidity_module":
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return err
			}
			zone.LiquidityModule = boolValue

		case "unbonding_enabled":
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return err
			}
			zone.UnbondingEnabled = boolValue

		case "deposits_enabled":
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return err
			}
			zone.DepositsEnabled = boolValue

		case "return_to_sender":
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return err
			}
			zone.ReturnToSender = boolValue

		case "messages_per_tx":
			intVal, err := strconv.Atoi(change.Value)
			if err != nil {
				return err
			}
			if intVal < 1 {
				return fmt.Errorf("invalid value for messages_per_tx: %d", intVal)
			}
			zone.MessagesPerTx = int64(intVal)

		case "account_prefix":
			zone.AccountPrefix = change.Value

		case "is_118":
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return err
			}
			zone.Is_118 = boolValue

		case "connection_id":
			if !strings.HasPrefix(change.Value, "connection-") {
				return errors.New("unexpected connection format")
			}
			if zone.DepositAddress != nil || zone.DelegationAddress != nil || zone.PerformanceAddress != nil || zone.WithdrawalAddress != nil {
				return errors.New("zone already intialised, cannot update connection_id")
			}
			if k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.IsPositive() {
				return errors.New("zone has assets minted, cannot update connection_id without potentially losing assets")
			}

			connection, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, change.Value)
			if !found {
				return errors.New("unable to fetch connection")
			}

			clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
			if !found {
				return errors.New("unable to fetch client state")
			}

			tmClientState, ok := clientState.(*tmclienttypes.ClientState)
			if !ok {
				return errors.New("error unmarshaling client state")
			}

			if tmClientState.Status(ctx, k.IBCKeeper.ClientKeeper.ClientStore(ctx, connection.ClientId), k.IBCKeeper.Codec()) != ibcexported.Active {
				return errors.New("new connection client state is not active")
			}

			zone.ConnectionId = change.Value

			k.SetZone(ctx, &zone)

			// generate deposit account
			portOwner := zone.ChainId + ".deposit"
			if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
				return err
			}

			// generate withdrawal account
			portOwner = zone.ChainId + ".withdrawal"
			if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
				return err
			}

			// generate perf account
			portOwner = zone.ChainId + ".performance"
			if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
				return err
			}

			// generate delegate accounts
			portOwner = zone.ChainId + ".delegate"
			if err := k.registerInterchainAccount(ctx, zone.ConnectionId, portOwner); err != nil {
				return err
			}

			period := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))
			query := stakingTypes.QueryValidatorsRequest{}
			err := k.EmitValSetQuery(ctx, zone.ConnectionId, zone.ChainId, query, sdkmath.NewInt(period))
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unexpected key '%s'", change.Key)
		}
	}
	k.SetZone(ctx, &zone)

	k.Logger(ctx).Info("applied changes to zone", "changes", p.Changes, "zone", zone.ChainId)

	return nil
}
