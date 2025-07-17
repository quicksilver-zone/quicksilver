package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/mocks"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/services"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
	"github.com/stretchr/testify/assert"
)

func TestAssetsHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		mockResponse   *types.Response
		mockErrors     map[string]error
		expectedStatus int
	}{
		{
			name:    "successful assets request",
			address: "test-address",
			mockResponse: &types.Response{
				Messages: []prewards.MsgSubmitClaim{
					{UserAddress: "test-address"},
				},
				Assets: map[string][]types.Asset{
					"chain1": {{Denom: "token1", Amount: "100"}},
				},
			},
			mockErrors:     make(map[string]error),
			expectedStatus: http.StatusOK,
		},
		{
			name:    "assets request with errors",
			address: "test-address",
			mockResponse: &types.Response{
				Messages: []prewards.MsgSubmitClaim{},
				Assets:   map[string][]types.Asset{},
			},
			mockErrors: map[string]error{
				"OsmosisClaim": assert.AnError,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock cache manager
			mockCacheManager := &mocks.MockCacheManager{
				GetConnectionsFunc: func(ctx context.Context) ([]prewards.ConnectionProtocolData, error) {
					return make([]prewards.ConnectionProtocolData, 0), nil
				},
				GetOsmosisParamsFunc: func(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error) {
					return make([]prewards.OsmosisParamsProtocolData, 0), nil
				},
				GetUmeeParamsFunc: func(ctx context.Context) ([]prewards.UmeeParamsProtocolData, error) {
					return make([]prewards.UmeeParamsProtocolData, 0), nil
				},
			}

			// Create mock claims service
			mockClaimsService := &mocks.MockClaimsService{}

			// Create assets service
			cfg := types.Config{}
			heights := map[string]int64{}
			assetsService := services.NewAssetsService(cfg, mockCacheManager, mockClaimsService, heights)

			// Create output function
			outputFunc := func(w http.ResponseWriter, response *types.Response, errors map[string]error) {
				// Mock output function
			}

			// Create handler
			handler := NewAssetsHandler(assetsService, outputFunc)

			// Create request with mux vars
			req := httptest.NewRequest(http.MethodGet, "/assets/"+tt.address, nil)
			vars := map[string]string{
				"address": tt.address,
			}
			req = mux.SetURLVars(req, vars)
			w := httptest.NewRecorder()

			// Call handler
			handler.Handle(w, req)

			// Assert results
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}

func TestNewAssetsHandler(t *testing.T) {
	mockCacheManager := &mocks.MockCacheManager{}
	mockClaimsService := &mocks.MockClaimsService{}
	cfg := types.Config{}
	heights := map[string]int64{}
	assetsService := services.NewAssetsService(cfg, mockCacheManager, mockClaimsService, heights)

	outputFunc := func(w http.ResponseWriter, response *types.Response, errors map[string]error) {}
	handler := NewAssetsHandler(assetsService, outputFunc)

	assert.NotNil(t, handler)
	assert.Equal(t, assetsService, handler.assetsService)
	assert.Equal(t, outputFunc, handler.outputFunc)
}

func TestGetAssetsHandler(t *testing.T) {
	mockCacheManager := &mocks.MockCacheManager{}
	mockClaimsService := &mocks.MockClaimsService{}
	cfg := types.Config{}
	heights := map[string]int64{}
	outputFunc := func(w http.ResponseWriter, response *types.Response, errors map[string]error) {}

	handlerFunc := GetAssetsHandler(context.Background(), cfg, mockCacheManager, mockClaimsService, heights, outputFunc)

	assert.NotNil(t, handlerFunc)

	// Test that the returned function works
	req := httptest.NewRequest(http.MethodGet, "/assets/test-address", nil)
	vars := map[string]string{
		"address": "test-address",
	}
	req = mux.SetURLVars(req, vars)
	w := httptest.NewRecorder()

	handlerFunc(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
