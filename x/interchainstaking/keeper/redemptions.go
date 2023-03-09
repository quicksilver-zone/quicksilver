package keeper

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	lsmstakingtypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// processRedemptionForLsm will determine based on user intent, the tokens to return to the user, generate Redeem message and send them.
func (k *Keeper) processRedemptionForLsm(ctx sdk.Context, zone *types.Zone, sender sdk.AccAddress, destination string, nativeTokens math.Int, burnAmount sdk.Coin, hash string) error {
	intent, found := k.GetIntent(ctx, zone, sender.String(), false)
	// msgs is slice of MsgTokenizeShares, so we can handle dust allocation later.
	msgs := make([]*lsmstakingtypes.MsgTokenizeShares, 0)
	intents := intent.Intents
	if !found || len(intents) == 0 {
		// if user has no intent set (this can happen if redeeming tokens that were obtained offchain), use global intent.
		// Note: this can be improved; user will receive a bunch of tokens.
		intents = zone.GetAggregateIntentOrDefault()
	}
	outstanding := nativeTokens
	distribution := make(map[string]uint64, 0)

	availablePerValidator, _ := k.GetUnlockedTokensForZone(ctx, zone)

	for _, intent := range intents.Sort() {
		thisAmount := intent.Weight.MulInt(nativeTokens).TruncateInt()
		if thisAmount.GT(availablePerValidator[intent.ValoperAddress]) {
			return errors.New("unable to satisfy unbond request; delegations may be locked")
		}
		distribution[intent.ValoperAddress] = thisAmount.Uint64()
		outstanding = outstanding.Sub(thisAmount)
	}

	distribution[intents[0].ValoperAddress] += outstanding.Uint64()

	for _, valoper := range utils.Keys(distribution) {
		msgs = append(msgs, &lsmstakingtypes.MsgTokenizeShares{
			DelegatorAddress:    zone.DelegationAddress.Address,
			ValidatorAddress:    valoper,
			Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewIntFromUint64(distribution[valoper])),
			TokenizedShareOwner: destination,
		})
	}
	// add unallocated dust.
	msgs[0].Amount = msgs[0].Amount.AddAmount(outstanding)
	sdkMsgs := make([]sdk.Msg, 0)
	for _, msg := range msgs {
		sdkMsgs = append(sdkMsgs, sdk.Msg(msg))
	}
	k.AddWithdrawalRecord(ctx, zone.ChainId, sender.String(), []*types.Distribution{}, destination, sdk.Coins{}, burnAmount, hash, WithdrawStatusTokenize, time.Unix(0, 0))

	return k.SubmitTx(ctx, sdkMsgs, zone.DelegationAddress, hash)
}

// queueRedemption will determine based on zone intent, the tokens to unbond, and add a withdrawal record with status QUEUED.
func (k *Keeper) queueRedemption(
	ctx sdk.Context,
	zone *types.Zone,
	sender sdk.AccAddress,
	destination string,
	nativeTokens math.Int,
	burnAmount sdk.Coin,
	hash string,
) error { //nolint:unparam // we know that the error is always nil
	distribution := make([]*types.Distribution, 0)

	amount := sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, nativeTokens))
	k.AddWithdrawalRecord(
		ctx,
		zone.ChainId,
		sender.String(),
		distribution,
		destination,
		amount,
		burnAmount,
		hash,
		WithdrawStatusQueued,
		time.Time{},
	)

	return nil
}

// GetUnlockedTokensForZone will iterate over all delegation records for a zone, and then remove the
// locked tokens (those actively being redelegated), returning a slice of int64 staking tokens that
// are unlocked and free to redelegate or unbond.
func (k *Keeper) GetUnlockedTokensForZone(ctx sdk.Context, zone *types.Zone) (map[string]math.Int, math.Int) {
	availablePerValidator := make(map[string]math.Int, len(zone.Validators))
	total := sdk.ZeroInt()
	for _, delegation := range k.GetAllDelegations(ctx, zone) {
		thisAvailable, found := availablePerValidator[delegation.ValidatorAddress]
		if !found {
			thisAvailable = sdk.ZeroInt()
		}
		availablePerValidator[delegation.ValidatorAddress] = thisAvailable.Add(delegation.Amount.Amount)
		total = total.Add(delegation.Amount.Amount)
	}
	for _, redelegation := range k.ZoneRedelegationRecords(ctx, zone.ChainId) {
		thisAvailable, found := availablePerValidator[redelegation.Destination]
		if found {
			availablePerValidator[redelegation.Destination] = thisAvailable.Sub(sdk.NewInt(redelegation.Amount))
			if availablePerValidator[redelegation.Destination].LT(sdk.ZeroInt()) {
				panic("negative available amount; unable to continue")
			}
			total = total.Sub(sdk.NewInt(redelegation.Amount))
		}
	}

	return availablePerValidator, total
}

