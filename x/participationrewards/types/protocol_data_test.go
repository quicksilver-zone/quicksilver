package types_test

import (
	"encoding/json"
	"testing"

	liquiditytypes "github.com/ingenuity-build/quicksilver/crescent-types/liquidity/types"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"

	"github.com/stretchr/testify/require"
)

// tests that {} is an invalid string, and that an error is thrown when unmarshalled.
// see: https://github.com/ingenuity-build/quicksilver/issues/214
func TestUnmarshalProtocolDataRejectsZeroLengthJson(t *testing.T) {
	_, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisPool, []byte("{}"))
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
			cpd := types.ConnectionProtocolData{
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

func TestLiquidProtocolData_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		pd      types.LiquidAllowedDenomProtocolData
		wantErr bool
	}{
		{
			"liquid_data",
			types.LiquidAllowedDenomProtocolData{
				ChainID:               "somechain-1",
				IbcDenom:              "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
				QAssetDenom:           "uqstake",
				RegisteredZoneChainID: "testzone-1",
			},
			false,
		},
		{
			"liquid_invalid_chainid",
			types.LiquidAllowedDenomProtocolData{
				ChainID:               "badzone",
				IbcDenom:              "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
				QAssetDenom:           "uqstake",
				RegisteredZoneChainID: "testzone-1",
			},
			true,
		},
		{
			"liquid_invalid_rzchainid",
			types.LiquidAllowedDenomProtocolData{
				ChainID:               "somechain-1",
				IbcDenom:              "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
				QAssetDenom:           "uqstake",
				RegisteredZoneChainID: "badzone",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pd.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func marshalledUmeeData[V types.UmeeInterestScalarProtocolData | types.UmeeUTokenSupplyProtocolData | types.UmeeLeverageModuleBalanceProtocolData | types.UmeeReservesProtocolData | types.UmeeTotalBorrowsProtocolData](data types.UmeeProtocolData) []byte {
	result, _ := json.Marshal(&V{UmeeProtocolData: data})
	return result
}

func TestUnmarshalProtocolData(t *testing.T) {
	testUmeeData := types.UmeeProtocolData{Denom: "test", Data: []byte{0x6e, 0x75, 0x6c, 0x6c}}
	testAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")

	type args struct {
		datatype types.ProtocolDataType
		data     json.RawMessage
	}
	tests := []struct {
		name    string
		args    args
		want    types.ProtocolDataI
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
				datatype: types.ProtocolDataTypeConnection,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"connection_data",
			args{
				datatype: types.ProtocolDataTypeConnection,
				data:     []byte(`{"connectionid": "connection-0","chainid": "somechain","lastepoch": 10000}`),
			},
			&types.ConnectionProtocolData{
				ConnectionID: "connection-0",
				ChainID:      "somechain",
				LastEpoch:    10000,
			},
			false,
		},
		{
			"liquid_data_empty",
			args{
				datatype: types.ProtocolDataTypeLiquidToken,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"liquid_data",
			args{
				datatype: types.ProtocolDataTypeLiquidToken,
				data:     []byte(`{"chainid": "somechain-1","ibcdenom": "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3","qassetdenom": "uqstake", "registeredzonechainid": "registeredzone-1"}`),
			},
			&types.LiquidAllowedDenomProtocolData{
				ChainID:               "somechain-1",
				IbcDenom:              "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
				QAssetDenom:           "uqstake",
				RegisteredZoneChainID: "registeredzone-1",
			},
			false,
		},
		{
			"osmosispool_data_empty",
			args{
				datatype: types.ProtocolDataTypeOsmosisPool,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"osmosispool_data",
			args{
				datatype: types.ProtocolDataTypeOsmosisPool,
				data:     []byte(`{"poolid": 1, "poolname": "atom/osmo","denoms": {"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {"denom": "uatom", "chainid": "cosmoshub-4"}, "uosmo": {"denom": "uosmo", "chainid": "osmosis-1"}}}`),
			},
			&types.OsmosisPoolProtocolData{
				PoolID:   1,
				PoolName: "atom/osmo",
				Denoms: map[string]types.DenomWithZone{
					"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {Denom: "uatom", ChainID: "cosmoshub-4"},
					"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
				},
			},
			false,
		},
		{
			"osmosis_params_empty",
			args{
				datatype: types.ProtocolDataTypeOsmosisParams,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"osmosis_params",
			args{
				datatype: types.ProtocolDataTypeOsmosisParams,
				data:     []byte(`{"ChainID": "test-01"}`),
			},
			&types.OsmosisParamsProtocolData{
				ChainID: "test-01",
			},
			false,
		},
		{
			"umee_params_empty",
			args{
				datatype: types.ProtocolDataTypeUmeeParams,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"umee_params",
			args{
				datatype: types.ProtocolDataTypeUmeeParams,
				data:     []byte(`{"ChainID": "test-01"}`),
			},
			&types.UmeeParamsProtocolData{
				ChainID: "test-01",
			},
			false,
		},
		{
			"umee_interest_scalar_empty",
			args{
				datatype: types.ProtocolDataTypeUmeeInterestScalar,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"umee_interest_scalar",
			args{
				datatype: types.ProtocolDataTypeUmeeInterestScalar,
				data:     marshalledUmeeData[types.UmeeInterestScalarProtocolData](testUmeeData),
			},
			&types.UmeeInterestScalarProtocolData{testUmeeData},
			false,
		},
		{
			"umee_utoken_supply_empty",
			args{
				datatype: types.ProtocolDataTypeUmeeUTokenSupply,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"umee_utoken_supply",
			args{
				datatype: types.ProtocolDataTypeUmeeUTokenSupply,
				data:     marshalledUmeeData[types.UmeeUTokenSupplyProtocolData](testUmeeData),
			},
			&types.UmeeUTokenSupplyProtocolData{testUmeeData},
			false,
		},
		{
			"umee_leverage_module_balance_empty",
			args{
				datatype: types.ProtocolDataTypeUmeeLeverageModuleBalance,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"umee_leverage_module_balance",
			args{
				datatype: types.ProtocolDataTypeUmeeLeverageModuleBalance,
				data:     marshalledUmeeData[types.UmeeLeverageModuleBalanceProtocolData](testUmeeData),
			},
			&types.UmeeLeverageModuleBalanceProtocolData{testUmeeData},
			false,
		},
		{
			"umee_reserves_data_empty",
			args{
				datatype: types.ProtocolDataTypeUmeeReserves,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"umee_reserves_data",
			args{
				datatype: types.ProtocolDataTypeUmeeReserves,
				data:     marshalledUmeeData[types.UmeeReservesProtocolData](testUmeeData),
			},
			&types.UmeeReservesProtocolData{testUmeeData},
			false,
		},
		{
			"umee_total_borrows_empty",
			args{
				datatype: types.ProtocolDataTypeUmeeTotalBorrows,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"umee_total_borrows",
			args{
				datatype: types.ProtocolDataTypeUmeeTotalBorrows,
				data:     marshalledUmeeData[types.UmeeTotalBorrowsProtocolData](testUmeeData),
			},
			&types.UmeeTotalBorrowsProtocolData{testUmeeData},
			false,
		},
		{
			"crescent_pool_empty",
			args{
				datatype: types.ProtocolDataTypeCrescentPool,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"crescent_pool",
			args{
				datatype: types.ProtocolDataTypeCrescentPool,
				data: func() []byte {
					pool := &liquiditytypes.Pool{
						PoolCoinDenom: "pool1",
						Id:            1,
					}
					pooldata, _ := json.Marshal(pool)
					pd := types.CrescentPoolProtocolData{
						PoolID: 1,
						Denoms: map[string]types.DenomWithZone{
							"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3": {ChainID: "cosmoshub-4", Denom: "uatom"},
							"ibc/15E9C5CF5969080539DB395FA7D9C0868265217EFC528433671AAF9B1912D159": {ChainID: "osmosis-1", Denom: "uosmo"},
						},
						PoolData: pooldata,
					}
					data, _ := json.Marshal(&pd)
					return data
				}(),
			},
			&types.CrescentPoolProtocolData{
				PoolID: 1,
				Denoms: map[string]types.DenomWithZone{
					"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3": {ChainID: "cosmoshub-4", Denom: "uatom"},
					"ibc/15E9C5CF5969080539DB395FA7D9C0868265217EFC528433671AAF9B1912D159": {ChainID: "osmosis-1", Denom: "uosmo"},
				},
				PoolData: func() []byte {
					pool := &liquiditytypes.Pool{
						PoolCoinDenom: "pool1",
						Id:            1,
					}
					pooldata, _ := json.Marshal(pool)
					return pooldata
				}(),
			},
			false,
		},
		{
			"crescent_pool_coin_supply_empty",
			args{
				datatype: types.ProtocolDataTypeCrescentPoolCoinSupply,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"crescent_pool_coin_supply",
			args{
				datatype: types.ProtocolDataTypeCrescentPoolCoinSupply,
				data: func() []byte {
					pd := &types.CrescentPoolCoinSupplyProtocolData{
						PoolCoinDenom: "pool1",
						Supply:        []byte{0x6e, 0x75, 0x6c, 0x6c},
					}
					data, _ := json.Marshal(pd)
					return data
				}(),
			},
			&types.CrescentPoolCoinSupplyProtocolData{
				PoolCoinDenom: "pool1",
				Supply:        []byte{0x6e, 0x75, 0x6c, 0x6c},
			},
			false,
		},
		{
			"crescent_reserve_address_balance_empty",
			args{
				datatype: types.ProtocolDataTypeCrescentReserveAddressBalance,
				data:     []byte(`{}`),
			},
			nil,
			true,
		},
		{
			"crescent_reserve_address_balance",
			args{
				datatype: types.ProtocolDataTypeCrescentReserveAddressBalance,
				data: func() []byte {
					pd := &types.CrescentReserveAddressBalanceProtocolData{
						ReserveAddress: testAddress,
						Balance:        []byte{0x6e, 0x75, 0x6c, 0x6c},
					}
					data, _ := json.Marshal(pd)
					return data
				}(),
			},
			&types.CrescentReserveAddressBalanceProtocolData{
				ReserveAddress: testAddress,
				Balance:        []byte{0x6e, 0x75, 0x6c, 0x6c},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := types.UnmarshalProtocolData(tt.args.datatype, tt.args.data)
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
