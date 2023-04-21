package keeper_test

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (s *KeeperTestSuite) TestKeeper_DelegationStore() {
	s.SetupTest()
	s.setupTestZones()

	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	// get test zone
	zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	zoneValidatorAddresses := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

	performanceDelegations := icsKeeper.GetAllPerformanceDelegations(ctx, &zone)
	s.Require().Len(performanceDelegations, 4)

	performanceDelegationPointers := icsKeeper.GetAllPerformanceDelegationsAsPointer(ctx, &zone)
	for i, pdp := range performanceDelegationPointers {
		s.Require().Equal(performanceDelegations[i], *pdp)
	}

	// update performance delegation
	updateDelegation, found := icsKeeper.GetPerformanceDelegation(ctx, &zone, zoneValidatorAddresses[0])
	s.Require().True(found)
	s.Require().Equal(uint64(0), updateDelegation.Amount.Amount.Uint64())

	updateDelegation.Amount.Amount = sdkmath.NewInt(10000)
	icsKeeper.SetPerformanceDelegation(ctx, &zone, updateDelegation)

	updatedDelegation, found := icsKeeper.GetPerformanceDelegation(ctx, &zone, zoneValidatorAddresses[0])
	s.Require().True(found)
	s.Require().Equal(updateDelegation, updatedDelegation)

	// check that there are no delegations
	delegations := icsKeeper.GetAllDelegations(ctx, &zone)
	s.Require().Len(delegations, 0)

	// set delegations
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		types.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[0],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000)),
		),
	)
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		types.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[1],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000000)),
		),
	)
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		types.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[2],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(20000000)),
		),
	)

	// check for delegations set above
	delegations = icsKeeper.GetAllDelegations(ctx, &zone)
	s.Require().Len(delegations, 3)

	// load and match pointers
	delegationPointers := icsKeeper.GetAllDelegationsAsPointer(ctx, &zone)
	for i, dp := range delegationPointers {
		s.Require().Equal(delegations[i], *dp)
	}

	// get delegations for delegation address and match
	addr, err := sdk.AccAddressFromBech32(zone.DelegationAddress.GetAddress())
	s.Require().NoError(err)
	dds := icsKeeper.GetDelegatorDelegations(ctx, &zone, addr)
	s.Require().Len(dds, 3)
	s.Require().Equal(delegations, dds)
}

type delegationUpdate struct {
	delegation types.Delegation
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
		delegation *types.Delegation
		updates    []delegationUpdate
		expected   types.Delegation
	}{
		{
			"single update, relative increase +3000",
			&types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(6000))},
		},
		{
			"single update, relative increase +3000",
			&types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   true,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
		},
		{
			"multi update, relative increase +3000, +2000",
			&types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
					absolute:   false,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(8000))},
		},
		{
			"multi update, relative +3000, absolute +2000",
			&types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
					absolute:   true,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
		},
		{
			"new delegation, relative increase +10000",
			nil,
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val5.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(10000))},
					absolute:   false,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val5.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(10000))},
		},
		{
			"new delegation, absolute increase +15000",
			nil,
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val6.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(15000))},
					absolute:   true,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val6.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(15000))},
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			qApp := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			zone, found := qApp.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			if tt.delegation != nil {
				qApp.InterchainstakingKeeper.SetDelegation(ctx, &zone, *tt.delegation)
			}

			for _, update := range tt.updates {
				err := qApp.InterchainstakingKeeper.UpdateDelegationRecordForAddress(ctx, update.delegation.DelegationAddress, update.delegation.ValidatorAddress, update.delegation.Amount, &zone, update.absolute)
				s.Require().NoError(err)
			}

			actual, found := qApp.InterchainstakingKeeper.GetDelegation(ctx, &zone, tt.expected.DelegationAddress, tt.expected.ValidatorAddress)
			s.Require().True(found)
			s.Require().Equal(tt.expected, actual)
		})
	}
}

