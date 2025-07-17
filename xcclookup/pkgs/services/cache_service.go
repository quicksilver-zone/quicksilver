package services

import (
	"context"
	"encoding/json"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// CacheOutput represents the output structure for cache data
type CacheOutput struct {
	Connections    []prewards.ConnectionProtocolData
	OsmosisPools   []prewards.OsmosisPoolProtocolData
	OsmosisClPools []prewards.OsmosisClPoolProtocolData
	OsmosisParams  []prewards.OsmosisParamsProtocolData
	Tokens         []prewards.LiquidAllowedDenomProtocolData
}

// CacheService handles cache-related operations
type CacheService struct {
	cacheManager types.CacheManagerInterface
}

// NewCacheService creates a new cache service
func NewCacheService(cacheManager types.CacheManagerInterface) *CacheService {
	return &CacheService{
		cacheManager: cacheManager,
	}
}

// GetCacheData retrieves all cache data and returns it as JSON
func (s *CacheService) GetCacheData(ctx context.Context) ([]byte, error) {
	connections, err := s.cacheManager.GetConnections(ctx)
	if err != nil {
		return nil, err
	}

	osmosisPools, err := s.cacheManager.GetOsmosisPools(ctx)
	if err != nil {
		return nil, err
	}

	osmosisParams, err := s.cacheManager.GetOsmosisParams(ctx)
	if err != nil {
		return nil, err
	}

	osmosisClPools, err := s.cacheManager.GetOsmosisClPools(ctx)
	if err != nil {
		return nil, err
	}

	liquidAllowedDenoms, err := s.cacheManager.GetLiquidAllowedDenoms(ctx)
	if err != nil {
		return nil, err
	}

	out := CacheOutput{
		Connections:    connections,
		OsmosisPools:   osmosisPools,
		OsmosisParams:  osmosisParams,
		OsmosisClPools: osmosisClPools,
		Tokens:         liquidAllowedDenoms,
	}

	return json.Marshal(out)
}
