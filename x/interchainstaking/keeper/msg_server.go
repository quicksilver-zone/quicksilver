package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	tmclienttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
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
	var (
		baseZone types.Zone
		found    bool
	)

	ctx := sdk.UnwrapSDKContext(goCtx)
	// checking msg authority is the gov module address
	if k.Keeper.GetGovAuthority() != msg.Authority {
		return &types.MsgRegisterZoneResponse{},
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.Keeper.GetGovAuthority(), msg.Authority,
			)
	}

	// get chain id from connection
	chainID, err := k.GetChainID(ctx, msg.ConnectionID)
	if err != nil {
		return &types.MsgRegisterZoneResponse{}, fmt.Errorf("unable to obtain chain id: %w", err)
	}

	// if subzone
	if msg.SubzoneInfo != nil {
		if chainID != msg.SubzoneInfo.BaseChainID {
			return &types.MsgRegisterZoneResponse{}, fmt.Errorf("incorrect ID \"%s\" for subzone \"%s\"", chainID, msg.SubzoneInfo.BaseChainID)
		}

		// get zone
		baseZone, found = k.GetZone(ctx, msg.SubzoneInfo.BaseChainID)
		if !found {
			return &types.MsgRegisterZoneResponse{}, fmt.Errorf("unable to find base chain \"%s\" for subzone \"%s\"", chainID, msg.SubzoneInfo.BaseChainID)
		}

		// check if subzone ID already is taken
		_, found = k.GetZone(ctx, msg.SubzoneInfo.ChainID)
		if found {
			return &types.MsgRegisterZoneResponse{}, fmt.Errorf("subzone ID already exists \"%s\"", msg.SubzoneInfo.ChainID)
		}

		// set chainID to be specified unique ID
		chainID = msg.SubzoneInfo.ChainID
	}

	// get zone
	_, found = k.GetZone(ctx, chainID)
	if found {
		return &types.MsgRegisterZoneResponse{}, fmt.Errorf("invalid chain id, zone for \"%s\" already registered", chainID)
	}

	connection, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, msg.ConnectionID)
	if !found {
		return &types.MsgRegisterZoneResponse{}, errors.New("unable to fetch connection")
	}

	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return &types.MsgRegisterZoneResponse{}, errors.New("unable to fetch client state")
	}

	tmClientState, ok := clientState.(*tmclienttypes.ClientState)
	if !ok {
		return &types.MsgRegisterZoneResponse{}, errors.New("error unmarshaling client state")
	}

	if tmClientState.Status(ctx, k.IBCKeeper.ClientKeeper.ClientStore(ctx, connection.ClientId), k.IBCKeeper.Codec()) != ibcexported.Active {
		return &types.MsgRegisterZoneResponse{}, errors.New("client state is not active")
	}

	zone := &types.Zone{
		ChainId:            chainID,
		ConnectionId:       msg.ConnectionID,
		LocalDenom:         msg.LocalDenom,
		BaseDenom:          msg.BaseDenom,
		AccountPrefix:      msg.AccountPrefix,
		RedemptionRate:     sdk.NewDec(1),
		LastRedemptionRate: sdk.NewDec(1),
		UnbondingEnabled:   msg.UnbondingEnabled,
		ReturnToSender:     msg.ReturnToSender,
		LiquidityModule:    msg.LiquidityModule,
		DepositsEnabled:    msg.DepositsEnabled,
		Decimals:           msg.Decimals,
		UnbondingPeriod:    int64(tmClientState.UnbondingPeriod),
		MessagesPerTx:      msg.MessagesPerTx,
		Is_118:             msg.Is_118,
		SubzoneInfo:        msg.SubzoneInfo,
	}

	// verify subzone if setting
	if zone.IsSubzone() {
		if err := types.ValidateSubzoneForBasezone(*zone, baseZone); err != nil {
			return &types.MsgRegisterZoneResponse{}, err
		}
	}

	k.SetZone(ctx, zone)

	// generate deposit account
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, zone.DepositPortOwner()); err != nil {
		return &types.MsgRegisterZoneResponse{}, err
	}

	// generate the withdrawal account
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, zone.WithdrawalPortOwner()); err != nil {
		return nil, err
	}

	// generate the perf account
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, zone.PerformancePortOwner()); err != nil {
		return &types.MsgRegisterZoneResponse{}, err
	}

	// generate delegate accounts
	if err := k.registerInterchainAccount(ctx, zone.ConnectionId, zone.DelegatePortOwner()); err != nil {
		return &types.MsgRegisterZoneResponse{}, err
	}

	// query val set for base zone
	if !zone.IsSubzone() {
		period := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))
		query := stakingTypes.QueryValidatorsRequest{}
		err = k.EmitValSetQuery(ctx, zone.ConnectionId, zone, query, sdkmath.NewInt(period))
		if err != nil {
			return &types.MsgRegisterZoneResponse{}, err
		}
	}

	err = k.hooks.AfterZoneCreated(ctx, zone.ConnectionId, zone.ChainId, zone.AccountPrefix)
	if err != nil {
		return &types.MsgRegisterZoneResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRegisterZone,
			sdk.NewAttribute(types.AttributeKeyConnectionID, msg.ConnectionID),
			sdk.NewAttribute(types.AttributeKeyChainID, chainID),
		),
	})

	return &types.MsgRegisterZoneResponse{
		ZoneID: chainID,
	}, nil
}

