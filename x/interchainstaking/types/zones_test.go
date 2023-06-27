package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestIsDelegateAddress(t *testing.T) {
	acc := addressutils.GenerateAccAddressForTest()
	acc2 := addressutils.GenerateAccAddressForTest()
	bech32 := addressutils.MustEncodeAddressToBech32("cosmos", acc)
	bech322 := addressutils.MustEncodeAddressToBech32("cosmos", acc2)
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", DelegationAddress: &types.ICAAccount{Address: bech32}, Is_118: true}
	require.True(t, zone.IsDelegateAddress(bech32))
	require.False(t, zone.IsDelegateAddress(bech322))
}

func TestGetDelegationAccount(t *testing.T) {
	acc := addressutils.GenerateAccAddressForTest()
	bech32 := addressutils.MustEncodeAddressToBech32("cosmos", acc)
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", DelegationAddress: &types.ICAAccount{Address: bech32}, Is_118: true}
	zone2 := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}

	delegateAccount, err := zone.GetDelegationAccount()
	require.NoError(t, err)
	require.Equal(t, bech32, delegateAccount.Address)

	acc2, err2 := zone2.GetDelegationAccount()
	require.Error(t, err2)
	require.Nil(t, acc2)
}

func TestValidateCoinsForZone(t *testing.T) {
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", Is_118: true}
	valAddresses := []string{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7"}
	require.NoError(t, zone.ValidateCoinsForZone(sdk.NewCoins(sdk.NewCoin("cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy1", sdk.OneInt())), valAddresses))
	require.Errorf(t, zone.ValidateCoinsForZone(sdk.NewCoins(sdk.NewCoin("cosmosvaloper18ldc09yx4aua9g8mkl3sj526hgydzzyehcyjjr1", sdk.OneInt())), valAddresses), "invalid denom for zone: cosmosvaloper18ldc09yx4aua9g8mkl3sj526hgydzzyehcyjjr1")
}

func TestCoinsToIntent(t *testing.T) {
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", Is_118: true}
	valAddresses := []string{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7"}

	testCases := []struct {
		amount         sdk.Coins
		expectedIntent map[string]sdk.Dec
	}{
		{
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(45)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf16", sdk.NewInt(55)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(45),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(55),
			},
		},
		{
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(350)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf16", sdk.NewInt(350)),
				sdk.NewCoin("cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy6", sdk.NewInt(300)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(350),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(350),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDec(300),
			},
		},
		{
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(3900)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf16", sdk.NewInt(5500)),
				sdk.NewCoin("cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy6", sdk.NewInt(3000)),
				sdk.NewCoin("cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll2", sdk.NewInt(500)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(3900),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(5500),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDec(3000),
				"cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll": sdk.NewDec(500),
			},
		},
	}

	for _, tc := range testCases {
		out := zone.ConvertCoinsToOrdinalIntents(tc.amount, valAddresses)
		for _, v := range out {
			if !tc.expectedIntent[v.ValoperAddress].Equal(v.Weight) {
				t.Errorf("Got %v expected %v", v.Weight, tc.expectedIntent[v.ValoperAddress])
			}
		}
	}
}

func TestDecodeMemo(t *testing.T) {
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", Is_118: true}
	zone.Validators = append(zone.Validators,
		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)

	testCases := []struct {
		name               string
		memo               string
		amount             int
		expectedIntent     map[string]sdk.Dec
		expectedMemoFields types.MemoFields
		wantErr            bool
	}{
		{
			memo:   "AipahL/4TH3a0Ry4wHOG6RkoxWdcpLxuppAElPH3PNriuvHIuI/1/AuKM5w=",
			amount: 100,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(45),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(55),
			},
			expectedMemoFields: types.MemoFields{
				2: types.MemoField{ID: 2, Data: []uint8{0x5a, 0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc, 0x6e, 0xa6, 0x90, 0x4, 0x94, 0xf1, 0xf7, 0x3c, 0xda, 0xe2, 0xba, 0xf1, 0xc8, 0xb8, 0x8f, 0xf5, 0xfc, 0xb, 0x8a, 0x33, 0x9c}},
			},
		},
		{
			memo:   "Aj9GhL/4TH3a0Ry4wHOG6RkoxWdcpLxGppAElPH3PNriuvHIuI/1/AuKM5w8r/n1pxbN1wEwTq5vx/QsgP3upYQ=",
			amount: 1000,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(350),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(350),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDec(300),
			},

			expectedMemoFields: types.MemoFields{
				2: types.MemoField{ID: 2, Data: []uint8{0x46, 0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc, 0x46, 0xa6, 0x90, 0x4, 0x94, 0xf1, 0xf7, 0x3c, 0xda, 0xe2, 0xba, 0xf1, 0xc8, 0xb8, 0x8f, 0xf5, 0xfc, 0xb, 0x8a, 0x33, 0x9c, 0x3c, 0xaf, 0xf9, 0xf5, 0xa7, 0x16, 0xcd, 0xd7, 0x1, 0x30, 0x4e, 0xae, 0x6f, 0xc7, 0xf4, 0x2c, 0x80, 0xfd, 0xee, 0xa5, 0x84}},
			},
		},
		{
			memo:   "AlROhL/4TH3a0Ry4wHOG6RkoxWdcpLw0ppAElPH3PNriuvHIuI/1/AuKM5w8r/n1pxbN1wEwTq5vx/QsgP3upYQK7EkpebEEzVgFDJYdIBcOLCYxjvw=",
			amount: 10,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(39, 1),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(26, 1),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDec(3),
				"cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll": sdk.NewDecWithPrec(5, 1),
			},
			expectedMemoFields: types.MemoFields{
				2: types.MemoField{ID: 2, Data: []uint8{0x4e, 0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc, 0x34, 0xa6, 0x90, 0x4, 0x94, 0xf1, 0xf7, 0x3c, 0xda, 0xe2, 0xba, 0xf1, 0xc8, 0xb8, 0x8f, 0xf5, 0xfc, 0xb, 0x8a, 0x33, 0x9c, 0x3c, 0xaf, 0xf9, 0xf5, 0xa7, 0x16, 0xcd, 0xd7, 0x1, 0x30, 0x4e, 0xae, 0x6f, 0xc7, 0xf4, 0x2c, 0x80, 0xfd, 0xee, 0xa5, 0x84, 0xa, 0xec, 0x49, 0x29, 0x79, 0xb1, 0x4, 0xcd, 0x58, 0x5, 0xc, 0x96, 0x1d, 0x20, 0x17, 0xe, 0x2c, 0x26, 0x31, 0x8e, 0xfc}},
			},
		},
		{
			name:   "val intents and memo fields",
			memo:   "AipahL/4TH3a0Ry4wHOG6RkoxWdcpLxuppAElPH3PNriuvHIuI/1/AuKM5wAAgEC",
			amount: 100,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(45),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(55),
			},
			expectedMemoFields: types.MemoFields{
				0: types.MemoField{ID: 0, Data: []uint8{0x1, 0x2}},
				2: types.MemoField{ID: 2, Data: []uint8{0x5a, 0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc, 0x6e, 0xa6, 0x90, 0x4, 0x94, 0xf1, 0xf7, 0x3c, 0xda, 0xe2, 0xba, 0xf1, 0xc8, 0xb8, 0x8f, 0xf5, 0xfc, 0xb, 0x8a, 0x33, 0x9c}},
			},
		},
		{
			name:    "empty memo",
			memo:    "",
			wantErr: false,
		},
		{
			name:    "invalid length",
			memo:    "ToS/+Ex92tEcuMBzhukZKMVnXKS8NKaQBJTx9zza4rrxyLiP9fwLijOcPK/59acWzdcBME6ub8f0LID97qWECuxJKXmxBM1YBQyWHSAXDiwmMY78K",
			wantErr: true,
		},
		{
			name:    "invalid base64",
			memo:    "\xFF",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			memoFields, err := zone.DecodeMemo(tc.memo)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			validatorIntents, found := memoFields.Intent(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(int64(tc.amount)))), &zone)
			if len(tc.expectedIntent) > 0 {
				require.True(t, found)
				for _, v := range validatorIntents {
					if !tc.expectedIntent[v.ValoperAddress].Equal(v.Weight) {
						t.Errorf("Got %v expected %v", v.Weight, tc.expectedIntent[v.ValoperAddress])
					}
				}
			} else {
				require.False(t, found)
			}

			require.Equal(t, tc.expectedMemoFields, memoFields)
		})
	}
}

