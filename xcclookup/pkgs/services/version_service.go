package services

import (
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// VersionService handles version-related operations
type VersionService struct {
	versionService types.VersionServiceInterface
}

// NewVersionService creates a new version service
func NewVersionService(versionService types.VersionServiceInterface) *VersionService {
	return &VersionService{
		versionService: versionService,
	}
}

// GetVersion retrieves the version information
func (s *VersionService) GetVersion() ([]byte, error) {
	return s.versionService.GetVersion()
}
