package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func TestGenesisState(t *testing.T) {
	// test default genesis state
	testGenesisState := types.GenesisState{
		Params: types.Params{
			DistributionProportions: types.DistributionProportions{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.34"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("0.33"),
				LockupAllocation:             sdk.MustNewDecFromStr("0.33"),
			},
		},
	}
	defaultGenesisState := types.DefaultGenesisState()
	require.Equal(t, *defaultGenesisState, testGenesisState)
	// test new genesis state
	newGenesisState := types.NewGenesisState(
		types.Params{
			DistributionProportions: types.DistributionProportions{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.5"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("0.3"),
				LockupAllocation:             sdk.MustNewDecFromStr("0.2"),
			},
		},
	)
	testGenesisState = types.GenesisState{
		Params: types.Params{
			DistributionProportions: types.DistributionProportions{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.5"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("0.3"),
				LockupAllocation:             sdk.MustNewDecFromStr("0.2"),
			},
		},
	}
	require.Equal(t, *newGenesisState, testGenesisState)
}

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		Params       types.Params
		ProtocolData []*types.KeyedProtocolData
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"defaults",
			fields{},
			true,
		},
		{
			"invalid_protocolData",
			fields{
				Params: types.DefaultParams(),
				ProtocolData: []*types.KeyedProtocolData{
					{
						"liquid",
						&types.ProtocolData{
							Type: "liquidtoken",
							Data: []byte("{}"),
						},
					},
				},
			},
			true,
		},
		{
			"valid_protocolData",
			fields{
				Params: types.DefaultParams(),
				ProtocolData: []*types.KeyedProtocolData{
					{
						"liquid",
						&types.ProtocolData{
							Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)],
							Data: []byte(validLiquidData),
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := types.GenesisState{
				Params:       tt.fields.Params,
				ProtocolData: tt.fields.ProtocolData,
			}

			err := gs.Validate()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
