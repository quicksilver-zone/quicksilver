package runner

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	stdlog "log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quicksilver-zone/quicksilver/v7/icq-relayer/pkg/config"
	"github.com/quicksilver-zone/quicksilver/v7/icq-relayer/prommetrics"

	"github.com/go-kit/log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmquery "github.com/cometbft/cometbft/libs/pubsub/query"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	qstypes "github.com/quicksilver-zone/quicksilver/v7/x/interchainquery/types"
	lensclient "github.com/strangelove-ventures/lens/client"
	lensquery "github.com/strangelove-ventures/lens/client/query"
	"google.golang.org/grpc/metadata"

	cmtjson "github.com/cometbft/cometbft/libs/json"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	jsonrpctypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	tmtypes "github.com/cometbft/cometbft/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/dgraph-io/ristretto"
)

type Clients []*lensclient.ChainClient

const VERSION = "icq/v0.10.0"

var (
	WaitInterval          = time.Second * 6
	HistoricQueryInterval = time.Second * 15
	MaxHistoricQueries    = 12
	MaxTxMsgs             = 12
	ctx                   = context.Background()
	sendQueue             = map[string]chan sdk.Msg{}
	cache                 *ristretto.Cache
	globalCfg             *config.Config
)

func (clients Clients) GetForChainId(chainId string) *lensclient.ChainClient {
	for _, chainClient := range clients {
		if chainClient.Config.ChainID == chainId {
			return chainClient
		}
	}
	return nil
}

func Run(cfg *config.Config, home string) error {
	globalCfg = cfg
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	_ = logger.Log("worker", "init", "msg", "starting icq relayer", "version", VERSION)
	_ = logger.Log("worker", "init", "msg", "permitted queries", "queries", strings.Join(globalCfg.AllowedQueries, ","))

	reg := prometheus.NewRegistry()
	metrics := *prommetrics.NewMetrics(reg)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // Num keys to track frequency of (10M).
		MaxCost:     1 << 30, // Maximum cost of cache (1GB).
		BufferItems: 64,      // Number of keys per Get buffer.
	})
	if err != nil {
		panic("unable to start ristretto cache")
	}

	http.Handle("/metrics", promHandler)
	go func() {
		stdlog.Fatal(http.ListenAndServe(":2112", nil))
	}()

	defer func() {
		err := Close()
		if err != nil {
			stdlog.Fatal("Error in Closing the routine")
		}
	}()
	for _, c := range cfg.Chains {
		cfg.Cl[c.ChainID], err = lensclient.NewChainClient(nil, c, home, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}

		err = logger.Log("worker", "init", "msg", "configured chain", "chain", c.ChainID)
		if err != nil {
			return err
		}
		sendQueue[c.ChainID] = make(chan sdk.Msg)
		metrics.SendQueue.WithLabelValues("send-queue").Set(float64(len(sendQueue[c.ChainID])))
	}

	query := tmquery.MustParse(fmt.Sprintf("message.module='%s'", "interchainquery"))

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	defaultClient, ok := cfg.Cl[cfg.DefaultChain]
	if !ok {
		panic("unable to create default chainClient; Client is nil")
	}
	err = defaultClient.RPCClient.Start()
	if err != nil {
		_ = logger.Log("error", err.Error())
	}

	_ = logger.Log("worker", "init", "msg", "configuring subscription on default chainClient", "chain", defaultClient.Config.ChainID)

	ch, err := defaultClient.RPCClient.Subscribe(ctx, defaultClient.Config.ChainID+"-icq", query.String())
	if err != nil {
		_ = logger.Log("error", err.Error())
		return err
	}
	wg.Add(1)
	go func(chainId string, ch <-chan coretypes.ResultEvent) {
		defer wg.Done()
		for v := range ch {
			v.Events["source"] = []string{chainId}
			// why does this always trigger twice? messages are deduped later, but this causes 2x queries to trigger.
			time.Sleep(75 * time.Millisecond) // try to avoid thundering herd.
			go handleEvent(v, log.With(logger, "worker", "chainClient", "chain", defaultClient.Config.ChainID), metrics)
		}
	}(defaultClient.Config.ChainID, ch)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := FlushSendQueue(defaultClient.Config.ChainID, log.With(logger, "worker", "flusher", "chain", defaultClient.Config.ChainID), metrics)
		if err != nil {
			_ = logger.Log("Flush Go-routine Bailing")
			panic(err)
		}
	}()

	for _, chainClient := range globalCfg.Cl {
		if chainClient.Config.ChainID != cfg.DefaultChain {
			wg.Add(1)
			go func(defaultClient *lensclient.ChainClient, srcClient *lensclient.ChainClient, logger log.Logger) {
				defer wg.Done()
			CNT:
				for {
					time.Sleep(HistoricQueryInterval)
					req := &qstypes.QueryRequestsRequest{
						Pagination: &querytypes.PageRequest{Limit: 500},
						ChainId:    srcClient.Config.ChainID,
					}

					bz := defaultClient.Codec.Marshaler.MustMarshal(req)
					metrics.HistoricQueryRequests.WithLabelValues("historic_requests").Inc()
					res, err := defaultClient.RPCClient.ABCIQuery(ctx, "/quicksilver.interchainquery.v1.QuerySrvr/Queries", bz)
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							err := logger.Log("error", fmt.Sprintf("timeout: %s", err.Error()))
							if err != nil {
								return
							}
							continue CNT
						}
						panic(fmt.Sprintf("panic(3): %v", err))
					}
					out := &qstypes.QueryRequestsResponse{}
					err = defaultClient.Codec.Marshaler.Unmarshal(res.Response.Value, out)
					if err != nil {
						err := logger.Log("msg", "Error: Unable to unmarshal: ", "error", err)
						if err != nil {
							return
						}
						continue CNT
					}
					_ = logger.Log("worker", "chainClient", "msg", "fetched historic queries for chain", "count", len(out.Queries))

					if len(out.Queries) > 0 {
						go handleHistoricRequests(out.Queries, defaultClient.Config.ChainID, log.With(logger, "worker", "historic"), metrics)
					}
				}
			}(defaultClient, chainClient, log.With(logger, "chain", defaultClient.Config.ChainID, "src_chain", chainClient.Config.ChainID))
		}
	}

	return nil
}

