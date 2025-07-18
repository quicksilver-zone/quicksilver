package claims

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	rpcclient "github.com/cometbft/cometbft/rpc/client"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/logger"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// MembranePosition represents a position in the Membrane protocol
type MembranePosition struct {
	PositionID       math.Int                  `json:"position_id"`
	CollateralAssets []MembraneCollateralAsset `json:"collateral_assets"`
	CreditAmount     math.Int                  `json:"credit_amount"`
}

// MembraneCollateralAsset represents a collateral asset in a Membrane position
type MembraneCollateralAsset struct {
	Asset        MembraneCollateralAssetAsset `json:"asset"`
	MaxBorrowLTV sdk.Dec                      `json:"max_borrow_LTV"`
	MaxLTV       sdk.Dec                      `json:"max_LTV"`
	RateIndex    sdk.Dec                      `json:"rate_index"`
}

// MembraneCollateralAssetAsset represents the asset information in a collateral asset
type MembraneCollateralAssetAsset struct {
	Info   MembraneCollateralAssetInfo `json:"info"`
	Amount math.Int                    `json:"amount"`
}

// MembraneCollateralAssetInfo represents the info field in a collateral asset
type MembraneCollateralAssetInfo struct {
	NativeToken MembraneNativeToken `json:"native_token"`
}

// MembraneNativeToken represents a native token
type MembraneNativeToken struct {
	Denom string `json:"denom"`
}

// MembraneClaim performs the reverse operation of the participation rewards module
// It queries the Membrane contract for user positions and validates them against
// the allowed liquid tokens for the Osmosis chain
func MembraneClaim(
	ctx context.Context,
	cfg types.Config,
	cacheMgr *types.CacheManager,
	address string,
	submitAddress string,
	chain string,
	height int64,
) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	log := logger.FromContext(ctx)

	// Get cached data
	membraneParamsCache, err := types.GetCache[prewards.MembraneProtocolData](ctx, cacheMgr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get membrane params cache: %w", err)
	}

	osmosisParamsCache, err := types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get osmosis params cache: %w", err)
	}

	liquidAllowedDenomsCache, err := types.GetCache[prewards.LiquidAllowedDenomProtocolData](ctx, cacheMgr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get liquid allowed denoms cache: %w", err)
	}

	zoneCache, err := types.GetCache[icstypes.Zone](ctx, cacheMgr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get zone cache: %w", err)
	}

	// Find the membrane params for the chain
	var membraneParams prewards.MembraneProtocolData
	found := false
	for _, params := range membraneParamsCache {
		// For now, we assume there's only one membrane params entry
		// In the future, we might need to match by chain ID
		membraneParams = params
		found = true
		break
	}

	if !found {
		log.Warn("No membrane params found in cache")
		return nil, nil, nil
	}

	// Find the osmosis params for the chain
	found = false
	for _, params := range osmosisParamsCache {
		if params.ChainID == chain {
			found = true
			break
		}
	}

	if !found {
		log.Warn("No osmosis params found for chain", "chain", chain)
		return nil, nil, nil
	}

	// Get the host endpoint for the chain
	host, ok := cfg.Chains[chain]
	if !ok {
		log.Warn("No endpoint found for chain", "chain", chain)
		return nil, nil, nil
	}

	// Create RPC client
	client, err := types.NewRPCClient(host, time.Duration(cfg.Timeout)*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create RPC client: %w", err)
	}

	// Setup codec
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	cmtypes.RegisterInterfaces(interfaceRegistry)

	// Convert address to Osmosis format
	addrBytes, err := addressutils.AccAddressFromBech32(address, "")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid address: %w", err)
	}

	addr, err := addressutils.EncodeAddressToBech32("osmo", addrBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode address: %w", err)
	}
	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}

	// Generate the key for the membrane contract query using the proper key format
	// The key should be a CW namespaced key: 0x03 + contract_address + 0x00 + "positions" + user_address
	contractAddrBytes, err := addressutils.AccAddressFromBech32(membraneParams.ContractAddress, "osmo")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid membrane contract address: %w", err)
	}

	// Build the CW namespaced key
	key := make([]byte, 0)
	key = append(key, 0x03)                   // prefix
	key = append(key, contractAddrBytes...)   // contract address
	key = append(key, 0x00)                   // null terminator
	key = append(key, byte(len("positions"))) // length of "positions"
	key = append(key, []byte("positions")...) // "positions"
	key = append(key, []byte(addr)...)        // user address

	// Query the membrane contract
	abciquery, err := client.ABCIQueryWithOptions(
		ctx,
		"/store/wasm/key",
		key,
		rpcclient.ABCIQueryOptions{Height: height, Prove: true},
	)
	if err != nil {
		log.Debug("Failed to query membrane contract", "address", addr, "error", err)
		return nil, nil, fmt.Errorf("failed to query membrane contract: %w", err)
	}

	// Parse the response
	var positions []MembranePosition
	if err := json.Unmarshal(abciquery.Response.Value, &positions); err != nil {
		log.Debug("Failed to unmarshal membrane positions", "address", addr, "error", err)
		return nil, nil, fmt.Errorf("failed to unmarshal membrane positions: %w", err)
	}

	// Process each position
	for _, position := range positions {
		for _, collateralAsset := range position.CollateralAssets {
			// Check if this denom is in the allowed liquid tokens
			for _, liquidToken := range liquidAllowedDenomsCache {
				if liquidToken.ChainID == chain && liquidToken.IbcDenom == collateralAsset.Asset.Info.NativeToken.Denom {
					// Find the corresponding zone
					for _, zone := range zoneCache {
						if zone.ChainId == liquidToken.RegisteredZoneChainID {
							// Create or update the claim message
							if _, ok := msg[zone.ChainId]; !ok {
								msg[zone.ChainId] = prewards.MsgSubmitClaim{
									UserAddress: submitAddress,
									Zone:        zone.ChainId,
									SrcZone:     chain,
									ClaimType:   cmtypes.ClaimTypeMembrane,
									Proofs:      make([]*cmtypes.Proof, 0),
								}
							}

							// Add the proof
							proof := cmtypes.Proof{
								Data:      abciquery.Response.Value,
								Key:       abciquery.Response.Key,
								ProofOps:  abciquery.Response.ProofOps,
								Height:    abciquery.Response.Height,
								ProofType: "membrane",
							}

							chainMsg := msg[zone.ChainId]
							chainMsg.Proofs = append(chainMsg.Proofs, &proof)
							msg[zone.ChainId] = chainMsg

							// Add to assets
							coin := sdk.NewCoin(liquidToken.QAssetDenom, collateralAsset.Asset.Amount)
							assets[chain] = assets[chain].Add(coin)

							break
						}
					}
					break
				}
			}
		}
	}

	log.Debug("Membrane claim processing completed", "address", address, "chain", chain, "positions_count", len(msg))
	return msg, assets, nil
}
