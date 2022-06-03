package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// ? callback concurrency issues on map ?
type zoneScore struct {
	zoneId           string // chainId
	totalVotingPower sdk.Int
	validatorScores  map[string]*validator
}

type validator struct {
	powerPercentage   sdk.Dec
	performanceScore  sdk.Dec
	distributionScore sdk.Dec
	overallScore      sdk.Dec

	*icstypes.Validator
}

func (k Keeper) allocateValidatorSelectionRewards(ctx sdk.Context, allocation sdk.Coins) error {
	k.Logger(ctx).Info("allocateValidatorChoiceRewards", "allocation", allocation)

	var rewardscb Callback = func(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
		zs, err := k.zoneCallback(ctx, response, query)
		if err != nil {
			return err
		}

		k.Logger(ctx).Info(
			"Callback Zone Score",
			"zone", zs.zoneId,
			"total voting power", zs.totalVotingPower,
			"validator scores", zs.validatorScores,
		)

		// TODO: distribute zone allocation

		return nil
	}

	for i, zone := range k.icsKeeper.AllRegisteredZones(ctx) {
		k.Logger(ctx).Info("Zones", "i", i, "zone", zone.ChainId, "performance address", zone.PerformanceAddress.GetAddress())

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
			rewardscb,
		)
	}

	// TODO: distribute zone allocation on callback
	// We burn for now to ensure sound accounting for testing purposes
	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, allocation)
	if err != nil {
		return err
	}

	return fmt.Errorf("allocateValidatorChoiceRewards not implemented")
}

func (k Keeper) zoneCallback(ctx sdk.Context, response []byte, query icqtypes.Query) (*zoneScore, error) {
	delegatorRewards := distrtypes.QueryDelegationTotalRewardsResponse{}
	err := k.cdc.Unmarshal(response, &delegatorRewards)
	if err != nil {
		return nil, err
	}

	zone, found := k.icsKeeper.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return nil, fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	k.Logger(ctx).Info(
		"Performance Rewards Zone Callback Response",
		"zone", zone.ChainId,
		"rewards", delegatorRewards,
	)

	zs := zoneScore{
		zoneId:           zone.ChainId,
		totalVotingPower: sdk.NewInt(0),
		validatorScores:  make(map[string]*validator),
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
	k.Logger(ctx).Info("Calculate distribution scores", "zone", zone.ChainId)

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
		zs.totalVotingPower = zs.totalVotingPower.Add(val.VotingPower)
		if _, exists := zs.validatorScores[val.ValoperAddress]; !exists {
			zs.validatorScores[val.ValoperAddress] = &validator{Validator: val}
		}

		// Set max/min
		if max.LT(val.VotingPower) {
			max = val.VotingPower
			k.Logger(ctx).Info("New power max", "max", max, "validator", val.ValoperAddress)
		}
		if min.GT(val.VotingPower) {
			min = val.VotingPower
			k.Logger(ctx).Info("New power min", "min", min, "validator", val.ValoperAddress)
		}
	}

	k.Logger(ctx).Info("zone voting power", "zone", zone.ChainId, "total voting power", zs.totalVotingPower)

	if zs.totalVotingPower.IsZero() {
		k.Logger(ctx).Error("Zone invalid, zero voting power", "zone", zone)
		panic("This should never happen!")
	}

	// calculate power percentage and normalized distribution scores
	maxp := max.ToDec().Quo(zs.totalVotingPower.ToDec())
	minp := min.ToDec().Quo(zs.totalVotingPower.ToDec())
	for _, vs := range zs.validatorScores {
		// calculate power percentage
		vs.powerPercentage = vs.VotingPower.ToDec().Quo(zs.totalVotingPower.ToDec())

		// calculate normalized distribution score
		vs.distributionScore = sdk.NewDec(1).Sub(
			vs.powerPercentage.Sub(minp).Mul(
				sdk.NewDec(1).Quo(maxp),
			),
		)

		k.Logger(ctx).Info(
			"Validator Score",
			"validator", vs.ValoperAddress,
			"power percentage", vs.powerPercentage,
			"distribution score", vs.distributionScore,
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
	k.Logger(ctx).Info("Calculate performance & overall scores")

	rewards := delegatorRewards.GetRewards()
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
		vs, exists := zs.validatorScores[reward.ValidatorAddress]
		if !exists {
			k.Logger(ctx).Info("validator may have been removed from active set", "validator", reward.ValidatorAddress)
			continue
		}

		vs.performanceScore = reward.Reward.AmountOf(zone.BaseDenom).Quo(expected)
		if vs.performanceScore.GT(limit) {
			vs.performanceScore = limit
		}
		k.Logger(ctx).Info("Performance Score", "validator", vs.ValoperAddress, "performance", vs.performanceScore)

		// calculate overall score
		vs.overallScore = vs.distributionScore.Mul(vs.performanceScore)
		k.Logger(ctx).Info("Overall Score", "validator", vs.ValoperAddress, "overall", vs.overallScore)

		// prepare validator performance withdrawal msg
		msg := &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: zone.PerformanceAddress.GetAddress(),
			ValidatorAddress: vs.ValoperAddress,
		}
		msgs = append(msgs, msg)
	}

	// submit rewards withdrawals to reset zone performance for next epoch
	k.Logger(ctx).Info("Send performance rewards withdrawal messages to reset scores for next epoch")
	if len(msgs) > 0 {
		if err := k.icsKeeper.SubmitTx(ctx, msgs, zone.PerformanceAddress); err != nil {
			return err
		}
	}

	return nil
}
