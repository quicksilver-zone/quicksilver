package keeper

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

// AllocateValidatorSelectionRewards utilizes IBC to query the performance
// rewards account for each zone to determine validator performance and
// corresponding rewards allocations. Each zone's response is dealt with
// individually in a callback.
func (k Keeper) AllocateValidatorSelectionRewards(ctx sdk.Context) {
	k.icsKeeper.IterateZones(ctx, func(_ int64, zone *icstypes.Zone) (stop bool) {
		if zone.PerformanceAddress != nil {
			k.Logger(ctx).Info("zones", "chain_id", zone.ChainId, "performance address", zone.PerformanceAddress.Address)

			// obtain zone performance account rewards
			rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: zone.PerformanceAddress.Address}
			bz := k.cdc.MustMarshal(&rewardsQuery)

			k.IcqKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				"cosmos.distribution.v1beta1.Query/DelegationTotalRewards",
				bz,
				sdkmath.NewInt(-1),
				types.ModuleName,
				ValidatorSelectionRewardsCallbackID,
				0,
			)
		}
		return false
	})
}

// getZoneScores returns an instance of zoneScore containing the calculated
// zone validator scores.
func (k Keeper) getZoneScores(
	ctx sdk.Context,
	zone icstypes.Zone,
	delegatorRewards distrtypes.QueryDelegationTotalRewardsResponse,
) (*types.ZoneScore, error) {
	k.Logger(ctx).Info(
		"performance rewards zone callback response",
		"zone", zone.ChainId,
		"rewards", delegatorRewards,
	)

	zs := types.ZoneScore{
		ZoneID:           zone.ChainId,
		TotalVotingPower: sdkmath.NewInt(0),
		ValidatorScores:  make(map[string]*types.Validator),
	}

	if err := k.CalcDistributionScores(ctx, zone, &zs); err != nil {
		return nil, err
	}

	if err := k.CalcOverallScores(ctx, zone, delegatorRewards, &zs); err != nil {
		return nil, err
	}

	return &zs, nil
}

// CalcDistributionScores calculates the validator distribution scores for the
// given zone based on the normalized voting power of the validators; scoring
// favours smaller validators for decentraliztion purposes.
func (k Keeper) CalcDistributionScores(ctx sdk.Context, zone icstypes.Zone, zs *types.ZoneScore) error {
	k.Logger(ctx).Info("calculate distribution scores", "zone", zone.ChainId)

	zoneValidators := k.icsKeeper.GetValidators(ctx, zone.ChainId)
	if len(zoneValidators) == 0 {
		return fmt.Errorf("zone %v has no validators", zone.ChainId)
	}

	// calculate total voting power
	// and determine min/max voting power for zone
	max := sdkmath.NewInt(0)
	min := sdkmath.NewInt(999999999999999999)
	for _, zoneVal := range zoneValidators {
		val := zoneVal
		if val.VotingPower.IsNegative() {
			return fmt.Errorf("unexpected negative voting power for %s", val.ValoperAddress)
		}
		// compute zone total voting power
		zs.TotalVotingPower = zs.TotalVotingPower.Add(val.VotingPower)
		if _, exists := zs.ValidatorScores[val.ValoperAddress]; !exists {
			zs.ValidatorScores[val.ValoperAddress] = &types.Validator{Validator: &val}
		}

		// Set max/min
		if max.LT(val.VotingPower) {
			max = val.VotingPower
			k.Logger(ctx).Info("new power max", "max", max, "validator", val.ValoperAddress)
		}
		if min.GT(val.VotingPower) {
			min = val.VotingPower
			k.Logger(ctx).Info("new power min", "min", min, "validator", val.ValoperAddress)
		}
	}

	k.Logger(ctx).Info("zone voting power", "zone", zone.ChainId, "total voting power", zs.TotalVotingPower)

	if zs.TotalVotingPower.IsZero() {
		err := errors.New("invalid zone, zero voting power")
		k.Logger(ctx).Error(err.Error(), "zone", zone)
		return err
	}

	// calculate power percentage and normalized distribution scores
	maxp := sdkmath.LegacyNewDecFromInt(max).Quo(sdkmath.LegacyNewDecFromInt(zs.TotalVotingPower))
	minp := sdkmath.LegacyNewDecFromInt(min).Quo(sdkmath.LegacyNewDecFromInt(zs.TotalVotingPower))
	for _, vs := range zs.ValidatorScores {
		// calculate power percentage
		vs.PowerPercentage = sdkmath.LegacyNewDecFromInt(vs.VotingPower).Quo(sdkmath.LegacyNewDecFromInt(zs.TotalVotingPower))

		// calculate normalized distribution score
		vs.DistributionScore = sdkmath.LegacyNewDec(1).Sub(
			vs.PowerPercentage.Sub(minp).Mul(
				sdkmath.LegacyNewDec(1).Quo(maxp),
			),
		)

		k.Logger(ctx).Debug(
			"validator score",
			"validator", vs.ValoperAddress,
			"power percentage", vs.PowerPercentage,
			"distribution score", vs.DistributionScore,
		)
	}

	return nil
}

