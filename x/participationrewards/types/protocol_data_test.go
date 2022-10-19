package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// tests that {} is an invalid string, and that an error is thrown when unmarshalled.
// see: https://github.com/ingenuity-build/quicksilver/issues/214
func TestUnmarshalProtocolDataRejectsZeroLengthJson(t *testing.T) {
	_, err := UnmarshalProtocolData(ProtocolDataTypeOsmosisPool, []byte("{}"))
	require.Error(t, err)
}

func TestConnectionProtocolData_ValidateBasic(t *testing.T) {
	type fields struct {
		ConnectionID string
		ChainID      string
		LastEpoch    int64
		Prefix       string
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
			"valid",
			fields{
				ConnectionID: "connection-0",
				ChainID:      "testchain-1",
				LastEpoch:    30000,
				Prefix:       "cosmos",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpd := ConnectionProtocolData{
				ConnectionID: tt.fields.ConnectionID,
				ChainID:      tt.fields.ChainID,
				LastEpoch:    tt.fields.LastEpoch,
				Prefix:       tt.fields.Prefix,
			}
			err := cpd.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUnmarshalProtocolData(t *testing.T) {
	type args struct {
		datatype ProtocolDataType
		data     json.RawMessage
	}
	tests := []struct {
		name    string
		args    args
		want    ProtocolDataI
		wantErr bool
	}{
		{
			"blank",
			args{},
			nil,
			true,
		},
		{
			"unknown_protocol_type",
			args{
				datatype: 99999,
				data:     []byte{},
			},
			nil,
			true,
		},
		{
			"connection_data_empty",
			args{
				datatype: 0,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"connection_data",
			args{
				datatype: ProtocolDataTypeConnection,
				data:     []byte(`{"connectionid": "connection-0","chainid": "somechain","lastepoch": 10000}`),
			},
			ConnectionProtocolData{
				ConnectionID: "connection-0",
				ChainID:      "somechain",
				LastEpoch:    10000,
			},
			false,
		},
		{
			"liquid_data_empty",
			args{
				datatype: ProtocolDataTypeLiquidToken,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"liquid_data",
			args{
				datatype: ProtocolDataTypeLiquidToken,
				data:     []byte(`{"chainid": "somechain","localdenom": "lstake","denom": "qstake"}`),
			},
			LiquidAllowedDenomProtocolData{
				ChainID:    "somechain",
				Denom:      "qstake",
				LocalDenom: "lstake",
			},
			false,
		},
		{
			"osmosispool_data_empty",
			args{
				datatype: ProtocolDataTypeOsmosisPool,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"osmosispool_data",
			args{
				datatype: ProtocolDataTypeOsmosisPool,
				data:     []byte(`{"poolid": 1, "poolname": "atom/osmo","zones": {"cosmoshub-4": "IBC/atom_denom", "osmosis-1": "IBC/osmo_denom"}}`),
			},
			OsmosisPoolProtocolData{
				PoolID:   1,
				PoolName: "atom/osmo",
				Zones:    map[string]string{"cosmoshub-4": "IBC/atom_denom", "osmosis-1": "IBC/osmo_denom"},
			},
			false,
		},
		{
			"osmosis_params_empty",
			args{
				datatype: ProtocolDataTypeOsmosisParams,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"osmosis_params",
			args{
				datatype: ProtocolDataTypeOsmosisParams,
				data:     []byte(`{"ChainID": "test-01"}`),
			},
			OsmosisParamsProtocolData{
				ChainID: "test-01",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalProtocolData(tt.args.datatype, tt.args.data)
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