type Query struct {
	SourceChainId string
	ConnectionId  string
	ChainId       string
	QueryId       string
	Type          string
	Height        int64
	Request       []byte
}

func handleHistoricRequests(queries []qstypes.Query, sourceChainId string, logger log.Logger, metrics prommetrics.Metrics) {
	metrics.HistoricQueries.WithLabelValues("historic-queries").Set(float64(len(queries)))

	if len(queries) == 0 {
		return
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(queries), func(i, j int) { queries[i], queries[j] = queries[j], queries[i] })

	sort.SliceStable(queries, func(i, j int) bool {
		return queries[i].CallbackId == "allbalances" || queries[i].CallbackId == "depositinterval" || queries[i].CallbackId == "deposittx" || queries[i].LastEmission.GT(queries[j].LastEmission)
	})

	for _, query := range queries[0:int(math.Min(float64(len(queries)), float64(MaxHistoricQueries)))] {
		_, ok := globalCfg.Cl[query.ChainId]
		if !ok {
			continue
		}

		q := Query{}
		q.SourceChainId = sourceChainId
		q.ChainId = query.ChainId
		q.ConnectionId = query.ConnectionId
		q.QueryId = query.Id
		q.Request = query.Request
		q.Type = query.QueryType

		if _, found := cache.Get("query/" + q.QueryId); found {
			// break if this is in the cache
			continue
		}

		currentheight, found := cache.Get("currentblock/" + q.ChainId)
		if !found {
			block, err := globalCfg.Cl[q.ChainId].RPCClient.Block(ctx, nil)
			if err != nil {
				panic(fmt.Sprintf("panic(6): %v", err))
			}
			currentheight = block.Block.LastCommit.Height - 1
			cache.SetWithTTL("currentblock/"+q.ChainId, currentheight, 1, 6*time.Second)
			logger.Log("msg", "caching currentblock", "height", currentheight)
		} else {
			logger.Log("msg", "using cached currentblock", "height", currentheight)
		}
		q.Height = currentheight.(int64)

		handle := false
		if len(globalCfg.AllowedQueries) == 0 {
			handle = true
		} else {
			for _, msgType := range globalCfg.AllowedQueries {
				if q.Type == msgType {
					handle = true
					break
				}
			}
		}

		if !handle {
			_ = logger.Log("msg", "Ignoring existing query; not a permitted type", "id", query.Id, "type", q.Type)
			continue
		}
		_ = logger.Log("msg", "Handling existing query", "id", query.Id)

		time.Sleep(75 * time.Millisecond) // try to avoid thundering herd.

		go doRequestWithMetrics(q, logger, metrics)
	}
}

