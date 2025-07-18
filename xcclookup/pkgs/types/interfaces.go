package types

import (
	"context"
	"errors"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

// ErrUnsupportedCacheManager is returned when the cache manager type is not supported
var ErrUnsupportedCacheManager = errors.New("unsupported cache manager type")

// RPCClientInterface defines the interface for RPC client operations
type RPCClientInterface interface {
	ABCIQuery(ctx context.Context, path string, data []byte) (*ABCIQueryResponse, error)
	ABCIQueryWithOptions(ctx context.Context, path string, data []byte, opts ABCIQueryOptions) (*ABCIQueryResponse, error)
}

// CacheManagerInterface defines the interface for cache operations
type CacheManagerInterface interface {
	GetConnections(ctx context.Context) ([]prewards.ConnectionProtocolData, error)
	GetOsmosisParams(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error)
	GetOsmosisPools(ctx context.Context) ([]prewards.OsmosisPoolProtocolData, error)
	GetOsmosisClPools(ctx context.Context) ([]prewards.OsmosisClPoolProtocolData, error)
	GetLiquidAllowedDenoms(ctx context.Context) ([]prewards.LiquidAllowedDenomProtocolData, error)
	GetUmeeParams(ctx context.Context) ([]prewards.UmeeParamsProtocolData, error)
	GetMembraneParams(ctx context.Context) ([]prewards.MembraneProtocolData, error)
	GetZones(ctx context.Context) ([]icstypes.Zone, error)
	AddMocks(ctx context.Context, mocks interface{}) error
}

// VersionServiceInterface defines the interface for version operations
type VersionServiceInterface interface {
	GetVersion() ([]byte, error)
}

// OsmosisResult represents the result of Osmosis claim operations
type OsmosisResult struct {
	Err           error
	OsmosisPool   OsmosisPool
	OsmosisClPool OsmosisClPool
}

// OsmosisPool represents Osmosis pool data
type OsmosisPool struct {
	Msg    map[string]prewards.MsgSubmitClaim
	Assets map[string]sdk.Coins
	Err    error
}

// OsmosisClPool represents Osmosis CL pool data
type OsmosisClPool struct {
	Msg    map[string]prewards.MsgSubmitClaim
	Assets map[string]sdk.Coins
	Err    error
}

// ClaimsServiceInterface defines the interface for claims operations
type ClaimsServiceInterface interface {
	OsmosisClaim(ctx context.Context, address, submitAddress, chain string, height int64) (OsmosisResult, error)
	UmeeClaim(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error)
	LiquidClaim(ctx context.Context, address, submitAddress string, connection prewards.ConnectionProtocolData, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error)
	MembraneClaim(ctx context.Context, address, submitAddress, chain string, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error)
}

// OutputFunction defines the interface for output functions
type OutputFunction func(http.ResponseWriter, *Response, map[string]error)

// ABCIQueryResponse represents the response from ABCI queries
type ABCIQueryResponse struct {
	Response ABCIResponse
}

// ABCIResponse represents the ABCI response data
type ABCIResponse struct {
	Value    []byte
	Key      []byte
	ProofOps interface{}
	Height   int64
}

// ABCIQueryOptions represents options for ABCI queries
type ABCIQueryOptions struct {
	Height int64
	Prove  bool
}
