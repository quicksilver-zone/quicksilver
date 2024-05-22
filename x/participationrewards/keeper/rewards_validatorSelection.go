package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

// QueryValidatorDelegationPerformance utilizes IBC to query the performance
// rewards account for each zone to determine validator performance and
// corresponding rewards allocations. Each zone's response is dealt with
// individually in a callback.
func (k Keeper) QueryValidatorDelegationPerformance(ctx sdk.Context, zone *icstypes.Zone) {
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
			sdk.NewInt(-1),
			types.ModuleName,
			ValidatorSelectionRewardsCallbackID,
			0,
		)
		k.EventManagerKeeper.AddEvent(ctx, types.ModuleName, zone.ChainId, "validator_performance", "", emtypes.EventTypeICQQueryDelegations, emtypes.EventStatusActive, nil, nil)
	}
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
		TotalVotingPower: sdk.NewInt(0),
		ValidatorScores:  make(map[string]*types.ValidatorScore),
	}

	if err := k.CalcDistributionScores(ctx, zone, &zs); err != nil {
		return nil, err
	}

	if err := k.CalcPerformanceScores(ctx, zone, delegatorRewards, &zs); err != nil {
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

	vps := map[string]math.Int{}
	for _, zoneVal := range zoneValidators {
		val := zoneVal
		if val.VotingPower.IsNegative() {
			continue
		}
		// compute zone total voting power
		zs.TotalVotingPower = zs.TotalVotingPower.Add(val.VotingPower)
		vps[val.ValoperAddress] = val.VotingPower
	}

	k.Logger(ctx).Info("zone voting power", "zone", zone.ChainId, "total voting power", zs.TotalVotingPower)

	if zs.TotalVotingPower.IsZero() {
		err := errors.New("invalid zone, zero voting power")
		k.Logger(ctx).Error(err.Error(), "zone", zone)
		return err
	}

	// calculate power percentage and normalized distribution scores
	for _, valoper := range utils.Keys[math.Int](vps) {
		vpPercent := sdk.NewDecFromInt(vps[valoper]).Quo(sdk.NewDecFromInt(zs.TotalVotingPower))
		zs.ValidatorScores[valoper] = &types.ValidatorScore{DistributionScore: sdk.NewDec(1).Quo(vpPercent)}

		k.Logger(ctx).Debug(
			"validator score",
			"validator", valoper,
			"power percentage", vpPercent,
			"distribution score", zs.ValidatorScores[valoper].DistributionScore,
		)
	}

	return nil
}

// CalcPerformanceScores calculates he overall validator scores for the given zone
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
func (k Keeper) CalcPerformanceScores(
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

	expected := total.Quo(sdk.NewDec(int64(len(rewards))))

	k.Logger(ctx).Info(
		"performance account rewards",
		"rewards", rewards,
		"total", total,
		"expected", expected,
	)

	maxScore := sdk.ZeroDec()

	msgs := make([]sdk.Msg, 0)
	for _, reward := range rewards {
		vs, exists := zs.ValidatorScores[reward.ValidatorAddress]
		if !exists {
			k.Logger(ctx).Info("validator may have been removed from active set", "validator", reward.ValidatorAddress)
			continue
		}

		rootScore := reward.Reward.AmountOf(zone.BaseDenom).Quo(expected)
		vs.PerformanceScore = rootScore.Mul(rootScore)
		if vs.PerformanceScore.GT(maxScore) {
			maxScore = vs.PerformanceScore
		}
	}

	for _, reward := range rewards {
		vs := zs.ValidatorScores[reward.ValidatorAddress]
		vs.PerformanceScore = vs.PerformanceScore.Quo(maxScore)
		k.Logger(ctx).Info("overall score", "validator", reward.ValidatorAddress, "distribution", vs.DistributionScore, "performance", vs.PerformanceScore, "total", vs.TotalScore())
		val, found := k.icsKeeper.GetValidator(ctx, zone.ChainId, addressutils.MustValAddressFromBech32(reward.ValidatorAddress, ""))
		if !found {
			k.Logger(ctx).Error("unable to find validator", "validator", reward.ValidatorAddress)
		} else {
			val.Score = vs.TotalScore()
			if err := k.icsKeeper.SetValidator(ctx, zone.ChainId, val); err != nil {
				k.Logger(ctx).Error("unable to set score for validator", "validator", reward.ValidatorAddress)
			}
		}

		// prepare validator performance withdrawal msg
		msg := &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: zone.PerformanceAddress.GetAddress(),
			ValidatorAddress: reward.ValidatorAddress,
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

	sum := sdk.NewDec(0)
	userScores := make([]types.UserScore, 0)
	// obtain snapshotted intents of last epoch boundary
	k.icsKeeper.IterateDelegatorIntents(ctx, zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
		uSum := sdk.NewDec(0)
		for _, intent := range di.GetIntents() {
			// calc overall user score
			score := sdk.ZeroDec()
			if vs, exists := zs.ValidatorScores[intent.ValoperAddress]; exists {
				if !vs.TotalScore().IsNil() {
					score = intent.Weight.Mul(vs.TotalScore())
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

	allocation := sdk.NewDecFromInt(sdk.NewIntFromUint64(zone.ValidatorSelectionAllocation))
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