func handleEvent(event coretypes.ResultEvent, logger log.Logger, metrics prommetrics.Metrics) {
	queries := []Query{}
	source := event.Events["source"]
	connections := event.Events["message.connection_id"]
	chains := event.Events["message.chain_id"]
	queryIds := event.Events["message.query_id"]
	types := event.Events["message.type"]
	request := event.Events["message.request"]
	height := event.Events["message.height"]

	items := len(queryIds)

	for i := 0; i < items; i++ {
		_, ok := globalCfg.Cl[chains[i]]
		if !ok {
			continue
		}

		req, err := hex.DecodeString(request[i])
		if err != nil {
			panic(fmt.Sprintf("panic(4): %v", err))
		}
		h, err := strconv.ParseInt(height[i], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("panic(5): %v", err))
		}

		handle := false
		if len(globalCfg.AllowedQueries) == 0 {
			handle = true
		} else {
			for _, msgType := range globalCfg.AllowedQueries {
				if types[i] == msgType {
					handle = true
					break
				}
			}
		}

		if !handle {
			_ = logger.Log("msg", "Ignoring current query; not a permitted type", "id", queryIds[i], "type", types[i])
			continue
		}

		if _, found := cache.Get("query/" + queryIds[i]); found {
			// break if this is in the cache
			fmt.Println("avoiding duplicate")
			continue
		}

		if h == 0 {
			currentheight, found := cache.Get("currentblock/" + chains[i])
			if !found {
				block, err := globalCfg.Cl[chains[i]].RPCClient.Block(ctx, nil)
				if err != nil {
					panic(fmt.Sprintf("panic(6): %v", err))
				}
				currentheight = block.Block.LastCommit.Height - 1
				cache.SetWithTTL("currentblock/"+chains[i], currentheight, 1, 6*time.Second)
				logger.Log("msg", "caching currentblock", "height", currentheight)
			} else {
				logger.Log("msg", "using cached currentblock", "height", currentheight)
			}
			h = currentheight.(int64)
		}

		cache.SetWithTTL("query/"+queryIds[i], true, 0, 10*time.Second) // just long enough to not duplicate.
		queries = append(queries, Query{source[0], connections[i], chains[i], queryIds[i], types[i], h, req})
	}

	for _, q := range queries {
		go doRequestWithMetrics(q, log.With(logger, "src_chain", q.ChainId), metrics)
	}
}

func RunGRPCQuery(ctx context.Context, client *lensclient.ChainClient, method string, reqBz []byte, md metadata.MD, metrics prommetrics.Metrics) (abcitypes.ResponseQuery, metadata.MD, error) {
	// parse height header
	height, err := lensclient.GetHeightFromMetadata(md)
	if err != nil {
		return abcitypes.ResponseQuery{}, nil, err
	}

	prove, err := lensclient.GetProveFromMetadata(md)
	if err != nil {
		return abcitypes.ResponseQuery{}, nil, err
	}

	abciReq := abcitypes.RequestQuery{
		Path:   method,
		Data:   reqBz,
		Height: height,
		Prove:  prove,
	}

	metrics.ABCIRequests.WithLabelValues("abci_requests", method).Inc()

	abciRes, err := client.QueryABCI(ctx, abciReq)
	if err != nil {
		return abcitypes.ResponseQuery{}, nil, err
	}
	return abciRes, md, nil
}

