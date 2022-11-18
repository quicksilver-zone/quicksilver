package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	"github.com/ingenuity-build/xcclookup/internal/multierror"
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
	Errors   multierror.Errors         `json:"errors,omitempty"`
}

type Asset struct {
	Type   string
	Amount sdk.Coins
}
