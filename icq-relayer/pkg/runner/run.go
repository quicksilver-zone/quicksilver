package runner

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	stdlog "log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	sdk "github.com/cosmos/cosmos-sdk/types"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtxtypes "github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	tmclient "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
	"github.com/dgraph-io/ristretto"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/types"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/prommetrics"
	qstypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	VERSION             = "icq-relayer"
	QUICKSILVER_VERSION = ""
	COMMIT              = ""
)

type ClientUpdateRequirement struct {
	ConnectionId string
	ChainId      string
	Height       int64
}

type Message struct {
	Msg          sdk.Msg
	ClientUpdate *ClientUpdateRequirement
}

var (
	MaxTxMsgs             int
	WaitInterval          = time.Second * 6
	HistoricQueryInterval = time.Second * 15
	TxMsgs                = MaxTxMsgs
	ctx                   = context.Background()
	sendQueue             = make(chan Message)
	cache                 *ristretto.Cache
	LastReduced           = time.Now()

	// Variables used for retries
	RtyAttNum = uint(5)
	RtyAtt    = retry.Attempts(RtyAttNum)
	RtyDel    = retry.Delay(time.Millisecond * 800)
	RtyErr    = retry.LastErrorOnly(true)
)

func Run(ctx context.Context, cfg *types.Config, errHandler func(error)) error {

	MaxTxMsgs = cfg.MaxMsgsPerTx
	if MaxTxMsgs == 0 {
		MaxTxMsgs = 30
	}
	TxMsgs = MaxTxMsgs

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	_ = logger.Log("worker", "init", "msg", "starting icq relayer", "version", VERSION)
	_ = logger.Log("worker", "init", "msg", "permitted queries", "queries", strings.Join(cfg.AllowedQueries, ","))

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
		stdlog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.BindPort), nil))
	}()

	defer func() {
		err := Close(cfg)
		if err != nil {
			logger.Log("worker", "init", "msg", "error in closing the routine")
		}
	}()

	if err := cfg.DefaultChain.Init(cfg.ProtoCodec, cache); err != nil {
		fmt.Println(err)
		return err
	}

	for _, c := range cfg.Chains {
		if err := c.Init(cfg.ProtoCodec, cache); err != nil {
			return err
		}

		_ = logger.Log("worker", "init", "msg", "configured chain", "chain", c.ChainID)
	}

	query := tmquery.MustParse(fmt.Sprintf("message.module='%s'", "interchainquery"))

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	_ = logger.Log("worker", "init", "msg", "configuring subscription on default chainClient", "chain", cfg.DefaultChain.ChainID)

	ch, err := cfg.DefaultChain.GetClient().Subscribe(ctx, cfg.DefaultChain.ChainID+"-icq", query.String())
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
			go handleEvent(cfg, v, log.With(logger, "worker", "chainClient", "chain", cfg.DefaultChain.ChainID), metrics)
		}
	}(cfg.DefaultChain.ChainID, ch)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := FlushSendQueue(cfg, log.With(logger, "worker", "flusher", "chain", cfg.DefaultChain.ChainID), metrics)
		if err != nil {
			_ = logger.Log("Flush Go-routine Bailing")
			panic(err)
		}
	}()

	for _, chainCfg := range cfg.Chains {
		wg.Add(1)
		go func(defaultClient *types.ChainConfig, srcClient *types.ReadOnlyChainConfig, logger log.Logger) {
			defer wg.Done()
		CNT:
			for {
				time.Sleep(HistoricQueryInterval)
				req := &qstypes.QueryRequestsRequest{
					Pagination: &querytypes.PageRequest{Limit: 500},
					ChainId:    srcClient.ChainID,
				}

				bz := cfg.ProtoCodec.MustMarshal(req)
				metrics.HistoricQueryRequests.WithLabelValues("historic_requests").Inc()
				var res *coretypes.ResultABCIQuery
				if err = retry.Do(func() error {
					res, err = defaultClient.Client.ABCIQuery(ctx, "/quicksilver.interchainquery.v1.QuerySrvr/Queries", bz)
					return err
				}, RtyAtt, RtyDel, RtyErr); err != nil {
					continue CNT
				}
				out := qstypes.QueryRequestsResponse{}
				err = cfg.ProtoCodec.Unmarshal(res.Response.Value, &out)
				if err != nil {
					err := logger.Log("msg", "Error: Unable to unmarshal: ", "error", err)
					if err != nil {
						return
					}
					continue CNT
				}
				_ = logger.Log("worker", "chainClient", "msg", "fetched historic queries for chain", "count", len(out.Queries))

				if len(out.Queries) > 0 {
					go handleHistoricRequests(cfg, srcClient, out.Queries, cfg.DefaultChain.ChainID, log.With(logger, "worker", "historic"), metrics)
				}
			}
		}(cfg.DefaultChain, chainCfg, log.With(logger, "chain", cfg.DefaultChain.ChainID, "src_chain", chainCfg.ChainID))
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

func handleHistoricRequests(cfg *types.Config, queryClient *types.ReadOnlyChainConfig, queries []qstypes.Query, sourceChainId string, logger log.Logger, metrics prommetrics.Metrics) {
	metrics.HistoricQueries.WithLabelValues("historic-queries").Set(float64(len(queries)))

	if len(queries) == 0 {
		return
	}

	rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(queries), func(i, j int) { queries[i], queries[j] = queries[j], queries[i] })

	sort.SliceStable(queries, func(i, j int) bool {
		return queries[i].CallbackId == "allbalances" || queries[i].CallbackId == "depositinterval" || queries[i].CallbackId == "deposittx" || queries[i].LastEmission.GT(queries[j].LastEmission)
	})

	for _, query := range queries[0:int(math.Min(float64(len(queries)), float64(TxMsgs)))] {
		q := Query{}
		q.SourceChainId = sourceChainId
		q.ChainId = query.ChainId
		q.ConnectionId = query.ConnectionId
		q.QueryId = query.Id
		q.Request = query.Request
		q.Type = query.QueryType

		if _, found := cache.Get("query/" + q.QueryId); found {
			logger.Log("msg", "Query already in cache", "id", q.QueryId)
			continue
		}

		if _, found := cache.Get("ignore/" + q.QueryId); found {
			logger.Log("msg", "Query already in ignore cache", "id", q.QueryId)
			continue
		}

		var err error
		q.Height, err = queryClient.GetCurrentHeight(ctx, cache, log.With(logger, "chain", queryClient.ChainID))
		if err != nil {
			_ = logger.Log("msg", "Error getting current height", "error", err)
			continue
		}

		handle := false
		if len(cfg.AllowedQueries) == 0 {
			handle = true
		} else {
			for _, msgType := range cfg.AllowedQueries {
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
		//_ = logger.Log("msg", "Handling existing query", "id", query.Id)

		time.Sleep(75 * time.Millisecond) // try to avoid thundering herd.

		cache.Set("query/"+q.QueryId, true, 0)

		go doRequestWithMetrics(cfg, q, logger, metrics)
	}
}

func handleEvent(cfg *types.Config, event coretypes.ResultEvent, logger log.Logger, metrics prommetrics.Metrics) {
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
		client, ok := cfg.Chains[chains[i]]
		if !ok {
			continue
		}
		req, err := hex.DecodeString(request[i])
		if err != nil {
			_ = logger.Log("worker", "handler", "msg", err.Error())
			continue
		}
		h, err := strconv.ParseInt(height[i], 10, 64)
		if err != nil {
			_ = logger.Log("worker", "handler", "msg", err.Error())
			continue
		}

		handle := false
		if len(cfg.AllowedQueries) == 0 {
			handle = true
		} else {
			for _, msgType := range cfg.AllowedQueries {
				if types[i] == msgType {
					handle = true
					break
				}
			}
		}

		if !handle {
			_ = logger.Log("worker", "handler", "msg", "Ignoring current query; not a permitted type", "id", queryIds[i], "type", types[i])
			continue
		}
		if _, found := cache.Get("query/" + queryIds[i]); found {
			// break if this is in the cache
			continue
		}

		if _, found := cache.Get("ignore/" + queryIds[i]); found {
			logger.Log("msg", "Query already in ignore cache", "id", queryIds[i])
			// break if this is in the cache
			continue
		}

		if h == 0 {
			currentheight, err := client.GetCurrentHeight(ctx, cache, log.With(logger, "chain", chains[i]))
			if err != nil {
				logger.Log("worker", "handler", "msg", "error getting current block", "height", currentheight)
				continue
			}
		}
		cache.Set("query/"+queryIds[i], true, 0) // just long enough to not duplicate.
		queries = append(queries, Query{source[0], connections[i], chains[i], queryIds[i], types[i], h, req})
	}

	for _, q := range queries {
		go doRequestWithMetrics(cfg, q, log.With(logger, "src_chain", q.ChainId), metrics)
		time.Sleep(75 * time.Millisecond) // try to avoid thundering herd.
	}
}

var lightblockMutexesMutex = sync.Mutex{}
var lightblockMutexes = make(map[string]map[int64]*sync.Mutex)

func retryLightblock(ctx context.Context, chain *types.ReadOnlyChainConfig, height int64, metrics prommetrics.Metrics) (*tmtypes.LightBlock, error) {
	// use an outer mutex to block concurrent read/writes tofrom the mutex map.
	lightblockMutexesMutex.Lock()
	// use a mutex to block any previous request for this height, so we don't get many concurrent requests for the same height.
	if lightblockMutexes[chain.ChainID] == nil {
		lightblockMutexes[chain.ChainID] = make(map[int64]*sync.Mutex)
	}
	if lightblockMutexes[chain.ChainID][height] == nil {
		lightblockMutexes[chain.ChainID][height] = &sync.Mutex{}
	}
	lightblockMutexes[chain.ChainID][height].Lock()
	defer lightblockMutexes[chain.ChainID][height].Unlock()
	lightblockMutexesMutex.Unlock()

	lightBlock, found := cache.Get("lightblock/" + chain.ChainID + "/" + fmt.Sprintf("%d", height))
	var err error
	if !found {
		if err = retry.Do(func() error {
			ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(chain.QueryTimeoutSeconds))
			defer cancel()
			lightBlock, err = chain.LightBlock(ctx, height)
			metrics.LightBlockRequests.WithLabelValues("lightblock_requests").Inc()
			return err
		}, RtyAtt, RtyDel, RtyErr); err != nil {
			return nil, errors.New("unable to query light block, max interval exceeded")
		}
		cache.Set("lightblock/"+chain.ChainID+"/"+fmt.Sprintf("%d", height), lightBlock, 5)
	}
	return lightBlock.(*tmtypes.LightBlock), err
}

