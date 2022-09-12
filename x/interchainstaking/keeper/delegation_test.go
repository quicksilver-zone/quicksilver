package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/require"
)

func TestDetermineAllocations(t *testing.T) {
	// we auto generate the validator addresses in these tests. any dust gets allocated to the first validator in the list
	// once sorted alphabetically on valoper. So we track the `dust` as a separate value in each test case and allocate to
	// the first validator.

	val1 := utils.GenerateValAddressForTest()
	val2 := utils.GenerateValAddressForTest()
	val3 := utils.GenerateValAddressForTest()
	val4 := utils.GenerateValAddressForTest()

	tc := []struct {
		current  map[string]sdk.Int
		sum      sdk.Int
		target   map[string]*types.ValidatorIntent
		inAmount sdk.Coins
		expected map[string]sdk.Int
		dust     sdk.Int
	}{
		{
			current: map[string]sdk.Int{
				val1.String(): sdk.NewInt(350000),
				val2.String(): sdk.NewInt(650000),
				val3.String(): sdk.NewInt(75000),
			},
			sum: sdk.NewInt(1075000),
			target: map[string]*types.ValidatorIntent{
				val1.String(): {ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(30, 2)},
				val2.String(): {ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(63, 2)},
				val3.String(): {ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(7, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(50000))),
			expected: map[string]sdk.Int{
				val1.String(): sdk.ZeroInt(),
				val2.String(): sdk.NewInt(33182),
				val3.String(): sdk.NewInt(16818),
			},
		},
		{
			current: map[string]sdk.Int{
				val1.String(): sdk.NewInt(52),
				val2.String(): sdk.NewInt(24),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(4),
			},
			sum: sdk.NewInt(100),
			target: map[string]*types.ValidatorIntent{
				val1.String(): {ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				val2.String(): {ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				val3.String(): {ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				val4.String(): {ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(20))),
			expected: map[string]sdk.Int{
				val4.String(): sdk.NewInt(11),
				val3.String(): sdk.ZeroInt(),
				val2.String(): sdk.NewInt(6),
				val1.String(): sdk.NewInt(3),
			},
		},
		{
			current: map[string]sdk.Int{
				val1.String(): sdk.NewInt(52),
				val2.String(): sdk.NewInt(24),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(4),
			},
			sum: sdk.NewInt(100),
			target: map[string]*types.ValidatorIntent{
				val1.String(): {ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				val2.String(): {ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				val3.String(): {ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				val4.String(): {ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(50))),
			expected: map[string]sdk.Int{
				val4.String(): sdk.NewInt(20),
				val2.String(): sdk.NewInt(13),
				val1.String(): sdk.NewInt(10),
				val3.String(): sdk.NewInt(7),
			},
		},

		// test to check for div-by-zero when no existing allocations exist. this testcase needs to add 2 dust onto the
		// alphabetically first validator.

	}

	for caseNumber, val := range tc {
		allocations := icskeeper.DetermineAllocationsForDelegation(val.current, val.sum, val.target, val.inAmount)
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
			s.SetupZones()

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
