package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/services"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// CacheHandler handles cache-related HTTP requests
type CacheHandler struct {
	cacheService *services.CacheService
}

// NewCacheHandler creates a new cache handler
func NewCacheHandler(cacheService *services.CacheService) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
	}
}

// Handle handles cache requests
func (h *CacheHandler) Handle(w http.ResponseWriter, r *http.Request) {
	jsonOut, err := h.cacheService.GetCacheData(r.Context())
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprint(w, string(jsonOut))
}

// GetCacheHandler returns a function that creates a cache handler
func GetCacheHandler(
	ctx context.Context,
	_ types.Config,
	cacheMgr types.CacheManagerInterface,
) func(http.ResponseWriter, *http.Request) {
	cacheService := services.NewCacheService(cacheMgr)
	handler := NewCacheHandler(cacheService)
	return handler.Handle
}
