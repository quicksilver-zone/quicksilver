package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type zoneScore struct {
	zoneId           string // chainId
	totalVotingPower sdk.Dec
	validatorScores  []validatorScore
	allocation       sdk.Coins
}

type validatorScore struct {
	zindex            int
	expectedRewards   sdk.Dec
	actualRewards     sdk.Dec
	performanceScore  sdk.Dec
	distributionScore sdk.Dec
	overallScore      sdk.Dec
}

func (k Keeper) allocateValidatorChoiceRewards(ctx sdk.Context, allocation sdk.Coins) error {
	k.Logger(ctx).Info("allocateValidatorChoiceRewards", "allocation", allocation)
	// DEVTEST:
	if ctx.Context().Value("DEVTEST") == "DEVTEST" {
		fmt.Printf("\t\tAllocate Validator Choice Rewards:\t%v\n", allocation)
	}

	var zoneScores []zoneScore

	// DEVTEST TODO: setup registered zones for testing...
	for i, zone := range k.icsKeeper.AllRegisteredZones(ctx) {
		k.Logger(ctx).Info("Zones", "i", i, "zone", zone.ChainId)
		// DEVTEST:
		if ctx.Context().Value("DEVTEST") == "DEVTEST" {
			fmt.Printf("\t\t\tZone [%d]:\t%v\n", i, zone.ChainId)
		}

		zoneScores = append(zoneScores, k.getZoneScore(ctx, zone))
	}

	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, allocation)
	if err != nil {
		return err
	}

	return fmt.Errorf("allocateValidatorChoiceRewards not implemented")
}

func (k Keeper) getZoneScore(ctx sdk.Context, zone icstypes.RegisteredZone) zoneScore {
	// DEVTEST: Set special context for DEVTEST output
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), "DEVTEST", "DEVTEST"))

	zs := zoneScore{
		zoneId:           zone.ChainId,
		totalVotingPower: sdk.NewDec(0),
	}

	// create score struct and calculate total voting power
	for i, val := range zone.GetValidatorsSorted() {
		// DEVTEST:
		if ctx.Context().Value("DEVTEST") == "DEVTEST" {
			fmt.Printf("\t\t\t\tValidator [%d]: %s, VotingPower = %v\n", i, val.GetValoperAddress(), val.VotingPower)
		}

		zs.totalVotingPower = zs.totalVotingPower.Add(val.VotingPower)
		zs.validatorScores = append(
			zs.validatorScores,
			validatorScore{
				zindex:            i,
				expectedRewards:   sdk.NewDec(0),
				actualRewards:     sdk.NewDec(0),
				performanceScore:  sdk.NewDec(0),
				distributionScore: sdk.NewDec(0),
				overallScore:      sdk.NewDec(0),
			},
		)
	}

	// calculate expected scores and min/max (percentage) for zone
	max := sdk.NewDec(0)
	min := sdk.NewDec(100)
	for i, val := range zone.Validators {
		zs.validatorScores[i].expectedRewards = val.VotingPower.Quo(zs.totalVotingPower).Mul(sdk.NewDec(100))
		// DEVTEST:
		if ctx.Context().Value("DEVTEST") == "DEVTEST" {
			fmt.Printf("\t\t\t\texpectedRewards: %v\n", zs.validatorScores[i].expectedRewards)
		}

		// Set max/min
		if max.LT(zs.validatorScores[i].expectedRewards) {
			max = zs.validatorScores[i].expectedRewards
			// DEVTEST:
			if ctx.Context().Value("DEVTEST") == "DEVTEST" {
				fmt.Printf("\t\t\t\t\tNew Max: %v\n", max)
			}
		}
		if min.GT(zs.validatorScores[i].expectedRewards) {
			min = zs.validatorScores[i].expectedRewards
			// DEVTEST:
			if ctx.Context().Value("DEVTEST") == "DEVTEST" {
				fmt.Printf("\t\t\t\t\tNew Min: %v\n", min)
			}
		}
	}

	// calculate scores
	for i := 0; i < len(zs.validatorScores); i++ {
		zs.validatorScores[i].distributionScore = sdk.NewDec(1).Sub(
			zs.validatorScores[i].expectedRewards.Sub(min).Mul(
				sdk.NewDec(1).Quo(max),
			),
		)

		// DEVTEST:
		if ctx.Context().Value("DEVTEST") == "DEVTEST" {
			fmt.Printf("\t\t\t\tDistribution Score: %v\n", zs.validatorScores[i].distributionScore)
		}
	}

	return zs
}
