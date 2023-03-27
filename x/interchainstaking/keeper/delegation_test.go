package keeper_test

import (
	"fmt"
	"testing"
	"time"

	cosmosmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_DelegationStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	zoneValidatorAddresses := zone.GetValidatorsAddressesAsSlice()

	performanceDelegations := icsKeeper.GetAllPerformanceDelegations(ctx, &zone)
	suite.Require().Len(performanceDelegations, 4)

	performanceDelegationPointers := icsKeeper.GetAllPerformanceDelegationsAsPointer(ctx, &zone)
	for i, pdp := range performanceDelegationPointers {
		suite.Require().Equal(performanceDelegations[i], *pdp)
	}

	// update performance delegation
	updateDelegation, found := icsKeeper.GetPerformanceDelegation(ctx, &zone, zoneValidatorAddresses[0])
	suite.Require().True(found)
	suite.Require().Equal(uint64(0), updateDelegation.Amount.Amount.Uint64())

	updateDelegation.Amount.Amount = cosmosmath.NewInt(10000)
	icsKeeper.SetPerformanceDelegation(ctx, &zone, updateDelegation)

	updatedDelegation, found := icsKeeper.GetPerformanceDelegation(ctx, &zone, zoneValidatorAddresses[0])
	suite.Require().True(found)
	suite.Require().Equal(updateDelegation, updatedDelegation)

	// check that there are no delegations
	delegations := icsKeeper.GetAllDelegations(ctx, &zone)
	suite.Require().Len(delegations, 0)

	// set delegations
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		icstypes.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[0],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000)),
		),
	)
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		icstypes.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[1],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000000)),
		),
	)
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		icstypes.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[2],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(20000000)),
		),
	)

	// check for delegations set above
	delegations = icsKeeper.GetAllDelegations(ctx, &zone)
	suite.Require().Len(delegations, 3)

	// load and match pointers
	delegationPointers := icsKeeper.GetAllDelegationsAsPointer(ctx, &zone)
	for i, dp := range delegationPointers {
		suite.Require().Equal(delegations[i], *dp)
	}

	// get delegations for delegation address and match
	addr, err := sdk.AccAddressFromBech32(zone.DelegationAddress.GetAddress())
	suite.Require().NoError(err)
	dds := icsKeeper.GetDelegatorDelegations(ctx, &zone, addr)
	suite.Require().Len(dds, 3)
	suite.Require().Equal(delegations, dds)
}

func TestDetermineAllocationsForDelegation(t *testing.T) {
	// we auto generate the validator addresses in these tests. any dust gets allocated to the first validator in the list
	// once sorted alphabetically on valoper.

	val1 := utils.GenerateValAddressForTest()
	val2 := utils.GenerateValAddressForTest()
	val3 := utils.GenerateValAddressForTest()
	val4 := utils.GenerateValAddressForTest()

	tc := []struct {
		current  map[string]cosmosmath.Int
		target   icstypes.ValidatorIntents
		inAmount sdk.Coins
		expected map[string]cosmosmath.Int
		dust     cosmosmath.Int
	}{
		{
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(350000),
				val2.String(): sdk.NewInt(650000),
				val3.String(): sdk.NewInt(75000),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(30, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(63, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(7, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(50000))),
			expected: map[string]cosmosmath.Int{
				val1.String(): sdk.ZeroInt(),
				val2.String(): sdk.NewInt(33182),
				val3.String(): sdk.NewInt(16818),
			},
		},
		{
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(52),
				val2.String(): sdk.NewInt(24),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(4),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(20))),
			expected: map[string]cosmosmath.Int{
				val4.String(): sdk.NewInt(11),
				val3.String(): sdk.ZeroInt(),
				val2.String(): sdk.NewInt(6),
				val1.String(): sdk.NewInt(3),
			},
		},
		{
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(52),
				val2.String(): sdk.NewInt(24),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(4),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(50))),
			expected: map[string]cosmosmath.Int{
				val4.String(): sdk.NewInt(20),
				val2.String(): sdk.NewInt(13),
				val1.String(): sdk.NewInt(10),
				val3.String(): sdk.NewInt(7),
			},
		},
	}

	for caseNumber, val := range tc {
		sum := cosmosmath.ZeroInt()
		for _, amount := range val.current {
			sum = sum.Add(amount)
		}
		allocations := icskeeper.DetermineAllocationsForDelegation(val.current, sum, val.target, val.inAmount)
		require.Equal(t, len(val.expected), len(allocations))
		for valoper := range val.expected {
			ex, ok := val.expected[valoper]
			require.True(t, ok)
			ac, ok := allocations[valoper]
			require.True(t, ok)
			require.True(t, ex.Equal(ac), fmt.Sprintf("Test Case #%d failed; allocations did not equal expected output - expected %s, got %s.", caseNumber, val.expected[valoper], allocations[valoper]))
		}
	}
}

