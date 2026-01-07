package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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

	if !k.GetUnbondingEnabled(ctx) {
		return nil, errors.New("unbonding is currently disabled")
	}

	zone := k.GetZoneByLocalDenom(ctx, msg.Value.Denom)

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
	blockHeight := ctx.BlockHeight()
	if blockHeight < 0 {
		return nil, fmt.Errorf("block height is negative: %d", blockHeight)
	}

	binary.BigEndian.PutUint64(heightBytes, uint64(blockHeight))
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

func (k msgServer) CancelRedemption(goCtx context.Context, msg *types.MsgCancelRedemption) (*types.MsgCancelRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	record, found := k.GetWithdrawalRecord(ctx, msg.ChainId, msg.Hash, types.WithdrawStatusQueued)
	// QUEUED records can be cancelled at any time.
	if !found {
		// check for errored unbond in UNBONDING status
		record, found = k.GetWithdrawalRecord(ctx, msg.ChainId, msg.Hash, types.WithdrawStatusUnbond)
		if !found {
			return nil, fmt.Errorf("no queued record with hash %q found", msg.Hash)
		}
		if record.SendErrors == 0 {
			return nil, fmt.Errorf("cannot cancel unbond %q with no errors", msg.Hash)
		}
	}

	if record.Delegator != msg.FromAddress && k.GetGovAuthority(ctx) != msg.FromAddress {
		return nil, fmt.Errorf("incorrect user for record with hash %q", msg.Hash)
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
			sdk.NewAttribute(types.AttributeKeyUser, record.Delegator),
			sdk.NewAttribute(types.AttributeKeyHash, msg.Hash),
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
		),
	})

	return &types.MsgCancelRedemptionResponse{Returned: record.BurnAmount}, nil
}

func (k msgServer) RequeueRedemption(goCtx context.Context, msg *types.MsgRequeueRedemption) (*types.MsgRequeueRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check for errored unbond in UNBONDING status
	record, found := k.GetWithdrawalRecord(ctx, msg.ChainId, msg.Hash, types.WithdrawStatusUnbond)
	if !found {
		return nil, fmt.Errorf("no unbonding record with hash %q found", msg.Hash)
	}
	if record.SendErrors == 0 {
		return nil, fmt.Errorf("cannot requeue unbond %q with no errors", msg.Hash)
	}

	if record.Delegator != msg.FromAddress && k.GetGovAuthority(ctx) != msg.FromAddress {
		return nil, fmt.Errorf("incorrect user for record with hash %q", msg.Hash)
	}

	// all good. update sendErrors to zero, nil the distributions and amount (as this we be recalculated when processed), and update the state to queued.
	record.SendErrors = 0
	record.Amount = nil
	record.Distribution = nil
	record.CompletionTime = time.Time{}
	k.UpdateWithdrawalRecordStatus(ctx, &record, types.WithdrawStatusQueued)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRedemptionRequeue,
			sdk.NewAttribute(types.AttributeKeyUser, record.Delegator),
			sdk.NewAttribute(types.AttributeKeyHash, msg.Hash),
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
		),
	})

	return &types.MsgRequeueRedemptionResponse{}, nil
}

func (k msgServer) UpdateRedemption(goCtx context.Context, msg *types.MsgUpdateRedemption) (*types.MsgUpdateRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetGovAuthority(ctx) != msg.FromAddress {
		return nil, errors.New("MsgUpdateRedemption may only be executed by the gov authority")
	}

	switch msg.NewStatus {
	case types.WithdrawStatusTokenize: // intentionally removed as not currently supported, but included here for completeness.
		return nil, errors.New("new status WithdrawStatusTokenize not supported")
	case types.WithdrawStatusQueued:
	case types.WithdrawStatusUnbond:
	case types.WithdrawStatusSend: // send is not a valid state for recovery, included here for completeness.
		return nil, errors.New("new status WithdrawStatusSend not supported")
	case types.WithdrawStatusCompleted:
	default:
		return nil, errors.New("new status not provided or invalid")
	}

	var r *types.WithdrawalRecord

	k.IteratePrefixedWithdrawalRecords(ctx, []byte(msg.ChainId), func(index int64, record types.WithdrawalRecord) (stop bool) {
		if record.Txhash == msg.Hash {
			r = &record
			return true
		}
		return false
	})

	if r == nil {
		return nil, fmt.Errorf("no unbonding record with hash %q found", msg.Hash)
	}

	if msg.NewStatus == types.WithdrawStatusQueued {
		// update sendErrors to zero, nil the distributions and amount (as this we be recalculated when processed), and update the state to queued.
		r.SendErrors = 0
		r.Amount = nil
		r.Distribution = nil
		r.CompletionTime = time.Time{}
		r.Acknowledged = false
	}

	k.UpdateWithdrawalRecordStatus(ctx, r, msg.NewStatus)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRedemptionRequeue,
			sdk.NewAttribute(types.AttributeKeyUser, r.Delegator),
			sdk.NewAttribute(types.AttributeKeyHash, msg.Hash),
			sdk.NewAttribute(types.AttributeKeyNewStatus, string(msg.NewStatus)),
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
		),
	})

	return &types.MsgUpdateRedemptionResponse{}, nil
}

