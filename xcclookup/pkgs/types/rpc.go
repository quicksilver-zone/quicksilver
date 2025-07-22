package types

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmhttp "github.com/cometbft/cometbft/rpc/client/http"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"

	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/logger"
)

// ErrorString is a custom type that implements json.Marshaler to serialize errors as strings
type ErrorString struct {
	error
}

func (e ErrorString) MarshalJSON() ([]byte, error) {
	if e.error == nil {
		return []byte("null"), nil
	}
	return json.Marshal(e.Error())
}

func (e *ErrorString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		e.error = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	e.error = fmt.Errorf("%s", s)
	return nil
}

func NewRPCClient(addr string, timeout time.Duration) (*tmhttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = timeout
	rpcClient, err := tmhttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}
	return rpcClient, nil
}

type Response struct {
	Messages []prewards.MsgSubmitClaim `json:"messages"`
	Assets   map[string][]Asset        `json:"assets"`
	Errors   *ErrorString              `json:"errors,omitempty"`
	mu       sync.RWMutex
}

func (response *Response) Update(ctx context.Context, messages map[string]prewards.MsgSubmitClaim, assets map[string]sdk.Coins, assetType string) {
	log := logger.FromContext(ctx)
	response.mu.Lock()
	defer response.mu.Unlock()

	for _, message := range messages {
		if message.ClaimType == types.ClaimTypeUndefined {
			log.Error("Skipping message with undefined claim", "message", message)
			continue
		}
		log.Debug("Adding message", "claim_type", message.ClaimType, "zone", message.Zone, "src_zone", message.SrcZone, "user_address", message.UserAddress)
		response.Messages = append(response.Messages, message)

	}

	for chainID, asset := range assets {
		if asset.IsZero() {
			log.Error("Skipping asset with no value", "chain_id", chainID, "asset", asset)
			continue
		}
		log.Debug("Adding asset", "chain_id", chainID, "asset", asset)
		response.Assets[chainID] = append(response.Assets[chainID], Asset{Type: assetType, Amount: asset})
	}
}

// GetAssets returns a copy of the assets map for thread-safe reading
func (response *Response) GetAssets() map[string][]Asset {
	response.mu.RLock()
	defer response.mu.RUnlock()

	// Create a deep copy to avoid race conditions
	result := make(map[string][]Asset)
	for k, v := range response.Assets {
		result[k] = append([]Asset{}, v...)
	}
	return result
}

// GetMessages returns a copy of the messages slice for thread-safe reading
func (response *Response) GetMessages() []prewards.MsgSubmitClaim {
	response.mu.RLock()
	defer response.mu.RUnlock()

	// Create a copy to avoid race conditions
	result := make([]prewards.MsgSubmitClaim, len(response.Messages))
	copy(result, response.Messages)
	return result
}

type Asset struct {
	Type   string
	Amount sdk.Coins
}