func doRequestWithMetrics(cfg *types.Config, query Query, logger log.Logger, metrics prommetrics.Metrics) {
	startTime := time.Now()
	metrics.Requests.WithLabelValues("requests", query.Type).Inc()
	doRequest(cfg, query, logger, metrics)
	endTime := time.Now()
	metrics.RequestsLatency.WithLabelValues("request-latency", query.Type).Observe(endTime.Sub(startTime).Seconds())
}

func doRequest(cfg *types.Config, query Query, logger log.Logger, metrics prommetrics.Metrics) {
	var err error
	client, ok := cfg.Chains[query.ChainId]
	if !ok {
		return
	}

	//_ = logger.Log("msg", "Handling request", "type", query.Type, "id", query.QueryId, "height", query.Height)

	pathParts := strings.Split(query.Type, "/")
	prove := pathParts[len(pathParts)-1] == "key" // fetch proof if the query is 'key'

	var res abcitypes.ResponseQuery

	switch query.Type {
	// until we fix ordering and pagination in the binary, we can override the query here.
	case "cosmos.tx.v1beta1.Service/GetTxsEvent":
		request := qstypes.GetTxsEventRequest{}
		err = cfg.ProtoCodec.Unmarshal(query.Request, &request)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in Unmarshalling Request", "type", query.Type, "id", query.QueryId, "height", query.Height)
			return
		}

		request.OrderBy = txtypes.OrderBy_ORDER_BY_DESC
		request.Limit = cfg.MaxTxsPerQuery
		request.Pagination.Limit = cfg.MaxTxsPerQuery

		query.Request, err = cfg.ProtoCodec.Marshal(&request)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in Marshalling Request", "type", query.Type, "id", query.QueryId, "height", query.Height)
			return
		}

		_ = logger.Log("msg", "Handling GetTxsEvents", "id", query.QueryId, "height", query.Height)
		res, err = client.RunABCIQuery(ctx, "/"+query.Type, query.Request, query.Height, prove, metrics)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in RunGRPCQuery", "type", query.Type, "id", query.QueryId, "height", query.Height)
			return
		}

		txDecoder := authtxtypes.DefaultTxDecoder(cfg.ProtoCodec)
		resp := txtypes.GetTxsEventResponse{}
		err = cfg.ProtoCodec.Unmarshal(res.Value, &resp)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in Unmarshalling Response", "type", query.Type, "id", query.QueryId, "height", query.Height)
		}
		_ = logger.Log("msg", "Got transactions", "num", len(resp.TxResponses), "id", query.QueryId, "height", query.Height)
		txResponses := make([]*sdk.TxResponse, 0, len(resp.TxResponses))
		for _, txr := range resp.TxResponses {
			tx, err := txDecoder(txr.Tx.Value)

			if err != nil {
				_ = logger.Log("msg", "Error: Failed in Unpacking Any", "type", query.Type, "id", query.QueryId, "height", query.Height)
			}
			msgs := tx.GetMsgs()
			for _, msg := range msgs {
				_, ok := msg.(*banktypes.MsgSend)
				if ok {
					txResponses = append(txResponses, txr)
				} else {
					fmt.Println("found non MsgSend message; removing")
				}
			}
		}
		resp.TxResponses = txResponses
		res.Value, err = cfg.ProtoCodec.Marshal(&resp)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in Marshalling Response", "type", query.Type, "id", query.QueryId, "height", query.Height)
			return
		}

	case "tendermint.Tx":
		// custom request type for fetching a tx with proof.
		req := txtypes.GetTxRequest{}
		cfg.ProtoCodec.MustUnmarshal(query.Request, &req)
		hashBytes, err := hex.DecodeString(req.GetHash())
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get decode hash %s", err))
			return
		}

		proofAny, height, err := client.Tx(hashBytes)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not fetch proof %s", err))
			return
		}

		clientId, err := cfg.DefaultChain.GetClientId(ctx, query.ConnectionId, logger, metrics)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get client id %s", err))
			return
		}

		header, err := getHeader(ctx, cfg, client, clientId, height-1, logger, true, metrics)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get header %s", err))
			return
		}

		resp := qstypes.GetTxWithProofResponse{Header: header, ProofAny: proofAny}
		res.Value = cfg.ProtoCodec.MustMarshal(&resp)

	case "ibc.ClientUpdate":
		// return a dummy message to settle the query.
		msg := &qstypes.MsgSubmitQueryResponse{ChainId: query.ChainId, QueryId: query.QueryId, Result: []byte{}, Height: int64(sdk.BigEndianToUint64(query.Request)), ProofOps: &crypto.ProofOps{}, FromAddress: cfg.DefaultChain.GetAddress()}
		sendQueue <- Message{Msg: msg, ClientUpdate: &ClientUpdateRequirement{ConnectionId: query.ConnectionId, ChainId: query.ChainId, Height: int64(sdk.BigEndianToUint64(query.Request))}}
		return
	default:
		res, err = client.RunABCIQuery(ctx, "/"+query.Type, query.Request, query.Height, prove, metrics)
		if err != nil {
			_ = logger.Log("msg", "Error: Failed in RunGRPCQuery", "type", query.Type, "id", query.QueryId, "height", query.Height)
			return
		}
	}

	msg := &qstypes.MsgSubmitQueryResponse{ChainId: query.ChainId, QueryId: query.QueryId, Result: res.Value, Height: res.Height, ProofOps: res.ProofOps, FromAddress: cfg.DefaultChain.GetAddress()}
	var clientUpdate *ClientUpdateRequirement

	if prove {
		clientUpdate = &ClientUpdateRequirement{ConnectionId: query.ConnectionId, ChainId: query.ChainId, Height: res.Height}
	}

	sendQueue <- Message{Msg: msg, ClientUpdate: clientUpdate}
}

