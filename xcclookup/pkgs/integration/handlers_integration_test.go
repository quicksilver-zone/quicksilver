package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/handlers"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/mocks"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/services"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
	"github.com/stretchr/testify/assert"
)

func TestHandlersIntegration(t *testing.T) {
	t.Run("version handler integration", func(t *testing.T) {
		// Setup
		mockVersionService := &mocks.MockVersionService{
			GetVersionFunc: func() ([]byte, error) {
				return []byte(`{"version":"1.0.0","build":"test"}`), nil
			},
		}
		versionService := services.NewVersionService(mockVersionService)
		handler := handlers.NewVersionHandler(versionService)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/version", nil)
		w := httptest.NewRecorder()

		// Execute
		handler.Handle(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "1.0.0", response["version"])
		assert.Equal(t, "test", response["build"])
	})

	t.Run("cache handler integration", func(t *testing.T) {
		// Setup
		mockCacheManager := &mocks.MockCacheManager{
			GetConnectionsFunc: func(ctx context.Context) ([]prewards.ConnectionProtocolData, error) {
				return make([]prewards.ConnectionProtocolData, 0), nil
			},
			GetOsmosisPoolsFunc: func(ctx context.Context) ([]prewards.OsmosisPoolProtocolData, error) {
				return make([]prewards.OsmosisPoolProtocolData, 0), nil
			},
			GetOsmosisParamsFunc: func(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error) {
				return make([]prewards.OsmosisParamsProtocolData, 0), nil
			},
			GetOsmosisClPoolsFunc: func(ctx context.Context) ([]prewards.OsmosisClPoolProtocolData, error) {
				return make([]prewards.OsmosisClPoolProtocolData, 0), nil
			},
			GetLiquidAllowedDenomsFunc: func(ctx context.Context) ([]prewards.LiquidAllowedDenomProtocolData, error) {
				return make([]prewards.LiquidAllowedDenomProtocolData, 0), nil
			},
		}
		cacheService := services.NewCacheService(mockCacheManager)
		handler := handlers.NewCacheHandler(cacheService)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/cache", nil)
		w := httptest.NewRecorder()

		// Execute
		handler.Handle(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response services.CacheOutput
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	t.Run("assets handler integration", func(t *testing.T) {
		// Setup
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

		mockClaimsService := &mocks.MockClaimsService{
			OsmosisClaimFunc: func(ctx context.Context, address, submitAddress, chain string, height int64) (types.OsmosisResult, error) {
				return types.OsmosisResult{
					OsmosisPool: types.OsmosisPool{
						Msg: map[string]prewards.MsgSubmitClaim{
							"chain1": {UserAddress: address},
						},
						Assets: map[string]sdk.Coins{
							"chain1": sdk.NewCoins(sdk.NewCoin("token1", sdk.NewInt(100))),
						},
					},
				}, nil
			},
			UmeeClaimFunc: func(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
				return map[string]prewards.MsgSubmitClaim{
						"umee-1": {UserAddress: address},
					}, map[string]sdk.Coins{
						"umee-1": sdk.NewCoins(sdk.NewCoin("utoken1", sdk.NewInt(200))),
					}, nil
			},
			LiquidClaimFunc: func(ctx context.Context, address, submitAddress string, connection prewards.ConnectionProtocolData, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
				return map[string]prewards.MsgSubmitClaim{
						"liquid-1": {UserAddress: address},
					}, map[string]sdk.Coins{
						"liquid-1": sdk.NewCoins(sdk.NewCoin("liquid1", sdk.NewInt(300))),
					}, nil
			},
		}

		cfg := types.Config{
			Chains: map[string]string{
				"osmosis-1": "http://osmosis:26657",
				"umee-1":    "http://umee:26657",
			},
		}
		heights := map[string]int64{
			"osmosis-1": 1000,
			"umee-1":    2000,
		}

		assetsService := services.NewAssetsService(cfg, mockCacheManager, mockClaimsService, heights)

		outputFunc := func(w http.ResponseWriter, response *types.Response, errors map[string]error) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		handler := handlers.NewAssetsHandler(assetsService, outputFunc)

		// Create request with mux vars
		req := httptest.NewRequest(http.MethodGet, "/assets/test-address", nil)
		vars := map[string]string{
			"address": "test-address",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()

		// Execute
		handler.Handle(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})
}
