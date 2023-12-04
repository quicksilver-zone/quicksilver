package keeper

import (
	"errors"
	"fmt"
	"sort"
	"time"

	lsmstakingtypes "github.com/quicksilver-zone/quicksilver/x/lsmtypes"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
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

	sdkMsgs := make([]sdk.Msg, 0)
	for _, msg := range msgs {
		sdkMsgs = append(sdkMsgs, sdk.Msg(msg))
	}
	distributions := make([]*types.Distribution, 0)

	for valoper, amount := range distribution {
		newDistribution := types.Distribution{
			Valoper: valoper,
			Amount:  amount,
		}
		distributions = append(distributions, &newDistribution)
	}

	k.AddWithdrawalRecord(ctx, zone.ChainId, sender.String(), distributions, destination, sdk.Coins{}, burnAmount, hash, types.WithdrawStatusTokenize, time.Unix(0, 0), k.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch)

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
	distributions := make([]*types.Distribution, 0)
	amount := sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, nativeTokens))

	k.AddWithdrawalRecord(
		ctx,
		zone.ChainId,
		sender.String(),
		distributions,
		destination,
		amount,
		burnAmount,
		hash,
		types.WithdrawStatusQueued,
		time.Time{},
		k.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch,
	)

	return nil
}

// GetUnlockedTokensForZone will iterate over all delegation records for a zone, and then remove the
// locked tokens (those actively being redelegated), returning a slice of int64 staking tokens that
// are unlocked and free to redelegate or unbond.
func (k *Keeper) GetUnlockedTokensForZone(ctx sdk.Context, zone *types.Zone) (map[string]math.Int, math.Int, error) {
	availablePerValidator := make(map[string]math.Int, len(zone.Validators))
	total := sdk.ZeroInt()
	for _, delegation := range k.GetAllDelegations(ctx, zone.ChainId) {
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
	coinsOutPerValidator := make(map[string]sdk.Coin, 0)
	// list of withdrawal tx hashes per validator
	txHashesPerValidator := make(map[string][]string, 0)

	// total amount coins to withdraw
	totalToWithdraw := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())

	// map of distributions per withdrawal
	distributionsPerWithdrawal := make(map[string][]*types.Distribution, 0)

	// map of coins per withdrawal
	amountToWithdrawPerWithdrawal := make(map[string]sdk.Coin, 0)

	// find total number of unlockedTokens (not locked by previous redelegations) for the given zone
	_, totalAvailable, err := k.GetUnlockedTokensForZone(ctx, zone)
	if err != nil {
		return err
	}

	// iterate all withdrawal records for the zone in the QUEUED state.
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusQueued, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Info("handling queued withdrawal request", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", withdrawal.Amount)
		if len(withdrawal.Amount) != 1 { // native unbonding can only unbond the baseDenom
			k.Logger(ctx).Error("withdrawal %s has no amount set; cannot process...", withdrawal.Txhash)
			return false
		}

		if !withdrawal.Amount[0].IsPositive() {
			k.Logger(ctx).Error("withdrawal %s attempting to withdraw non-positive amount; cannot process...", withdrawal.Txhash)
			return false
		}

		// native unbonding can only unbond the baseDenom
		if withdrawal.Amount[0].Denom != zone.BaseDenom {
			k.Logger(ctx).Error("withdrawal %s attempting to withdraw invalid amount; cannot process...", withdrawal.Txhash)
			return false
		}

		// check whether the running total of withdrawals can be satisfied by the available unlocked tokens.
		// if not return true to stop iterating and return all records up until now.
		if totalAvailable.LT(totalToWithdraw.Amount.Add(withdrawal.Amount[0].Amount)) {
			k.Logger(ctx).Error("unable to satisfy further unbondings this epoch")
			// do not process this or subsequent withdrawals this epoch.
			return true
		}

		// increment total to withdraw by the withdrawal amount
		totalToWithdraw = totalToWithdraw.Add(withdrawal.Amount[0])

		// set per withdrawal amount
		amountToWithdrawPerWithdrawal[withdrawal.Txhash] = withdrawal.Amount[0]

		// initialise empty distribution slice per withdrawal
		distributionsPerWithdrawal[withdrawal.Txhash] = make([]*types.Distribution, 0)
		return false
	})

	// no undelegations to attempt
	if len(amountToWithdrawPerWithdrawal) == 0 {
		return nil
	}

	tokensAllocatedForWithdrawalPerValidator, err := k.DeterminePlanForUndelegation(ctx, zone, sdk.NewCoins(totalToWithdraw))
	if err != nil {
		return err
	}
	valopers := utils.Keys(tokensAllocatedForWithdrawalPerValidator)
	// set current source validator to zero.
	vidx := 0
	v := valopers[vidx]
