package keeper

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	lsmstakingtypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"
)

// processRedemptionForLsm will determine based on user intent, the tokens to return to the user, generate Redeem message and send them.
func (k *Keeper) processRedemptionForLsm(ctx sdk.Context, zone types.Zone, sender sdk.AccAddress, destination string, nativeTokens math.Int, burnAmount sdk.Coin, hash string) error {
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

	availablePerValidator := k.GetUnlockedTokensForZone(ctx, &zone)

	for _, intent := range intents.Sort() {
		thisAmount := intent.Weight.MulInt(nativeTokens).TruncateInt()
		if thisAmount.Int64() > availablePerValidator[intent.ValoperAddress] {
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
	zone types.Zone,
	sender sdk.AccAddress,
	destination string,
	nativeTokens math.Int,
	burnAmount sdk.Coin,
	hash string,
) error { //nolint:unparam // we know that the error is always nil
	distribution := make([]*types.Distribution, 0)
	outstanding := nativeTokens

	aggregateIntent := zone.GetAggregateIntentOrDefault()
	for _, intent := range aggregateIntent {
		thisAmount := intent.Weight.MulInt(nativeTokens).TruncateInt()
		outstanding = outstanding.Sub(thisAmount)
		dist := types.Distribution{
			Valoper: intent.ValoperAddress,
			Amount:  thisAmount.Uint64(),
		}

		distribution = append(distribution, &dist)
	}
	// handle dust ? ok to do uint64 calc here or do we use math.Int (just more verbose) ?
	distribution[0].Amount += outstanding.Uint64()

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
func (k *Keeper) GetUnlockedTokensForZone(ctx sdk.Context, zone *types.Zone) map[string]int64 {
	availablePerValidator := map[string]int64{}
	for _, delegation := range k.GetAllDelegations(ctx, zone) {
		thisAvailable, found := availablePerValidator[delegation.ValidatorAddress]
		if !found {
			thisAvailable = 0
		}
		availablePerValidator[delegation.ValidatorAddress] = thisAvailable + delegation.Amount.Amount.Int64()
	}
	for _, redelegation := range k.ZoneRedelegationRecords(ctx, zone.ChainId) {
		thisAvailable, found := availablePerValidator[redelegation.Destination]
		if found {
			availablePerValidator[redelegation.Destination] = thisAvailable - redelegation.Amount
		}
	}
	return availablePerValidator
}

// handle queued unbondings is called once per epoch to aggregate all queued unbondings into
// a single unbond transaction per delegation.
func (k *Keeper) HandleQueuedUnbondings(ctx sdk.Context, zone *types.Zone, epoch int64) error {
	// out here will only ever be in native bond denom
	out := make(map[string]sdk.Coin, 0)
	txhashes := make(map[string][]string, 0)

	availablePerValidator := k.GetUnlockedTokensForZone(ctx, zone)

	var err error
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, WithdrawStatusQueued, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		// copy this so we can rollback on fail
		thisAvail := availablePerValidator
		thisOut := make(map[string]sdk.Coin, 0)
		k.Logger(ctx).Info("unbonding funds", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", withdrawal.Amount)
		for _, dist := range withdrawal.Distribution {
			if thisAvail[dist.Valoper] < int64(dist.Amount) {
				// we cannot satisfy this unbond this epoch.
				k.Logger(ctx).Error("unable to satisfy unbonding for this epoch, due to locked tokens.", "txhash", withdrawal.Txhash, "user", withdrawal.Delegator, "chain", zone.ChainId, "validator", dist.Valoper, "avail", thisAvail[dist.Valoper], "wanted", int64(dist.Amount))
				return false
			}
			thisOut[dist.Valoper] = sdk.NewCoin(zone.BaseDenom, math.NewIntFromUint64(dist.Amount))
			thisAvail[dist.Valoper] -= int64(dist.Amount)

			// if the validator has been historically slashed, and delegatorShares does not match tokens, then we end up with 'clipping'.
			// clipping is the truncation of the expected unbonding amount because of the need to have whole integer tokens.
			// the amount unbonded is emitted as an event, but not in the response, so we never _know_ this has happened.
			// as such, if we know the validator has hisotrical slashing, we remove 1 utoken from the distribution for this validator, with
			// the expectation that clipping will occur. We do not reduce the amount requested to unbond.
			val, found := zone.GetValidatorByValoper(dist.Valoper)
			if !found {
				// something kooky is going on...
				err = fmt.Errorf("unable to find a validator we expected to exist [%s]", dist.Valoper)
				return true
			}
			if !val.DelegatorShares.Equal(sdk.NewDecFromInt(val.VotingPower)) && dist.Amount > 0 {
				dist.Amount--
			}
		}

		// update record of available balances.
		availablePerValidator = thisAvail

		for valoper, amount := range thisOut {
			existing, found := out[valoper]
			if !found {
				out[valoper] = amount
				txhashes[valoper] = []string{withdrawal.Txhash}

			} else {
				out[valoper] = existing.Add(amount)
				txhashes[valoper] = append(txhashes[valoper], withdrawal.Txhash)

			}
		}

		k.UpdateWithdrawalRecordStatus(ctx, &withdrawal, WithdrawStatusUnbond)
		return false
	})
	if err != nil {
		return err
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
