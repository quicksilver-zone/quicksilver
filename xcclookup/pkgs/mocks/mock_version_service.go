package mocks

import (
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// MockVersionService is a mock implementation of VersionServiceInterface
type MockVersionService struct {
	GetVersionFunc func() ([]byte, error)
}

// GetVersion calls the mock function
func (m *MockVersionService) GetVersion() ([]byte, error) {
	if m.GetVersionFunc != nil {
		return m.GetVersionFunc()
	}
	return []byte(`{"version":"test"}`), nil
}

// Ensure MockVersionService implements VersionServiceInterface
var _ types.VersionServiceInterface = &MockVersionService{}
