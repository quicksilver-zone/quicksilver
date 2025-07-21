package mocks

import (
	"context"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// MockCacheManager is a mock implementation of CacheManagerInterface
type MockCacheManager struct {
	GetConnectionsFunc         func(ctx context.Context) ([]prewards.ConnectionProtocolData, error)
	GetOsmosisParamsFunc       func(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error)
	GetOsmosisPoolsFunc        func(ctx context.Context) ([]prewards.OsmosisPoolProtocolData, error)
	GetOsmosisClPoolsFunc      func(ctx context.Context) ([]prewards.OsmosisClPoolProtocolData, error)
	GetLiquidAllowedDenomsFunc func(ctx context.Context) ([]prewards.LiquidAllowedDenomProtocolData, error)
	GetUmeeParamsFunc          func(ctx context.Context) ([]prewards.UmeeParamsProtocolData, error)
	GetMembraneParamsFunc      func(ctx context.Context) ([]prewards.MembraneProtocolData, error)
	GetZonesFunc               func(ctx context.Context) ([]icstypes.Zone, error)
	AddMocksFunc               func(ctx context.Context, mocks interface{}) error
}

// GetConnections calls the mock function
func (m *MockCacheManager) GetConnections(ctx context.Context) ([]prewards.ConnectionProtocolData, error) {
	if m.GetConnectionsFunc != nil {
		return m.GetConnectionsFunc(ctx)
	}
	return make([]prewards.ConnectionProtocolData, 0), nil
}

// GetOsmosisParams calls the mock function
func (m *MockCacheManager) GetOsmosisParams(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error) {
	if m.GetOsmosisParamsFunc != nil {
		return m.GetOsmosisParamsFunc(ctx)
	}
	return make([]prewards.OsmosisParamsProtocolData, 0), nil
}

// GetOsmosisPools calls the mock function
func (m *MockCacheManager) GetOsmosisPools(ctx context.Context) ([]prewards.OsmosisPoolProtocolData, error) {
	if m.GetOsmosisPoolsFunc != nil {
		return m.GetOsmosisPoolsFunc(ctx)
	}
	return make([]prewards.OsmosisPoolProtocolData, 0), nil
}

// GetOsmosisClPools calls the mock function
func (m *MockCacheManager) GetOsmosisClPools(ctx context.Context) ([]prewards.OsmosisClPoolProtocolData, error) {
	if m.GetOsmosisClPoolsFunc != nil {
		return m.GetOsmosisClPoolsFunc(ctx)
	}
	return make([]prewards.OsmosisClPoolProtocolData, 0), nil
}

// GetLiquidAllowedDenoms calls the mock function
func (m *MockCacheManager) GetLiquidAllowedDenoms(ctx context.Context) ([]prewards.LiquidAllowedDenomProtocolData, error) {
	if m.GetLiquidAllowedDenomsFunc != nil {
		return m.GetLiquidAllowedDenomsFunc(ctx)
	}
	return make([]prewards.LiquidAllowedDenomProtocolData, 0), nil
}

// GetUmeeParams calls the mock function
func (m *MockCacheManager) GetUmeeParams(ctx context.Context) ([]prewards.UmeeParamsProtocolData, error) {
	if m.GetUmeeParamsFunc != nil {
		return m.GetUmeeParamsFunc(ctx)
	}
	return make([]prewards.UmeeParamsProtocolData, 0), nil
}

// GetMembraneParams calls the mock function
func (m *MockCacheManager) GetMembraneParams(ctx context.Context) ([]prewards.MembraneProtocolData, error) {
	if m.GetMembraneParamsFunc != nil {
		return m.GetMembraneParamsFunc(ctx)
	}
	return make([]prewards.MembraneProtocolData, 0), nil
}

// GetZones calls the mock function
func (m *MockCacheManager) GetZones(ctx context.Context) ([]icstypes.Zone, error) {
	if m.GetZonesFunc != nil {
		return m.GetZonesFunc(ctx)
	}
	return make([]icstypes.Zone, 0), nil
}

// AddMocks calls the mock function
func (m *MockCacheManager) AddMocks(ctx context.Context, mocks interface{}) error {
	if m.AddMocksFunc != nil {
		return m.AddMocksFunc(ctx, mocks)
	}
	return nil
}

// Ensure MockCacheManager implements CacheManagerInterface
var _ types.CacheManagerInterface = &MockCacheManager{}
