package types_test

import (
	"testing"

	"github.com/quicksilver-zone/quicksilver/x/mint/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisValidate(t *testing.T) {
	invalidParams := types.DefaultGenesis()
	invalidParams.Params.MintDenom = "" // cannot be empty

	invalidMinter := types.DefaultGenesis()
	invalidMinter.Minter.EpochProvisions = sdk.NewDec(-1) // cannot be empty

	tests := []struct {
		name    string
		genesis *types.GenesisState
		isValid bool
	}{
		{
			name:    "valid genesis",
			genesis: types.DefaultGenesis(),
			isValid: true,
		},
		{
			name:    "invalid params",
			genesis: invalidParams,
			isValid: false,
		},
		{
			name:    "invalid minter",
			genesis: invalidMinter,
			isValid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.Validate()
			if !tc.isValid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