type delegationUpdate struct {
	delegation icstypes.Delegation
	absolute   bool
}

func (s *KeeperTestSuite) TestUpdateDelegation() {
	del1 := utils.GenerateAccAddressForTest()

	val1 := utils.GenerateValAddressForTest()
	val2 := utils.GenerateValAddressForTest()
	val3 := utils.GenerateValAddressForTest()
	val4 := utils.GenerateValAddressForTest()
	val5 := utils.GenerateValAddressForTest()
	val6 := utils.GenerateValAddressForTest()

	tests := []struct {
		name       string
		delegation *icstypes.Delegation
		updates    []delegationUpdate
		expected   icstypes.Delegation
	}{
		{
			"single update, relative increase +3000",
			&icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
			},
			icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(6000))},
		},
		{
			"single update, relative increase +3000",
			&icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   true,
				},
			},
			icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
		},
		{
			"multi update, relative increase +3000, +2000",
			&icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
				{
					delegation: icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
					absolute:   false,
				},
			},
			icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(8000))},
		},
		{
			"multi update, relative +3000, absolute +2000",
			&icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
				{
					delegation: icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
					absolute:   true,
				},
			},
			icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
		},
		{
			"new delegation, relative increase +10000",
			nil,
			[]delegationUpdate{
				{
					delegation: icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val5.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(10000))},
					absolute:   false,
				},
			},
			icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val5.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(10000))},
		},
		{
			"new delegation, absolute increase +15000",
			nil,
			[]delegationUpdate{
				{
					delegation: icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val6.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(15000))},
					absolute:   true,
				},
			},
			icstypes.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val6.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(15000))},
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			app := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			if tt.delegation != nil {
				app.InterchainstakingKeeper.SetDelegation(ctx, &zone, *tt.delegation)
			}

			for _, update := range tt.updates {
				app.InterchainstakingKeeper.UpdateDelegationRecordForAddress(ctx, update.delegation.DelegationAddress, update.delegation.ValidatorAddress, update.delegation.Amount, &zone, update.absolute)
			}

			actual, found := app.InterchainstakingKeeper.GetDelegation(ctx, &zone, tt.expected.DelegationAddress, tt.expected.ValidatorAddress)
			s.Require().True(found)
			s.Require().Equal(tt.expected, actual)
		})
	}
}

