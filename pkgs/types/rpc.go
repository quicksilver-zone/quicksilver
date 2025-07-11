package types

import (
	"fmt"
	"log"
	"sync"
	"time"

	tmhttp "github.com/cometbft/cometbft/rpc/client/http"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/multierror"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

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
	Errors   multierror.MultiError     `json:"errors,omitempty"`
}

func (response *Response) Update(messages map[string]prewards.MsgSubmitClaim, assets map[string]sdk.Coins, assetType string) {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	for _, message := range messages {
		if message.ClaimType == types.ClaimTypeUndefined {
			log.Default().Println("ERROR: skipping message with undefined claim: ", message)
			continue
		}
		fmt.Printf("adding message. claim_type: %s, zone: %s, src_zone: %s, user_address: %s\n", message.ClaimType, message.Zone, message.SrcZone, message.UserAddress)
		response.Messages = append(response.Messages, message)

	}

	for chainID, asset := range assets {
		if asset.IsZero() {
			log.Default().Printf("ERROR: skipping asset with no value: chain_id: %s, asset: %s\n", chainID, asset)
			continue
		}
		fmt.Printf("Adding asset: chain_id: %s, asset: %s\n", chainID, asset)
		response.Assets[chainID] = append(response.Assets[chainID], Asset{Type: assetType, Amount: asset})
	}
}

type Asset struct {
	Type   string
	Amount sdk.Coins
}
