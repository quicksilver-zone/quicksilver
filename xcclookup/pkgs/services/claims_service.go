package services

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/claims"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/lookup"
)

// ClaimsService is a concrete implementation of ClaimsServiceInterface
type ClaimsService struct {
	cfg      lookup.Config
	cacheMgr lookup.CacheManagerInterface
}

// NewClaimsService creates a new claims service
func NewClaimsService(cfg lookup.Config, cacheMgr lookup.CacheManagerInterface) *ClaimsService {
	return &ClaimsService{
		cfg:      cfg,
		cacheMgr: cacheMgr,
	}
}

// OsmosisClaim implements ClaimsServiceInterface
func (c *ClaimsService) OsmosisClaim(ctx context.Context, address, submitAddress, chain string, height int64) (lookup.OsmosisResult, error) {
	// For now, we'll need to cast the cache manager to the concrete type
	// This is a limitation of the current design - we need to refactor the claims package
	// to use interfaces instead of concrete types
	if concreteCacheMgr, ok := c.cacheMgr.(*lookup.CacheManager); ok {
		result := claims.OsmosisClaim(ctx, c.cfg, concreteCacheMgr, address, submitAddress, chain, height)
		return lookup.OsmosisResult{
			Err:           result.Err,
			OsmosisPool:   lookup.OsmosisPool(result.OsmosisPool),
			OsmosisClPool: lookup.OsmosisClPool(result.OsmosisClPool),
		}, nil
	}
	return lookup.OsmosisResult{}, lookup.ErrUnsupportedCacheManager
}

// UmeeClaim implements ClaimsServiceInterface
func (c *ClaimsService) UmeeClaim(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	if concreteCacheMgr, ok := c.cacheMgr.(*lookup.CacheManager); ok {
		return claims.UmeeClaim(ctx, c.cfg, concreteCacheMgr, address, submitAddress, chain, height)
	}
	return nil, nil, lookup.ErrUnsupportedCacheManager
}

// LiquidClaim implements ClaimsServiceInterface
func (c *ClaimsService) LiquidClaim(ctx context.Context, address, submitAddress string, connection prewards.ConnectionProtocolData, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	if concreteCacheMgr, ok := c.cacheMgr.(*lookup.CacheManager); ok {
		return claims.LiquidClaim(ctx, c.cfg, concreteCacheMgr, address, submitAddress, connection, height)
	}
	return nil, nil, lookup.ErrUnsupportedCacheManager
}

// MembraneClaim implements ClaimsServiceInterface
func (c *ClaimsService) MembraneClaim(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	if concreteCacheMgr, ok := c.cacheMgr.(*lookup.CacheManager); ok {
		return claims.MembraneClaim(ctx, c.cfg, concreteCacheMgr, address, submitAddress, chain, height)
	}
	return nil, nil, lookup.ErrUnsupportedCacheManager
}
