package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/mocks"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/services"
	"github.com/stretchr/testify/assert"
)

func TestVersionHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		mockVersion    []byte
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful version request",
			mockVersion:    []byte(`{"version":"1.0.0"}`),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"version":"1.0.0"}`,
		},
		{
			name:           "version service error",
			mockVersion:    nil,
			mockError:      errors.New("version service unavailable"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Error: version service unavailable",
		},
		{
			name:           "write error",
			mockVersion:    []byte(`{"version":"1.0.0"}`),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"version":"1.0.0"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock version service
			mockVersionService := &mocks.MockVersionService{
				GetVersionFunc: func() ([]byte, error) {
					return tt.mockVersion, tt.mockError
				},
			}

			// Create version service
			versionService := services.NewVersionService(mockVersionService)

			// Create handler
			handler := NewVersionHandler(versionService)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/version", nil)
			w := httptest.NewRecorder()

			// Call handler
			handler.Handle(w, req)

			// Assert results
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}

func TestNewVersionHandler(t *testing.T) {
	mockVersionService := &mocks.MockVersionService{}
	versionService := services.NewVersionService(mockVersionService)
	handler := NewVersionHandler(versionService)

	assert.NotNil(t, handler)
	assert.Equal(t, versionService, handler.versionService)
}

func TestGetVersionHandler(t *testing.T) {
	mockVersionService := &mocks.MockVersionService{}
	versionService := services.NewVersionService(mockVersionService)

	handlerFunc := GetVersionHandler(versionService)

	assert.NotNil(t, handlerFunc)

	// Test that the returned function works
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	handlerFunc(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