// func (s *KeeperTestSuite) TestCalculateDeltas() {

// 	// we auto generate the validator addresses in these tests. any dust gets allocated to the first validator in the list
// 	// once sorted alphabetically on valoper.

// 	val1 := utils.GenerateValAddressForTest()
// 	val2 := utils.GenerateValAddressForTest()
// 	val3 := utils.GenerateValAddressForTest()
// 	val4 := utils.GenerateValAddressForTest()

// 	zone := types.Zone{Validators: []*types.Validator{
// 		{ValoperAddress: val1.String(), CommissionRate: sdk.NewDecWithPrec(30, 2), Status: stakingtypes.BondStatusBonded},
// 		{ValoperAddress: val2.String(), CommissionRate: sdk.NewDecWithPrec(25, 2), Status: stakingtypes.BondStatusBonded},
// 		{ValoperAddress: val3.String(), CommissionRate: sdk.NewDecWithPrec(10, 2), Status: stakingtypes.BondStatusBonded},
// 		{ValoperAddress: val4.String(), CommissionRate: sdk.NewDecWithPrec(12, 2), Status: stakingtypes.BondStatusBonded},
// 	}}

// 	zone2 := types.Zone{Validators: []*types.Validator{
// 		{ValoperAddress: val1.String(), CommissionRate: sdk.NewDecWithPrec(30, 2), Status: stakingtypes.BondStatusBonded},
// 		{ValoperAddress: val2.String(), CommissionRate: sdk.NewDecWithPrec(25, 2), Status: stakingtypes.BondStatusBonded},
// 		{ValoperAddress: val3.String(), CommissionRate: sdk.NewDecWithPrec(10, 2), Status: stakingtypes.BondStatusBonded},
// 		{ValoperAddress: val4.String(), CommissionRate: sdk.NewDecWithPrec(75, 2), Status: stakingtypes.BondStatusBonded},
// 	}}

