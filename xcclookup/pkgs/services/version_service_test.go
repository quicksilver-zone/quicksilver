package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/mocks"
)

func TestVersionService_GetVersion(t *testing.T) {
	tests := []struct {
		name           string
		mockVersion    []byte
		mockError      error
		expectedResult []byte
		expectedError  error
	}{
		{
			name:           "successful version retrieval",
			mockVersion:    []byte(`{"version":"1.0.0"}`),
			mockError:      nil,
			expectedResult: []byte(`{"version":"1.0.0"}`),
			expectedError:  nil,
		},
		{
			name:           "version service error",
			mockVersion:    nil,
			mockError:      errors.New("version service unavailable"),
			expectedResult: nil,
			expectedError:  errors.New("version service unavailable"),
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

			// Create version service with mock
			service := NewVersionService(mockVersionService)

			// Call the method
			result, err := service.GetVersion()

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestNewVersionService(t *testing.T) {
	mockVersionService := &mocks.MockVersionService{}
	service := NewVersionService(mockVersionService)

	assert.NotNil(t, service)
	assert.Equal(t, mockVersionService, service.versionService)
}
