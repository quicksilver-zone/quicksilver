package runner

import (
	"context"
	"encoding/hex"
	"fmt"
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
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	tmclient "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/types"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/prommetrics"
	qstypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	VERSION = "icq-relayer/v1.0.0-alpha.0"
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
	LastReduced           time.Time

	// Variables used for retries
	RtyAttNum = uint(5)
	RtyAtt    = retry.Attempts(RtyAttNum)
	RtyDel    = retry.Delay(time.Millisecond * 800)
	RtyErr    = retry.LastErrorOnly(true)
)

func Run(ctx context.Context, cfg *types.Config, errHandler func(error), cmd *cobra.Command) error {
	MaxTxMsgs = cfg.MaxMsgsPerTx
	if MaxTxMsgs == 0 {
		MaxTxMsgs = 40
	}
	TxMsgs = MaxTxMsgs

	// Configure zerolog to log to os.Stderr
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Add global context fields
	logger = logger.With().Timestamp().Caller().Logger()

	// Log messages
	logger.Info().
		Str("worker", "init").
		Str("msg", "starting icq relayer").
		Str("version", VERSION).
		Msg("")

	logger.Info().
		Str("worker", "init").
		Str("msg", "permitted queries").
		Str("queries", strings.Join(cfg.AllowedQueries, ",")).
		Msg("")

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
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.BindPort), nil)
		logger.Fatal().Err(err).Msg("HTTP metrics failed to start")
	}()

	defer func() {
		err := Close(cfg)
		if err != nil {
			logger.Error().Str("worker", "init").Str("msg", "error in closing the routine").Msg("")
		}
	}()

	if err := cfg.DefaultChain.Init(cfg.ProtoCodec, cache); err != nil {
		logger.Error().Err(err).Msg("Failed to initialize default chain")
		return err
	}

	for _, c := range cfg.Chains {
		if err := c.Init(cfg.ProtoCodec, cache); err != nil {
			return err
		}

		logger.Info().
			Str("worker", "init").
			Str("msg", "configured chain").
			Str("chain", c.ChainID).
			Msg("")
	}

	query := tmquery.MustParse(fmt.Sprintf("message.module='%s'", "interchainquery"))

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	logger.Info().
		Str("worker", "init").
		Str("msg", "configuring subscription on default chainClient").
		Str("chain", cfg.DefaultChain.ChainID).
		Msg("")

	ch, err := cfg.DefaultChain.GetClient().Subscribe(ctx, cfg.DefaultChain.ChainID+"-icq", query.String())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to subscribe to default chain client")
		return err
	}
	wg.Add(1)
	go func(chainId string, ch <-chan coretypes.ResultEvent) {
		defer wg.Done()
		for v := range ch {
			v.Events["source"] = []string{chainId}
			// why does this always trigger twice? messages are deduped later, but this causes 2x queries to trigger.
			time.Sleep(75 * time.Millisecond) // try to avoid thundering herd.
			go handleEvent(cfg, v, logger.With().Str("worker", "chainClient").Str("chain", cfg.DefaultChain.ChainID).Logger(), metrics, cmd)
		}
	}(cfg.DefaultChain.ChainID, ch)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := FlushSendQueue(cfg, logger.With().Str("worker", "flusher").Str("chain", cfg.DefaultChain.ChainID).Logger(), metrics, cmd)
		if err != nil {
			logger.Error().Msg("Flush Go-routine Bailing")
			panic(err)
		}
	}()

	for _, chainCfg := range cfg.Chains {
		wg.Add(1)
		go func(defaultClient *types.ChainConfig, srcClient *types.ReadOnlyChainConfig, logger zerolog.Logger) {
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
					logger.Error().Err(err).Str("msg", "Error: Unable to unmarshal").Msg("")
					if err != nil {
						return
					}
					continue CNT
				}
				logger.Info().
					Str("worker", "chainClient").
					Str("msg", "fetched historic queries for chain").
					Int("count", len(out.Queries)).
					Msg("")

				if len(out.Queries) > 0 {
					go handleHistoricRequests(cfg, srcClient, out.Queries,
						cfg.DefaultChain.ChainID, logger.With().Str("worker",
							"historic").Logger(), metrics, cmd)
				}
			}
		}(cfg.DefaultChain, chainCfg, logger.With().Str("chain",
			cfg.DefaultChain.ChainID).Str("src_chain", chainCfg.ChainID).Logger())
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

