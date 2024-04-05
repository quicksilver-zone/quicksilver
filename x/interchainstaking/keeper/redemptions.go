package keeper

import (
	"errors"
	"fmt"
	"sort"
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// processRedemptionForLsm will determine based on user intent, the tokens to return to the user, generate Redeem message and send them.
// func (k *Keeper) processRedemptionForLsm(ctx sdk.Context, zone *types.Zone, sender sdk.AccAddress, destination string, nativeTokens sdkmath.Int, burnAmount sdk.Coin, hash string) error {
// 	intent, found := k.GetDelegatorIntent(ctx, zone, sender.String(), false)
// 	// msgs is slice of MsgTokenizeShares, so we can handle dust allocation later.
// 	msgs := make([]*lsmstakingtypes.MsgTokenizeShares, 0)
// 	var err error
// 	intents := intent.Intents

// 	if !found || len(intents) == 0 {
// 		// if user has no intent set (this can happen if redeeming tokens that were obtained offchain), use global intent.
// 		// Note: this can be improved; user will receive a bunch of tokens.
// 		intents, err = k.GetAggregateIntentOrDefault(ctx, zone)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	outstanding := nativeTokens
// 	distribution := make(map[string]uint64, 0)

// 	availablePerValidator, _, err := k.GetUnlockedTokensForZone(ctx, zone)
// 	if err != nil {
// 		return err
// 	}
// 	for _, intent := range intents.Sort() {
// 		thisAmount := intent.Weight.MulInt(nativeTokens).TruncateInt()
// 		if thisAmount.GT(availablePerValidator[intent.ValoperAddress]) {
// 			return errors.New("unable to satisfy unbond request; delegations may be locked")
// 		}
// 		distribution[intent.ValoperAddress] = thisAmount.Uint64()
// 		outstanding = outstanding.Sub(thisAmount)
// 	}

// 	distribution[intents[0].ValoperAddress] += outstanding.Uint64()

// 	if distribution[intents[0].ValoperAddress] > availablePerValidator[intents[0].ValoperAddress].Uint64() {
// 		return errors.New("unable to satisfy unbond request (2); delegations may be locked")
// 	}

// 	for _, valoper := range utils.Keys(distribution) {
// 		msgs = append(msgs, &lsmstakingtypes.MsgTokenizeShares{
// 			DelegatorAddress:    zone.DelegationAddress.Address,
// 			ValidatorAddress:    valoper,
// 			Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewIntFromUint64(distribution[valoper])),
// 			TokenizedShareOwner: destination,
// 		})
// 	}

// 	sdkMsgs := make([]sdk.Msg, 0)
// 	for _, msg := range msgs {
// 		sdkMsgs = append(sdkMsgs, sdk.Msg(msg))
// 	}
// 	distributions := make([]*types.Distribution, 0)

// 	for valoper, amount := range distribution {
// 		newDistribution := types.Distribution{
// 			Valoper: valoper,
// 			Amount:  amount,
// 		}
// 		distributions = append(distributions, &newDistribution)
// 	}

// 	k.AddWithdrawalRecord(ctx, zone.ChainId, sender.String(), distributions, destination, sdk.Coins{}, burnAmount, hash, types.WithdrawStatusTokenize, time.Unix(0, 0), k.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch)

// 	return k.SubmitTx(ctx, sdkMsgs, zone.DelegationAddress, hash, zone.MessagesPerTx)
// }

// queueRedemption will determine based on zone intent, the tokens to unbond, and add a withdrawal record with status QUEUED.
func (k *Keeper) queueRedemption(
	ctx sdk.Context,
	zone *types.Zone,
	sender sdk.AccAddress,
	destination string,
	burnAmount sdk.Coin,
	hash string,
) error { //nolint:unparam // we know that the error is always nil
	distributions := make([]*types.Distribution, 0)

	_ = k.AddWithdrawalRecord(
		ctx,
		zone.ChainId,
		sender.String(),
		distributions,
		destination,
		burnAmount,
		hash,
		types.WithdrawStatusQueued,
		time.Time{},
		k.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch,
	)

	return nil
}

// GetUnlockedTokensForZone will iterate over all validators for a zone, summing delegated amounts,
// and then remove the locked tokens (those actively being redelegated), returning a slice of int64
// staking tokens that are unlocked and free to redelegate or unbond.
func (k *Keeper) GetUnlockedTokensForZone(ctx sdk.Context, zone *types.Zone) (map[string]sdkmath.Int, sdkmath.Int, error) {
	validators := k.GetValidators(ctx, zone.ChainId)

	availablePerValidator := make(map[string]sdkmath.Int, len(validators))
	total := sdk.ZeroInt()
	// for each validator, fetch delegated amount.
	for _, validator := range validators {
		delegation, found := k.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, validator.ValoperAddress)
		if !found {
			availablePerValidator[validator.ValoperAddress] = sdk.ZeroInt()
		} else {
			availablePerValidator[validator.ValoperAddress] = delegation.Amount.Amount
			total = total.Add(delegation.Amount.Amount)
		}
	}

	// for each redelegation, remove the amount being redelegated to from the destination,
	// as this cannot be available for unbonding or redelegation.
	for _, redelegation := range k.ZoneRedelegationRecords(ctx, zone.ChainId) {
		thisAvailable, found := availablePerValidator[redelegation.Destination]
		if found {
			availablePerValidator[redelegation.Destination] = thisAvailable.Sub(redelegation.Amount)
			if availablePerValidator[redelegation.Destination].LT(sdk.ZeroInt()) {
				return map[string]sdkmath.Int{}, sdk.ZeroInt(), fmt.Errorf("negative available amount [chain: %s, validator: %s, amount: %s]; unable to continue", zone.ChainId, redelegation.Destination, availablePerValidator[redelegation.Destination].String())
			}
			total = total.Sub(redelegation.Amount)
		}
	}

	return availablePerValidator, total, nil
}