func TestCalculateDeltas(t *testing.T) {
	// we auto generate the validator addresses in these tests. any dust gets allocated to the first validator in the list
	// once sorted alphabetically on valoper.

	val1 := utils.GenerateValAddressForTest()
	val2 := utils.GenerateValAddressForTest()
	val3 := utils.GenerateValAddressForTest()
	val4 := utils.GenerateValAddressForTest()

	zone := icstypes.Zone{Validators: []*icstypes.Validator{
		{ValoperAddress: val1.String(), CommissionRate: sdk.NewDecWithPrec(30, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val2.String(), CommissionRate: sdk.NewDecWithPrec(25, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val3.String(), CommissionRate: sdk.NewDecWithPrec(10, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val4.String(), CommissionRate: sdk.NewDecWithPrec(12, 2), Status: stakingtypes.BondStatusBonded},
	}}

	zone2 := icstypes.Zone{Validators: []*icstypes.Validator{
		{ValoperAddress: val1.String(), CommissionRate: sdk.NewDecWithPrec(30, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val2.String(), CommissionRate: sdk.NewDecWithPrec(25, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val3.String(), CommissionRate: sdk.NewDecWithPrec(10, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val4.String(), CommissionRate: sdk.NewDecWithPrec(75, 2), Status: stakingtypes.BondStatusBonded},
	}}

	tc := []struct {
		current  map[string]cosmosmath.Int
		target   icstypes.ValidatorIntents
		expected icstypes.ValidatorIntents
	}{
		{
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(350000),
				val2.String(): sdk.NewInt(650000),
				val3.String(): sdk.NewInt(75000),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(30, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(63, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(7, 2)},
			},
			expected: icstypes.ValidatorIntents{
				{ValoperAddress: val2.String(), Weight: sdk.NewDec(27250)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDec(250)},
				{ValoperAddress: val1.String(), Weight: sdk.NewDec(-27500)},
			},
		},
		{
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(53),
				val2.String(): sdk.NewInt(26),
				val3.String(): sdk.NewInt(14),
				val4.String(): sdk.NewInt(7),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(28, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(12, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			expected: icstypes.ValidatorIntents{
				{ValoperAddress: val4.String(), Weight: sdk.NewDec(3)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDec(2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDec(-2)},
				{ValoperAddress: val1.String(), Weight: sdk.NewDec(-3)},
			},
		},
		{
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(30),
				val2.String(): sdk.NewInt(30),
				val3.String(): sdk.NewInt(60),
				val4.String(): sdk.NewInt(180),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			expected: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDec(120)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDec(45)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDec(-15)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDec(-150)},
			},
		},
		// default intent -- all equal
		{
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(15),
				val2.String(): sdk.NewInt(5),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone.GetAggregateIntentOrDefault(),
			expected: icstypes.ValidatorIntents{
				{ValoperAddress: val2.String(), Weight: sdk.NewDec(20)},
				{ValoperAddress: val1.String(), Weight: sdk.NewDec(10)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDec(5)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDec(-35)},
			},
		},
		{
			// GetAggregateIntentOrDefault will preclude val4 on account on high commission.
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDec(25)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDec(21)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDec(13)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDec(-60)},
			},
		},
	}

	for caseNumber, val := range tc {
		sum := cosmosmath.ZeroInt()
		for _, amount := range val.current {
			sum = sum.Add(amount)
		}
		deltas := icskeeper.CalculateDeltas(val.current, sum, val.target)
		fmt.Println("Deltas", deltas)
		require.Equal(t, len(val.expected), len(deltas), fmt.Sprintf("expected %d RebalanceTargets in case %d, got %d", len(val.expected), caseNumber, len(deltas)))
		for idx, expected := range val.expected {
			require.Equal(t, expected, deltas[idx], fmt.Sprintf("case %d, idx %d: Expected %v, got %v", caseNumber, idx, expected, deltas[idx]))
		}

	}
}

func TestDetermineAllocationsForRebalance(t *testing.T) {
	// we auto generate the validator addresses in these tests. any dust gets allocated to the first validator in the list
	// once sorted alphabetically on valoper.

	val1 := utils.GenerateValAddressForTest()
	val2 := utils.GenerateValAddressForTest()
	val3 := utils.GenerateValAddressForTest()
	val4 := utils.GenerateValAddressForTest()

	zone := icstypes.Zone{Validators: []*icstypes.Validator{
		{ValoperAddress: val1.String(), CommissionRate: sdk.NewDecWithPrec(30, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val2.String(), CommissionRate: sdk.NewDecWithPrec(25, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val3.String(), CommissionRate: sdk.NewDecWithPrec(10, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val4.String(), CommissionRate: sdk.NewDecWithPrec(12, 2), Status: stakingtypes.BondStatusBonded},
	}}

	zone2 := icstypes.Zone{Validators: []*icstypes.Validator{
		{ValoperAddress: val1.String(), CommissionRate: sdk.NewDecWithPrec(30, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val2.String(), CommissionRate: sdk.NewDecWithPrec(25, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val3.String(), CommissionRate: sdk.NewDecWithPrec(10, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val4.String(), CommissionRate: sdk.NewDecWithPrec(75, 2), Status: stakingtypes.BondStatusBonded},
	}}

	tc := []struct {
		name          string
		current       map[string]cosmosmath.Int
		target        icstypes.ValidatorIntents
		expected      []icskeeper.RebalanceTarget
		dust          cosmosmath.Int
		redelegations []types.RedelegationRecord
	}{
		{
			name: "case 1",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(350000),
				val2.String(): sdk.NewInt(650000),
				val3.String(): sdk.NewInt(75000),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(30, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(63, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(7, 2)},
			},
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(27250), Source: val1.String(), Target: val2.String()},
				{Amount: cosmosmath.NewInt(250), Source: val1.String(), Target: val3.String()},
			},
		},
		{
			name: "case 2",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(56),
				val2.String(): sdk.NewInt(24),
				val3.String(): sdk.NewInt(14),
				val4.String(): sdk.NewInt(5),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(28, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(12, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(4), Source: val1.String(), Target: val4.String()},
				{Amount: cosmosmath.NewInt(3), Source: val1.String(), Target: val2.String()},
			},
			redelegations: []types.RedelegationRecord{},
		},
		{
			name: "case 3",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(30),
				val2.String(): sdk.NewInt(30),
				val3.String(): sdk.NewInt(60),
				val4.String(): sdk.NewInt(180),
			},
			target: icstypes.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(42), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				//{Amount: cosmosmath.NewInt(10), Source: val4.String(), Target: val2.String()},
			},
			redelegations: []types.RedelegationRecord{},
		},
		// default intent -- all equal
		{
			name: "case 4 - default intent, all equal",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(15),
				val2.String(): sdk.NewInt(5),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone.GetAggregateIntentOrDefault(),
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(14), Source: val4.String(), Target: val2.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				//{Amount: cosmosmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				//{Amount: cosmosmath.NewInt(5), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{},
		},
		//
		{
			name: "case 5 - default intent with val4 high commission; truncate rebalance to 50% of tvl",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(14), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				//{Amount: cosmosmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				//{Amount: cosmosmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{},
		},
		{
			name: "case 6 - includes redelegation, no impact",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(14), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				//{Amount: cosmosmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				//{Amount: cosmosmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 30, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
		{
			name: "case 7 - includes redelegation, truncated delegation",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				//{Amount: cosmosmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				//{Amount: cosmosmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 50, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
		{
			name: "case 8 - includes redelegation, truncated delegation",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				//{Amount: cosmosmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				//{Amount: cosmosmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 50, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
		{
			name: "case 8 - includes redelegation, truncated delegation overflow",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(2),
				val2.String(): sdk.NewInt(8),
				val3.String(): sdk.NewInt(30),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []icskeeper.RebalanceTarget{
				{Amount: cosmosmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				//{Amount: cosmosmath.NewInt(4), Source: val3.String(), Target: val1.String()},  // joe: I would expect this to be included...
				//{Amount: cosmosmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 50, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
		{
			name: "case 9 - includes redelegation, zero delegation",
			current: map[string]cosmosmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target:   zone2.GetAggregateIntentOrDefault(),
			expected: []icskeeper.RebalanceTarget{
				//{Amount: cosmosmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				//{Amount: cosmosmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				//{Amount: cosmosmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 60, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
	}

	for _, val := range tc {
		sum := cosmosmath.ZeroInt()
		for _, amount := range val.current {
			sum = sum.Add(amount)
		}
		allocations := icskeeper.DetermineAllocationsForRebalancing(val.current, sum, val.target, val.redelegations)
		require.Equal(t, len(val.expected), len(allocations), fmt.Sprintf("expected %d RebalanceTargets in '%s', got %d", len(val.expected), val.name, len(allocations)))
		for idx, rebalance := range val.expected {
			require.Equal(t, rebalance, allocations[idx], fmt.Sprintf("%s, idx %d: Expected %v, got %v", val.name, idx, rebalance, allocations[idx]))
		}
	}
}

func (s *KeeperTestSuite) TestStoreGetDeleteDelegation() {
	s.Run("delegation - store / get / delete", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		ctx := s.chainA.GetContext()

		zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		s.Require().True(found)

		delegator := utils.GenerateAccAddressForTest()
		validator := utils.GenerateValAddressForTest()

		_, found = app.InterchainstakingKeeper.GetDelegation(ctx, &zone, delegator.String(), validator.String())
		s.Require().False(found)

		newDelegation := icstypes.NewDelegation(delegator.String(), validator.String(), sdk.NewCoin("uatom", sdk.NewInt(5000)))
		app.InterchainstakingKeeper.SetDelegation(ctx, &zone, newDelegation)

		fetchedDelegation, found := app.InterchainstakingKeeper.GetDelegation(ctx, &zone, delegator.String(), validator.String())
		s.Require().True(found)
		s.Require().Equal(newDelegation, fetchedDelegation)

		allDelegations := app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)
		s.Require().Len(allDelegations, 1)

		err := app.InterchainstakingKeeper.RemoveDelegation(ctx, &zone, newDelegation)
		s.Require().NoError(err)

		allDelegations2 := app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)
		s.Require().Len(allDelegations2, 0)
	})
}
