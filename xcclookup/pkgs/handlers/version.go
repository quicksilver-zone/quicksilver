package handlers

import (
	"fmt"
	"net/http"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/services"
)

// VersionHandler handles version-related HTTP requests
type VersionHandler struct {
	versionService *services.VersionService
}

// NewVersionHandler creates a new version handler
func NewVersionHandler(versionService *services.VersionService) *VersionHandler {
	return &VersionHandler{
		versionService: versionService,
	}
}

// Handle handles version requests
func (h *VersionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	version, err := h.versionService.GetVersion()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(version)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
}

// GetVersionHandler returns a function that creates a version handler
func GetVersionHandler(versionService *services.VersionService) func(http.ResponseWriter, *http.Request) {
	handler := NewVersionHandler(versionService)
	return handler.Handle
}
