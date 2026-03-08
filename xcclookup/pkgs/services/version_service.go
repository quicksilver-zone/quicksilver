package services

import (
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/lookup"
)

// VersionService handles version-related operations
type VersionService struct {
	versionService lookup.VersionServiceInterface
}

// NewVersionService creates a new version service
func NewVersionService(versionService lookup.VersionServiceInterface) *VersionService {
	return &VersionService{
		versionService: versionService,
	}
}

// GetVersion retrieves the version information
func (s *VersionService) GetVersion() ([]byte, error) {
	return s.versionService.GetVersion()
}
