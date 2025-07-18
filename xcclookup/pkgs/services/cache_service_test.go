package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/mocks"
)

func TestCacheService_GetCacheData(t *testing.T) {
	tests := []struct {
		name               string
		mockConnections    []prewards.ConnectionProtocolData
		mockOsmosisPools   []prewards.OsmosisPoolProtocolData
		mockOsmosisParams  []prewards.OsmosisParamsProtocolData
		mockOsmosisClPools []prewards.OsmosisClPoolProtocolData
		mockTokens         []prewards.LiquidAllowedDenomProtocolData
		mockError          error
		expectedError      error
	}{
		{
			name: "successful cache data retrieval",
			mockConnections: []prewards.ConnectionProtocolData{
				{ChainID: "test-chain-1", LastEpoch: 1},
				{ChainID: "test-chain-2", LastEpoch: 2},
			},
			mockOsmosisPools: []prewards.OsmosisPoolProtocolData{
				{PoolID: 1},
				{PoolID: 2},
			},
			mockOsmosisParams: []prewards.OsmosisParamsProtocolData{
				{ChainID: "osmosis-1"},
			},
			mockOsmosisClPools: []prewards.OsmosisClPoolProtocolData{
				{PoolID: 1},
			},
			mockTokens: []prewards.LiquidAllowedDenomProtocolData{
				{QAssetDenom: "token-1"},
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "cache error on connections",
			mockError:     errors.New("cache error"),
			expectedError: errors.New("cache error"),
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
				GetOsmosisPoolsFunc: func(ctx context.Context) ([]prewards.OsmosisPoolProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockOsmosisPools, nil
				},
				GetOsmosisParamsFunc: func(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockOsmosisParams, nil
				},
				GetOsmosisClPoolsFunc: func(ctx context.Context) ([]prewards.OsmosisClPoolProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockOsmosisClPools, nil
				},
				GetLiquidAllowedDenomsFunc: func(ctx context.Context) ([]prewards.LiquidAllowedDenomProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockTokens, nil
				},
			}

			// Create cache service with mock
			service := NewCacheService(mockCacheManager)

			// Call the method
			result, err := service.GetCacheData(t.Context())

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Verify JSON structure
				var output CacheOutput
				err = json.Unmarshal(result, &output)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockConnections), len(output.Connections))
				assert.Equal(t, len(tt.mockOsmosisPools), len(output.OsmosisPools))
				assert.Equal(t, len(tt.mockOsmosisParams), len(output.OsmosisParams))
				assert.Equal(t, len(tt.mockOsmosisClPools), len(output.OsmosisClPools))
				assert.Equal(t, len(tt.mockTokens), len(output.Tokens))
			}
		})
	}
}

func TestNewCacheService(t *testing.T) {
	mockCacheManager := &mocks.MockCacheManager{}
	service := NewCacheService(mockCacheManager)

	assert.NotNil(t, service)
	assert.Equal(t, mockCacheManager, service.cacheManager)
}
