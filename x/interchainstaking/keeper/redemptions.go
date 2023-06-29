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
	intent, found := k.GetDelegatorIntent(ctx, zone, sender.String(), false)
	// msgs is slice of MsgTokenizeShares, so we can handle dust allocation later.
	msgs := make([]*lsmstakingtypes.MsgTokenizeShares, 0)
	var err error
	intents := intent.Intents
	if !found || len(intents) == 0 {
		// if user has no intent set (this can happen if redeeming tokens that were obtained offchain), use global intent.
		// Note: this can be improved; user will receive a bunch of tokens.
		intents, err = k.GetAggregateIntentOrDefault(ctx, zone)
		if err != nil {
			return err
		}
	}
	outstanding := nativeTokens
	distribution := make(map[string]uint64, 0)

	availablePerValidator, _, err := k.GetUnlockedTokensForZone(ctx, zone)
	if err != nil {
		return err
	}

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
	k.AddWithdrawalRecord(ctx, zone.ChainId, sender.String(), []*types.Distribution{}, destination, sdk.Coins{}, burnAmount, hash, types.WithdrawStatusTokenize, time.Unix(0, 0))

	return k.SubmitTx(ctx, sdkMsgs, zone.DelegationAddress, hash, zone.MessagesPerTx)
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
		types.WithdrawStatusQueued,
		time.Time{},
	)

	return nil
}

// GetUnlockedTokensForZone will iterate over all delegation records for a zone, and then remove the
// locked tokens (those actively being redelegated), returning a slice of int64 staking tokens that
// are unlocked and free to redelegate or unbond.
func (k *Keeper) GetUnlockedTokensForZone(ctx sdk.Context, zone *types.Zone) (map[string]math.Int, math.Int, error) {
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
				return map[string]math.Int{}, sdk.ZeroInt(), fmt.Errorf("negative available amount [chain: %s, validator: %s, amount: %s]; unable to continue", zone.ChainId, redelegation.Destination, availablePerValidator[redelegation.Destination].String())
			}
			total = total.Sub(sdk.NewInt(redelegation.Amount))
		}
	}

	return availablePerValidator, total, nil
}

// HandleQueuedUnbondings is called once per epoch to aggregate all queued unbondings into
// a single unbond transaction per delegation.
func (k *Keeper) HandleQueuedUnbondings(ctx sdk.Context, zone *types.Zone, epoch int64) error {
	// out here will only ever be in native bond denom
	valOutCoinsMap := make(map[string]sdk.Coin, 0)
	txHashes := make(map[string][]string, 0)

	totalToWithdraw := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	txDistrsMap := make(map[string][]*types.Distribution, 0)
	txCoinMap := make(map[string]sdk.Coin, 0)
	_, totalAvailable, err := k.GetUnlockedTokensForZone(ctx, zone)
	if err != nil {
		return err
	}

	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusQueued, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Info("handling queued withdrawal request", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", withdrawal.Amount)
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

		txCoinMap[withdrawal.Txhash] = withdrawal.Amount[0]
		txDistrsMap[withdrawal.Txhash] = make([]*types.Distribution, 0)
		return false
	})

	// no undelegations to attempt
	if totalToWithdraw.IsZero() {
		return nil
	}

	allocationsMap, err := k.DeterminePlanForUndelegation(ctx, zone, sdk.NewCoins(totalToWithdraw))
	if err != nil {
		return err
	}
	valopers := utils.Keys(allocationsMap)
	vidx := 0
	v := valopers[vidx]