func retryLightblock(ctx context.Context, client *lensclient.ChainClient, height int64, maxTime int, logger log.Logger, metrics prommetrics.Metrics) (*tmtypes.LightBlock, error) {
	lightBlock, found := cache.Get("lightblock/" + client.Config.ChainID + "/" + fmt.Sprintf("%d", height))
	var err error
	if !found {
		interval := 1
		_ = logger.Log("msg", "Querying lightblock", "attempt", interval)
		lightBlock, err = client.LightProvider.LightBlock(ctx, height)
		metrics.LightBlockRequests.WithLabelValues("lightblock_requests").Inc()

		if err != nil {
			for {
				time.Sleep(time.Duration(interval) * time.Second)
				_ = logger.Log("msg", "Requerying lightblock", "attempt", interval)
				lightBlock, err = client.LightProvider.LightBlock(ctx, height)
				metrics.LightBlockRequests.WithLabelValues("lightblock_requests").Inc()
				interval = interval + 1
				if err == nil {
					break
				} else if interval > maxTime {
					return nil, fmt.Errorf("unable to query light block, max interval exceeded")
				}
			}
		}
		cache.Set("lightblock/"+client.Config.ChainID+"/"+fmt.Sprintf("%d", height), lightBlock, 5)
	} else {
		_ = logger.Log("msg", "got lightblock from cache")
	}
	return lightBlock.(*tmtypes.LightBlock), err
}

func doRequestWithMetrics(query Query, logger log.Logger, metrics prommetrics.Metrics) {
	startTime := time.Now()
	metrics.Requests.WithLabelValues("requests", query.Type).Inc()
	doRequest(query, logger, metrics)
	endTime := time.Now()
	metrics.RequestsLatency.WithLabelValues("request-latency", query.Type).Observe(endTime.Sub(startTime).Seconds())
}

func doRequest(query Query, logger log.Logger, metrics prommetrics.Metrics) {
	var err error
	client := globalCfg.Cl[query.ChainId]
	if client == nil {
		return
	}

	_ = logger.Log("msg", "Handling request", "type", query.Type, "id", query.QueryId, "height", query.Height)

	newCtx := lensclient.SetHeightOnContext(ctx, query.Height)
	pathParts := strings.Split(query.Type, "/")
	if pathParts[len(pathParts)-1] == "key" { // fetch proof if the query is 'key'
		newCtx = lensclient.SetProveOnContext(newCtx, true)
	}
	inMd, ok := metadata.FromOutgoingContext(newCtx)
	if !ok {
		panic("failed on not ok")
	}

	var res abcitypes.ResponseQuery
	submitClient := globalCfg.Cl[query.SourceChainId]

	switch query.Type {
	// until we fix ordering and pagination in the binary, we can override the query here.
	case "cosmos.tx.v1beta1.Service/GetTxsEvent":
		request := txtypes.GetTxsEventRequest{}
		err = client.Codec.Marshaler.Unmarshal(query.Request, &request)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in Unmarshalling Request", "type", query.Type, "id", query.QueryId, "height", query.Height)
			panic(fmt.Sprintf("panic(7a): %v", err))
		}
		request.OrderBy = txtypes.OrderBy_ORDER_BY_DESC
		request.Limit = 200
		request.Pagination.Limit = 200

		query.Request, err = client.Codec.Marshaler.Marshal(&request)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in Marshalling Request", "type", query.Type, "id", query.QueryId, "height", query.Height)
			panic(fmt.Sprintf("panic(7b): %v", err))
		}

		_ = logger.Log("msg", "Handling GetTxsEvents", "id", query.QueryId, "height", query.Height)
		res, _, err = RunGRPCQuery(ctx, client, "/"+query.Type, query.Request, inMd, metrics)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in RunGRPCQuery", "type", query.Type, "id", query.QueryId, "height", query.Height)
			panic(fmt.Sprintf("panic(7c): %v", err))
		}

	case "tendermint.Tx":
		req := txtypes.GetTxRequest{}
		client.Codec.Marshaler.MustUnmarshal(query.Request, &req)
		hashBytes, err := hex.DecodeString(req.GetHash())
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get decode hash %s", err))
			return
		}
		txRes, height, err := Tx(client, hashBytes)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not fetch proof %s", err))
			return
		}

		protoProof := txRes.ToProto()

		submitQuerier := lensquery.Query{Client: submitClient, Options: lensquery.DefaultOptions()}
		clientId, found := cache.Get("clientId/" + query.ConnectionId)
		if !found {
			connection, err := submitQuerier.Ibc_Connection(query.ConnectionId)
			if err != nil {
				_ = logger.Log("msg", fmt.Sprintf("Error: Could not get connection from chain %s", err))
				return
			}
			clientId = connection.Connection.ClientId
			cache.Set("clientId/"+query.ConnectionId, clientId, 1)
		}

		header, err := getHeader(ctx, client, submitClient, clientId.(string), height-1, logger, true, metrics)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get header %s", err))
			return
		}

		resp := qstypes.GetTxWithProofResponse{Proof: &protoProof, Header: header}
		res.Value = client.Codec.Marshaler.MustMarshal(&resp)

	case "ibc.ClientUpdate":
		submitClientUpdate(client, submitClient, query, int64(sdk.BigEndianToUint64(query.Request)), logger, metrics)
		// return a dummy message to settle the query.
		from, _ := submitClient.GetKeyAddress()
		msg := &qstypes.MsgSubmitQueryResponse{ChainId: query.ChainId, QueryId: query.QueryId, Result: []byte{}, Height: int64(sdk.BigEndianToUint64(query.Request)), ProofOps: &crypto.ProofOps{}, FromAddress: submitClient.MustEncodeAccAddr(from)}
		sendQueue[query.SourceChainId] <- msg
		return
	default:
		res, _, err = RunGRPCQuery(ctx, client, "/"+query.Type, query.Request, inMd, metrics)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in RunGRPCQuery", "type", query.Type, "id", query.QueryId, "height", query.Height)
			panic(fmt.Sprintf("panic(7): %v", err))
		}
	}

	// submit tx to queue
	from, _ := submitClient.GetKeyAddress()
	if pathParts[len(pathParts)-1] == "key" {
		submitClientUpdate(client, submitClient, query, res.Height, logger, metrics)
	}

	msg := &qstypes.MsgSubmitQueryResponse{ChainId: query.ChainId, QueryId: query.QueryId, Result: res.Value, Height: res.Height, ProofOps: res.ProofOps, FromAddress: submitClient.MustEncodeAccAddr(from)}
	sendQueue[query.SourceChainId] <- msg
}