func asyncCacheClientUpdate(ctx context.Context, cfg *types.Config, client *types.ReadOnlyChainConfig, query Query, height int64, logger log.Logger, metrics prommetrics.Metrics) error {
	cacheKey := fmt.Sprintf("cu/%s-%d", query.ConnectionId, height)
	queryKey := fmt.Sprintf("cuquery/%s-%d", query.ConnectionId, height)

	_, ok := cache.Get("cu/" + cacheKey)
	if ok {
		fmt.Println("cache found for ", cacheKey)
		return nil
	} else {
		_, ok := cache.Get(queryKey) // lock on querying the same block
		if ok {
			//_ = logger.Log("msg", "cache miss, but cuquery/"+cacheKey+" exists, skipping")
			return nil
		}
		cache.SetWithTTL(queryKey, true, 1, 10*time.Second)
		defer cache.Del(queryKey)

		clientId, err := cfg.DefaultChain.GetClientId(ctx, query.ConnectionId, logger, metrics)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get clientId %s", err))
			return err
		}
		header, err := getHeader(ctx, cfg, client, clientId, height, logger, false, metrics)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get header %s", err))
			return err
		}

		msg, err := clienttypes.NewMsgUpdateClient(clientId, header, cfg.DefaultChain.GetAddress())
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not create msg update: %s", err))
			return err
		}
		cache.SetWithTTL(cacheKey, msg, 5, 10*time.Minute)
		return nil
	}
}