WITHDRAWAL:
	for _, hash := range utils.Keys(txCoinMap) {
		for {
			if txCoinMap[hash].Amount.IsZero() {
				continue WITHDRAWAL
			}
			if allocationsMap[v].GT(txCoinMap[hash].Amount) {
				allocationsMap[v] = allocationsMap[v].Sub(txCoinMap[hash].Amount)
				txDistrsMap[hash] = append(txDistrsMap[hash], &types.Distribution{Valoper: v, Amount: txCoinMap[hash].Amount.Uint64()})
				existing, found := valOutCoinsMap[v]
				if !found {
					valOutCoinsMap[v] = txCoinMap[hash]
					txHashes[v] = []string{hash}

				} else {
					valOutCoinsMap[v] = existing.Add(txCoinMap[hash])
					txHashes[v] = append(txHashes[v], hash)
				}
				txCoinMap[hash] = sdk.NewCoin(txCoinMap[hash].Denom, sdk.ZeroInt())
				continue WITHDRAWAL
			}

			txDistrsMap[hash] = append(txDistrsMap[hash], &types.Distribution{Valoper: v, Amount: allocationsMap[v].Uint64()})
			txCoinMap[hash] = sdk.NewCoin(txCoinMap[hash].Denom, txCoinMap[hash].Amount.Sub(allocationsMap[v]))
			existing, found := valOutCoinsMap[v]
			if !found {
				valOutCoinsMap[v] = sdk.NewCoin(zone.BaseDenom, allocationsMap[v])
				txHashes[v] = []string{hash}

			} else {
				valOutCoinsMap[v] = existing.Add(sdk.NewCoin(zone.BaseDenom, allocationsMap[v]))
				txHashes[v] = append(txHashes[v], hash)
			}

			allocationsMap[v] = sdk.ZeroInt()
			if allocationsMap[v].IsZero() {
				if len(valopers) > vidx+1 {
					vidx++
					v = valopers[vidx]
				} else {
					if !txCoinMap[hash].Amount.IsZero() {
						return fmt.Errorf("unable to satisfy unbonding")
					}
					continue WITHDRAWAL
				}
			}
		}
	}

	for _, hash := range utils.Keys(txDistrsMap) {
		record, found := k.GetWithdrawalRecord(ctx, zone.ChainId, hash, types.WithdrawStatusQueued)
		if !found {
			return errors.New("unable to find withdrawal record")
		}
		record.Distribution = txDistrsMap[hash]
		k.UpdateWithdrawalRecordStatus(ctx, &record, types.WithdrawStatusUnbond)
	}

	if len(txHashes) == 0 {
		// no records to handle.
		return nil
	}

	var msgs []sdk.Msg
	for _, valoper := range utils.Keys(valOutCoinsMap) {
		if !valOutCoinsMap[valoper].Amount.IsZero() {
			msgs = append(msgs, &stakingtypes.MsgUndelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: valoper, Amount: valOutCoinsMap[valoper]})
		}
	}

	k.Logger(ctx).Info("unbonding messages to send", "msg", msgs)

	err = k.SubmitTx(ctx, msgs, zone.DelegationAddress, types.EpochWithdrawalMemo(epoch), zone.MessagesPerTx)
	if err != nil {
		return err
	}

	for _, valoper := range utils.Keys(valOutCoinsMap) {
		if !valOutCoinsMap[valoper].Amount.IsZero() {
			sort.Strings(txHashes[valoper])
			k.SetUnbondingRecord(ctx, types.UnbondingRecord{ChainId: zone.ChainId, EpochNumber: epoch, Validator: valoper, RelatedTxhash: txHashes[valoper]})
		}
	}

	return nil
}

func (k *Keeper) GCCompletedUnbondings(ctx sdk.Context, zone *types.Zone) error {
	var err error

	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusCompleted, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		if ctx.BlockTime().After(withdrawal.CompletionTime.Add(24 * time.Hour)) {
			k.Logger(ctx).Info("garbage collecting completed unbondings")
			k.DeleteWithdrawalRecord(ctx, zone.ChainId, withdrawal.Txhash, types.WithdrawStatusCompleted)
		}
		return false
	})

	return err
}

func (k *Keeper) DeterminePlanForUndelegation(ctx sdk.Context, zone *types.Zone, amount sdk.Coins) (map[string]math.Int, error) {
	currentAllocations, currentSum, _ := k.GetDelegationMap(ctx, zone)
	availablePerValidator, _, err := k.GetUnlockedTokensForZone(ctx, zone)
	if err != nil {
		return nil, err
	}
	targetAllocations, err := k.GetAggregateIntentOrDefault(ctx, zone)
	if err != nil {
		return nil, err
	}
	allocations := types.DetermineAllocationsForUndelegation(currentAllocations, currentSum, targetAllocations, availablePerValidator, amount)
	return allocations, nil
}