// 	tc := []struct {
// 		current  map[string]sdkmath.Int
// 		target   func(qs *app.Quicksilver) types.ValidatorIntents {}
// 		expected types.ValidatorIntents
// 	}{
// 		{
// 			current: map[string]sdkmath.Int{
// 				val1.String(): sdk.NewInt(350000),
// 				val2.String(): sdk.NewInt(650000),
// 				val3.String(): sdk.NewInt(75000),
// 			},
// 			target: func(qs *app.Quicksilver) types.ValidatorIntents {
// 				return types.ValidatorIntents{
// 					{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(30, 2)},
// 					{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(63, 2)},
// 					{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(7, 2)},
// 				}
// 			},
// 			expected: types.ValidatorIntents{
// 				{ValoperAddress: val2.String(), Weight: sdk.NewDec(27250)},
// 				{ValoperAddress: val3.String(), Weight: sdk.NewDec(250)},
// 				{ValoperAddress: val1.String(), Weight: sdk.NewDec(-27500)},
// 			},
// 		},
// 		{
// 			current: map[string]sdkmath.Int{
// 				val1.String(): sdk.NewInt(53),
// 				val2.String(): sdk.NewInt(26),
// 				val3.String(): sdk.NewInt(14),
// 				val4.String(): sdk.NewInt(7),
// 			},
// 			target: func(qs *app.Quicksilver) types.ValidatorIntents {
// 				return types.ValidatorIntents{
// 					{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
// 					{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(28, 2)},
// 					{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(12, 2)},
// 					{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
// 				}
// 			},
// 			expected: types.ValidatorIntents{
// 				{ValoperAddress: val4.String(), Weight: sdk.NewDec(3)},
// 				{ValoperAddress: val2.String(), Weight: sdk.NewDec(2)},
// 				{ValoperAddress: val3.String(), Weight: sdk.NewDec(-2)},
// 				{ValoperAddress: val1.String(), Weight: sdk.NewDec(-3)},
// 			},
// 		},
// 		{
// 			current: map[string]sdkmath.Int{
// 				val1.String(): sdk.NewInt(30),
// 				val2.String(): sdk.NewInt(30),
// 				val3.String(): sdk.NewInt(60),
// 				val4.String(): sdk.NewInt(180),
// 			},
// 			target: func(qs *app.Quicksilver) types.ValidatorIntents {
// 				return types.ValidatorIntents{
// 					{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
// 					{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
// 					{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
// 					{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
// 				}
// 			},
// 			expected: types.ValidatorIntents{
// 				{ValoperAddress: val1.String(), Weight: sdk.NewDec(120)},
// 				{ValoperAddress: val2.String(), Weight: sdk.NewDec(45)},
// 				{ValoperAddress: val3.String(), Weight: sdk.NewDec(-15)},
// 				{ValoperAddress: val4.String(), Weight: sdk.NewDec(-150)},
// 			},
// 		},
// 		// default intent -- all equal
// 		{
// 			current: map[string]sdkmath.Int{
// 				val1.String(): sdk.NewInt(15),
// 				val2.String(): sdk.NewInt(5),
// 				val3.String(): sdk.NewInt(20),
// 				val4.String(): sdk.NewInt(60),
// 			},
// 			target: func(qs *app.Quicksilver) types.ValidatorIntents {
// 				return qs.InterchainstakingKeeper.GetAggregateIntentOrDefault(ctx, zone)
// 			},
// 			expected: types.ValidatorIntents{
// 				{ValoperAddress: val2.String(), Weight: sdk.NewDec(20)},
// 				{ValoperAddress: val1.String(), Weight: sdk.NewDec(10)},
// 				{ValoperAddress: val3.String(), Weight: sdk.NewDec(5)},
// 				{ValoperAddress: val4.String(), Weight: sdk.NewDec(-35)},
// 			},
// 		},
// 		{
// 			// GetAggregateIntentOrDefault will preclude val4 on account on high commission.
// 			current: map[string]sdkmath.Int{
// 				val1.String(): sdk.NewInt(8),
// 				val2.String(): sdk.NewInt(12),
// 				val3.String(): sdk.NewInt(20),
// 				val4.String(): sdk.NewInt(60),
// 			},
// 			target: func(qs *app.Quicksilver) types.ValidatorIntents {
// 				return qs.InterchainstakingKeeper.GetAggregateIntentOrDefault(ctx, zone2)
// 			},
// 			expected: types.ValidatorIntents{
// 				{ValoperAddress: val1.String(), Weight: sdk.NewDec(25)},
// 				{ValoperAddress: val2.String(), Weight: sdk.NewDec(21)},
// 				{ValoperAddress: val3.String(), Weight: sdk.NewDec(13)},
// 				{ValoperAddress: val4.String(), Weight: sdk.NewDec(-60)},
// 			},
// 		},
// 	}

// 	for caseNumber, val := range tc {
// 		s.Run(fmt.Sprint("case %d", caseNumber), func() {
// 			s.SetupTest()
// 			s.setupTestZones()

// 			app := s.GetQuicksilverApp(s.chainA)
// 			ctx := s.chainA.GetContext()

// 			sum := sdkmath.ZeroInt()
// 			for _, amount := range val.current {
// 				sum = sum.Add(amount)
// 			}
// 			deltas := types.CalculateDeltas(val.current, sum, val.(app))
// 			s.Require().Equal(len(val.expected), len(deltas), fmt.Sprintf("expected %d RebalanceTargets in case %d, got %d", len(val.expected), caseNumber, len(deltas)))
// 			for idx, expected := range val.expected {
// 				s.Require().Equal(expected, deltas[idx], fmt.Sprintf("case %d, idx %d: Expected %v, got %v", caseNumber, idx, expected, deltas[idx]))
// 			}
// 		})
// 	}
// }

