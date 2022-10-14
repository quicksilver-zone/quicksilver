package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/stretchr/testify/require"
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
			dp := DistributionProportions{
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

func TestClaim_ValidateBasic(t *testing.T) {
	type fields struct {
		UserAddress string
		ChainId     string
		Amount      uint64
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
			"invalid_address",
			fields{
				UserAddress: "cosmos1234567890",
				ChainId:     "testzone-1",
				Amount:      10000,
			},
			true,
		},
		{
			"invalid_chain_id",
			fields{
				UserAddress: utils.GenerateAccAddressForTest().String(),
				ChainId:     "",
				Amount:      10000,
			},
			true,
		},
		{
			"invalid_chain_id",
			fields{
				UserAddress: utils.GenerateAccAddressForTest().String(),
				ChainId:     "",
				Amount:      10000,
			},
			true,
		},
		{
			"invalid_amount",
			fields{
				UserAddress: utils.GenerateAccAddressForTest().String(),
				ChainId:     "testzone-1",
				Amount:      0,
			},
			true,
		},
		{
			"valid",
			fields{
				UserAddress: utils.GenerateAccAddressForTest().String(),
				ChainId:     "testzone-1",
				Amount:      1000000,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Claim{
				UserAddress: tt.fields.UserAddress,
				ChainId:     tt.fields.ChainId,
				Amount:      tt.fields.Amount,
			}
			err := c.ValidateBasic()
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
	"zones": {
		"zone_id": "IBC/zone_denom"
	}
}`
	validLiquidData := `{
	"chainid": "somechain",
	"localdenom": "lstake",
	"denom": "qstake"
}`
	type fields struct {
		Key          string
		ProtocolData *ProtocolData
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
				&ProtocolData{},
			},
			true,
		},
		{
			"pd_osmosis_nil_data",
			fields{
				"osmosispools/1",
				&ProtocolData{
					Type: ProtocolDataType_name[int32(ProtocolDataTypeOsmosisPool)],
					Data: nil,
				},
			},
			true,
		},
		{
			"pd_osmosis_empty_data",
			fields{
				"osmosispools/1",
				&ProtocolData{
					Type: ProtocolDataType_name[int32(ProtocolDataTypeOsmosisPool)],
					Data: []byte("{}"),
				},
			},
			true,
		},
		{
			"pd_osmosis_invalid",
			fields{
				"osmosispools/1",
				&ProtocolData{
					Type: ProtocolDataType_name[int32(ProtocolDataTypeOsmosisPool)],
					Data: []byte(invalidOsmosisData),
				},
			},
			true,
		},
		{
			"pd_osmosis_valid",
			fields{
				"osmosispools/1",
				&ProtocolData{
					Type: ProtocolDataType_name[int32(ProtocolDataTypeOsmosisPool)],
					Data: []byte(validOsmosisData),
				},
			},
			false,
		},
		{
			"pd_liquid_invalid",
			fields{
				"liquid",
				&ProtocolData{
					Type: ProtocolDataType_name[int32(ProtocolDataTypeLiquidToken)],
					Data: []byte("{}"),
				},
			},
			true,
		},
		{
			"pd_liquid_valid",
			fields{
				"liquid",
				&ProtocolData{
					Type: ProtocolDataType_name[int32(ProtocolDataTypeLiquidToken)],
					Data: []byte(validLiquidData),
				},
			},
			false,
		},
		{
			"pd_unknown",
			fields{
				"unknown",
				&ProtocolData{
					Type: "unknown",
					Data: []byte("{}"),
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kpd := KeyedProtocolData{
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