func TestUpdateIntentWithMemo(t *testing.T) {
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", Is_118: true}
	zone.Validators = append(zone.Validators,
		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)

	testCases := []struct {
		baseAmount     int
		originalIntent map[string]sdk.Dec
		memo           string
		amount         int
		expectedIntent map[string]sdk.Dec
	}{
		{
			baseAmount: 100,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
			memo: "AipahL/4TH3a0Ry4wHOG6RkoxWdcpLxuppAElPH3PNriuvHIuI/1/AuKM5w=",

			amount: 100,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
		},
		{
			baseAmount: 100,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
			memo:   "AipahL/4TH3a0Ry4wHOG6RkoxWdcpLxuppAElPH3PNriuvHIuI/1/AuKM5w=",
			amount: 1000,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
		},
		{
			baseAmount: 100,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(25, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(75, 2),
			},
			memo:   "AipahL/4TH3a0Ry4wHOG6RkoxWdcpLxuppAElPH3PNriuvHIuI/1/AuKM5w=",
			amount: 100,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(35, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(65, 2),
			},
		},
		{
			baseAmount: 1000,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(25, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(75, 2),
			},
			memo:   "Aj9GhL/4TH3a0Ry4wHOG6RkoxWdcpLxGppAElPH3PNriuvHIuI/1/AuKM5w8r/n1pxbN1wEwTq5vx/QsgP3upYQ=",
			amount: 1000,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(30, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDecWithPrec(15, 2),
			},
		},
	}

	for caseidx, tc := range testCases {
		memoFields, err := zone.DecodeMemo(tc.memo)
		require.NoError(t, err)
		memoIntent, found := memoFields.Intent(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(int64(tc.amount)))), &zone)
		require.True(t, found)
		intent := zone.UpdateZoneIntentWithMemo(memoIntent, intentFromDecSlice(tc.originalIntent), sdk.NewDec(int64(tc.baseAmount)))
		for idx, v := range intent.Intents.Sort() {
			if !tc.expectedIntent[v.ValoperAddress].Equal(v.Weight) {
				t.Errorf("Case [%d:%d] -> Got %v expected %v", caseidx, idx, v.Weight, tc.expectedIntent[v.ValoperAddress])
			}
		}
	}
}