/*func TestDetermineAllocationsForRebalance(t *testing.T) {
	// we auto generate the validator addresses in these tests. any dust gets allocated to the first validator in the list
	// once sorted alphabetically on valoper.

	val1 := utils.GenerateValAddressForTest()
	val2 := utils.GenerateValAddressForTest()
	val3 := utils.GenerateValAddressForTest()
	val4 := utils.GenerateValAddressForTest()

	zone := types.Zone{Validators: []*types.Validator{
		{ValoperAddress: val1.String(), CommissionRate: sdk.NewDecWithPrec(30, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val2.String(), CommissionRate: sdk.NewDecWithPrec(25, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val3.String(), CommissionRate: sdk.NewDecWithPrec(10, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val4.String(), CommissionRate: sdk.NewDecWithPrec(12, 2), Status: stakingtypes.BondStatusBonded},
	}}

	zone2 := types.Zone{Validators: []*types.Validator{
		{ValoperAddress: val1.String(), CommissionRate: sdk.NewDecWithPrec(30, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val2.String(), CommissionRate: sdk.NewDecWithPrec(25, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val3.String(), CommissionRate: sdk.NewDecWithPrec(10, 2), Status: stakingtypes.BondStatusBonded},
		{ValoperAddress: val4.String(), CommissionRate: sdk.NewDecWithPrec(75, 2), Status: stakingtypes.BondStatusBonded},
	}}

	tc := []struct {
		name          string
		current       map[string]sdkmath.Int
		locked        map[string]bool
		target        types.ValidatorIntents
		expected      []types.RebalanceTarget
		dust          sdkmath.Int
		redelegations []types.RedelegationRecord
	}{
		{
			name: "case 1",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(350000),
				val2.String(): sdk.NewInt(650000),
				val3.String(): sdk.NewInt(75000),
			},
			target: types.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(30, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(63, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(7, 2)},
			},
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(27250), Source: val1.String(), Target: val2.String()},
				{Amount: sdkmath.NewInt(250), Source: val1.String(), Target: val3.String()},
			},
		},
		{
			name: "case 2",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(56),
				val2.String(): sdk.NewInt(24),
				val3.String(): sdk.NewInt(14),
				val4.String(): sdk.NewInt(5),
			},
			target: types.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(28, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(12, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(4), Source: val1.String(), Target: val4.String()},
				{Amount: sdkmath.NewInt(3), Source: val1.String(), Target: val2.String()},
			},
			redelegations: []types.RedelegationRecord{},
		},
		{
			name: "case 3",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(30),
				val2.String(): sdk.NewInt(30),
				val3.String(): sdk.NewInt(60),
				val4.String(): sdk.NewInt(180),
			},
			target: types.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(42), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				// {Amount: sdkmath.NewInt(10), Source: val4.String(), Target: val2.String()},
			},
			redelegations: []types.RedelegationRecord{},
		},
		// default intent -- all equal
		{
			name: "case 4 - default intent, all equal",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(15),
				val2.String(): sdk.NewInt(5),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone.GetAggregateIntentOrDefault(),
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(14), Source: val4.String(), Target: val2.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				// {Amount: sdkmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// {Amount: sdkmath.NewInt(5), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{},
		},
		//
		{
			name: "case 5 - default intent with val4 high commission; truncate rebalance to 50% of tvl",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(14), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				// {Amount: sdkmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				// {Amount: sdkmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{},
		},
		{
			name: "case 6 - includes redelegation, no impact",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(14), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				// {Amount: sdkmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				// {Amount: sdkmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 30, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
		{
			name: "case 7 - includes redelegation, truncated delegation",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				// {Amount: sdkmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				// {Amount: sdkmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 50, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
		{
			name: "case 8 - includes redelegation, truncated delegation",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				// {Amount: sdkmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				// {Amount: sdkmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 50, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
		{
			name: "case 9 - includes redelegation, truncated delegation overflow",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(2),
				val2.String(): sdk.NewInt(8),
				val3.String(): sdk.NewInt(30),
				val4.String(): sdk.NewInt(60),
			},
			target: zone2.GetAggregateIntentOrDefault(),
			expected: []types.RebalanceTarget{
				{Amount: sdkmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				// {Amount: sdkmath.NewInt(4), Source: val3.String(), Target: val1.String()}, // joe: I would expect this to be included...
				// {Amount: sdkmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 50, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
		{
			name: "case 10 - includes redelegation, zero delegation",
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(8),
				val2.String(): sdk.NewInt(12),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(60),
			},
			target:   zone2.GetAggregateIntentOrDefault(),
			expected: []types.RebalanceTarget{
				// {Amount: sdkmath.NewInt(10), Source: val4.String(), Target: val1.String()},
				// below values _would_ applied, if we weren't limited by a max of total/7
				// {Amount: sdkmath.NewInt(21), Source: val4.String(), Target: val2.String()},
				// {Amount: sdkmath.NewInt(4), Source: val4.String(), Target: val3.String()},
			},
			redelegations: []types.RedelegationRecord{
				{ChainId: "test-1", EpochNumber: 1, Source: val2.String(), Destination: val4.String(), Amount: 60, CompletionTime: time.Now().Add(time.Hour)},
			},
		},
	}

	for _, val := range tc {
		sum := sdkmath.ZeroInt()
		for _, amount := range val.current {
			sum = sum.Add(amount)
		}
		allocations := types.DetermineAllocationsForRebalancing(val.current, val.locked, sum, val.target, val.redelegations, nil)
		require.Equal(t, len(val.expected), len(allocations), fmt.Sprintf("expected %d RebalanceTargets in '%s', got %d", len(val.expected), val.name, len(allocations)))
		for idx, rebalance := range val.expected {
			require.Equal(t, rebalance, allocations[idx], fmt.Sprintf("%s, idx %d: Expected %v, got %v", val.name, idx, rebalance, allocations[idx]))
		}
	}
}*/

