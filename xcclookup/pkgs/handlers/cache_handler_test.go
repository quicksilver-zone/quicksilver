package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/mocks"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/services"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
	"github.com/stretchr/testify/assert"
)

func TestCacheHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		mockCacheData  []byte
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful cache request",
			mockCacheData: func() []byte {
				data := services.CacheOutput{
					Connections: []prewards.ConnectionProtocolData{
						{ChainID: "test-chain-1", LastEpoch: 1},
					},
					OsmosisPools: []prewards.OsmosisPoolProtocolData{
						{PoolID: "pool-1"},
					},
				}
				jsonData, _ := json.Marshal(data)
				return jsonData
			}(),
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "cache service error",
			mockCacheData:  nil,
			mockError:      errors.New("cache service unavailable"),
			expectedStatus: http.StatusOK, // The handler doesn't set status code on error
			expectedBody:   "Error: cache service unavailable",
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
					return make([]prewards.ConnectionProtocolData, 0), nil
				},
				GetOsmosisPoolsFunc: func(ctx context.Context) ([]prewards.OsmosisPoolProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return make([]prewards.OsmosisPoolProtocolData, 0), nil
				},
				GetOsmosisParamsFunc: func(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return make([]prewards.OsmosisParamsProtocolData, 0), nil
				},
				GetOsmosisClPoolsFunc: func(ctx context.Context) ([]prewards.OsmosisClPoolProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return make([]prewards.OsmosisClPoolProtocolData, 0), nil
				},
				GetLiquidAllowedDenomsFunc: func(ctx context.Context) ([]prewards.LiquidAllowedDenomProtocolData, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return make([]prewards.LiquidAllowedDenomProtocolData, 0), nil
				},
			}

			// Create cache service
			cacheService := services.NewCacheService(mockCacheManager)

			// Create handler
			handler := NewCacheHandler(cacheService)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/cache", nil)
			w := httptest.NewRecorder()

			// Call handler
			handler.Handle(w, req)

			// Assert results
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockError != nil {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			} else {
				assert.NotEmpty(t, w.Body.String())
			}
		})
	}
}

func TestNewCacheHandler(t *testing.T) {
	mockCacheManager := &mocks.MockCacheManager{}
	cacheService := services.NewCacheService(mockCacheManager)
	handler := NewCacheHandler(cacheService)

	assert.NotNil(t, handler)
	assert.Equal(t, cacheService, handler.cacheService)
}

func TestGetCacheHandler(t *testing.T) {
	mockCacheManager := &mocks.MockCacheManager{}
	cfg := types.Config{}

	handlerFunc := GetCacheHandler(context.Background(), cfg, mockCacheManager)

	assert.NotNil(t, handlerFunc)

	// Test that the returned function works
	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	w := httptest.NewRecorder()

	handlerFunc(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
