package types

import (
	"encoding/json"
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
	type fields struct {
		Key          string
		ProtocolData *ProtocolData
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kpd := KeyedProtocolData{
				Key:          tt.fields.Key,
				ProtocolData: tt.fields.ProtocolData,
			}
			if err := kpd.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("KeyedProtocolData.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProtocolData_ValidateBasic(t *testing.T) {
	type fields struct {
		Protocol string
		Type     string
		Data     json.RawMessage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pd := ProtocolData{
				Protocol: tt.fields.Protocol,
				Type:     tt.fields.Type,
				Data:     tt.fields.Data,
			}
			if err := pd.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("ProtocolData.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