// HandleQueuedUnbondings is called once per epoch to aggregate all queued unbondings into
// a single unbond transaction per delegation.
func (k *Keeper) HandleQueuedUnbondings(ctx sdk.Context, zone *types.Zone, epoch int64) error {
	// out here will only ever be in native bond denom
	out := make(map[string]sdk.Coin, 0)
	txhashes := make(map[string][]string, 0)

	totalToWithdraw := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	distributions := make(map[string][]*types.Distribution, 0)
	amounts := make(map[string]sdk.Coin, 0)
	_, totalAvailable := k.GetUnlockedTokensForZone(ctx, zone)

	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, WithdrawStatusQueued, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Info("handling queued withdrawal request", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", withdrawal.Amount, "chain_id", zone.ChainId)
		if len(withdrawal.Amount) == 0 {
			k.Logger(ctx).Error("withdrawal %s has no amount set; cannot process...", withdrawal.Txhash)
			return false
		}
		if totalAvailable.LT(totalToWithdraw.Amount.Add(withdrawal.Amount[0].Amount)) {
			k.Logger(ctx).Error("unable to satisfy further unbondings this epoch")
			// do not process this or subsequent withdrawals this epoch.
			return true
		}
		totalToWithdraw = totalToWithdraw.Add(withdrawal.Amount[0])

		amounts[withdrawal.Txhash] = withdrawal.Amount[0]
		distributions[withdrawal.Txhash] = make([]*types.Distribution, 0)
		return false
	})

	// no undelegations to attempt
	if totalToWithdraw.IsZero() {
		return nil
	}

	allocations := k.DeterminePlanForUndelegation(ctx, zone, sdk.NewCoins(totalToWithdraw))
	valopers := utils.Keys(allocations)
	vidx := 0
	v := valopers[vidx]
WITHDRAWAL:
	for _, hash := range utils.Keys(amounts) {
		for {
			fmt.Println(amounts[hash].Amount)
			if amounts[hash].Amount.IsZero() {
				continue WITHDRAWAL
			}
			if allocations[v].GT(amounts[hash].Amount) {
				allocations[v] = allocations[v].Sub(amounts[hash].Amount)
				distributions[hash] = append(distributions[hash], &types.Distribution{Valoper: v, Amount: amounts[hash].Amount.Uint64()})
				existing, found := out[v]
				if !found {
					out[v] = amounts[hash]
					txhashes[v] = []string{hash}

				} else {
					out[v] = existing.Add(amounts[hash])
					txhashes[v] = append(txhashes[v], hash)
				}
				amounts[hash] = sdk.NewCoin(amounts[hash].Denom, sdk.ZeroInt())
				continue WITHDRAWAL
			} else {
				distributions[hash] = append(distributions[hash], &types.Distribution{Valoper: v, Amount: allocations[v].Uint64()})
				amounts[hash] = sdk.NewCoin(amounts[hash].Denom, amounts[hash].Amount.Sub(allocations[v]))
				existing, found := out[v]
				if !found {
					out[v] = sdk.NewCoin(zone.BaseDenom, allocations[v])
					txhashes[v] = []string{hash}

				} else {
					out[v] = existing.Add(sdk.NewCoin(zone.BaseDenom, allocations[v]))
					txhashes[v] = append(txhashes[v], hash)
				}
				allocations[v] = sdk.ZeroInt()
			}
			if allocations[v].IsZero() {
				fmt.Println("valopers len", len(valopers))
				fmt.Println("vidx+1", vidx+1)
				if len(valopers) > vidx+1 {
					vidx++
					v = valopers[vidx]
				} else {
					if !amounts[hash].Amount.IsZero() {
						return fmt.Errorf("unable to satisfy unbonding")
					}
					continue WITHDRAWAL
				}
			}
		}
	}

	for _, hash := range utils.Keys(distributions) {
		record, found := k.GetWithdrawalRecord(ctx, zone.ChainId, hash, WithdrawStatusQueued)
		if !found {
			return errors.New("unable to find withdrawal record")
		}
		record.Distribution = distributions[hash]
		k.UpdateWithdrawalRecordStatus(ctx, &record, WithdrawStatusUnbond)
	}

	if len(txhashes) == 0 {
		// no records to handle.
		return nil
	}

	var msgs []sdk.Msg
	for _, valoper := range utils.Keys(out) {
		if !out[valoper].Amount.IsZero() {
			sort.Strings(txhashes[valoper])
			k.SetUnbondingRecord(ctx, types.UnbondingRecord{ChainId: zone.ChainId, EpochNumber: epoch, Validator: valoper, RelatedTxhash: txhashes[valoper]})
			msgs = append(msgs, &stakingtypes.MsgUndelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: valoper, Amount: out[valoper]})
		}
	}

	k.Logger(ctx).Info("unbonding messages to send", "msg", msgs)

	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, fmt.Sprintf("withdrawal/%d", epoch))
}

func (k *Keeper) GCCompletedUnbondings(ctx sdk.Context, zone *types.Zone) error {
	var err error

	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, WithdrawStatusCompleted, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		if ctx.BlockTime().After(withdrawal.CompletionTime.Add(24 * time.Hour)) {
			k.Logger(ctx).Info("garbage collecting completed unbondings")
			k.DeleteWithdrawalRecord(ctx, zone.ChainId, withdrawal.Txhash, WithdrawStatusCompleted)
		}
		return false
	})

	return err
}

func (k *Keeper) DeterminePlanForUndelegation(ctx sdk.Context, zone *types.Zone, amount sdk.Coins) map[string]math.Int {
	currentAllocations, currentSum, _ := k.GetDelegationMap(ctx, zone)
	availablePerValidator, _ := k.GetUnlockedTokensForZone(ctx, zone)
	targetAllocations := zone.GetAggregateIntentOrDefault()
	allocations := types.DetermineAllocationsForUndelegation(currentAllocations, currentSum, targetAllocations, availablePerValidator, amount)
	return allocations
}