WITHDRAWAL:
	for _, hash := range utils.Keys(amountToWithdrawPerWithdrawal) {
		for {
			// if amountToWithdrawPerWithdrawal has been satisified, then continue.
			if amountToWithdrawPerWithdrawal[hash].Amount.IsZero() {
				continue WITHDRAWAL
			}

			// if current selected validator allocation for withdrawal can satisfy this withdrawal in totality...
			if tokensAllocatedForWithdrawalPerValidator[v].GTE(amountToWithdrawPerWithdrawal[hash].Amount) {
				// sub current withdrawal amount from allocation.
				tokensAllocatedForWithdrawalPerValidator[v] = tokensAllocatedForWithdrawalPerValidator[v].Sub(amountToWithdrawPerWithdrawal[hash].Amount)
				// create a distribution from this validator for the withdrawal
				distributionsPerWithdrawal[hash] = append(distributionsPerWithdrawal[hash], &types.Distribution{Valoper: v, Amount: amountToWithdrawPerWithdrawal[hash].Amount.Uint64()})

				// add the amount and hash to per validator records
				existing, found := coinsOutPerValidator[v]
				if !found {
					coinsOutPerValidator[v] = amountToWithdrawPerWithdrawal[hash]
					txHashesPerValidator[v] = []string{hash}

				} else {
					coinsOutPerValidator[v] = existing.Add(amountToWithdrawPerWithdrawal[hash])
					txHashesPerValidator[v] = append(txHashesPerValidator[v], hash)
				}

				// set withdrawal amount to zero, and continue to outer loop (next withdrawal record).
				amountToWithdrawPerWithdrawal[hash] = sdk.NewCoin(amountToWithdrawPerWithdrawal[hash].Denom, sdk.ZeroInt())
				continue WITHDRAWAL
			}

			// otherwise (current validator allocation cannot wholly satisfy current record), allocate entire allocation to this withdrawal.
			distributionsPerWithdrawal[hash] = append(distributionsPerWithdrawal[hash], &types.Distribution{Valoper: v, Amount: tokensAllocatedForWithdrawalPerValidator[v].Uint64()})
			amountToWithdrawPerWithdrawal[hash] = sdk.NewCoin(amountToWithdrawPerWithdrawal[hash].Denom, amountToWithdrawPerWithdrawal[hash].Amount.Sub(tokensAllocatedForWithdrawalPerValidator[v]))
			existing, found := coinsOutPerValidator[v]
			if !found {
				coinsOutPerValidator[v] = sdk.NewCoin(zone.BaseDenom, tokensAllocatedForWithdrawalPerValidator[v])
				txHashesPerValidator[v] = []string{hash}
			} else {
				coinsOutPerValidator[v] = existing.Add(sdk.NewCoin(zone.BaseDenom, tokensAllocatedForWithdrawalPerValidator[v]))
				txHashesPerValidator[v] = append(txHashesPerValidator[v], hash)
			}

			// set current val to zero.
			tokensAllocatedForWithdrawalPerValidator[v] = sdk.ZeroInt()
			// next validator
			if len(valopers) > vidx+1 {
				vidx++
				v = valopers[vidx]
			} else if !amountToWithdrawPerWithdrawal[hash].Amount.IsZero() {
				return fmt.Errorf("unable to satisfy unbonding")
			}
		}
	}

	for _, hash := range utils.Keys(distributionsPerWithdrawal) {
		record, found := k.GetWithdrawalRecord(ctx, zone.ChainId, hash, types.WithdrawStatusQueued)
		if !found {
			return errors.New("unable to find withdrawal record")
		}
		record.Distribution = distributionsPerWithdrawal[hash]
		k.UpdateWithdrawalRecordStatus(ctx, &record, types.WithdrawStatusUnbond)
	}

	if len(txHashesPerValidator) == 0 {
		// no records to handle.
		return nil
	}

	var msgs []sdk.Msg
	for _, valoper := range utils.Keys(coinsOutPerValidator) {
		if !coinsOutPerValidator[valoper].Amount.IsZero() {
			msgs = append(msgs, &stakingtypes.MsgUndelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: valoper, Amount: coinsOutPerValidator[valoper]})
		}
	}

	k.Logger(ctx).Info("unbonding messages to send", "msg", msgs)

	err = k.SubmitTx(ctx, msgs, zone.DelegationAddress, types.EpochWithdrawalMemo(epoch), zone.MessagesPerTx)
	if err != nil {
		return err
	}

	for _, valoper := range utils.Keys(coinsOutPerValidator) {
		if !coinsOutPerValidator[valoper].Amount.IsZero() {
			sort.Strings(txHashesPerValidator[valoper])
			k.SetUnbondingRecord(ctx, types.UnbondingRecord{ChainId: zone.ChainId, EpochNumber: epoch, Validator: valoper, RelatedTxhash: txHashesPerValidator[valoper]})
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
	currentAllocations, currentSum, _, _ := k.GetDelegationMap(ctx, zone.ChainId)
	availablePerValidator, _, err := k.GetUnlockedTokensForZone(ctx, zone)
	if err != nil {
		return nil, err
	}
	targetAllocations, err := k.GetAggregateIntentOrDefault(ctx, zone)
	if err != nil {
		return nil, err
	}
	allocations := types.DetermineAllocationsForUndelegation(currentAllocations, map[string]bool{}, currentSum, targetAllocations, availablePerValidator, amount)
	return allocations, nil
}