func getCachedClientUpdate(connectionId string, height int64) (sdk.Msg, error) {
	cacheKey := fmt.Sprintf("cu/%s-%d", connectionId, height)
	cu, ok := cache.Get(cacheKey)
	if ok {
		fmt.Printf("cache hit for %s-%d\n", connectionId, height)
		return cu.(sdk.Msg), nil
	}
	fmt.Printf("cache miss for %s-%d\n", connectionId, height)
	return nil, errors.New("client update not found")
}

func getHeader(ctx context.Context, cfg *types.Config, client *types.ReadOnlyChainConfig, clientId string, requestHeight int64, logger log.Logger, historicOk bool, metrics prommetrics.Metrics) (*tmclient.Header, error) {
	state, err := cfg.DefaultChain.GetClientState(ctx, clientId, logger, metrics)
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
		return nil, errors.New("error: Could coerce trusted height")
	}

	if !historicOk && clientHeight.RevisionHeight >= uint64(requestHeight+1) {
		oldHeights, err := cfg.DefaultChain.GetClientStateHeights(ctx, clientId, client.ChainID, uint64(requestHeight-500), logger, metrics, 0)
		if err != nil {
			return nil, fmt.Errorf("error: Could not get old heights: %w", err)
		}
		clientHeight = oldHeights[0]
	}

	newBlock, err := retryLightblock(ctx, client, int64(requestHeight+1), metrics)
	if err != nil {
		return nil, fmt.Errorf("error: Could not fetch new light block from chain: %v", err)
	}
	trustedBlock, err := retryLightblock(ctx, client, int64(clientHeight.RevisionHeight)+1, metrics)
	if err != nil {
		return nil, fmt.Errorf("error: Could not fetch trusted light block from chain: %v", err)
	}

	valSet := tmtypes.NewValidatorSet(newBlock.ValidatorSet.Validators)
	trustedValSet := tmtypes.NewValidatorSet(trustedBlock.ValidatorSet.Validators)
	protoVal, err := valSet.ToProto()
	if err != nil {
		return nil, fmt.Errorf("error: Could not fetch new valset from chain: %v", err)
	}

	trustedProtoVal, err := trustedValSet.ToProto()
	if err != nil {
		return nil, fmt.Errorf("error: Could not fetch trusted valset from chain: %v", err)
	}

	header := &tmclient.Header{
		SignedHeader:      newBlock.SignedHeader.ToProto(),
		ValidatorSet:      protoVal,
		TrustedHeight:     clientHeight,
		TrustedValidators: trustedProtoVal,
	}

	return header, nil
}

