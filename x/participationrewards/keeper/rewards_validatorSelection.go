package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// zoneScore is an internal struct to track transient state for the calculation
// of zone scores. It specifically tallies the total zone voting power used in
// calculations to determine validator voting power percentages;
type zoneScore struct {
	ZoneId           string // chainId
	TotalVotingPower sdk.Int
	ValidatorScores  map[string]*validator
}

// validator is an internal struct to track transient state for the calculation
// of zone scores. It contains all relevant validator scoring metrics with a
// pointer reference to the actual validator (embedded).
type validator struct {
	PowerPercentage   sdk.Dec
	PerformanceScore  sdk.Dec
	DistributionScore sdk.Dec

	*icstypes.Validator
}

// userAllocation is an internal struct to track transient state for rewards
// distribution. It contains the user address and the coins that are allocated
// to it.
type userAllocation struct {
	Address string
	Coins   sdk.Coins
}

// allocateValidatorSelectionRewards utilizes IBC to query the performance
// rewards account for each zone to determine validator performance and
// corresponding rewards allocations. Each zone's response is dealt with
// individually in a callback.
func (k Keeper) allocateValidatorSelectionRewards(ctx sdk.Context) error {
	k.Logger(ctx).Info("allocateValidatorChoiceRewards")

	for i, zone := range k.icsKeeper.AllRegisteredZones(ctx) {
		k.Logger(ctx).Info("zones", "i", i, "zone", zone.ChainId, "performance address", zone.PerformanceAddress.GetAddress())

		// obtain zone performance account rewards
		rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: zone.PerformanceAddress.GetAddress()}
		bz := k.cdc.MustMarshal(&rewardsQuery)

		k.icqKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			"cosmos.distribution.v1beta1.Query/DelegationTotalRewards",
			bz,
			sdk.NewInt(-1),
			types.ModuleName,
			"validatorselectionrewards",
			0,
		)
	}

	return nil
}

// getZoneScores returns an instance of zoneScore containing the calculated
// zone validator scores.
func (k Keeper) getZoneScores(
	ctx sdk.Context,
	zone icstypes.RegisteredZone,
	delegatorRewards distrtypes.QueryDelegationTotalRewardsResponse,
) (*zoneScore, error) {
	k.Logger(ctx).Info(
		"performance rewards zone callback response",
		"zone", zone.ChainId,
		"rewards", delegatorRewards,
	)

	zs := zoneScore{
		ZoneId:           zone.ChainId,
		TotalVotingPower: sdk.NewInt(0),
		ValidatorScores:  make(map[string]*validator),
	}

	if err := k.calcDistributionScores(ctx, zone, &zs); err != nil {
		return nil, err
	}

	if err := k.calcOverallScores(ctx, zone, delegatorRewards, &zs); err != nil {
		return nil, err
	}

	return &zs, nil
}

