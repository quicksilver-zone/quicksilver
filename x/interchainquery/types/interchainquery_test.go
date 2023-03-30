package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestQuery_ValidateBasic(t *testing.T) {
	type fields struct {
		ID           string
		ConnectionID string
		ChainID      string
		QueryType    string
		Request      []byte
		Period       sdkmath.Int
		LastHeight   sdkmath.Int
		CallbackID   string
		TTL          uint64
		LastEmission sdkmath.Int
	}
	var tests []struct {
		name    string
		fields  fields
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Query{
				Id:           tt.fields.ID,
				ConnectionId: tt.fields.ConnectionID,
				ChainId:      tt.fields.ChainID,
				QueryType:    tt.fields.QueryType,
				Request:      tt.fields.Request,
				Period:       tt.fields.Period,
				LastHeight:   tt.fields.LastHeight,
				CallbackId:   tt.fields.CallbackID,
				Ttl:          tt.fields.TTL,
				LastEmission: tt.fields.LastEmission,
			}

			err := q.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
		})
	}
}

func TestDataPoint_ValidateBasic(t *testing.T) {
	type fields struct {
		ID           string
		RemoteHeight sdkmath.Int
		LocalHeight  sdkmath.Int
		Value        []byte
	}
	var tests []struct {
		name    string
		fields  fields
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := DataPoint{
				Id:           tt.fields.ID,
				RemoteHeight: tt.fields.RemoteHeight,
				LocalHeight:  tt.fields.LocalHeight,
				Value:        tt.fields.Value,
			}

			err := dp.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.NoError(t, err)
				return
			}
		})
	}
}