func TestUpdateIntentWithMemoBad(t *testing.T) {
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", Is_118: true}
	zone.Validators = append(zone.Validators,
		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)

	testCases := []struct {
		baseAmount     int
		originalIntent map[string]sdk.Dec
		memo           string
		amount         int
		errorMsg       string
	}{
		{
			baseAmount: 100,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
			memo:     "WoS/+Ex92tEcuMBzhukZKMVnXKS8bqaQBJT",
			amount:   100,
			errorMsg: "unable to determine intent from memo: Failed to decode base64 message: illegal base64 data at input byte 32",
		},
	}

	for _, tc := range testCases {
		_, err := zone.DecodeMemo(tc.memo)
		require.Errorf(t, err, tc.errorMsg)
	}
}

func TestUpdateIntentWithCoins(t *testing.T) {
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", Is_118: true}
	zone.Validators = append(zone.Validators,
		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&types.Validator{ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)
	testCases := []struct {
		baseAmount     int
		originalIntent map[string]sdk.Dec
		amount         sdk.Coins
		expectedIntent map[string]sdk.Dec
	}{
		{
			baseAmount: 100,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(450)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf2", sdk.NewInt(550)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
		},
		{
			baseAmount: 100,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(45000)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf2", sdk.NewInt(55000)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(45, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
			},
		},
		{
			baseAmount: 100,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(25, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(75, 2),
			},
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(45)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf2", sdk.NewInt(55)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(35, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(65, 2),
			},
		},
		{
			baseAmount: 1000,
			originalIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(25, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(75, 2),
			},
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(350)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf2", sdk.NewInt(350)),
				sdk.NewCoin("cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy4", sdk.NewInt(300)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(30, 2),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(55, 2),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDecWithPrec(15, 2),
			},
		},
	}
	valAddresses := []string{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7"}

	for _, tc := range testCases {
		intent := zone.UpdateIntentWithCoins(intentFromDecSlice(tc.originalIntent), sdk.NewDec(int64(tc.baseAmount)), tc.amount, valAddresses)
		for _, v := range intent.Intents {
			if !tc.expectedIntent[v.ValoperAddress].Equal(v.Weight) {
				t.Errorf("Got %v expected %v", v.Weight, tc.expectedIntent[v.ValoperAddress])
			}
		}
	}
}