func (k msgServer) SignalIntent(goCtx context.Context, msg *types.MsgSignalIntent) (*types.MsgSignalIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get zone
	zone, ok := k.GetZone(ctx, msg.ChainId)
	if !ok {
		return nil, fmt.Errorf("invalid chain id %q", msg.ChainId)
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

	// checking msg authority is the gov module address
	if k.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(ctx), msg.Authority,
			)
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
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(ctx), msg.Authority,
			)
	}

	zone, found := k.GetZone(ctx, msg.ChainId)
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

func (k msgServer) GovAddValidatorDenyList(goCtx context.Context, msg *types.MsgGovAddValidatorDenyList) (*types.MsgGovAddValidatorDenyListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(ctx), msg.Authority,
			)
	}

	zone, found := k.GetZone(ctx, msg.ChainId)
	if !found {
		return nil,
			fmt.Errorf(
				"no zone found for: %s",
				msg.ChainId,
			)
	}
	valAddr, err := addressutils.ValAddressFromBech32(msg.OperatorAddress, zone.GetValoperPrefix())
	if err != nil {
		return nil, err
	}

	if err := k.SetZoneValidatorToDenyList(ctx, zone.ChainId, valAddr); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeAddValidatorDenyList,
			sdk.NewAttribute(types.AttributeKeyOperatorAddress, msg.OperatorAddress),
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
		),
	})

	return &types.MsgGovAddValidatorDenyListResponse{}, nil
}

func (k msgServer) GovRemoveValidatorDenyList(goCtx context.Context, msg *types.MsgGovRemoveValidatorDenyList) (*types.MsgGovRemoveValidatorDenyListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(ctx), msg.Authority,
			)
	}

	zone, found := k.GetZone(ctx, msg.ChainId)
	if !found {
		return nil,
			fmt.Errorf(
				"no zone found for: %s",
				msg.ChainId,
			)
	}
	valAddr, err := addressutils.ValAddressFromBech32(msg.OperatorAddress, zone.GetValoperPrefix())
	if err != nil {
		return nil, err
	}
	if found := k.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr); !found {
		return nil, fmt.Errorf("validator %s not found in deny list", msg.OperatorAddress)
	}

	k.RemoveValidatorFromDenyList(ctx, zone.ChainId, valAddr)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRemoveValidatorDenyList,
			sdk.NewAttribute(types.AttributeKeyOperatorAddress, msg.OperatorAddress),
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
		),
	})

	return &types.MsgGovRemoveValidatorDenyListResponse{}, nil
}

func (k msgServer) GovExecuteICATx(goCtx context.Context, msg *types.MsgGovExecuteICATx) (*types.MsgGovExecuteICATxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(ctx), msg.Authority,
			)
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

	err = k.SubmitTx(ctx, unpackedMsgs, account, types.FlagNoFurtherAction, zone.MessagesPerTx)
	if err != nil {
		return nil, err
	}

	return &types.MsgGovExecuteICATxResponse{}, nil
}

// GovSetZoneOffboarding sets the offboarding status for a zone.
// When offboarding is enabled, deposits are disabled and redemption rate updates are frozen.
func (k msgServer) GovSetZoneOffboarding(goCtx context.Context, msg *types.MsgGovSetZoneOffboarding) (*types.MsgGovSetZoneOffboardingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(ctx), msg.Authority,
			)
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
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(ctx), msg.Authority,
			)
	}

	zone, found := k.GetZone(ctx, msg.ChainId)
	if !found {
		return nil, fmt.Errorf("zone not found for chain id: %s", msg.ChainId)
	}

	var cancelledCount uint64
	refundedAmounts := sdk.NewCoins()
	recordsToDelete := []types.WithdrawalRecord{}
	amountsToRefund := sdk.NewCoins()

	// Iterate through all queued withdrawal records for this zone
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusQueued, func(idx int64, record types.WithdrawalRecord) bool {
		recordsToDelete = append(recordsToDelete, record)
		amountsToRefund = amountsToRefund.Add(record.BurnAmount)
		return false
	})

	// Also cancel records in UNBOND status (unbonding initiated but not complete)
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusUnbond, func(idx int64, record types.WithdrawalRecord) bool {
		recordsToDelete = append(recordsToDelete, record)
		amountsToRefund = amountsToRefund.Add(record.BurnAmount)
		return false
	})

	if len(amountsToRefund) > 0 {
		// check if the amounts to refund are greater than the escrow balance
		escrowBalance := k.BankKeeper.GetBalance(ctx, k.AccountKeeper.GetModuleAddress(types.EscrowModuleAccount), amountsToRefund[0].Denom)
		if escrowBalance.IsLT(amountsToRefund[0]) {
			return nil, fmt.Errorf("insufficient escrow balance to refund. Escrow balance: %s, amounts to refund: %s", escrowBalance.String(), amountsToRefund[0].String())
		}

		// Process each record
		for _, record := range recordsToDelete {
			userAccAddress, err := addressutils.AddressFromBech32(record.Delegator, "")
			if err != nil {
				k.Logger(ctx).Error("failed to parse delegator address", "address", record.Delegator, "error", err)
				return nil, err
			}

			// Refund the qAssets from escrow to user
			if err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.EscrowModuleAccount, userAccAddress, sdk.NewCoins(record.BurnAmount)); err != nil {
				k.Logger(ctx).Error("failed to refund qAssets", "user", record.Delegator, "amount", record.BurnAmount, "error", err)
				return nil, err
			}

			// Delete the withdrawal record
			k.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)

			cancelledCount++
			refundedAmounts = refundedAmounts.Add(record.BurnAmount)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeRedemptionCancellation,
					sdk.NewAttribute(types.AttributeKeyReturnedAmount, record.BurnAmount.String()),
					sdk.NewAttribute(types.AttributeKeyUser, record.Delegator),
					sdk.NewAttribute(types.AttributeKeyHash, record.Txhash),
					sdk.NewAttribute(types.AttributeKeyChainID, record.ChainId),
				),
			)
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeCancelAllPendingRedemptions,
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
			sdk.NewAttribute(types.AttributeKeyCancelledCount, fmt.Sprintf("%d", cancelledCount)),
			sdk.NewAttribute(types.AttributeKeyRefundedAmounts, refundedAmounts.String()),
		),
	})

	return &types.MsgGovCancelAllPendingRedemptionsResponse{
		CancelledCount:  cancelledCount,
		RefundedAmounts: refundedAmounts,
	}, nil
}