func (s *KeeperTestSuite) TestStoreGetDeleteDelegation() {
	s.Run("delegation - store / get / delete", func() {
		s.SetupTest()
		s.setupTestZones()

		qApp := s.GetQuicksilverApp(s.chainA)
		ctx := s.chainA.GetContext()

		zone, found := qApp.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		s.Require().True(found)

		delegator := utils.GenerateAccAddressForTest()
		validator := utils.GenerateValAddressForTest()

		_, found = qApp.InterchainstakingKeeper.GetDelegation(ctx, &zone, delegator.String(), validator.String())
		s.Require().False(found)

		newDelegation := types.NewDelegation(delegator.String(), validator.String(), sdk.NewCoin("uatom", sdk.NewInt(5000)))
		qApp.InterchainstakingKeeper.SetDelegation(ctx, &zone, newDelegation)

		fetchedDelegation, found := qApp.InterchainstakingKeeper.GetDelegation(ctx, &zone, delegator.String(), validator.String())
		s.Require().True(found)
		s.Require().Equal(newDelegation, fetchedDelegation)

		allDelegations := qApp.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)
		s.Require().Len(allDelegations, 1)

		err := qApp.InterchainstakingKeeper.RemoveDelegation(ctx, &zone, newDelegation)
		s.Require().NoError(err)

		allDelegations2 := qApp.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)
		s.Require().Len(allDelegations2, 0)
	})
}