// tm0.37 has a breaking change whereby tx events are no longer base64 encoded, so are represented as string and not bytes.
// As a result, we cannot use the RPCClient.Tx() method which attempts to unmarshal the Result, including the underlying Tx object.
// As such, we want to query the result directly, and unmarshal the json ourselves, to a representation of the result that conveniently
// does not contain the Tx object (that we don't use, because the TxProof already contains a byte representation of tx anyway!)
// Note: this function is compatible with 0.34 and 0.37 representations of transactions.
func Tx(client *lensclient.ChainClient, hash []byte) (tmtypes.TxProof, int64, error) {
	params := map[string]interface{}{
		"hash":  hash,
		"prove": true,
	}

	id := jsonrpctypes.JSONRPCIntID(0)

	request, err := jsonrpctypes.MapToRequest(id, "tx", params)
	if err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("failed to encode params: %w", err)
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	requestBuf := bytes.NewBuffer(requestBytes)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, client.Config.RPCAddr, requestBuf)
	if err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("request failed: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	httpClient, err := jsonrpcclient.DefaultHTTPClient(client.Config.RPCAddr)
	if err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("create client failed: %w", err)
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("post failed: %w", err)
	}

	defer httpResponse.Body.Close()
	defer httpClient.CloseIdleConnections()

	responseBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &jsonrpctypes.RPCResponse{}
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("error unmarshalling: %w", err)
	}

	if response.Error != nil {
		return tmtypes.TxProof{}, 0, response.Error
	}

	// Unmarshal the RawMessage into the result.
	result := TxResultMinimal{}
	if err := cmtjson.Unmarshal(response.Result, &result); err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("error unmarshalling result: %w", err)
	}

	height, err := strconv.Atoi(result.Height)
	if err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("failed to unmarshal tx height: %w", err)
	}

	return result.Proof, int64(height), nil
}

