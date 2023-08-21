package types_test

import (
	"encoding/json"
	"testing"

	liquiditytypes "github.com/ingenuity-build/quicksilver/third-party-chains/crescent-types/liquidity/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"

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
	testUmeeData := types.UmeeProtocolData{Denom: "test", Data: []byte{0x6e, 0x75, 0x6c, 0x6c}}
	testAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")

	invalidOsmosisData := `{
	"poolname": "osmosispools/1",
	"denoms": {	}
}`
	validOsmosisData := `{
	"poolid": 1,
	"poolname": "atom/osmo",
	"pooltype": "balancer",
	"denoms": {
		"uosmo": {"chainid": "osmosis-1", "denom": "uosmo"}
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
		{
			"umee_params_invalid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeParams)],
					Data: []byte("{}"),
				},
			},
			true,
		},
		{
			"umee_params_valid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeParams)],
					Data: []byte(`{"ChainID": "test-01"}`),
				},
			},
			false,
		},
		{
			"umee_reserves_invalid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeReserves)],
					Data: []byte(`{}`),
				},
			},
			true,
		},
		{
			"umee_reserves_valid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeReserves)],
					Data: marshalledUmeeData[types.UmeeReservesProtocolData](testUmeeData),
				},
			},
			false,
		},
		{
			"umee_interest_scalar_invalid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeInterestScalar)],
					Data: []byte(`{}`),
				},
			},
			true,
		},
		{
			"umee_interest_scalar_valid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeInterestScalar)],
					Data: marshalledUmeeData[types.UmeeInterestScalarProtocolData](testUmeeData),
				},
			},
			false,
		},
		{
			"umee_utoken_supply_invalid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeUTokenSupply)],
					Data: []byte(`{}`),
				},
			},
			true,
		},
		{
			"umee_utoken_supply_valid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeUTokenSupply)],
					Data: marshalledUmeeData[types.UmeeUTokenSupplyProtocolData](testUmeeData),
				},
			},
			false,
		},
		{
			"umee_leverage_module_balance_invalid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeLeverageModuleBalance)],
					Data: []byte(`{}`),
				},
			},
			true,
		},
		{
			"umee_leverage_module_balance_valid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeLeverageModuleBalance)],
					Data: marshalledUmeeData[types.UmeeLeverageModuleBalanceProtocolData](testUmeeData),
				},
			},
			false,
		},
		{
			"umee_total_borrows_invalid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeTotalBorrows)],
					Data: []byte(`{}`),
				},
			},
			true,
		},
		{
			"umee_total_borrows_valid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeUmeeTotalBorrows)],
					Data: marshalledUmeeData[types.UmeeTotalBorrowsProtocolData](testUmeeData),
				},
			},
			false,
		},
		{
			"crescent_pool_invalid",
			fields{
				"pool1",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeCrescentPool)],
					Data: []byte(`{}`),
				},
			},
			true,
		},
		{
			"crescent_pool_valid",
			fields{
				"pool1",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeCrescentPool)],
					Data: func() []byte {
						pool := &liquiditytypes.Pool{
							PoolCoinDenom: "pool1",
							Id:            1,
						}
						pooldata, _ := json.Marshal(pool)
						pd := types.CrescentPoolProtocolData{
							PoolID:   1,
							Denom:    "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
							PoolData: pooldata,
						}
						data, _ := json.Marshal(&pd)
						return data
					}(),
				},
			},
			false,
		},
		{
			"crescent_pool_coin_supply_invalid",
			fields{
				"pool1",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeCrescentPoolCoinSupply)],
					Data: []byte(`{}`),
				},
			},
			true,
		},
		{
			"crescent_pool_coin_supply_valid",
			fields{
				"pool1",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeCrescentPoolCoinSupply)],
					Data: func() []byte {
						pd := &types.CrescentPoolCoinSupplyProtocolData{
							PoolCoinDenom: "pool1",
							Supply:        []byte{0x6e, 0x75, 0x6c, 0x6c},
						}
						data, _ := json.Marshal(pd)
						return data
					}(),
				},
			},
			false,
		},
		{
			"crescent_reserve_address_balance_invalid",
			fields{
				"test",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeCrescentReserveAddressBalance)],
					Data: []byte(`{}`),
				},
			},
			true,
		},
		{
			"crescent_reserve_address_balance_valid",
			fields{
				"uosmo",
				&types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeCrescentReserveAddressBalance)],
					Data: func() []byte {
						pd := &types.CrescentReserveAddressBalanceProtocolData{
							ReserveAddress: testAddress,
							Denom:          "uosmo",
							Balance:        []byte{0x6e, 0x75, 0x6c, 0x6c},
						}
						data, _ := json.Marshal(pd)
						return data
					}(),
				},
			},
			false,
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