func (s *KeeperTestSuite) TestFlushOutstandingDelegations() {
	userAddress := utils.GenerateAccAddressForTest().String()
	denom := "uatom"
	tests := []struct {
		name               string
		setStatements      func(ctx sdk.Context, quicksilver *app.Quicksilver)
		delAddrBalance     sdk.Coin
		mockAck            bool
		expectedDelegation sdk.Coins
		assertStatements   func(ctx sdk.Context, quicksilver *app.Quicksilver) bool
	}{
		{
			name:           "zero pending delegations, no pending receipts, no exclusion receipts",
			setStatements:  func(ctx sdk.Context, quicksilver *app.Quicksilver) {},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.ZeroInt()),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				return true
			},
		},
		{
			name: "zero pending delegations, 2 pending receipts and no exclusion receipts",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				pendingRecieptTime := cutOffTime.Add(-2 * time.Hour)
				excludedReciptTime := cutOffTime.Add(-3 * time.Hour)

				rcpt1 := types.Receipt{
					ChainId: s.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000), // 20% deposit
						),
					),
					FirstSeen: &pendingRecieptTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: s.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100), // 20% deposit
						),
					),
					FirstSeen: &excludedReciptTime,
					Completed: nil,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(0)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}

					return false
				})

				s.Require().Equal(0, count)
				return true
			},
		},
		{
			name: "zero pending delegations, 1  pending receipt and 1 exclusion receipt",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				pendingRecieptTime := cutOffTime.Add(-2 * time.Hour)
				excludedReciptTime := cutOffTime.Add(2 * time.Hour)

				rcpt1 := types.Receipt{
					ChainId: s.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000), // 20% deposit
						),
					),
					FirstSeen: &pendingRecieptTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: s.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100), // 20% deposit
						),
					),
					FirstSeen: &excludedReciptTime,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(100)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}
					return false
				})
				s.Require().Equal(1, count)
				return true
			},
		},
		{
			name: "non-zero pending delegations, 1 pending receipts and 1 exclusion recipts ",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				pendingRecieptTime := cutOffTime.Add(-2 * time.Hour)
				excludedReciptTime := cutOffTime.Add(2 * time.Hour)

				rcpt1 := types.Receipt{
					ChainId: s.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000), // 20% deposit
						),
					),
					FirstSeen: &pendingRecieptTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: s.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100), // 20% deposit
						),
					),
					FirstSeen: &excludedReciptTime,
					Completed: nil,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(2000100)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}
					return false
				})
				s.Require().Equal(1, count)
				return true
			},
			mockAck:            true,
			expectedDelegation: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
		},
		{
			name: "non-zero pending delegations, 2 pending receipts",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				pendingRecieptTime := cutOffTime.Add(-2 * time.Hour)
				excludedReciptTime := cutOffTime.Add(-3 * time.Hour)

				rcpt1 := types.Receipt{
					ChainId: s.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000), // 20% deposit
						),
					),
					FirstSeen: &pendingRecieptTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: s.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100), // 20% deposit
						),
					),
					FirstSeen: &excludedReciptTime,
					Completed: nil,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(2000100)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}
					return false
				})
				s.Require().Equal(0, count)
				return true
			},
			mockAck:            true,
			expectedDelegation: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000100))),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.setupTestZones()

			quicksilver := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			test.setStatements(ctx, quicksilver)
			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)
			before := quicksilver.InterchainstakingKeeper.AllReceipts(ctx)
			err := quicksilver.InterchainstakingKeeper.FlushOutstandingDelegations(ctx, &zone, test.delAddrBalance)
			if test.mockAck {
				var msgs []sdk.Msg
				allocations, err := quicksilver.InterchainstakingKeeper.DeterminePlanForDelegation(ctx, &zone, test.expectedDelegation)
				s.Require().NoError(err)
				msgs = append(msgs, quicksilver.InterchainstakingKeeper.PrepareDelegationMessagesForCoins(&zone, allocations)...)
				for _, msg := range msgs {
					err := quicksilver.InterchainstakingKeeper.HandleDelegate(ctx, msg, "batch/1577836910")
					s.Require().NoError(err)
				}
			}
			after := quicksilver.InterchainstakingKeeper.AllReceipts(ctx)
			fmt.Println(before, after)
			s.Require().NoError(err)
			isCorrect := test.assertStatements(ctx, quicksilver)
			s.Require().True(isCorrect)
		})
	}
}