// a minimised representation of the Tx emitted by a Tx query, only containing Height and Proof and thus compatbiel with tm0.34 and tm0.37.
type TxResultMinimal struct {
	Height string          `json:"height"`
	Proof  tmtypes.TxProof `json:"proof"`
}

func submitClientUpdate(client, submitClient *lensclient.ChainClient, query Query, height int64, logger log.Logger, metrics prommetrics.Metrics) {
	from, _ := submitClient.GetKeyAddress()
	submitQuerier := lensquery.Query{Client: submitClient, Options: lensquery.DefaultOptions()}
	clientId, found := cache.Get("clientId/" + query.ConnectionId)
	if !found {
		connection, err := submitQuerier.Ibc_Connection(query.ConnectionId)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get connection from chain %s", err))
			return
		}
		clientId = connection.Connection.ClientId
		cache.Set("clientId/"+query.ConnectionId, clientId, 1)
	}

	header, err := getHeader(ctx, client, submitClient, clientId.(string), height, logger, false, metrics)
	if err != nil {
		_ = logger.Log("msg", fmt.Sprintf("Error: Could not get header %s", err))
		return
	}
	anyHeader, err := clienttypes.PackHeader(header)
	if err != nil {
		_ = logger.Log("msg", fmt.Sprintf("Error: Could not pack header %s", err))
		return
	}

	msg := &clienttypes.MsgUpdateClient{
		ClientId: clientId.(string), // needs to be passed in as part of request.
		Signer:   submitClient.MustEncodeAccAddr(from),
	}

	sendQueue[query.SourceChainId] <- msg
	metrics.SendQueue.WithLabelValues("send-queue").Set(float64(len(sendQueue)))
}

func getHeader(ctx context.Context, client, submitClient *lensclient.ChainClient, clientId string, requestHeight int64, logger log.Logger, historicOk bool, metrics prommetrics.Metrics) (*tmclient.Header, error) {
	submitQuerier := lensquery.Query{Client: submitClient, Options: lensquery.DefaultOptions()}
	state, err := submitQuerier.Ibc_ClientState(clientId) // pass in from request
	if err != nil {
		return nil, fmt.Errorf("error: Could not get state from chain: %q ", err.Error())
	}
	unpackedState, err := clienttypes.UnpackClientState(state.ClientState)
	if err != nil {
		return nil, fmt.Errorf("error: Could not unpack state from chain: %q ", err.Error())
	}

	trustedHeight := unpackedState.GetLatestHeight()
	clientHeight, ok := trustedHeight.(clienttypes.Height)
	if !ok {
		return nil, fmt.Errorf("error: Could coerce trusted height")
	}

	if !historicOk && clientHeight.RevisionHeight >= uint64(requestHeight+1) {
		return nil, fmt.Errorf("trusted height >= request height")
	}

	_ = logger.Log("msg", "Fetching client update for height", "height", requestHeight+1)
	newBlock, err := retryLightblock(ctx, client, int64(requestHeight+1), 5, logger, metrics)
	if err != nil {
		panic(fmt.Sprintf("Error: Could not fetch updated LC from chain - bailing: %v", err))
	}

	trustedBlock, err := retryLightblock(ctx, client, int64(clientHeight.RevisionHeight)+1, 5, logger, metrics)
	if err != nil {
		panic(fmt.Sprintf("Error: Could not fetch updated LC from chain - bailing (2): %v", err))
	}

	valSet := tmtypes.NewValidatorSet(newBlock.ValidatorSet.Validators)
	trustedValSet := tmtypes.NewValidatorSet(trustedBlock.ValidatorSet.Validators)
	protoVal, err := valSet.ToProto()
	if err != nil {
		panic(fmt.Sprintf("Error: Could not get valset from chain: %v", err))
	}
	trustedProtoVal, err := trustedValSet.ToProto()
	if err != nil {
		panic(fmt.Sprintf("Error: Could not get trusted valset from chain: %v", err))
	}

	header := &tmclient.Header{
		SignedHeader:      newBlock.SignedHeader.ToProto(),
		ValidatorSet:      protoVal,
		TrustedHeight:     clientHeight,
		TrustedValidators: trustedProtoVal,
	}

	return header, nil
}

