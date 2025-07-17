package claims

import (
	"testing"

	"github.com/stretchr/testify/assert"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

func TestGetTokenMap(t *testing.T) {
	tests := []struct {
		name           string
		liquidAllowed  []prewards.LiquidAllowedDenomProtocolData
		zones          []icstypes.Zone
		chain          string
		keyPrefix      string
		ignores        types.Ignores
		expectedResult map[string]TokenTuple
	}{
		{
			name: "successful token mapping",
			liquidAllowed: []prewards.LiquidAllowedDenomProtocolData{
				{
					QAssetDenom:           "qasset1",
					IbcDenom:              "ibc/denom1",
					ChainID:               "chain1",
					RegisteredZoneChainID: "zone1",
				},
				{
					QAssetDenom:           "qasset2",
					IbcDenom:              "ibc/denom2",
					ChainID:               "chain1",
					RegisteredZoneChainID: "zone2",
				},
			},
			zones: []icstypes.Zone{
				{ChainId: "zone1"},
				{ChainId: "zone2"},
			},
			chain:     "chain1",
			keyPrefix: "prefix",
			ignores:   types.Ignores{},
			expectedResult: map[string]TokenTuple{
				"prefixibc/denom1": {denom: "qasset1", chain: "zone1"},
				"prefixibc/denom2": {denom: "qasset2", chain: "zone2"},
			},
		},
		{
			name: "ignored token",
			liquidAllowed: []prewards.LiquidAllowedDenomProtocolData{
				{
					QAssetDenom:           "qasset1",
					IbcDenom:              "ibc/denom1",
					ChainID:               "chain1",
					RegisteredZoneChainID: "zone1",
				},
			},
			zones: []icstypes.Zone{
				{ChainId: "zone1"},
			},
			chain:     "chain1",
			keyPrefix: "",
			ignores: types.Ignores{
				{Type: "liquid", Key: "qasset1"},
			},
			expectedResult: map[string]TokenTuple{},
		},
		{
			name: "zone not onboarded",
			liquidAllowed: []prewards.LiquidAllowedDenomProtocolData{
				{
					QAssetDenom:           "qasset1",
					IbcDenom:              "ibc/denom1",
					ChainID:               "chain1",
					RegisteredZoneChainID: "zone1",
				},
			},
			zones: []icstypes.Zone{
				{ChainId: "zone2"}, // Different zone
			},
			chain:          "chain1",
			keyPrefix:      "",
			ignores:        types.Ignores{},
			expectedResult: map[string]TokenTuple{},
		},
		{
			name: "wrong chain",
			liquidAllowed: []prewards.LiquidAllowedDenomProtocolData{
				{
					QAssetDenom:           "qasset1",
					IbcDenom:              "ibc/denom1",
					ChainID:               "chain1",
					RegisteredZoneChainID: "zone1",
				},
			},
			zones: []icstypes.Zone{
				{ChainId: "zone1"},
			},
			chain:          "chain2", // Different chain
			keyPrefix:      "",
			ignores:        types.Ignores{},
			expectedResult: map[string]TokenTuple{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTokenMap(tt.liquidAllowed, tt.zones, tt.chain, tt.keyPrefix, tt.ignores)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestZoneOnboarded(t *testing.T) {
	tests := []struct {
		name   string
		zones  []icstypes.Zone
		token  prewards.LiquidAllowedDenomProtocolData
		result bool
	}{
		{
			name: "zone is onboarded",
			zones: []icstypes.Zone{
				{ChainId: "zone1"},
				{ChainId: "zone2"},
			},
			token: prewards.LiquidAllowedDenomProtocolData{
				RegisteredZoneChainID: "zone1",
			},
			result: true,
		},
		{
			name: "zone is not onboarded",
			zones: []icstypes.Zone{
				{ChainId: "zone1"},
				{ChainId: "zone2"},
			},
			token: prewards.LiquidAllowedDenomProtocolData{
				RegisteredZoneChainID: "zone3",
			},
			result: false,
		},
		{
			name:  "empty zones list",
			zones: []icstypes.Zone{},
			token: prewards.LiquidAllowedDenomProtocolData{
				RegisteredZoneChainID: "zone1",
			},
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZoneOnboarded(tt.zones, tt.token)
			assert.Equal(t, tt.result, result)
		})
	}
}