func handleHistoricRequests(cfg *types.Config,
	queryClient *types.ReadOnlyChainConfig, queries []qstypes.Query,
	sourceChainId string,
	logger zerolog.Logger,
	metrics prommetrics.Metrics,
	cmd *cobra.Command) {
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
			logger.Info().Str("msg", "Query already in cache").Str("id", q.QueryId).Msg("")
			continue
		}

		if _, found := cache.Get("ignore/" + q.QueryId); found {
			logger.Info().Str("msg", "Query already in ignore cache").Str("id", q.QueryId).Msg("")
			continue
		}

		var err error
		q.Height, err = queryClient.GetCurrentHeight(ctx, cache, logger.With().Str("chain", queryClient.ChainID).Logger())
		if err != nil {
			logger.Error().Str("msg", "Error getting current height").Err(err).Msg("")
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
			logger.Info().Str("msg", "Ignoring existing query; not a permitted type").Str("id", query.Id).Str("type", q.Type).Msg("")
			continue
		}
		//_ = logger.Log("msg", "Handling existing query", "id", query.Id)

		time.Sleep(75 * time.Millisecond) // try to avoid thundering herd.

		cache.Set("query/"+q.QueryId, true, 0)

		go doRequestWithMetrics(cfg, q, logger, metrics, cmd)
	}
}

func handleEvent(cfg *types.Config, event coretypes.ResultEvent,
	logger zerolog.Logger, metrics prommetrics.Metrics, cmd *cobra.Command) {
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
			logger.Error().Str("worker", "handler").Str("msg", err.Error()).Msg("")
			continue
		}
		h, err := strconv.ParseInt(height[i], 10, 64)
		if err != nil {
			logger.Error().Str("worker", "handler").Str("msg", err.Error()).Msg("")
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
			logger.Info().Str("worker", "handler").Str("msg", "Ignoring current query; not a permitted type").Str("id", queryIds[i]).Str("type", types[i]).Msg("")
			continue
		}
		if _, found := cache.Get("query/" + queryIds[i]); found {
			// break if this is in the cache
			continue
		}

		if _, found := cache.Get("ignore/" + queryIds[i]); found {
			logger.Info().Str("msg", "Query already in ignore cache").Str("id", queryIds[i]).Msg("")
			// break if this is in the cache
			continue
		}

		if h == 0 {
			currentheight, err := client.GetCurrentHeight(ctx, cache, logger.With().Str("chain", chains[i]).Logger())
			if err != nil {
				logger.Error().Str("worker", "handler").Str("msg", "error getting current block").Int64("height", currentheight).Msg("")
				continue
			}
		}
		cache.Set("query/"+queryIds[i], true, 0) // just long enough to not duplicate.
		queries = append(queries, Query{source[0], connections[i], chains[i], queryIds[i], types[i], h, req})
	}

	for _, q := range queries {
		go doRequestWithMetrics(cfg, q, logger.With().Str("src_chain",
			q.ChainId).Logger(), metrics, cmd)
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
			return nil, fmt.Errorf("unable to query light block, max interval exceeded")
		}
		cache.Set("lightblock/"+chain.ChainID+"/"+fmt.Sprintf("%d", height), lightBlock, 5)
	}
	return lightBlock.(*tmtypes.LightBlock), err
}

