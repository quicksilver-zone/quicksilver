package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func TestParseStakingDelegationKeyValid(t *testing.T) {
	delAddr, err := sdk.AccAddressFromBech32("cosmos1zcuaqawcpzn7q9wmulagvjjv7f72qearnep4jt")
	require.NoError(t, err, "failed to parse delAddress from bech32")
	valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1zcuaqawcpzn7q9wmulagvjjv7f72qearkd4q7c")
	require.NoError(t, err, "failed to parse valAddress from bech32")
	key := stakingtypes.GetDelegationKey(delAddr, valAddr)
	del, val, err := types.ParseStakingDelegationKey(key)
	require.NoError(t, err, "expected no error in ParseStakingDelegationKey()")
	require.Equal(t, delAddr, del, "require original and parsed delegator addresses match")
	require.Equal(t, valAddr, val, "require original and parsed validator addresses match")
}

func TestParseStakingDelegationKeyInvalidLength(t *testing.T) {
	var key []byte
	_, _, err := types.ParseStakingDelegationKey(key)
	require.ErrorContains(t, err, "out of bounds reading byte 0")

	key = []byte{0x31}
	_, _, err = types.ParseStakingDelegationKey(key)
	require.ErrorContains(t, err, "out of bounds reading delegator address length")

	key = []byte{0x31, 0x42}
	_, _, err = types.ParseStakingDelegationKey(key)
	require.ErrorContains(t, err, "invalid delegator address length")
}

func TestParseStakingDelegationKeyInvalidPrefix(t *testing.T) {
	key := []byte{0x42}
	_, _, err := types.ParseStakingDelegationKey(key)
	require.Errorf(t, err, "not a valid delegation key")
}

func TestParseStakingDelegationKeyInvalidTruncated(t *testing.T) {
	delAddr, err := sdk.AccAddressFromBech32("cosmos1zcuaqawcpzn7q9wmulagvjjv7f72qearnep4jt")
	require.NoError(t, err, "failed to parse delAddress from bech32")
	valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1zcuaqawcpzn7q9wmulagvjjv7f72qearkd4q7c")
	require.NoError(t, err, "failed to parse valAddress from bech32")
	key := stakingtypes.GetDelegationKey(delAddr, valAddr)
	// truncate the last byte of the key.
	_, _, err = types.ParseStakingDelegationKey(key[:len(key)-1])
	require.Error(t, err, "out of bounds reading validator address")
}
