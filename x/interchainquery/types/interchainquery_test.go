package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestQuery_ValidateBasic(t *testing.T) {
	type fields struct {
		ID           string
		ConnectionID string
		ChainID      string
		QueryType    string
		Request      []byte
		Period       math.Int
		LastHeight   math.Int
		CallbackID   string
		TTL          uint64
		LastEmission math.Int
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
		RemoteHeight math.Int
		LocalHeight  math.Int
		Value        []byte
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
