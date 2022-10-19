package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState(t *testing.T) {
	// test default genesis state
	testGenesisState := GenesisState{
		Params{
			DistributionProportions: DistributionProportions{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.34"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("0.33"),
				LockupAllocation:             sdk.MustNewDecFromStr("0.33"),
			},
		},
		nil,
	}
	defaultGenesisState := DefaultGenesisState()
	require.Equal(t, *defaultGenesisState, testGenesisState)
	// test new genesis state
	newGenesisState := NewGenesisState(
		Params{
			DistributionProportions: DistributionProportions{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.5"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("0.3"),
				LockupAllocation:             sdk.MustNewDecFromStr("0.2"),
			},
		},
	)
	testGenesisState = GenesisState{
		Params{
			DistributionProportions: DistributionProportions{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.5"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("0.3"),
				LockupAllocation:             sdk.MustNewDecFromStr("0.2"),
			},
		},
		nil,
	}
	require.Equal(t, *newGenesisState, testGenesisState)
}

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		Params       Params
		ProtocolData []*KeyedProtocolData
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
				Params: DefaultParams(),
				ProtocolData: []*KeyedProtocolData{
					{
						"liquid",
						&ProtocolData{
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
				Params: DefaultParams(),
				ProtocolData: []*KeyedProtocolData{
					{
						"liquid",
						&ProtocolData{
							Type: ProtocolDataType_name[int32(ProtocolDataTypeLiquidToken)],
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
			gs := GenesisState{
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