func intentFromDecSlice(in map[string]sdk.Dec) types.DelegatorIntent {
	out := types.DelegatorIntent{
		Delegator: addressutils.GenerateAccAddressForTest().String(),
		Intents:   []*types.ValidatorIntent{},
	}
	for addr, weight := range in {
		out.Intents = append(out.Intents, &types.ValidatorIntent{ValoperAddress: addr, Weight: weight})
	}
	return out
}

// func TestDetermineStateIntentDiff(t *testing.T) {
// 	zone := types.Zone{}
// 	d1 := []*types.Delegation{}
// 	d1 = append(d1, &types.Delegation{DelegationAddress: "cosmos1user1234", ValidatorAddress: "cosmos12345667890", Amount: sdk.NewDec(1000)})
// 	d1 = append(d1, &types.Delegation{DelegationAddress: "cosmos1user1235", ValidatorAddress: "cosmos12345667890", Amount: sdk.NewDec(500)})
// 	d1 = append(d1, &types.Delegation{DelegationAddress: "cosmos1user1236", ValidatorAddress: "cosmos12345667890", Amount: sdk.NewDec(300)})
// 	d1 = append(d1, &types.Delegation{DelegationAddress: "cosmos1user1237", ValidatorAddress: "cosmos12345667890", Amount: sdk.NewDec(200)})

// 	i1 := []types.DelegatorIntent{}
// 	i1 = append(i1, types.DelegatorIntent{Delegator: "cosmos1user1234", Intents: []*types.ValidatorIntent{{ValoperAddress: "cosmos12345667890", Weight: sdk.MustNewDecFromStr("0.5")}, {ValoperAddress: "cosmos987654321", Weight: sdk.MustNewDecFromStr("0.5")}}})
// 	i1 = append(i1, types.DelegatorIntent{Delegator: "cosmos1user1235", Intents: []*types.ValidatorIntent{{ValoperAddress: "cosmos12345667890", Weight: sdk.NewDec(1)}}})
// 	i1 = append(i1, types.DelegatorIntent{Delegator: "cosmos1user1236", Intents: []*types.ValidatorIntent{{ValoperAddress: "cosmos12345667890", Weight: sdk.NewDec(1)}}})
// 	i1 = append(i1, types.DelegatorIntent{Delegator: "cosmos1user1237", Intents: []*types.ValidatorIntent{{ValoperAddress: "cosmos12345667890", Weight: sdk.NewDec(1)}}})

// 	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmos12345667890", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewDec(2000), Delegations: d1})

// 	require.Equal(t, 0, 0)
// }

// func TestApplyDiffsToDistribution(t *testing.T) {
// 	testCases := []struct {
// 		distribution         map[string]sdk.Coin
// 		diff                 map[string]cosmosmath.Int
// 		expectedDistribution map[string]sdk.Coin
// 		expectedRemainder    cosmosmath.Int
// 	}{
// 		{
// 			distribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 3),
// 				"val2": sdk.NewInt64Coin("uatom", 3),
// 			},
// 			diff: map[string]cosmosmath.Int{
// 				"val1": sdk.NewInt(-1),
// 				"val2": sdk.NewInt(1),
// 			},
// 			expectedDistribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 4),
// 				"val2": sdk.NewInt64Coin("uatom", 2),
// 			},
// 			expectedRemainder: sdk.ZeroInt(),
// 		},