func doRequestWithMetrics(cfg *types.Config, query Query, logger zerolog.Logger,
	metrics prommetrics.Metrics, cmd *cobra.Command) {
	startTime := time.Now()
	metrics.Requests.WithLabelValues("requests", query.Type).Inc()
	doRequest(cfg, query, logger, metrics, cmd)
	endTime := time.Now()
	metrics.RequestsLatency.WithLabelValues("request-latency", query.Type).Observe(endTime.Sub(startTime).Seconds())
}

func doRequest(cfg *types.Config, query Query, logger zerolog.Logger,
	metrics prommetrics.Metrics, cmd *cobra.Command) {
	var err error
	client, ok := cfg.Chains[query.ChainId]
	if !ok {
		return
	}

	//_ = logger.Log("msg", "Handling request", "type", query.Type, "id", query.QueryId, "height", query.Height)

	pathParts := strings.Split(query.Type, "/")
	prove := pathParts[len(pathParts)-1] == "key" // fetch proof if the query is 'key'

	var res abcitypes.ResponseQuery
	clientCtx, err := sdkClient.GetClientTxContext(cmd)
	if err != nil {
		return
	}
	fromAddr := clientCtx.GetFromAddress()

	switch query.Type {
	// until we fix ordering and pagination in the binary, we can override the query here.
	case "cosmos.tx.v1beta1.Service/GetTxsEvent":
		request := qstypes.GetTxsEventRequest{}
		err = cfg.ProtoCodec.Unmarshal(query.Request, &request)
		if err != nil {
			logger.Error().Str("msg", "Error: Failed in Unmarshalling Request").Str("type", query.Type).Str("id", query.QueryId).Int64("height", query.Height).Msg("")
			return
		}
		request.OrderBy = txtypes.OrderBy_ORDER_BY_DESC
		request.Limit = 200
		request.Pagination.Limit = 200
		request.Query = request.Events[0]

		query.Request, err = cfg.ProtoCodec.Marshal(&request)
		if err != nil {
			logger.Error().Str("msg", "Error: Failed in Marshalling Request").Str("type", query.Type).Str("id", query.QueryId).Int64("height", query.Height).Msg("")
			return
		}

		logger.Info().Str("msg", "Handling GetTxsEvents").Str("id", query.QueryId).Int64("height", query.Height).Msg("")
		res, err = client.RunABCIQuery(ctx, "/"+query.Type, query.Request, query.Height, prove, metrics)
		if err != nil {
			logger.Error().Str("msg", "Error: Failed in RunGRPCQuery").Str("type", query.Type).Str("id", query.QueryId).Int64("height", query.Height).Msg("")
			return
		}

	case "tendermint.Tx":
		// custom request type for fetching a tx with proof.
		req := txtypes.GetTxRequest{}
		cfg.ProtoCodec.MustUnmarshal(query.Request, &req)
		hashBytes, err := hex.DecodeString(req.GetHash())
		if err != nil {
			logger.Error().Str("msg", fmt.Sprintf("Error: Could not decode hash %s", err)).Msg("")
			return
		}
		txRes, height, err := client.Tx(hashBytes)
		if err != nil {
			logger.Error().Str("msg", fmt.Sprintf("Error: Could not fetch proof %s", err)).Msg("")
			return
		}

		protoProof := txRes.ToProto()

		clientId, err := cfg.DefaultChain.GetClientId(ctx, query.ConnectionId, logger, metrics)
		if err != nil {
			logger.Error().Str("msg", fmt.Sprintf("Error: Could not get client id %s", err)).Msg("")
			return
		}

		header, err := getHeader(ctx, cfg, client, clientId, height-1, logger, true, metrics)
		if err != nil {
			logger.Error().Str("msg", fmt.Sprintf("Error: Could not get header %s", err)).Msg("")
			return
		}

		resp := qstypes.GetTxWithProofResponse{Proof: &protoProof, Header: header}
		res.Value = cfg.ProtoCodec.MustMarshal(&resp)

	case "ibc.ClientUpdate":
		// return a dummy message to settle the query.
		msg := &qstypes.MsgSubmitQueryResponse{
			ChainId:     query.ChainId,
			QueryId:     query.QueryId,
			Result:      []byte{},
			Height:      int64(sdk.BigEndianToUint64(query.Request)),
			ProofOps:    &crypto.ProofOps{},
			FromAddress: fromAddr.String(),
		}
		sendQueue <- Message{Msg: msg, ClientUpdate: &ClientUpdateRequirement{ConnectionId: query.ConnectionId, ChainId: query.ChainId, Height: int64(sdk.BigEndianToUint64(query.Request))}}
		return
	default:
		res, err = client.RunABCIQuery(ctx, "/"+query.Type, query.Request, query.Height, prove, metrics)
		if err != nil {
			logger.Error().Str("msg", "Error: Failed in RunGRPCQuery").Str("type", query.Type).Str("id", query.QueryId).Int64("height", query.Height).Msg("")
			return
		}
	}

	msg := &qstypes.MsgSubmitQueryResponse{
		ChainId:     query.ChainId,
		QueryId:     query.QueryId,
		Result:      res.Value,
		Height:      res.Height,
		ProofOps:    res.ProofOps,
		FromAddress: fromAddr.String()}
	var clientUpdate *ClientUpdateRequirement

	if prove {
		clientUpdate = &ClientUpdateRequirement{ConnectionId: query.ConnectionId, ChainId: query.ChainId, Height: res.Height}
	}

	sendQueue <- Message{Msg: msg, ClientUpdate: clientUpdate}
}

