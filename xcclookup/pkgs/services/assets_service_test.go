package services

import (
	"context"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/mocks"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
	"github.com/stretchr/testify/assert"
)

func TestAssetsService_GetAssets(t *testing.T) {
	// Store original function
	origGetMappedAddresses := types.GetMappedAddresses

	// Create a variable to hold the mock function
	var mockGetMappedAddresses func(context.Context, string, []prewards.ConnectionProtocolData, *types.Config) (map[string]string, error)
	mockGetMappedAddresses = func(ctx context.Context, address string, connections []prewards.ConnectionProtocolData, config *types.Config) (map[string]string, error) {
		return map[string]string{}, nil
	}

	// Replace the function temporarily
	types.GetMappedAddresses = mockGetMappedAddresses
	defer func() { types.GetMappedAddresses = origGetMappedAddresses }()

	tests := []struct {
		name              string
		address           string
		mockConnections   []prewards.ConnectionProtocolData
		mockOsmosisParams []prewards.OsmosisParamsProtocolData
		mockUmeeParams    []prewards.UmeeParamsProtocolData
		mockOsmosisResult types.OsmosisResult
		mockUmeeResult    map[string]prewards.MsgSubmitClaim
		mockUmeeAssets    map[string]sdk.Coins
		mockLiquidResult  map[string]prewards.MsgSubmitClaim
		mockLiquidAssets  map[string]sdk.Coins
		mockError         error
		expectedErrors    map[string]error
	}{
		{
			name:    "successful assets retrieval",
			address: "test-address",
			mockConnections: []prewards.ConnectionProtocolData{
				{ChainID: "test-chain-1", LastEpoch: 1},
			},
			mockOsmosisParams: []prewards.OsmosisParamsProtocolData{
				{ChainID: "osmosis-1"},
			},
			mockUmeeParams: []prewards.UmeeParamsProtocolData{
				{ChainID: "umee-1"},
			},
			mockOsmosisResult: types.OsmosisResult{
				OsmosisPool: types.OsmosisPool{
					Msg: map[string]prewards.MsgSubmitClaim{
						"chain1": {
							UserAddress: "test-address",
							ClaimType:   cmtypes.ClaimTypeOsmosisPool,
						},
					},
					Assets: map[string]sdk.Coins{
						"chain1": sdk.NewCoins(sdk.NewCoin("token1", sdk.NewInt(100))),
					},
				},
			},
			mockUmeeResult: map[string]prewards.MsgSubmitClaim{
				"umee-1": {
					UserAddress: "test-address",
					ClaimType:   cmtypes.ClaimTypeUmeeToken,
				},
			},
			mockUmeeAssets: map[string]sdk.Coins{
				"umee-1": sdk.NewCoins(sdk.NewCoin("utoken1", sdk.NewInt(200))),
			},
			mockLiquidResult: map[string]prewards.MsgSubmitClaim{
				"liquid-1": {
					UserAddress: "test-address",
					ClaimType:   cmtypes.ClaimTypeLiquidToken,
				},
			},
			mockLiquidAssets: map[string]sdk.Coins{
				"liquid-1": sdk.NewCoins(sdk.NewCoin("liquid1", sdk.NewInt(300))),
			},
			expectedErrors: make(map[string]error),
		},
		{
			name:      "cache error on connections",
			address:   "test-address",
			mockError: errors.New("cache error"),
			expectedErrors: map[string]error{
				"Connections": errors.New("cache error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock cache manager
			mockCacheManager := &mocks.MockCacheManager{
				GetConnectionsFunc: func(ctx context.Context) ([]prewards.ConnectionProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockConnections, nil
				},
				GetOsmosisParamsFunc: func(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockOsmosisParams, nil
				},
				GetUmeeParamsFunc: func(ctx context.Context) ([]prewards.UmeeParamsProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockUmeeParams, nil
				},
			}

			// Create mock claims service
			mockClaimsService := &mocks.MockClaimsService{
				OsmosisClaimFunc: func(ctx context.Context, address, submitAddress, chain string, height int64) (types.OsmosisResult, error) {
					return tt.mockOsmosisResult, nil
				},
				UmeeClaimFunc: func(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
					return tt.mockUmeeResult, tt.mockUmeeAssets, nil
				},
				LiquidClaimFunc: func(ctx context.Context, address, submitAddress string, connection prewards.ConnectionProtocolData, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
					return tt.mockLiquidResult, tt.mockLiquidAssets, nil
				},
			}

			// Create config
			cfg := types.Config{
				SourceChain: "quicksilver-1",
				Chains: map[string]string{
					"quicksilver-1": "http://quicksilver:26657",
					"osmosis-1":     "http://osmosis:26657",
					"umee-1":        "http://umee:26657",
				},
			}

			// Create heights map
			heights := map[string]int64{
				"osmosis-1": 1000,
				"umee-1":    2000,
			}

			// Create assets service with mocks
			service := NewAssetsService(cfg, mockCacheManager, mockClaimsService, heights)

			// Call the method
			response, errs := service.GetAssets(context.Background(), tt.address)

			// Debug: print out any errors
			if len(errs) > 0 {
				t.Logf("Generated errors: %+v", errs)
				for key, err := range errs {
					t.Logf("Error key: %s, Error: %v", key, err)
				}
			}

			// Assert results
			if len(tt.expectedErrors) > 0 {
				assert.Equal(t, len(tt.expectedErrors), len(errs))
				for key, expectedErr := range tt.expectedErrors {
					assert.Error(t, errs[key])
					assert.Equal(t, expectedErr.Error(), errs[key].Error())
				}
			} else {
				assert.NotNil(t, response)
				assert.Equal(t, 0, len(errs))
			}
		})
	}
}

func TestAssetsService_GetAssets_WithMappedAddresses(t *testing.T) {
	// Store original function
	origGetMappedAddresses := types.GetMappedAddresses

	// Mock mapped addresses for specific chains
	mockMappedAddresses := map[string]string{
		"osmosis-1": "osmo1mappedaddress",
		"umee-1":    "umee1mappedaddress",
	}

	// Create a variable to hold the mock function
	var mockGetMappedAddresses func(context.Context, string, []prewards.ConnectionProtocolData, *types.Config) (map[string]string, error)
	mockGetMappedAddresses = func(ctx context.Context, address string, connections []prewards.ConnectionProtocolData, config *types.Config) (map[string]string, error) {
		return mockMappedAddresses, nil
	}

	// Replace the function temporarily
	types.GetMappedAddresses = mockGetMappedAddresses
	defer func() { types.GetMappedAddresses = origGetMappedAddresses }()

	// Track which addresses were used in claims service calls
	var osmosisAddressUsed, umeeAddressUsed, liquidAddressUsed string

	// Create mock cache manager
	mockCacheManager := &mocks.MockCacheManager{
		GetConnectionsFunc: func(ctx context.Context) ([]prewards.ConnectionProtocolData, error) {
			return []prewards.ConnectionProtocolData{
				{ChainID: "liquid-1", LastEpoch: 1, Prefix: "liquid"},
			}, nil
		},
		GetOsmosisParamsFunc: func(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error) {
			return []prewards.OsmosisParamsProtocolData{
				{ChainID: "osmosis-1"},
			}, nil
		},
		GetUmeeParamsFunc: func(ctx context.Context) ([]prewards.UmeeParamsProtocolData, error) {
			return []prewards.UmeeParamsProtocolData{
				{ChainID: "umee-1"},
			}, nil
		},
	}

	// Create mock claims service that tracks the addresses used
	mockClaimsService := &mocks.MockClaimsService{
		OsmosisClaimFunc: func(ctx context.Context, address, submitAddress, chain string, height int64) (types.OsmosisResult, error) {
			osmosisAddressUsed = address
			return types.OsmosisResult{
				OsmosisPool: types.OsmosisPool{
					Msg: map[string]prewards.MsgSubmitClaim{
						"chain1": {
							UserAddress: submitAddress,
							ClaimType:   cmtypes.ClaimTypeOsmosisPool,
						},
					},
					Assets: map[string]sdk.Coins{
						"chain1": sdk.NewCoins(sdk.NewCoin("token1", sdk.NewInt(100))),
					},
				},
			}, nil
		},
		UmeeClaimFunc: func(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
			umeeAddressUsed = address
			return map[string]prewards.MsgSubmitClaim{
					"umee-1": {
						UserAddress: submitAddress,
						ClaimType:   cmtypes.ClaimTypeUmeeToken,
					},
				}, map[string]sdk.Coins{
					"umee-1": sdk.NewCoins(sdk.NewCoin("utoken1", sdk.NewInt(200))),
				}, nil
		},
		LiquidClaimFunc: func(ctx context.Context, address, submitAddress string, connection prewards.ConnectionProtocolData, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
			liquidAddressUsed = address
			return map[string]prewards.MsgSubmitClaim{
					"liquid-1": {
						UserAddress: submitAddress,
						ClaimType:   cmtypes.ClaimTypeLiquidToken,
					},
				}, map[string]sdk.Coins{
					"liquid-1": sdk.NewCoins(sdk.NewCoin("liquid1", sdk.NewInt(300))),
				}, nil
		},
	}

	// Create config
	cfg := types.Config{
		SourceChain: "quicksilver-1",
		Chains: map[string]string{
			"quicksilver-1": "http://quicksilver:26657",
			"osmosis-1":     "http://osmosis:26657",
			"umee-1":        "http://umee:26657",
		},
	}

	// Create heights map
	heights := map[string]int64{
		"osmosis-1": 1000,
		"umee-1":    2000,
		"liquid-1":  3000,
	}

	// Create assets service with mocks
	service := NewAssetsService(cfg, mockCacheManager, mockClaimsService, heights)

	// Call the method
	originalAddress := "quick1originaladdress"
	response, errs := service.GetAssets(context.Background(), originalAddress)

	// Assert no errors
	assert.Equal(t, 0, len(errs))
	assert.NotNil(t, response)

	// Verify that mapped addresses were used for chains that have mappings
	assert.Equal(t, mockMappedAddresses["osmosis-1"], osmosisAddressUsed, "Osmosis claim should use mapped address")
	assert.Equal(t, mockMappedAddresses["umee-1"], umeeAddressUsed, "Umee claim should use mapped address")

	// Verify that original address was used for liquid claim (no mapping for liquid-1)
	assert.Equal(t, originalAddress, liquidAddressUsed, "Liquid claim should use original address when no mapping exists")

	// Verify that the submit address is always the original address
	for _, msg := range response.Messages {
		assert.Equal(t, originalAddress, msg.UserAddress, "Submit address should always be the original address")
	}
}

func TestNewAssetsService(t *testing.T) {
	mockCacheManager := &mocks.MockCacheManager{}
	mockClaimsService := &mocks.MockClaimsService{}
	cfg := types.Config{}
	heights := map[string]int64{}

	service := NewAssetsService(cfg, mockCacheManager, mockClaimsService, heights)

	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.cfg)
	assert.Equal(t, mockCacheManager, service.cacheManager)
	assert.Equal(t, mockClaimsService, service.claimsService)
	assert.Equal(t, heights, service.heights)
}
