package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// helper function for ICA tests.
func NewICA() *types.ICAAccount {
	addr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	accAddr := sdk.AccAddress(addr)
	ica, err := types.NewICAAccount(accAddr.String(), "mercury-1.deposit")
	if err != nil {
		panic("failed to create ICA")
	}
	return ica
}

func TestNewICAAccountBadAddr(t *testing.T) {
	ica, err := types.NewICAAccount("cosmos1ssrxxe4xsls57ehrkswlkhlk", "mercury-1.deposit")
	require.Nil(t, ica, "expecting a nil ICAAccount")
	require.NotNil(t, err, "expecting a non-nil error")
	require.Contains(t, err.Error(), "invalid checksum")
}

// TestAccountSetBalanceGood tests that the balance can be set to a valid coin (good denom + non-negative value).
func TestAccountSetBalanceGood(t *testing.T) {
	ica := NewICA()
	err := ica.SetBalance(sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(300))))
	require.NoError(t, err, "setbalance failed")
	require.True(t, ica.Balance.AmountOf("uqck").Equal(sdk.NewInt(300)))
}

// tests that the balance panics when set to an invalid denomination.
func TestAccountSetBalanceBadDenom(t *testing.T) {
	ica := NewICA()
	require.PanicsWithError(t, "invalid denom: _fail", func() { ica.SetBalance(sdk.NewCoins(sdk.NewCoin("_fail", sdk.NewInt(300)))) })
}

// tests that the balance panics when set to a negative number.
func TestAccountSetBalanceNegativeAmount(t *testing.T) {
	ica := NewICA()
	require.PanicsWithError(t, "negative coin amount: -300", func() { ica.SetBalance(sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(-300)))) })
}

// tests that the balance panics when set to a negative number.
func TestAccountSetBalanceNonSortedCoins(t *testing.T) {
	ica := NewICA()
	nonSortedCoins := sdk.Coins{
		sdk.NewCoin("uqck", sdk.NewInt(300)),
		sdk.NewCoin("uqck", sdk.NewInt(200)),
	}
	err := ica.SetBalance(nonSortedCoins)
	require.NotNil(t, err, "non sorted coins should return an error")
}

func TestAccountSetWithdrawalAddress(t *testing.T) {
	ica := NewICA()

	cases := []struct {
		name    string
		addr    string
		want    string
		wantErr string
	}{
		{"empty address", "    ", "", "empty address string"},
		{
			name: "valid address",
			addr: "cosmos1ssrxxe4xsls57ehrkswlkhlkcverf0p0fpgyhzqw0hfdqj92ynxsw29r6e",
			want: "cosmos1ssrxxe4xsls57ehrkswlkhlkcverf0p0fpgyhzqw0hfdqj92ynxsw29r6e",
		},
		{
			name:    "non-empty but invalid address",
			addr:    "cosmos1ssrxxe4xsls57ehrkswlkhlk",
			wantErr: "invalid checksum",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ica.SetWithdrawalAddress(tc.addr)
			if tc.wantErr != "" {
				require.NotNil(t, err, "expected a non-nil error")
				require.Contains(t, err.Error(), tc.wantErr)
				return
			}

			// Otherwise not expecting an error.
			require.Nil(t, err, "expected a nil error")
			require.Equal(t, ica.WithdrawalAddress, tc.want, "addresses must be the same")
		})
	}
}

// tests balance waitgroup increments and decrements as expected and errors on negative wg value.
func TestIncrementDecrementWg(t *testing.T) {
	ica := NewICA()
	oldWg := ica.BalanceWaitgroup
	ica.IncrementBalanceWaitgroup()
	firstWg := ica.BalanceWaitgroup
	require.Equal(t, oldWg+1, firstWg)
	require.NoError(t, ica.DecrementBalanceWaitgroup())
	secondWg := ica.BalanceWaitgroup
	require.Equal(t, firstWg-1, secondWg)
	require.Error(t, ica.DecrementBalanceWaitgroup())
}