func asyncCacheClientUpdate(
	ctx context.Context,
	cfg *types.Config,
	client *types.ReadOnlyChainConfig,
	query Query, height int64,
	logger zerolog.Logger,
	metrics prommetrics.Metrics,
	cmd *cobra.Command) error {
	cacheKey := fmt.Sprintf("cu/%s-%d", query.ConnectionId, height)
	queryKey := fmt.Sprintf("cuquery/%s-%d", query.ConnectionId, height)
	clientCtx, err := sdkClient.GetClientTxContext(cmd)
	if err != nil {
		return err
	}
	fromAddr := clientCtx.GetFromAddress()
	_, ok := cache.Get("cu/" + cacheKey)
	if ok {
		log.Info().Msgf("cache found for %s", cacheKey)
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
			logger.Error().Str("msg", "Error: Could not get clientId").Err(err).Msg("")
			return err
		}
		header, err := getHeader(ctx, cfg, client, clientId, height, logger, false, metrics)
		if err != nil {
			logger.Error().Str("msg", "Error: Could not get header").Err(err).Msg("")
			return err
		}

		msg, err := clienttypes.NewMsgUpdateClient(clientId, header, fromAddr.String())
		if err != nil {
			logger.Error().Str("msg", "Error: Could not create msg update").Err(err).Msg("")
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
	return nil, fmt.Errorf("client update not found")
}

func getHeader(ctx context.Context, cfg *types.Config, client *types.ReadOnlyChainConfig, clientId string, requestHeight int64, logger zerolog.Logger, historicOk bool, metrics prommetrics.Metrics) (*tmclient.Header, error) {
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
		return nil, fmt.Errorf("error: Could coerce trusted height")
	}

	if !historicOk && clientHeight.RevisionHeight >= uint64(requestHeight+1) {
		//return nil, fmt.Errorf("trusted height >= request height")
		oldHeights, err := cfg.DefaultChain.GetClientStateHeights(ctx, clientId, client.ChainID, uint64(requestHeight-200), logger, metrics, 0)
		if err != nil {
			return nil, fmt.Errorf("error: Could not get old heights: %w", err)
		}
		clientHeight = oldHeights[0]
	}

	//_ = logger.Log("msg", "Fetching client update for height", "height", requestHeight+1)
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

func FlushSendQueue(cfg *types.Config, logger zerolog.Logger,
	metrics prommetrics.Metrics, cmd *cobra.Command) error {
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
			logger.Info().Str("msg", "increased batch size").Int("size", TxMsgs).Msg("")
			LastReduced = time.Now()
		}

		if len(toSend) > 5*TxMsgs {
			flush(cfg, toSend, logger, metrics, cmd)
			toSend = []Message{}
		}
		select {
		case msg := <-ch:
			toSend = append(toSend, msg)
			metrics.SendQueue.WithLabelValues("send-queue").Set(float64(len(sendQueue)))
			if msg.ClientUpdate != nil {
				go func() {
					err := asyncCacheClientUpdate(ctx, cfg,
						cfg.Chains[msg.ClientUpdate.ChainId],
						Query{ConnectionId: msg.ClientUpdate.ConnectionId,
							Height: msg.ClientUpdate.Height},
						msg.ClientUpdate.Height, logger, metrics, cmd)
					if err != nil {
						logger.Error().Str("msg", "Error: Could not submit client update").Err(err).Msg("")
					}
				}()
			}

		case <-time.After(WaitInterval):
			flush(cfg, toSend, logger, metrics, cmd)
			metrics.SendQueue.WithLabelValues("send-queue").Set(float64(len(sendQueue)))
			toSend = []Message{}
		}
	}
}