// HandleQueuedUnbondings is called once per epoch to aggregate all queued unbondings into
// a single unbond transaction per delegation.
func (k *Keeper) HandleQueuedUnbondings(ctx sdk.Context, zone *types.Zone, epoch int64) error {
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

	// get min of LastRedemptionRate (N-1) and RedemptionRate (N)
	rate := sdk.MinDec(zone.LastRedemptionRate, zone.RedemptionRate)

	// iterate all withdrawal records for the zone in the QUEUED state.
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusQueued, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Info("handling queued withdrawal request", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", withdrawal.Amount)

		nativeTokens := sdk.NewDecFromInt(withdrawal.BurnAmount.Amount).Mul(rate).TruncateInt()
		amount := sdk.NewCoin(zone.BaseDenom, nativeTokens)
		k.Logger(ctx).Info("tokens to distribute", "amount", amount)

		if !amount.IsPositive() {
			k.Logger(ctx).Error("withdrawal %s attempting to withdraw non-positive amount; cannot process...", withdrawal.Txhash)
			return false
		}

		withdrawal.Amount = sdk.NewCoins(amount)
		err = k.SetWithdrawalRecord(ctx, withdrawal)
		if err != nil {
			k.Logger(ctx).Error("unable to set withdrawal record", "error", err)
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

	coinsOutPerValidator, txHashesPerValidator, distributionsPerWithdrawal, err := AllocateWithdrawalsFromValidators(zone.BaseDenom, tokensAllocatedForWithdrawalPerValidator, amountToWithdrawPerWithdrawal, distributionsPerWithdrawal)
	if err != nil {
		return err
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

			k.EventManagerKeeper.AddEvent(ctx,
				types.ModuleName,
				zone.ChainId,
				fmt.Sprintf("%s/%s", types.EpochWithdrawalMemo(epoch), valoper),
				"unbondAck",
				emtypes.EventTypeICAUnbond,
				emtypes.EventStatusActive,
				nil,
				nil,
			)
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
			k.SetUnbondingRecord(ctx, types.UnbondingRecord{ChainId: zone.ChainId, EpochNumber: epoch, Validator: valoper, RelatedTxhash: txHashesPerValidator[valoper], Amount: coinsOutPerValidator[valoper]})
		}
	}

	// if err = zone.IncrementWithdrawalWaitgroup(k.Logger(ctx), uint32(len(msgs)), "trigger unbonding messages"); err != nil {
	// 	return err
	// }
	//k.SetZone(ctx, zone)

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

func (k *Keeper) DeterminePlanForUndelegation(ctx sdk.Context, zone *types.Zone, amount sdk.Coins) (map[string]sdkmath.Int, error) {
	currentAllocations, currentSum, _, _ := k.GetDelegationMap(ctx, zone.ChainId)
	availablePerValidator, _, err := k.GetUnlockedTokensForZone(ctx, zone)
	if err != nil {
		return nil, err
	}
	targetAllocations, err := k.GetAggregateIntentOrDefault(ctx, zone)
	if err != nil {
		return nil, err
	}
	return types.DetermineAllocationsForUndelegation(currentAllocations, map[string]bool{}, currentSum, targetAllocations, availablePerValidator, amount)
}

// AllocateWithdrawalsFromValidators, given a mao of tokens that can be withdrawn from validators
// and a map of withdrawal records, distributes one to the other.
//
// Returns: map of coins removed from each val, map of withdrawal hashes to allocate to each unbonding message, map of distributions for each withdrawal.
func AllocateWithdrawalsFromValidators(
	denom string,
	tokensAllocatedForWithdrawalPerValidator map[string]sdkmath.Int, // map of amounts that can be unbonded from each val
	amountToWithdrawPerWithdrawal map[string]sdk.Coin, // map of amounts to withdraw per queued withdrawal_record
	distributionsPerWithdrawal map[string][]*types.Distribution, // empty map of distributions
) (
	map[string]sdk.Coin, // map of coins to be removed from each val (does this just end up matching tokensAllocatedForWithdrawalPerValidator?)
	map[string][]string, // map of withdrawal_records txhashes allocated to each unbonding
	map[string][]*types.Distribution, // filled map of distributions
	error,
) {
	_amountToWithdrawPerWithdrawal := make(map[string]sdk.Coin, len(amountToWithdrawPerWithdrawal))
	_tokensAllocatedForWithdrawalPerValidator := make(map[string]sdkmath.Int, len(tokensAllocatedForWithdrawalPerValidator))
	for k, v := range amountToWithdrawPerWithdrawal {
		_amountToWithdrawPerWithdrawal[k] = v
	}
	for k, v := range tokensAllocatedForWithdrawalPerValidator {
		_tokensAllocatedForWithdrawalPerValidator[k] = v
	}

	// out here will only ever be in native bond denom
	coinsOutPerValidator := make(map[string]sdk.Coin, 0)
	// list of withdrawal tx hashes per validator
	txHashesPerValidator := make(map[string][]string, 0)

	valopers := utils.Keys(tokensAllocatedForWithdrawalPerValidator)
	// set current source validator to zero.
	vidx := 0
	v := valopers[vidx]
WITHDRAWAL:
	for _, hash := range utils.Keys(amountToWithdrawPerWithdrawal) {
		for {
			// if amountToWithdrawPerWithdrawal has been satisified, then continue.
			if amountToWithdrawPerWithdrawal[hash].IsZero() {
				continue WITHDRAWAL
			}

			// if current selected validator allocation for withdrawal can satisfy this withdrawal in totality...
			if tokensAllocatedForWithdrawalPerValidator[v].GTE(amountToWithdrawPerWithdrawal[hash].Amount) {
				// sub current withdrawal amount from allocation.
				tokensAllocatedForWithdrawalPerValidator[v] = tokensAllocatedForWithdrawalPerValidator[v].Sub(amountToWithdrawPerWithdrawal[hash].Amount)
				// create a distribution from this validator for the withdrawal
				distributionsPerWithdrawal[hash] = append(distributionsPerWithdrawal[hash], &types.Distribution{Valoper: v, Amount: amountToWithdrawPerWithdrawal[hash].Amount})

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
			distributionsPerWithdrawal[hash] = append(distributionsPerWithdrawal[hash], &types.Distribution{Valoper: v, Amount: tokensAllocatedForWithdrawalPerValidator[v]})
			amountToWithdrawPerWithdrawal[hash] = sdk.NewCoin(amountToWithdrawPerWithdrawal[hash].Denom, amountToWithdrawPerWithdrawal[hash].Amount.Sub(tokensAllocatedForWithdrawalPerValidator[v]))
			existing, found := coinsOutPerValidator[v]
			if !found {
				coinsOutPerValidator[v] = sdk.NewCoin(denom, tokensAllocatedForWithdrawalPerValidator[v])
				txHashesPerValidator[v] = []string{hash}
			} else {
				coinsOutPerValidator[v] = existing.Add(sdk.NewCoin(denom, tokensAllocatedForWithdrawalPerValidator[v]))
				txHashesPerValidator[v] = append(txHashesPerValidator[v], hash)
			}

			// set current val to zero.
			tokensAllocatedForWithdrawalPerValidator[v] = sdk.ZeroInt()
			// next validator
			if len(valopers) > vidx+1 {
				vidx++
				v = valopers[vidx]
			} else if !amountToWithdrawPerWithdrawal[hash].IsZero() {
				return nil, nil, nil, fmt.Errorf("unable to satisfy unbonding")
			}
		}
	}

	// sanity checks
	sumOut := sdk.NewCoin(denom, sdkmath.ZeroInt())
	for valoper, coinPerVal := range coinsOutPerValidator {
		if !coinPerVal.Amount.Equal(_tokensAllocatedForWithdrawalPerValidator[valoper]) {
			return nil, nil, nil, fmt.Errorf("allocation <-> coinOut mismatch for %s; in = %v, out = %v", valoper, _tokensAllocatedForWithdrawalPerValidator[valoper], coinPerVal)
		}
		sumOut = sumOut.Add(coinPerVal)
	}

	sumIn := sdk.NewCoin(denom, sdkmath.ZeroInt())
	for hash, tx := range _amountToWithdrawPerWithdrawal {
		sumIn = sumIn.Add(tx)
		dist := func(in []*types.Distribution) sdkmath.Int {
			sum := sdkmath.ZeroInt()
			for _, dist := range in {
				sum = sum.Add(dist.Amount)
			}
			return sum
		}(distributionsPerWithdrawal[hash])

		if !tx.Amount.Equal(dist) {
			return nil, nil, nil, fmt.Errorf("amountToWithdrawPerWithdrawal <-> distributionsPerWithdrawal mismatch for %s; tx = %v, dist = %v", hash, tx, dist)
		}
	}

	if !sumIn.Equal(sumOut) {
		return nil, nil, nil, fmt.Errorf("sumIn <-> sumOut mismatch; sumIn = %s, sumOut = %s", sumIn.Amount.String(), sumOut.Amount.String())
	}

	return coinsOutPerValidator, txHashesPerValidator, distributionsPerWithdrawal, nil
}