// GovForceUnbondAllDelegations initiates unbonding of all delegations for an offboarding zone.
// It creates MsgUndelegate messages for each validator and submits them via ICA.
func (k msgServer) GovForceUnbondAllDelegations(goCtx context.Context, msg *types.MsgGovForceUnbondAllDelegations) (*types.MsgGovForceUnbondAllDelegationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.GetGovAuthority(ctx) != msg.Authority {
		return nil,
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(ctx), msg.Authority,
			)
	}

	zone, found := k.GetZone(ctx, msg.ChainId)
	if !found {
		return nil, fmt.Errorf("zone not found for chain id: %s", msg.ChainId)
	}

	// Ensure zone is in offboarding mode
	if !zone.IsOffboarding {
		return nil, fmt.Errorf("zone %s is not in offboarding mode", msg.ChainId)
	}

	// Ensure delegation address exists
	if zone.DelegationAddress == nil {
		return nil, fmt.Errorf("zone %s has no delegation address", msg.ChainId)
	}

	// Get all delegations for this zone
	delegations := k.GetAllDelegations(ctx, zone.ChainId)
	if len(delegations) == 0 {
		return nil, fmt.Errorf("no delegations found for zone %s", msg.ChainId)
	}

	var msgs []sdk.Msg
	totalUnbonded := sdk.NewCoin(zone.BaseDenom, math.ZeroInt())
	var unbondingCount uint64

	// Create MsgUndelegate for each delegation
	for _, delegation := range delegations {
		if delegation.Amount.IsZero() {
			continue
		}

		undelegateMsg := &stakingtypes.MsgUndelegate{
			DelegatorAddress: zone.DelegationAddress.Address,
			ValidatorAddress: delegation.ValidatorAddress,
			Amount:           delegation.Amount,
		}
		msgs = append(msgs, undelegateMsg)
		totalUnbonded = totalUnbonded.Add(delegation.Amount)
		unbondingCount++
	}

	if len(msgs) == 0 {
		return nil, errors.New("no delegations to unbond")
	}

	// Use a special memo prefix for offboarding unbonds
	offboardingMemo := types.OffboardingUnbondMemo(ctx.BlockHeight())

	// Submit the unbonding transactions via ICA
	if err := k.SubmitTx(ctx, msgs, zone.DelegationAddress, offboardingMemo, zone.MessagesPerTx); err != nil {
		return nil, fmt.Errorf("failed to submit unbonding transactions: %w", err)
	}

	// Increment withdrawal waitgroup for the unbonding messages
	if err := zone.IncrementWithdrawalWaitgroup(k.Logger(ctx), uint32(len(msgs)), "offboarding unbond messages"); err != nil { //nolint:gosec
		return nil, err
	}
	k.SetZone(ctx, &zone)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeForceUnbondAllDelegations,
			sdk.NewAttribute(types.AttributeKeyChainID, msg.ChainId),
			sdk.NewAttribute(types.AttributeKeyUnbondingCount, fmt.Sprintf("%d", unbondingCount)),
			sdk.NewAttribute(types.AttributeKeyTotalUnbonded, totalUnbonded.String()),
		),
	})

	return &types.MsgGovForceUnbondAllDelegationsResponse{
		UnbondingCount: unbondingCount,
		TotalUnbonded:  totalUnbonded,
	}, nil
}
