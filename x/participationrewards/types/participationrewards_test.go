package types

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := DistributionProportions{
				ValidatorSelectionAllocation: tt.fields.ValidatorSelectionAllocation,
				HoldingsAllocation:           tt.fields.HoldingsAllocation,
				LockupAllocation:             tt.fields.LockupAllocation,
			}
			if err := dp.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("DistributionProportions.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParams_ValidateBasic(t *testing.T) {
	type fields struct {
		DistributionProportions DistributionProportions
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
			p := Params{
				DistributionProportions: tt.fields.DistributionProportions,
			}
			if err := p.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("Params.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Claim{
				UserAddress: tt.fields.UserAddress,
				ChainId:     tt.fields.ChainId,
				Amount:      tt.fields.Amount,
			}
			if err := c.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("Claim.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
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
