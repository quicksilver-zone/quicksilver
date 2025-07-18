package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/services"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// AssetsHandler handles assets-related HTTP requests
type AssetsHandler struct {
	assetsService *services.AssetsService
	outputFunc    types.OutputFunction
}

// NewAssetsHandler creates a new assets handler
func NewAssetsHandler(assetsService *services.AssetsService, outputFunc types.OutputFunction) *AssetsHandler {
	return &AssetsHandler{
		assetsService: assetsService,
		outputFunc:    outputFunc,
	}
}

// Handle handles assets requests
func (h *AssetsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response, errs := h.assetsService.GetAssets(r.Context(), address)

	// ensure a neatly formatted JSON response
	h.outputFunc(w, response, errs)
}

// GetAssetsHandler returns a function that creates an assets handler
func GetAssetsHandler(
	ctx context.Context,
	cfg types.Config,
	cacheMgr types.CacheManagerInterface,
	claimsService types.ClaimsServiceInterface,
	heights map[string]int64,
	outputFunc types.OutputFunction,
) func(http.ResponseWriter, *http.Request) {
	assetsService := services.NewAssetsService(cfg, cacheMgr, claimsService, heights)
	handler := NewAssetsHandler(assetsService, outputFunc)
	return handler.Handle
}
