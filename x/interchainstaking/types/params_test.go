package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func TestValidateParams(t *testing.T) {
	require.NoError(t, types.DefaultParams().Validate(), "default")
	require.NoError(t, types.NewParams(1, 1, sdkmath.LegacyNewDec(1), true).Validate(), "valid")
	require.Error(t, types.NewParams(0, 1, sdkmath.LegacyNewDec(1), true).Validate(), "0 deposit interval")
	require.Error(t, types.NewParams(1, 0, sdkmath.LegacyNewDec(1), true).Validate(), "0 valset interval")
	require.Error(t, types.NewParams(1, 1, sdkmath.LegacyNewDec(-1), true).Validate(), "negative commission rate")
}
