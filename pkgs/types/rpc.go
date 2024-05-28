package types

import (
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/multierror"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
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
		fmt.Println("Adding message: ", message)

		response.Messages = append(response.Messages, message)
	}

	for chainID, asset := range assets {
		fmt.Println("Adding asset: ", chainID, asset)
		response.Assets[chainID] = append(response.Assets[chainID], Asset{Type: assetType, Amount: asset})
	}
}

type Asset struct {
	Type   string
	Amount sdk.Coins
}
