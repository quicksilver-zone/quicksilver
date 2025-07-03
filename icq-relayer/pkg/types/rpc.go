package types

import (
	"context"
	"time"

	"github.com/cometbft/cometbft/libs/bytes"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cometbft/cometbft/types"
)

// RPCClient is a client for the Tendermint RPC interface. This is copy and pasted from github.com/cometbft/cometbft/rpc/client/http,
// which does not expost the baseRPCClient, and as such does not permit you to create a HTTP RPC Client without websocket connection.
type RPCClient struct {
	caller jsonrpcclient.Caller
}

type RPCClientI interface {
	rpcclient.ABCIClient
	rpcclient.HistoryClient
	rpcclient.NetworkClient
	rpcclient.SignClient
	rpcclient.StatusClient
}

var (
	_ RPCClientI = &RPCClient{}
	_ RPCClientI = &rpchttp.HTTP{}
)

// Create timeout enabled http client
func NewWithTimeout(remote string, timeout uint) (*RPCClient, error) {
	httpClient, err := jsonrpcclient.DefaultHTTPClient(remote)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = time.Duration(timeout) * time.Second

	rc, err := jsonrpcclient.NewWithHTTPClient(remote, httpClient)
	if err != nil {
		return nil, err
	}

	return &RPCClient{caller: rc}, nil
}