// calcDistributionScores calculates the validator distribution scores for the
// given zone based on the normalized voting power of the validators; scoring
// favours smaller validators for decentraliztion purposes.
func (k Keeper) calcDistributionScores(ctx sdk.Context, zone icstypes.RegisteredZone, zs *zoneScore) error {
	k.Logger(ctx).Info("calculate distribution scores", "zone", zone.ChainId)

	zoneValidators := zone.GetValidatorsSorted()
	if len(zoneValidators) == 0 {
		return fmt.Errorf("zone %v has no validators", zone.ChainId)
	}

	// calculate total voting power
	// and determine min/max voting power for zone
	max := sdk.NewInt(0)
	min := sdk.NewInt(999999999999999999)
	for _, val := range zoneValidators {
		// compute zone total voting power
		zs.TotalVotingPower = zs.TotalVotingPower.Add(val.VotingPower)
		if _, exists := zs.ValidatorScores[val.ValoperAddress]; !exists {
			zs.ValidatorScores[val.ValoperAddress] = &validator{Validator: val}
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
		k.Logger(ctx).Error("zone invalid, zero voting power", "zone", zone)
		panic("this should never happen!")
	}

	// calculate power percentage and normalized distribution scores
	maxp := sdk.NewDecFromInt(max).Quo(sdk.NewDecFromInt(zs.TotalVotingPower))
	minp := sdk.NewDecFromInt(min).Quo(sdk.NewDecFromInt(zs.TotalVotingPower))
	for _, vs := range zs.ValidatorScores {
		// calculate power percentage
		vs.PowerPercentage = sdk.NewDecFromInt(vs.VotingPower).Quo(sdk.NewDecFromInt(vs.VotingPower))

		// calculate normalized distribution score
		vs.DistributionScore = sdk.NewDec(1).Sub(
			vs.PowerPercentage.Sub(minp).Mul(
				sdk.NewDec(1).Quo(maxp),
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

// calcOverallScores calculates the overall validator scores for the given zone
// based on the combination of performance score and distribution score.
//
// The performance score is first calculated based on validator rewards earned
// from the zone performance account that delegates an exact amount to each
// validator. The total rewards earned by the performance account is divided
// by the number of active validators to obtain the expected rewards. The
// performance score for each validator is then simply the percentage of actual
// rewards compared to the expected rewards (capped at 100%).
//
// On completetion a msg is submitted to withdraw the zone performance rewards,
// resetting zone performance scoring for the next epoch.
func (k Keeper) calcOverallScores(
	ctx sdk.Context,
	zone icstypes.RegisteredZone,
	delegatorRewards distrtypes.QueryDelegationTotalRewardsResponse,
	zs *zoneScore,
) error {
	k.Logger(ctx).Info("calculate performance & overall scores")

	rewards := delegatorRewards.GetRewards()
	if rewards == nil {
		k.Logger(ctx).Error("No delegator rewards")
		return nil
	}

	total := delegatorRewards.GetTotal().AmountOf(zone.BaseDenom)
	expected := total.Quo(sdk.NewDec(int64(len(rewards))))

	k.Logger(ctx).Info(
		"performance account rewards",
		"rewards", rewards,
		"total", total,
		"expected", expected,
	)

	var msgs []sdk.Msg
	limit := sdk.NewDec(1.0)
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

		// calculate overall score
		vs.Score = vs.DistributionScore.Mul(vs.PerformanceScore)
		k.Logger(ctx).Info("overall score", "validator", vs.ValoperAddress, "overall", vs.Score)

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
		if err := k.icsKeeper.SubmitTx(ctx, msgs, zone.PerformanceAddress, ""); err != nil {
			return err
		}
	}

	// update zone with validator scores
	k.icsKeeper.SetRegisteredZone(ctx, zone)

	return nil
}

// calcUserValidatorSelectionAllocations returns a slice of userAllocation. It
// calculates individual user scores relative to overall zone score and then
// proportionally allocates rewards based on the individual zone allocation.
func (k Keeper) calcUserValidatorSelectionAllocations(
	ctx sdk.Context,
	zone icstypes.RegisteredZone,
	zs zoneScore,
) ([]userAllocation, error) {
	k.Logger(ctx).Info("calcUserAllocations", "zone", zone.ChainId, "scores", zs, "allocations", zone.ValidatorSelectionAllocation)

	userAllocations := make([]userAllocation, 0)

	if zone.ValidatorSelectionAllocation.IsZero() {
		k.Logger(ctx).Info("validator selection allocation is zero, nothing to allocate")
		return userAllocations, nil
	}

	type userScore struct {
		address string
		score   sdk.Dec
	}

	sum := sdk.NewDec(0)
	userScores := make([]userScore, 0)
	// obtain snapshotted intents of last epoch boundary
	for _, di := range k.icsKeeper.AllOrdinalizedIntents(ctx, zone, true) {
		uSum := sdk.NewDec(0)
		for _, intent := range di.GetIntents() {
			// calc overall user score
			score := sdk.ZeroDec()
			if vs, exists := zs.ValidatorScores[intent.ValoperAddress]; exists {
				if !vs.Score.IsNil() {
					score = intent.Weight.Mul(vs.Score)
				}
			}
			k.Logger(ctx).Info("user score for validator", "user", di.GetDelegator(), "validator", intent.ValoperAddress, "score", score)
			uSum = uSum.Add(score)
		}
		u := userScore{
			address: di.GetDelegator(),
			score:   uSum,
		}
		k.Logger(ctx).Info("user score for zone", "user", di.GetDelegator(), "zone", zs.ZoneId, "score", uSum)
		userScores = append(userScores, u)
		// calc overall zone score
		sum = sum.Add(uSum)
	}

	if sum.IsZero() {
		k.Logger(ctx).Info("zero sum score for zone", "zone", zone.ChainId)
		return userAllocations, nil
	}

	tokensPerPoint := sdk.NewDecFromInt(zone.ValidatorSelectionAllocation.AmountOfNoDenomValidation(k.stakingKeeper.BondDenom(ctx))).Quo(sum)
	k.Logger(ctx).Info("tokens per point", "zone", zs.ZoneId, "zone score", sum, "tpp", tokensPerPoint)
	for _, us := range userScores {
		ua := userAllocation{
			Address: us.address,
			Coins: sdk.NewCoins(
				sdk.NewCoin(
					k.stakingKeeper.BondDenom(ctx),
					us.score.Mul(tokensPerPoint).TruncateInt(),
				),
			),
		}
		userAllocations = append(userAllocations, ua)
	}

	return userAllocations, nil
}
