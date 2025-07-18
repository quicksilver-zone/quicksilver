package types

import (
	"encoding/json"
	"runtime"
)

// VersionService is a concrete implementation of VersionServiceInterface
type VersionService struct{}

// GetVersion implements VersionServiceInterface
func (v *VersionService) GetVersion() ([]byte, error) {
	version := map[string]string{
		"version": "1.0.0",
		"go":      runtime.Version(),
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	}

	return json.Marshal(version)
}

// GetVersion is a legacy function that uses the concrete implementation
func GetVersion() ([]byte, error) {
	vs := &VersionService{}
	return vs.GetVersion()
}