func (c *RPCClient) Status(ctx context.Context) (*ctypes.ResultStatus, error) {
	result := new(ctypes.ResultStatus)
	_, err := c.caller.Call(ctx, "status", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RPCClient) ABCIInfo(ctx context.Context) (*ctypes.ResultABCIInfo, error) {
	result := new(ctypes.ResultABCIInfo)
	_, err := c.caller.Call(ctx, "abci_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RPCClient) ABCIQuery(
	ctx context.Context,
	path string,
	data bytes.HexBytes,
) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryWithOptions(ctx, path, data, rpcclient.DefaultABCIQueryOptions)
}

func (c *RPCClient) ABCIQueryWithOptions(
	ctx context.Context,
	path string,
	data bytes.HexBytes,
	opts rpcclient.ABCIQueryOptions,
) (*ctypes.ResultABCIQuery, error) {
	result := new(ctypes.ResultABCIQuery)
	_, err := c.caller.Call(ctx, "abci_query",
		map[string]interface{}{"path": path, "data": data, "height": opts.Height, "prove": opts.Prove},
		result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RPCClient) BroadcastTxCommit(
	ctx context.Context,
	tx types.Tx,
) (*ctypes.ResultBroadcastTxCommit, error) {
	result := new(ctypes.ResultBroadcastTxCommit)
	_, err := c.caller.Call(ctx, "broadcast_tx_commit", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) BroadcastTxAsync(
	ctx context.Context,
	tx types.Tx,
) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX(ctx, "broadcast_tx_async", tx)
}

func (c *RPCClient) BroadcastTxSync(
	ctx context.Context,
	tx types.Tx,
) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX(ctx, "broadcast_tx_sync", tx)
}

func (c *RPCClient) broadcastTX(
	ctx context.Context,
	route string,
	tx types.Tx,
) (*ctypes.ResultBroadcastTx, error) {
	result := new(ctypes.ResultBroadcastTx)
	_, err := c.caller.Call(ctx, route, map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) UnconfirmedTxs(
	ctx context.Context,
	limit *int,
) (*ctypes.ResultUnconfirmedTxs, error) {
	result := new(ctypes.ResultUnconfirmedTxs)
	params := make(map[string]interface{})
	if limit != nil {
		params["limit"] = limit
	}
	_, err := c.caller.Call(ctx, "unconfirmed_txs", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) NumUnconfirmedTxs(ctx context.Context) (*ctypes.ResultUnconfirmedTxs, error) {
	result := new(ctypes.ResultUnconfirmedTxs)
	_, err := c.caller.Call(ctx, "num_unconfirmed_txs", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) CheckTx(ctx context.Context, tx types.Tx) (*ctypes.ResultCheckTx, error) {
	result := new(ctypes.ResultCheckTx)
	_, err := c.caller.Call(ctx, "check_tx", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) NetInfo(ctx context.Context) (*ctypes.ResultNetInfo, error) {
	result := new(ctypes.ResultNetInfo)
	_, err := c.caller.Call(ctx, "net_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) DumpConsensusState(ctx context.Context) (*ctypes.ResultDumpConsensusState, error) {
	result := new(ctypes.ResultDumpConsensusState)
	_, err := c.caller.Call(ctx, "dump_consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) ConsensusState(ctx context.Context) (*ctypes.ResultConsensusState, error) {
	result := new(ctypes.ResultConsensusState)
	_, err := c.caller.Call(ctx, "consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) ConsensusParams(
	ctx context.Context,
	height *int64,
) (*ctypes.ResultConsensusParams, error) {
	result := new(ctypes.ResultConsensusParams)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = height
	}
	_, err := c.caller.Call(ctx, "consensus_params", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) Health(ctx context.Context) (*ctypes.ResultHealth, error) {
	result := new(ctypes.ResultHealth)
	_, err := c.caller.Call(ctx, "health", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) BlockchainInfo(
	ctx context.Context,
	minHeight,
	maxHeight int64,
) (*ctypes.ResultBlockchainInfo, error) {
	result := new(ctypes.ResultBlockchainInfo)
	_, err := c.caller.Call(ctx, "blockchain",
		map[string]interface{}{"minHeight": minHeight, "maxHeight": maxHeight},
		result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) Genesis(ctx context.Context) (*ctypes.ResultGenesis, error) {
	result := new(ctypes.ResultGenesis)
	_, err := c.caller.Call(ctx, "genesis", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) GenesisChunked(ctx context.Context, id uint) (*ctypes.ResultGenesisChunk, error) {
	result := new(ctypes.ResultGenesisChunk)
	_, err := c.caller.Call(ctx, "genesis_chunked", map[string]interface{}{"chunk": id}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) Block(ctx context.Context, height *int64) (*ctypes.ResultBlock, error) {
	result := new(ctypes.ResultBlock)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = height
	}
	_, err := c.caller.Call(ctx, "block", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) BlockByHash(ctx context.Context, hash []byte) (*ctypes.ResultBlock, error) {
	result := new(ctypes.ResultBlock)
	params := map[string]interface{}{
		"hash": hash,
	}
	_, err := c.caller.Call(ctx, "block_by_hash", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) BlockResults(
	ctx context.Context,
	height *int64,
) (*ctypes.ResultBlockResults, error) {
	result := new(ctypes.ResultBlockResults)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = height
	}
	_, err := c.caller.Call(ctx, "block_results", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) Commit(ctx context.Context, height *int64) (*ctypes.ResultCommit, error) {
	result := new(ctypes.ResultCommit)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = height
	}
	_, err := c.caller.Call(ctx, "commit", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) Tx(ctx context.Context, hash []byte, prove bool) (*ctypes.ResultTx, error) {
	result := new(ctypes.ResultTx)
	params := map[string]interface{}{
		"hash":  hash,
		"prove": prove,
	}
	_, err := c.caller.Call(ctx, "tx", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) TxSearch(
	ctx context.Context,
	query string,
	prove bool,
	page,
	perPage *int,
	orderBy string,
) (*ctypes.ResultTxSearch, error) {
	result := new(ctypes.ResultTxSearch)
	params := map[string]interface{}{
		"query":    query,
		"prove":    prove,
		"order_by": orderBy,
	}

	if page != nil {
		params["page"] = page
	}
	if perPage != nil {
		params["per_page"] = perPage
	}

	_, err := c.caller.Call(ctx, "tx_search", params, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RPCClient) BlockSearch(
	ctx context.Context,
	query string,
	page, perPage *int,
	orderBy string,
) (*ctypes.ResultBlockSearch, error) {
	result := new(ctypes.ResultBlockSearch)
	params := map[string]interface{}{
		"query":    query,
		"order_by": orderBy,
	}

	if page != nil {
		params["page"] = page
	}
	if perPage != nil {
		params["per_page"] = perPage
	}

	_, err := c.caller.Call(ctx, "block_search", params, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RPCClient) Validators(
	ctx context.Context,
	height *int64,
	page,
	perPage *int,
) (*ctypes.ResultValidators, error) {
	result := new(ctypes.ResultValidators)
	params := make(map[string]interface{})
	if page != nil {
		params["page"] = page
	}
	if perPage != nil {
		params["per_page"] = perPage
	}
	if height != nil {
		params["height"] = height
	}
	_, err := c.caller.Call(ctx, "validators", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) BroadcastEvidence(
	ctx context.Context,
	ev types.Evidence,
) (*ctypes.ResultBroadcastEvidence, error) {
	result := new(ctypes.ResultBroadcastEvidence)
	_, err := c.caller.Call(ctx, "broadcast_evidence", map[string]interface{}{"evidence": ev}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) Header(ctx context.Context, height *int64) (*ctypes.ResultHeader, error) {
	result := new(ctypes.ResultHeader)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = height
	}
	_, err := c.caller.Call(ctx, "header", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) HeaderByHash(ctx context.Context, hash bytes.HexBytes) (*ctypes.ResultHeader, error) {
	result := new(ctypes.ResultHeader)
	params := map[string]interface{}{
		"hash": hash,
	}
	_, err := c.caller.Call(ctx, "header_by_hash", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