func FlushSendQueue(cfg *types.Config, logger log.Logger, metrics prommetrics.Metrics) error {
	time.Sleep(WaitInterval)
	toSend := []Message{}
	ch := sendQueue

	for {
		if LastReduced.Add(time.Second * 30).Before(time.Now()) {
			if 2*TxMsgs > MaxTxMsgs {
				TxMsgs = MaxTxMsgs
			} else {
				TxMsgs = 2 * TxMsgs
			}
			_ = logger.Log("msg", "increased batchsize", "size", TxMsgs)
			LastReduced = time.Now()
		}

		if len(toSend) > 5*TxMsgs {
			flush(cfg, toSend, logger, metrics)
			toSend = []Message{}
		}
		select {
		case msg := <-ch:
			toSend = append(toSend, msg)
			metrics.SendQueue.WithLabelValues("send-queue").Set(float64(len(sendQueue)))
			if msg.ClientUpdate != nil {
				go func() {
					err := asyncCacheClientUpdate(ctx, cfg, cfg.Chains[msg.ClientUpdate.ChainId], Query{ConnectionId: msg.ClientUpdate.ConnectionId, Height: msg.ClientUpdate.Height}, msg.ClientUpdate.Height, logger, metrics)
					if err != nil {
						_ = logger.Log("msg", fmt.Sprintf("Error: Could not submit client update %s", err))
					}
				}()
			}

		case <-time.After(WaitInterval):
			flush(cfg, toSend, logger, metrics)
			metrics.SendQueue.WithLabelValues("send-queue").Set(float64(len(sendQueue)))
			toSend = []Message{}
		}
	}
}