// TODO: refactor me!
func flush(cfg *types.Config, toSend []Message, logger zerolog.Logger,
	metrics prommetrics.Metrics, cmd *cobra.Command) {
	logger.Info().Msgf("Flush messages: %d", len(toSend))
	if len(toSend) > 0 {
		logger.Info().Msgf("Sending batch of %d messages", len(toSend))
		msgs := prepareMessages(toSend, logger)
		if len(msgs) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()
			code, err := cfg.DefaultChain.SignAndBroadcastMsgWithKey(ctx, cfg.ClientContext, msgs, VERSION, cmd)

			switch {
			case err == nil:
				logger.Info().Msgf("Sent batch of %d", len(msgs))
			case code == 12:
				logger.Warn().Msg("Not enough gas")
			case code == 19:
				logger.Warn().Msg("Tx already in mempool")
			case strings.Contains(err.Error(), "request body too large"):
				TxMsgs = TxMsgs / 4 * 3
				LastReduced = time.Now()
				logger.Warn().Msgf("Body too large: reduced batch size to %d", TxMsgs)
			case strings.Contains(err.Error(), "failed to execute message"):
				regex := regexp.MustCompile(`failed to execute message; message index: (\d+)`)
				match := regex.FindStringSubmatch(err.Error())
				idx, _ := strconv.Atoi(match[1])
				badMsg := msgs[idx].(*qstypes.MsgSubmitQueryResponse)
				cache.SetWithTTL("ignore/"+badMsg.QueryId, true, 1, time.Minute*5)
				logger.Error().Msgf("Failed to execute message; ignoring for five minutes (index: %d, error: %s)", idx, err.Error())
			case code == 65536:
				logger.Error().Err(err).Msg("Error in tx")
			default:
				logger.Error().Err(err).Msg("Failed to submit; we'll try again!")
				metrics.FailedTxs.WithLabelValues("failed_txs").Inc()
			}

		}
	}
}

func prepareMessages(msgSlice []Message, logger zerolog.Logger) []sdk.Msg {
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
			log.Info().Msg("unable to cast message to MsgSubmitQueryResponse")
			continue // unable to cast message to MsgSubmitQueryResponse
		}

		if _, ok := cache.Get("ignore/" + msg.QueryId); ok {
			logger.Info().Str("msg", "Query already in ignore cache").Str("id", msg.QueryId).Msg("")
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
				log.Info().Msg("client update ready; adding update and query response to send list")
				list = append(list, cu)
				list = append(list, entry.Msg)
				keys[msg.QueryId] = true
				keys[fmt.Sprintf("%s-%d", entry.ClientUpdate.ConnectionId, entry.ClientUpdate.Height)] = true
				cache.Del("query/" + msg.QueryId)
			}
		} else {
			log.Info().Msg("adding query response to send list")
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
