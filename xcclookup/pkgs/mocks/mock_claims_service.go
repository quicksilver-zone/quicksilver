package mocks

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// MockClaimsService is a mock implementation of ClaimsServiceInterface
type MockClaimsService struct {
	OsmosisClaimFunc func(ctx context.Context, address, submitAddress, chain string, height int64) (types.OsmosisResult, error)
	UmeeClaimFunc    func(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error)
	LiquidClaimFunc  func(ctx context.Context, address, submitAddress string, connection prewards.ConnectionProtocolData, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error)
}

// OsmosisClaim calls the mock function
func (m *MockClaimsService) OsmosisClaim(ctx context.Context, address, submitAddress, chain string, height int64) (types.OsmosisResult, error) {
	if m.OsmosisClaimFunc != nil {
		return m.OsmosisClaimFunc(ctx, address, submitAddress, chain, height)
	}
	return types.OsmosisResult{}, nil
}

// UmeeClaim calls the mock function
func (m *MockClaimsService) UmeeClaim(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	if m.UmeeClaimFunc != nil {
		return m.UmeeClaimFunc(ctx, address, submitAddress, chain, height)
	}
	return make(map[string]prewards.MsgSubmitClaim), make(map[string]sdk.Coins), nil
}

// LiquidClaim calls the mock function
func (m *MockClaimsService) LiquidClaim(ctx context.Context, address, submitAddress string, connection prewards.ConnectionProtocolData, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	if m.LiquidClaimFunc != nil {
		return m.LiquidClaimFunc(ctx, address, submitAddress, connection, height)
	}
	return make(map[string]prewards.MsgSubmitClaim), make(map[string]sdk.Coins), nil
}

// Ensure MockClaimsService implements ClaimsServiceInterface
var _ types.ClaimsServiceInterface = &MockClaimsService{}
