package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func TestDistributionProportions_ValidateBasic(t *testing.T) {
	type fields struct {
		ValidatorSelectionAllocation sdk.Dec
		HoldingsAllocation           sdk.Dec
		LockupAllocation             sdk.Dec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"blank",
			fields{},
			true,
		},
		{
			"invalid_proportions_gt",
			fields{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.5"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("0.5"),
				LockupAllocation:             sdk.MustNewDecFromStr("0.5"),
			},
			true,
		},
		{
			"invalid_proportions_lt",
			fields{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.3"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("0.3"),
				LockupAllocation:             sdk.MustNewDecFromStr("0.3"),
			},
			true,
		},
		{
			"invalid_proportions_negative",
			fields{
				ValidatorSelectionAllocation: sdk.MustNewDecFromStr("-0.4"),
				HoldingsAllocation:           sdk.MustNewDecFromStr("-0.3"),
				LockupAllocation:             sdk.MustNewDecFromStr("-0.3"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := types.DistributionProportions{
				ValidatorSelectionAllocation: tt.fields.ValidatorSelectionAllocation,
				HoldingsAllocation:           tt.fields.HoldingsAllocation,
				LockupAllocation:             tt.fields.LockupAllocation,
			}
			err := dp.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestKeyedProtocolData_ValidateBasic(t *testing.T) {
	invalidOsmosisData := `{
	"poolname": "osmosispools/1",
	"zones": {
		"": ""
	}
}`
	validOsmosisData := `{
	"poolid": 1,
	"poolname": "atom/osmo",
	"pooltype": "balancer",
	"zones": {
		"zone_id": "IBC/zone_denom"
	}
}`
	validLiquidData := `{
	"chainid": "somechain-1",
	"registeredzonechainid": "someotherchain-1",
	"ibcdenom": "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
	"qassetdenom": "uqstake"
}`
	type fields struct {
		Key          string
		ProtocolData *types.ProtocolData
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"blank",
			fields{},
			true,
		},
		{
			"blank_pd",
			fields{
				"somekey",
				&types.ProtocolData{},
			},
			true,
		},
		{
			"pd_osmosis_nil_data",
			fields{
				"osmosispools/1",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)],
					Data: nil,
				},
			},
			true,
		},
		{
			"pd_osmosis_empty_data",
			fields{
				"osmosispools/1",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)],
					Data: []byte("{}"),
				},
			},
			true,
		},
		{
			"pd_osmosis_invalid",
			fields{
				"osmosispools/1",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)],
					Data: []byte(invalidOsmosisData),
				},
			},
			true,
		},
		{
			"pd_osmosis_valid",
			fields{
				"osmosispools/1",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)],
					Data: []byte(validOsmosisData),
				},
			},
			false,
		},
		{
			"pd_liquid_invalid",
			fields{
				"liquid",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)],
					Data: []byte("{}"),
				},
			},
			true,
		},
		{
			"pd_liquid_valid",
			fields{
				"liquid",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)],
					Data: []byte(validLiquidData),
				},
			},
			false,
		},
		{
			"pd_unknown",
			fields{
				"unknown",
				&types.ProtocolData{
					Type: "unknown",
					Data: []byte("{}"),
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kpd := types.KeyedProtocolData{
				Key:          tt.fields.Key,
				ProtocolData: tt.fields.ProtocolData,
			}
			err := kpd.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
