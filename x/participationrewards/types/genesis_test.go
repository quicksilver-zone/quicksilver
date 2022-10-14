package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/utils"
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
		nil,
	}
	require.Equal(t, *newGenesisState, testGenesisState)
}

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		Params       Params
		Claims       []*Claim
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
			"invalid_claim",
			fields{
				Params: DefaultParams(),
				Claims: []*Claim{
					{
						UserAddress: utils.GenerateAccAddressForTest().String(),
						ChainId:     "testzone-1",
						Amount:      0,
					},
				},
				ProtocolData: nil,
			},
			true,
		},
		{
			"valid_claim",
			fields{
				Params: DefaultParams(),
				Claims: []*Claim{
					{
						UserAddress: utils.GenerateAccAddressForTest().String(),
						ChainId:     "testzone-1",
						Amount:      1000000,
					},
				},
				ProtocolData: nil,
			},
			false,
		},
		{
			"invalid_protocolData",
			fields{
				Params: DefaultParams(),
				Claims: nil,
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
				Claims: nil,
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
		{
			"valid_claim",
			fields{
				Params: DefaultParams(),
				Claims: []*Claim{
					{
						UserAddress: utils.GenerateAccAddressForTest().String(),
						ChainId:     "testzone-1",
						Amount:      1000000,
					},
				},
				ProtocolData: nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GenesisState{
				Params:       tt.fields.Params,
				Claims:       tt.fields.Claims,
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
