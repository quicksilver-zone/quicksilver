package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
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

func (k msgServer) RequestRedemption(goCtx context.Context, msg *types.MsgRequestRedemption) (*types.MsgRequestRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate coins are positive
	err := msg.Value.Validate()
	if err != nil {
		return nil, err
	}

	// validate recipient address
	if len(msg.DestinationAddress) == 0 {
		return nil, fmt.Errorf("recipient address not provided")
	}

	if _, _, err = bech32.DecodeAndConvert(msg.DestinationAddress); err != nil {
		return nil, err
	}

	var zone *types.Zone

	k.IterateZones(ctx, func(_ int64, thisZone types.Zone) bool {
		if thisZone.LocalDenom == msg.Value.GetDenom() {
			zone = &thisZone
			return true
		}
		return false
	})

	// does zone exist?
	if nil == zone {
		return nil, fmt.Errorf("unable to find matching zone for denom %s", msg.Value.GetDenom())
	}

	// does destination address match the prefix registered against the zone?
	if _, err := utils.AccAddressFromBech32(msg.DestinationAddress, zone.AccountPrefix); err != nil {
		return nil, fmt.Errorf("destination address %s does not match expected prefix %s", msg.DestinationAddress, zone.AccountPrefix)
	}

	sender, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	// does the user have sufficient assets to burn
	if !k.BankKeeper.HasBalance(ctx, sender, msg.Value) {
		return nil, fmt.Errorf("account has insufficient balance of qasset to burn")
	}

	// get min of LastRedemptionRate (N-1) and RedemptionRate (N)
	var rate sdk.Dec
	rate = zone.LastRedemptionRate
	if zone.RedemptionRate.LT(rate) {
		rate = zone.RedemptionRate
	}

	nativeTokens := sdk.NewDecFromInt(msg.Value.Amount).Mul(rate).TruncateInt()

	outTokens := sdk.NewCoin(zone.BaseDenom, nativeTokens)
	k.Logger(ctx).Error("outtokens", "o", outTokens)

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(ctx.BlockHeight()))
	hash := sha256.Sum256(append(msg.GetSignBytes(), heightBytes...))
	hashString := hex.EncodeToString(hash[:])

	// lock qAssets - how are we tracking this?
	if err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(msg.Value)); err != nil {
		return nil, err
	}

	userIntent, found := k.GetIntent(ctx, *zone, msg.FromAddress, false)

	if !found || len(userIntent.Intents) == 0 {
		vi := []*types.ValidatorIntent{}
		for _, v := range zone.GetAggregateIntentOrDefault() {
			vi = append(vi, v)
		}
		userIntent = types.DelegatorIntent{Delegator: msg.FromAddress, Intents: vi}
	}

	intentMap := userIntent.ToAllocations(sdk.NewDecFromInt(nativeTokens))

	targets, err := k.GetRedemptionTargets(ctx, *zone, intentMap) // map[string][string]sdk.Coin
	if err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("targets can never be zero length")
	}

	sumAmount := sdk.NewCoins()

	// redeemType := "tokenize"
	redeemType := "unbond" // TODO: revert to "tokenize"
	// does zone have LSM enabled?
	// if !zone.LiquidityModule {
	// 	// unbond workflow.
	// 	redeemType = "unbond"
	// }

	msgs := make(map[string][]sdk.Msg, 0)

	for _, target := range targets.Sorted() {
		if len(target.Value) == 1 {
			if _, ok := msgs[target.DelegatorAddress]; !ok {
				msgs[target.DelegatorAddress] = make([]sdk.Msg, 0)
			}
			if redeemType == "tokenize" {
				msgs[target.DelegatorAddress] = append(msgs[target.DelegatorAddress], &stakingtypes.MsgTokenizeShares{
					DelegatorAddress:    target.DelegatorAddress,
					ValidatorAddress:    target.ValidatorAddress,
					Amount:              target.Value[0],
					TokenizedShareOwner: msg.DestinationAddress,
				})
			} else {
				msgs[target.DelegatorAddress] = append(msgs[target.DelegatorAddress], &stakingtypes.MsgUndelegate{
					DelegatorAddress: target.DelegatorAddress,
					ValidatorAddress: target.ValidatorAddress,
					Amount:           target.Value[0],
				})
			}
			sumAmount = sumAmount.Add(target.Value[0])
			if _, found := k.GetWithdrawalRecord(ctx, zone, hashString, target.DelegatorAddress, target.ValidatorAddress); found {
				return nil, fmt.Errorf("cannot withdraw twice for the same delegator/validator tuple in a single transaction")
			}
			k.Logger(ctx).Info("Store", "del", target.DelegatorAddress, "val", target.ValidatorAddress, "hash", hashString, "chain", zone.ChainId)
			k.AddWithdrawalRecord(ctx, zone, target.DelegatorAddress, target.ValidatorAddress, msg.DestinationAddress, target.Value[0], msg.Value, hashString, time.Unix(0, 0))
		}
	}

	k.Logger(ctx).Error("messages", "m", msgs)

	delegators := make([]string, 0, len(msgs))
	for delegator := range msgs {
		delegators = append(delegators, delegator)
	}
	sort.Strings(delegators)

	for _, delegator := range delegators {
		icaAccount, err := zone.GetDelegationAccountByAddress(delegator)
		if err != nil {
			// panic here because something is terribly wrong if we can't find the delegation bucket here!!!
			panic(err)
		}
		err = k.SubmitTx(ctx, msgs[delegator], icaAccount, hashString)
		if err != nil {
			k.Logger(ctx).Error("error submitting tx", "err", err)
			return nil, err
		}
	}

	if !sumAmount.IsAllLTE(sdk.NewCoins(outTokens)) {
		k.Logger(ctx).Error("output coins > than expected!", "sum", sumAmount, "expected", outTokens)
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRedemptionRequest,
			sdk.NewAttribute(types.AttributeKeyBurnAmount, msg.Value.String()),
			sdk.NewAttribute(types.AttributeKeyRedeemAmount, sumAmount.String()),
			sdk.NewAttribute(types.AttributeKeyRecipientAddress, msg.DestinationAddress),
			sdk.NewAttribute(types.AttributeKeyRecipientChain, zone.ChainId),
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

	// validate intents (aggregated errors)
	if err := k.validateIntents(zone, msg.Intents); err != nil {
		return nil, err
	}

	intent := types.DelegatorIntent{
		Delegator: msg.FromAddress,
		Intents:   msg.Intents,
	}

	k.SetIntent(ctx, zone, intent, false)

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

func (k msgServer) validateIntents(zone types.Zone, intents []*types.ValidatorIntent) error {
	errors := make(map[string]error)

	for i, intent := range intents {
		_, found := zone.GetValidatorByValoper(intent.ValoperAddress)
		if !found {
			errors[fmt.Sprintf("intent[%v]", i)] = fmt.Errorf("unable to find valoper %s", intent.ValoperAddress)
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

func (k Keeper) EmitValsetRequery(ctx sdk.Context, connectionID string, chainID string) error {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz1, err := k.cdc.Marshal(&bondedQuery)
	if err != nil {
		return err
	}
	unbondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusUnbonded}
	bz2, err := k.cdc.Marshal(&unbondedQuery)
	if err != nil {
		return err
	}
	unbondingQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusUnbonding}
	bz3, err := k.cdc.Marshal(&unbondingQuery)
	if err != nil {
		return err
	}

	period := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))

	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		chainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz1,
		sdk.NewInt(period),
		types.ModuleName,
		"valset",
		0,
	)
	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		chainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz2,
		sdk.NewInt(period),
		types.ModuleName,
		"valset",
		0,
	)
	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		chainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz3,
		sdk.NewInt(period),
		types.ModuleName,
		"valset",
		0,
	)
	return nil
}