func (k msgServer) UpdateZone(goCtx context.Context, msg *types.MsgUpdateZone) (*types.MsgUpdateZoneResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// checking msg authority is the gov module address
	if k.Keeper.GetGovAuthority() != msg.Authority {
		return &types.MsgUpdateZoneResponse{},
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.Keeper.GetGovAuthority(), msg.Authority,
			)
	}

	zone, found := k.GetZone(ctx, msg.ZoneID)
	if !found {
		return &types.MsgUpdateZoneResponse{}, fmt.Errorf("unable to get registered zone for zone id: %s", msg.ZoneID)
	}

	for _, change := range msg.Changes {
		switch change.Key {
		case types.UpdateZoneKeyBaseDenom:
			if err := sdk.ValidateDenom(change.Value); err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}
			if k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.IsPositive() {
				return &types.MsgUpdateZoneResponse{}, errors.New("zone has assets minted, cannot update base_denom without potentially losing assets")
			}
			zone.BaseDenom = change.Value

		case types.UpdateZoneKeyLocalDenom:
			if err := sdk.ValidateDenom(change.Value); err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}
			if k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.IsPositive() {
				return &types.MsgUpdateZoneResponse{}, errors.New("zone has assets minted, cannot update local_denom without potentially losing assets")
			}
			zone.LocalDenom = change.Value

		case types.UpdateZoneKeyLiquidityModule:
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}
			zone.LiquidityModule = boolValue

		case types.UpdateZoneKeyUnbondingEnabled:
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}
			zone.UnbondingEnabled = boolValue

		case types.UpdateZoneKeyDepositsEnabled:
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}
			zone.DepositsEnabled = boolValue

		case types.UpdateZoneKeyReturnToSender:
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}
			zone.ReturnToSender = boolValue

		case types.UpdateZoneKeyMessagesPerTx:
			intVal, err := strconv.Atoi(change.Value)
			if err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}
			if intVal < 1 {
				return &types.MsgUpdateZoneResponse{}, fmt.Errorf("invalid value for messages_per_tx: %d", intVal)
			}
			zone.MessagesPerTx = int64(intVal)

		case types.UpdateZoneKeyAccountPrefix:
			zone.AccountPrefix = change.Value

		case types.UpdateZoneKeyIs118:
			boolValue, err := strconv.ParseBool(change.Value)
			if err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}
			zone.Is_118 = boolValue

		case types.UpdateZoneKeyConnectionID:
			if !strings.HasPrefix(change.Value, types.ConnectionPrefix) {
				return &types.MsgUpdateZoneResponse{}, errors.New("unexpected connection format")
			}
			if zone.DepositAddress != nil || zone.DelegationAddress != nil || zone.PerformanceAddress != nil || zone.WithdrawalAddress != nil {
				return &types.MsgUpdateZoneResponse{}, errors.New("zone already intialised, cannot update connection_id")
			}
			if k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.IsPositive() {
				return &types.MsgUpdateZoneResponse{}, errors.New("zone has assets minted, cannot update connection_id without potentially losing assets")
			}

			connection, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, change.Value)
			if !found {
				return &types.MsgUpdateZoneResponse{}, errors.New("unable to fetch connection")
			}

			clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
			if !found {
				return &types.MsgUpdateZoneResponse{}, errors.New("unable to fetch client state")
			}

			tmClientState, ok := clientState.(*tmclienttypes.ClientState)
			if !ok {
				return &types.MsgUpdateZoneResponse{}, errors.New("error unmarshaling client state")
			}

			if tmClientState.Status(ctx, k.IBCKeeper.ClientKeeper.ClientStore(ctx, connection.ClientId), k.IBCKeeper.Codec()) != ibcexported.Active {
				return &types.MsgUpdateZoneResponse{}, errors.New("new connection client state is not active")
			}

			zone.ConnectionId = change.Value

			k.SetZone(ctx, &zone)

			// generate deposit account
			if err := k.registerInterchainAccount(ctx, zone.ConnectionId, zone.DepositPortOwner()); err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}

			// generate withdrawal account
			if err := k.registerInterchainAccount(ctx, zone.ConnectionId, zone.WithdrawalPortOwner()); err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}

			// generate perf account
			if err := k.registerInterchainAccount(ctx, zone.ConnectionId, zone.PerformancePortOwner()); err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}

			// generate delegate accounts
			if err := k.registerInterchainAccount(ctx, zone.ConnectionId, zone.DelegatePortOwner()); err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}

			period := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))
			query := stakingTypes.QueryValidatorsRequest{}
			err := k.EmitValSetQuery(ctx, zone.ConnectionId, &zone, query, sdkmath.NewInt(period))
			if err != nil {
				return &types.MsgUpdateZoneResponse{}, err
			}

		default:
			return &types.MsgUpdateZoneResponse{}, fmt.Errorf("unexpected key '%s'", change.Key)
		}
	}
	k.SetZone(ctx, &zone)

	k.Logger(ctx).Info("applied changes to zone", "changes", msg.Changes, "zone", zone.ZoneID())

	return &types.MsgUpdateZoneResponse{}, nil
}

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
		return nil, fmt.Errorf("unbonding currently disabled for zone %s", zone.ZoneID())
	}

	// TODO check sub-zone
	if zone.IsSubzone() && (msg.FromAddress != zone.SubzoneInfo.Authority) {
		return nil, types.ErrInvalidSubzoneAuthority
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
			sdk.NewAttribute(types.AttributeKeyChainID, zone.ZoneID()),
			sdk.NewAttribute(types.AttributeKeyConnectionID, zone.ConnectionId),
		),
	})

	return &types.MsgRequestRedemptionResponse{}, nil
}

func (k msgServer) SignalIntent(goCtx context.Context, msg *types.MsgSignalIntent) (*types.MsgSignalIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get zone
	zone, ok := k.GetZone(ctx, msg.ChainId)
	if !ok {
		return nil, fmt.Errorf("invalid chain id \"%s\"", msg.ChainId)
	}

	// TODO check sub-zone
	if zone.IsSubzone() && (msg.FromAddress != zone.SubzoneInfo.Authority) {
		return nil, types.ErrInvalidSubzoneAuthority
	}

	// validate intents (aggregated errors)
	intents, err := types.IntentsFromString(msg.Intents)
	if err != nil {
		return nil, err
	}

	if err := k.validateValidatorIntents(ctx, &zone, intents); err != nil {
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

	// TODO handle for subzone?

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
	if k.Keeper.GetGovAuthority() != msg.Authority {
		return &types.MsgGovCloseChannelResponse{},
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.Keeper.GetGovAuthority(), msg.Authority,
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