func FlushSendQueue(chainId string, logger log.Logger, metrics prommetrics.Metrics) error {
	time.Sleep(WaitInterval)
	toSend := []sdk.Msg{}
	ch := sendQueue[chainId]

	for {
		if len(toSend) > MaxTxMsgs {
			flush(chainId, toSend, logger, metrics)
			toSend = []sdk.Msg{}
		}
		select {
		case msg := <-ch:
			toSend = append(toSend, msg)
			metrics.SendQueue.WithLabelValues("send-queue").Set(float64(len(sendQueue[chainId])))
		case <-time.After(WaitInterval):
			flush(chainId, toSend, logger, metrics)
			metrics.SendQueue.WithLabelValues("send-queue").Set(float64(len(sendQueue[chainId])))
			toSend = []sdk.Msg{}
		}
	}
}

// TODO: refactor me!
func flush(chainId string, toSend []sdk.Msg, logger log.Logger, metrics prommetrics.Metrics) {
	if len(toSend) > 0 {
		_ = logger.Log("msg", fmt.Sprintf("Sending batch of %d messages", len(toSend)))
		chainClient := globalCfg.Cl[chainId]
		if chainClient == nil {
			return
		}
		// dedupe on queryId
		msgs := unique(toSend, logger)
		if len(msgs) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()
			resp, err := chainClient.SendMsgs(ctx, msgs, VERSION)
			if err != nil {
				switch {
				case resp != nil && resp.Code == 19 && resp.Codespace == "sdk":
					_ = logger.Log("msg", "Tx already in mempool")
				case resp != nil && resp.Code == 12 && resp.Codespace == "sdk":
					_ = logger.Log("msg", "Not enough gas")
				case err.Error() == "context deadline exceeded":
					_ = logger.Log("msg", "Failed to submit in time, retrying")
					resp, err := chainClient.SendMsgs(ctx, msgs, VERSION)
					if err != nil {
						switch {
						case resp != nil && resp.Code == 19 && resp.Codespace == "sdk":
							_ = logger.Log("msg", "Tx already in mempool")
						case resp != nil && resp.Code == 12 && resp.Codespace == "sdk":
							_ = logger.Log("msg", "Not enough gas")
						case err.Error() == "context deadline exceeded":
							_ = logger.Log("msg", "Failed to submit in time, bailing")
						default:
							_ = logger.Log("msg", "Failed to submit after retry; nevermind, we'll try again!", "err", err)
							metrics.FailedTxs.WithLabelValues("failed_txs").Inc()
						}
					}
				default:
					_ = logger.Log("msg", "Failed to submit; nevermind, we'll try again!", "err", err)
					metrics.FailedTxs.WithLabelValues("failed_txs").Inc()
				}
			}
			_ = logger.Log("msg", fmt.Sprintf("Sent batch of %d (deduplicated) messages", len(msgs)))
		}
	}
}

func unique(msgSlice []sdk.Msg, logger log.Logger) []sdk.Msg {
	keys := make(map[string]bool)
	clientUpdateHeights := make(map[string]bool)

	list := []sdk.Msg{}
	for _, entry := range msgSlice {
		msg, ok := entry.(*clienttypes.MsgUpdateClient)
		if ok {
			header, _ := clienttypes.UnpackHeader(msg.Header)
			key := header.GetHeight().String()
			if _, value := clientUpdateHeights[key]; !value {
				clientUpdateHeights[key] = true
				list = append(list, entry)
				_ = logger.Log("msg", "Added ClientUpdate message", "height", key)
			}
			continue
		}
		msg2, ok2 := entry.(*qstypes.MsgSubmitQueryResponse)
		if ok2 {
			if _, value := keys[msg2.QueryId]; !value {
				keys[msg2.QueryId] = true
				list = append(list, entry)
				_ = logger.Log("msg", "Added SubmitResponse message", "id", msg2.QueryId)
			}
		}
	}

	return list
}

func Close() error {
	query := cmtquery.MustParse(fmt.Sprintf("message.module='%s'", "interchainquery"))

	for _, chainClient := range globalCfg.Cl {
		err := chainClient.RPCClient.Unsubscribe(ctx, chainClient.Config.ChainID+"-icq", query.String())
		if err != nil {
			return err
		}
	}
	return nil
}