// TODO: refactor me!
func flush(cfg *types.Config, toSend []Message, logger log.Logger, metrics prommetrics.Metrics) {
	fmt.Println("flush messages", len(toSend))
	if len(toSend) > 0 {
		_ = logger.Log("msg", fmt.Sprintf("Sending batch of %d messages", len(toSend)))
		msgs := prepareMessages(toSend, logger)
		if len(msgs) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()
			hash, code, err := cfg.DefaultChain.SignAndBroadcastMsg(ctx, cfg.ClientContext, msgs, VERSION)

			switch {
			case err == nil:
				_ = logger.Log("msg", fmt.Sprintf("Sent batch of %d (deduplicated) messages [hash: %s]", len(msgs), hash))
			case code == 12:
				_ = logger.Log("msg", "Not enough gas")
			case code == 19:
				_ = logger.Log("msg", "Tx already in mempool")
			case strings.Contains(err.Error(), "request body too large"):
				TxMsgs = max(1, TxMsgs*3/4)
				LastReduced = time.Now()
				_ = logger.Log("msg", "body too large: reduced batchsize", "size", TxMsgs)
			case strings.Contains(err.Error(), "failed to execute message"):
				regex := regexp.MustCompile(`failed to execute message; message index: (\d+)`)
				match := regex.FindStringSubmatch(err.Error())
				idx, _ := strconv.Atoi(match[1])
				badMsg, ok := msgs[idx].(*qstypes.MsgSubmitQueryResponse)
				if !ok {
					_ = logger.Log("msg", "Failed to execute non QueryResponse (probably submitting a ClientUpdate to a syncing node)", "index", match[1], "err", err.Error())
					metrics.FailedTxs.WithLabelValues("failed_txs").Inc()
					break
				}
				cache.SetWithTTL("ignore/"+badMsg.QueryId, true, 1, time.Minute*5)
				_ = logger.Log("msg", "Failed to execute message; ignoring for five minutes", "index", match[1], "err", err.Error())
			case code == 65536:
				_ = logger.Log("msg", "error in tx", "err", err.Error())
			default:
				_ = logger.Log("msg", "Failed to submit; nevermind, we'll try again!", "err", err)
				metrics.FailedTxs.WithLabelValues("failed_txs").Inc()
			}

		}
	}
}