// 		{
// 			distribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 1),
// 				"val2": sdk.NewInt64Coin("uatom", 5),
// 			},
// 			diff: map[string]cosmosmath.Int{
// 				"val1": sdk.NewInt(-1),
// 				"val2": sdk.NewInt(1),
// 			},
// 			expectedDistribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 2),
// 				"val2": sdk.NewInt64Coin("uatom", 4),
// 			},
// 			expectedRemainder: sdk.ZeroInt(),
// 		},
// 		{
// 			distribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 1),
// 				"val2": sdk.NewInt64Coin("uatom", 5),
// 			},
// 			diff: map[string]cosmosmath.Int{
// 				"val1": sdk.NewInt(2),
// 				"val2": sdk.NewInt(2),
// 				"val3": sdk.NewInt(-4),
// 				"val4": sdk.NewInt(0),
// 			},
// 			expectedDistribution: map[string]sdk.Coin{
// 				"val2": sdk.NewInt64Coin("uatom", 3),
// 			},
// 			expectedRemainder: sdk.NewInt(3),
// 		},
// 		{
// 			distribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 1),
// 				"val2": sdk.NewInt64Coin("uatom", 5),
// 				"val3": sdk.NewInt64Coin("uatom", 0),
// 			},
// 			diff: map[string]cosmosmath.Int{
// 				"val1": sdk.NewInt(2),
// 				"val2": sdk.NewInt(2),
// 				"val3": sdk.NewInt(-4),
// 				"val4": sdk.NewInt(0),
// 			},
// 			expectedDistribution: map[string]sdk.Coin{
// 				"val2": sdk.NewInt64Coin("uatom", 3),
// 				"val3": sdk.NewInt64Coin("uatom", 3),
// 			},
// 			expectedRemainder: sdk.ZeroInt(),
// 		},
// 	}

// 	zone := types.Zone{}

// 	for idx, i := range testCases {
// 		result, remainder := zone.ApplyDiffsToDistribution(i.distribution, i.diff)
// 		for k, v := range i.expectedDistribution {
// 			require.Truef(t, v.IsEqual(result[k]), "case %d: distribution %v does not match expected %v", idx, result[k], v)
// 		}
// 		require.Truef(t, i.expectedRemainder.Equal(remainder), "case %d: remainder %v does not match expected %v", idx, remainder, i.expectedRemainder)
// 	}

// }

func TestParseMemoFields(t *testing.T) {
	testCases := []struct {
		name               string
		fieldBytes         []byte
		expectedMemoFields types.MemoFields
		wantErr            bool
	}{
		{
			name:       "invalid no length data",
			fieldBytes: []byte{},
			wantErr:    true,
		},
		{
			name:       "invalid length field",
			fieldBytes: []byte{byte(types.FieldTypeAccountMap), 1, 0, 0},
			wantErr:    true,
		},
		{
			name: "invalid multiple",
			fieldBytes: []byte{
				byte(types.FieldTypeAccountMap), 3, 0, 0, // should be 2 for length field
				byte(types.FieldTypeReturnToSender), 4, 1, 1, 1, 3,
			},
			wantErr: true,
		},
		{
			name: "invalid address for account map",
			fieldBytes: []byte{
				byte(types.FieldTypeAccountMap), 0,
				byte(types.FieldTypeReturnToSender), 0,
			},
			wantErr: true,
		},
		{
			name: "invalid field id 3",
			fieldBytes: []byte{
				3, 2, 1, 1,
				byte(types.FieldTypeReturnToSender), 0,
			},
			wantErr: true,
		},
		{
			name:       "valid single",
			fieldBytes: []byte{byte(types.FieldTypeAccountMap), 2, 1, 1},
			expectedMemoFields: types.MemoFields{
				types.FieldTypeAccountMap: types.MemoField{
					ID:   0,
					Data: []byte{1, 1},
				},
			},
			wantErr: false,
		},
		{
			name: "valid multiple",
			fieldBytes: []byte{
				byte(types.FieldTypeAccountMap), 2, 1, 1,
				byte(types.FieldTypeReturnToSender), 0,
			},
			expectedMemoFields: types.MemoFields{
				types.FieldTypeAccountMap: {
					ID:   types.FieldTypeAccountMap,
					Data: []byte{1, 1},
				},
				types.FieldTypeReturnToSender: {
					ID:   types.FieldTypeReturnToSender,
					Data: nil,
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := types.ParseMemoFields(tc.fieldBytes)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedMemoFields, out)
		})
	}
}