// CalcOverallScores calculates the overall validator scores for the given zone
// based on the combination of performance score and distribution score.
//
// The performance score is first calculated based on validator rewards earned
// from the zone performance account that delegates an exact amount to each
// validator. The total rewards earned by the performance account is divided
// by the number of active validators to obtain the expected rewards. The
// performance score for each validator is then simply the percentage of actual
// rewards compared to the expected rewards (capped at 100%).
//
// On completion a msg is submitted to withdraw the zone performance rewards,
// resetting zone performance scoring for the next epoch.
func (k Keeper) CalcOverallScores(
	ctx sdk.Context,
	zone icstypes.Zone,
	delegatorRewards distrtypes.QueryDelegationTotalRewardsResponse,
	zs *types.ZoneScore,
) error {
	k.Logger(ctx).Info("calculate performance & overall scores")

	rewards := delegatorRewards.GetRewards()
	if rewards == nil {
		return nil
	}

	total := delegatorRewards.GetTotal().AmountOf(zone.BaseDenom)

	if total.IsZero() {
		return nil
	}

	expected := total.Quo(sdkmath.LegacyNewDec(int64(len(rewards))))

	k.Logger(ctx).Info(
		"performance account rewards",
		"rewards", rewards,
		"total", total,
		"expected", expected,
	)

	msgs := make([]sdk.Msg, 0)
	limit := sdkmath.LegacyNewDec(1.0)
	for _, reward := range rewards {
		vs, exists := zs.ValidatorScores[reward.ValidatorAddress]
		if !exists {
			k.Logger(ctx).Info("validator may have been removed from active set", "validator", reward.ValidatorAddress)
			continue
		}

		vs.PerformanceScore = reward.Reward.AmountOf(zone.BaseDenom).Quo(expected)
		if vs.PerformanceScore.GT(limit) {
			vs.PerformanceScore = limit
		}
		k.Logger(ctx).Info("performance score", "validator", vs.ValoperAddress, "performance", vs.PerformanceScore)

		// calculate and set overall score
		vs.Score = vs.DistributionScore.Mul(vs.PerformanceScore)
		k.Logger(ctx).Info("overall score", "validator", vs.ValoperAddress, "overall", vs.Score)
		if err := k.icsKeeper.SetValidator(ctx, zone.ChainId, *(vs.Validator)); err != nil {
			k.Logger(ctx).Error("unable to set score for validator", "validator", vs.ValoperAddress)
		}

		// prepare validator performance withdrawal msg
		msg := &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: zone.PerformanceAddress.GetAddress(),
			ValidatorAddress: vs.ValoperAddress,
		}
		msgs = append(msgs, msg)
	}

	// submit rewards withdrawals to reset zone performance for next epoch
	k.Logger(ctx).Info("send performance rewards withdrawal messages to reset scores for next epoch")
	if len(msgs) > 0 {
		if err := k.icsKeeper.SubmitTx(ctx, msgs, zone.PerformanceAddress, "", zone.MessagesPerTx); err != nil {
			return err
		}
	}

	// update zone with validator scores
	k.icsKeeper.SetZone(ctx, &zone)

	return nil
}

// CalcUserValidatorSelectionAllocations returns a slice of userAllocation. It
// calculates individual user scores relative to overall zone score and then
// proportionally allocates rewards based on the individual zone allocation.
func (k Keeper) CalcUserValidatorSelectionAllocations(
	ctx sdk.Context,
	zone *icstypes.Zone,
	zs types.ZoneScore,
) []types.UserAllocation {
	k.Logger(ctx).Info("calcUserValidatorSelectionAllocations", "zone", zone.ChainId, "scores", zs, "allocation", zone.ValidatorSelectionAllocation)

	userAllocations := make([]types.UserAllocation, 0)
	if zone.ValidatorSelectionAllocation == 0 {
		k.Logger(ctx).Info("validator selection allocation is zero, nothing to allocate")
		return userAllocations
	}

	sum := sdkmath.LegacyNewDec(0)
	userScores := make([]types.UserScore, 0)
	// obtain snapshotted intents of last epoch boundary
	k.icsKeeper.IterateDelegatorIntents(ctx, zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
		uSum := sdkmath.LegacyNewDec(0)
		for _, intent := range di.GetIntents() {
			// calc overall user score
			score := sdk.ZeroDec()
			if vs, exists := zs.ValidatorScores[intent.ValoperAddress]; exists {
				if !vs.Score.IsNil() {
					score = intent.Weight.Mul(vs.Score)
				}
			}
			k.Logger(ctx).Info("user score for validator", "user", di.GetDelegator(), "validator", intent.GetValoperAddress(), "score", score)
			uSum = uSum.Add(score)
		}
		u := types.UserScore{
			Address: di.GetDelegator(),
			Score:   uSum,
		}
		k.Logger(ctx).Info("user score for zone", "user", di.GetDelegator(), "zone", zs.ZoneID, "score", uSum)
		userScores = append(userScores, u)
		// calc overall zone score
		sum = sum.Add(uSum)
		return false
	})

	if sum.IsZero() {
		k.Logger(ctx).Info("zero sum score for zone", "zone", zone.ChainId)
		return userAllocations
	}

	allocation := sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(zone.ValidatorSelectionAllocation))
	tokensPerPoint := allocation.Quo(sum)
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	k.Logger(ctx).Info("tokens per point", "zone", zs.ZoneID, "zone score", sum, "tpp", tokensPerPoint)
	for _, us := range userScores {
		ua := types.UserAllocation{
			Address: us.Address,
			Amount:  sdk.NewCoin(bondDenom, us.Score.Mul(tokensPerPoint).TruncateInt()),
		}
		userAllocations = append(userAllocations, ua)
	}

	return userAllocations
}