func prepareMessages(msgSlice []Message, logger log.Logger) []sdk.Msg {
	keys := make(map[string]bool)

	list := []sdk.Msg{}

	exists := func(keys map[string]bool, key string) bool {
		_, ok := keys[key]
		return ok
	}

	for _, entry := range msgSlice {
		//fmt.Println("prepareMessages", idx)
		if len(list) > TxMsgs {
			//fmt.Println("transaction full; requeueing")
			go func(entry Message) { time.Sleep(time.Second * 2); sendQueue <- entry }(entry) // client update not ready; requeue.
			continue
		}

		msg, ok := entry.Msg.(*qstypes.MsgSubmitQueryResponse)
		if !ok {
			fmt.Println("unable to cast message to MsgSubmitQueryResponse")
			continue // unable to cast message to MsgSubmitQueryResponse
		}

		if _, ok := cache.Get("ignore/" + msg.QueryId); ok {
			logger.Log("msg", "Query already in ignore cache", "id", msg.QueryId)
			continue
		}

		if _, ok := keys[msg.QueryId]; ok {
			//fmt.Println("message already added")
			continue // message already added
		}

		if entry.ClientUpdate != nil && !exists(keys, fmt.Sprintf("%s-%d", entry.ClientUpdate.ConnectionId, entry.ClientUpdate.Height)) {
			//fmt.Println("client update required")
			cu, err := getCachedClientUpdate(entry.ClientUpdate.ConnectionId, entry.ClientUpdate.Height)
			if err != nil {
				//fmt.Println("client update not ready; requeueing")
				go func(entry Message) { time.Sleep(time.Second * 2); sendQueue <- entry }(entry) // client update not ready; requeue.
			} else {
				fmt.Println("client update ready; adding update and query response to send list")
				list = append(list, cu)
				list = append(list, entry.Msg)
				keys[msg.QueryId] = true
				keys[fmt.Sprintf("%s-%d", entry.ClientUpdate.ConnectionId, entry.ClientUpdate.Height)] = true
				cache.Del("query/" + msg.QueryId)
			}
		} else {
			fmt.Println("adding query response to send list")
			list = append(list, entry.Msg)
			keys[msg.QueryId] = true
			cache.Del("query/" + msg.QueryId)
		}

	}

	fmt.Printf("prepared %d messages\n", len(list))

	return list
}

func Close(cfg *types.Config) error {
	query := tmquery.MustParse(fmt.Sprintf("message.module='%s'", "interchainquery"))

	err := cfg.DefaultChain.GetClient().Unsubscribe(ctx, cfg.DefaultChain.ChainID+"-icq", query.String())
	if err != nil {
		return err
	}

	return cfg.DefaultChain.GetClient().Stop()
}
