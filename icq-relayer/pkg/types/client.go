package types

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"cosmossdk.io/math"

	"github.com/avast/retry-go/v4"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/dgraph-io/ristretto"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/prommetrics"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	cmtjson "github.com/tendermint/tendermint/libs/json"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	jsonrpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdkcryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	connectiontypes "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	log2 "github.com/go-kit/log"
	home "github.com/mitchellh/go-homedir"
	prov "github.com/tendermint/tendermint/light/provider"
	lighthttp "github.com/tendermint/tendermint/light/provider/http"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type ReadOnlyChainConfig struct {
	ChainID                     string
	RpcUrl                      string
	ConnectTimeoutSeconds       int
	QueryTimeoutSeconds         int
	QueryRetries                int
	QueryRetryDelayMilliseconds int
	Client                      RPCClientI        `toml:"-"`
	LightProvider               prov.Provider     `toml:"-"`
	Codec                       *codec.ProtoCodec `toml:"-"`
	Cache                       *ristretto.Cache  `toml:"-"`
}

type ChainConfig struct {
	*ReadOnlyChainConfig
	MnemonicPath           string
	Prefix                 string
	TxSubmitTimeoutSeconds int
	GasLimit               int
	GasPrice               string
	GasMultiplier          float64
	AddressBytes           sdktypes.AccAddress `toml:"-"`
}

func (r *ChainConfig) GetClient() *rpchttp.HTTP {
	return r.Client.(*rpchttp.HTTP)
}

func (r *ChainConfig) Init(codec *codec.ProtoCodec, cache *ristretto.Cache) error {
	var err error
	r.Client, err = rpchttp.NewWithTimeout(r.RpcUrl, "/websocket", uint(r.ConnectTimeoutSeconds))
	if err != nil {
		return err
	}

	err = r.GetClient().Start()
	if err != nil {
		return err
	}

	r.LightProvider, err = lighthttp.New(r.ChainID, r.RpcUrl)
	if err != nil {
		return err
	}

	r.Codec = codec
	r.Cache = cache

	return nil
}

func (r *ReadOnlyChainConfig) Init(codec *codec.ProtoCodec, cache *ristretto.Cache) error {
	var err error
	r.Client, err = NewWithTimeout(r.RpcUrl, uint(r.ConnectTimeoutSeconds))
	if err != nil {
		return err
	}
	// err = r.Client.Start()
	// if err != nil {
	// 	return err
	// }

	r.LightProvider, err = lighthttp.New(r.ChainID, r.RpcUrl)
	if err != nil {
		return err
	}

	r.Codec = codec
	r.Cache = cache

	return nil
}

func (r *ReadOnlyChainConfig) LightBlock(ctx context.Context, height int64) (*tmtypes.LightBlock, error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	lightBlock, err := r.LightProvider.LightBlock(ctx, height)
	if err != nil {
		fmt.Println("error getting light block", err)
		return nil, err
	}
	return lightBlock, nil
}

func (r *ReadOnlyChainConfig) GetClientState(ctx context.Context, clientId string, logger log2.Logger, metrics prommetrics.Metrics) (*clienttypes.QueryClientStateResponse, error) {
	clientStateQuery := clienttypes.QueryClientStateRequest{ClientId: clientId}
	bz := r.Codec.MustMarshal(&clientStateQuery)
	res, err := r.Client.ABCIQuery(ctx, "/ibc.core.client.v1.Query/ClientState", bz)
	if err != nil {
		_ = logger.Log("msg", fmt.Sprintf("Error: Could not get client state from client %s", err))
		return nil, err
	}
	clientStateResponse := clienttypes.QueryClientStateResponse{}
	err = r.Codec.Unmarshal(res.Response.Value, &clientStateResponse)
	if err != nil {
		_ = logger.Log("msg", fmt.Sprintf("Error: Could not unmarshal connection %s", err))
		return nil, err
	}

	return &clientStateResponse, nil
}

func (r *ReadOnlyChainConfig) GetClientStateHeights(ctx context.Context, clientId string, chainId string, height uint64, logger log2.Logger, metrics prommetrics.Metrics, depth int) ([]clienttypes.Height, error) {
	time.Sleep(100 * time.Millisecond)
	if depth > 500 {
		return nil, fmt.Errorf("reached max depth")
	}
	chainParts := strings.Split(chainId, "-")
	if len(chainParts) == 1 {
		chainParts = append(chainParts, "1")
	}
	key := fmt.Sprintf("%s-%d", chainParts[len(chainParts)-1], height)

	req := clienttypes.QueryConsensusStateHeightsRequest{ClientId: clientId, Pagination: &querytypes.PageRequest{Key: []byte(key)}}
	bz := r.Codec.MustMarshal(&req)
	res, err := r.Client.ABCIQuery(ctx, "/ibc.core.client.v1.Query/ConsensusStateHeights", bz)
	if err != nil {
		_ = logger.Log("msg", fmt.Sprintf("Error: Could not get consensus state heights from client %s", err))
		return nil, err
	}
	resp := clienttypes.QueryConsensusStateHeightsResponse{}
	err = r.Codec.Unmarshal(res.Response.Value, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.ConsensusStateHeights) == 0 {
		return r.GetClientStateHeights(ctx, clientId, chainId, height-500, logger, metrics, depth+1)
	}

	return resp.ConsensusStateHeights, nil
}

func (r *ReadOnlyChainConfig) GetClientId(ctx context.Context, connectionId string, logger log2.Logger, metrics prommetrics.Metrics) (string, error) {
	clientId, found := r.Cache.Get("clientId/" + connectionId)
	if !found {
		connectionQuery := connectiontypes.QueryConnectionRequest{ConnectionId: connectionId}
		bz := r.Codec.MustMarshal(&connectionQuery)
		res, err := r.Client.ABCIQuery(ctx, "/ibc.core.connection.v1.Query/Connection", bz)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not get connection from chain %s", err))
			return "", err
		}
		connectionResponse := connectiontypes.QueryConnectionResponse{}
		err = r.Codec.Unmarshal(res.Response.Value, &connectionResponse)
		if err != nil {
			_ = logger.Log("msg", fmt.Sprintf("Error: Could not unmarshal connection %s", err))
			return "", err
		}

		clientId = connectionResponse.Connection.ClientId
		r.Cache.Set("clientId/"+connectionId, clientId, 1)
	}
	return clientId.(string), nil
}

// tm0.37 has a breaking change whereby tx events are no longer base64 encoded, so are represented as string and not bytes.
// As a result, we cannot use the RPCClient.Tx() method which attempts to unmarshal the Result, including the underlying Tx object.
// As such, we want to query the result directly, and unmarshal the json ourselves, to a representation of the result that conveniently
// does not contain the Tx object (that we don't use, because the TxProof already contains a byte representation of tx anyway!)
// Note: this function is compatible with 0.34 and 0.37 representations of transactions.
func (r *ReadOnlyChainConfig) Tx(hash []byte) (tmtypes.TxProof, int64, error) {
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
	httpRequest, err := http.NewRequestWithContext(context.Background(), http.MethodPost, r.RpcUrl, requestBuf)
	if err != nil {
		return tmtypes.TxProof{}, 0, fmt.Errorf("request failed: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	httpClient, err := jsonrpcclient.DefaultHTTPClient(r.RpcUrl)
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

// define celestia based TxResultMinimal
type TxResultMinimalCelestia struct {
	Height string          `json:"height"`
	Proof  tmtypes.TxProof `json:"proof"`
}

type QueryClient interface {
	Init(codec *codec.ProtoCodec, cache *ristretto.Cache) error
	GetClientState(ctx context.Context, clientId string, logger log2.Logger, metrics prommetrics.Metrics) (*clienttypes.QueryClientStateResponse, error)
	GetClientStateHeights(ctx context.Context, clientId string, chainId string, height uint64, logger log2.Logger, metrics prommetrics.Metrics, depth int) ([]clienttypes.Height, error)
	GetClientId(ctx context.Context, connectionId string, logger log2.Logger, metrics prommetrics.Metrics) (string, error)
}

type TxClient interface {
	QueryClient
	SignAndBroadcastMsg(ctx context.Context, cliContext *client.Context, exec []sdktypes.Msg, memo string) (string, uint32, error)
}

var _ TxClient = &ChainConfig{}
var _ QueryClient = &ChainConfig{}

func SetSDKConfigPrefix(prefix string) {
	configuration := sdktypes.GetConfig()
	configuration.SetBech32PrefixForAccount(prefix, prefix+sdktypes.PrefixPublic)
	configuration.SetBech32PrefixForValidator(prefix, prefix+sdktypes.PrefixValidator+sdktypes.PrefixOperator)
	configuration.SetBech32PrefixForConsensusNode(prefix+sdktypes.PrefixValidator+sdktypes.PrefixConsensus, prefix+sdktypes.PrefixValidator+sdktypes.PrefixConsensus+sdktypes.PrefixPublic)
}

func (c *ChainConfig) GetAddress() string {
	addrString, err := Bech32ifyAddressBytes(c.Prefix, c.GetAddressBytes())
	if err != nil {
		panic(err)
	}
	return addrString
}

func (c *ChainConfig) GetAddressBytes() sdktypes.AccAddress {
	if c.AddressBytes == nil {
		_, addrBytes := c.GetPrivKeyAddress()
		c.AddressBytes = addrBytes
	}
	return c.AddressBytes
}

func (c *ChainConfig) GetPrivKeyAddress() (sdkcryptotypes.PrivKey, sdktypes.AccAddress) {
	kb, err := keyring.New("icq-relayer", keyring.BackendMemory, "", nil, c.Codec)
	if err != nil {
		log.Fatal(err)
	}

	path, err := home.Expand(c.MnemonicPath)
	if err != nil {
		log.Fatal(err)
	}
	mnemonicBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	mnemonic := strings.Trim(string(mnemonicBytes), "\n")

	keyringAlgos, _ := kb.SupportedAlgorithms()

	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		log.Fatal(err)
	}

	hdPath := hd.CreateHDPath(118, 0, 0)

	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath.String())
	if err != nil {
		log.Fatal(err)
	}

	privKey := algo.Generate()(derivedPriv)

	addrBytes := sdktypes.AccAddress(privKey.PubKey().Address())
	return privKey, addrBytes
}

func Bech32ifyAddressBytes(prefix string, address sdktypes.AccAddress) (string, error) {
	if address.Empty() {
		return "", nil
	}
	if len(address.Bytes()) == 0 {
		return "", nil
	}
	if len(prefix) == 0 {
		return "", errors.New("prefix cannot be empty")
	}
	return bech32.ConvertAndEncode(prefix, address.Bytes())
}

func (r *ReadOnlyChainConfig) RunABCIQuery(ctx context.Context, method string, reqBz []byte, height int64, prove bool, metrics prommetrics.Metrics) (abcitypes.ResponseQuery, error) {
	metrics.ABCIRequests.WithLabelValues("abci_requests", method).Inc()
	// metrics: query duration?
	var abciRes abcitypes.ResponseQuery
	if err := retry.Do(func() error {

		opts := rpcclient.ABCIQueryOptions{
			Height: height,
			Prove:  prove,
		}
		result, err := r.Client.ABCIQueryWithOptions(ctx, method, reqBz, opts)
		if err != nil {
			return err
		}
		abciRes = result.Response

		return err
	}, retry.Attempts(uint(r.QueryRetries)), retry.Delay(time.Duration(r.QueryRetryDelayMilliseconds)*time.Millisecond), retry.LastErrorOnly(true)); err != nil {
		return abcitypes.ResponseQuery{}, err
	}

	return abciRes, nil
}

func (r *ReadOnlyChainConfig) GetCurrentHeight(ctx context.Context, cache *ristretto.Cache, logger log2.Logger) (int64, error) {
	currentheight, found := cache.Get("currentblock/" + r.ChainID)
	if !found {
		var err error
		var block *coretypes.ResultBlock
		if err = retry.Do(func() error {
			block, err = r.Client.Block(ctx, nil)
			return err
		}, retry.Attempts(uint(r.QueryRetries)), retry.Delay(time.Duration(r.QueryRetryDelayMilliseconds)*time.Millisecond), retry.LastErrorOnly(true)); err != nil {
			return 0, err
		}

		fmt.Printf("block: %+v\n", block)

		if block == nil {
			return 0, fmt.Errorf("block is nil")
		}

		if block.Block == nil {
			return 0, fmt.Errorf("block.Block is nil")
		}

		if block.Block.LastCommit == nil {
			return 0, fmt.Errorf("block last commit is nil")
		}

		currentheight = block.Block.LastCommit.Height - 1
		if cache == nil {
			logger.Log("msg", "unable to cache current block height, no cache provided", "height", currentheight)
			return currentheight.(int64), nil
		}
		cache.SetWithTTL("currentblock/"+r.ChainID, currentheight, 1, 6*time.Second)
		//logger.Log("msg", "caching currentblock", "height", currentheight)
		//} else {
		//logger.Log("msg", "using cached currentblock", "height", currentheight)
	}
	return currentheight.(int64), nil
}

func (c *ChainConfig) prepareBytes(txBuilder client.TxBuilder, cliContext client.Context, signerData xauthsigning.SignerData, gas uint64, codec *codec.ProtoCodec, memo string) ([]byte, error) {
	txBuilder.SetGasLimit(gas)
	fee, err := sdktypes.ParseDecCoin(c.GasPrice)
	if err != nil {
		return []byte{}, err
	}
	fee.Amount = fee.Amount.MulInt(math.NewIntFromUint64(gas))
	feeCoin := sdktypes.NewCoin(fee.Denom, fee.Amount.TruncateInt())
	txBuilder.SetFeeAmount(sdktypes.NewCoins(feeCoin))
	txBuilder.SetMemo(memo)

	privKey, _ := c.GetPrivKeyAddress()
	sigv2, err := tx.SignWithPrivKey(
		cliContext.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, privKey, cliContext.TxConfig, signerData.Sequence)
	if err != nil {
		return []byte{}, err
	}

	err = txBuilder.SetSignatures(sigv2)
	if err != nil {
		return []byte{}, err
	}
	return cliContext.TxConfig.TxEncoder()(txBuilder.GetTx())
}

func (c *ChainConfig) SignAndBroadcastMsg(ctx context.Context, cliContext *client.Context, exec []sdktypes.Msg, version string) (string, uint32, error) {
	// Build the factory CLI
	// Create a new TxBuilder.
	txBuilder := cliContext.TxConfig.NewTxBuilder()

	err := txBuilder.SetMsgs(exec...)
	if err != nil {
		return "", 65536, err
	}

	SetSDKConfigPrefix(c.Prefix)
	ac, seq, err := cliContext.AccountRetriever.GetAccountNumberSequence(*cliContext, c.GetAddressBytes())
	if err != nil {
		return "", 65536, err
	}

	signerData := xauthsigning.SignerData{
		ChainID:       c.ChainID,
		AccountNumber: ac,
		Sequence:      seq,
	}

	txBytes, err := c.prepareBytes(txBuilder, *cliContext, signerData, 1000000, c.Codec, version)
	if err != nil {
		return "", 65536, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.TxSubmitTimeoutSeconds)*time.Second)
	defer cancel()

	serviceClient := txtypes.NewServiceClient(cliContext)

	simRes, err := serviceClient.Simulate(ctx, &txtypes.SimulateRequest{
		TxBytes: txBytes,
	})
	if err != nil {
		return "", 65536, err
	}

	gas := uint64(float64(simRes.GasInfo.GasUsed) * c.GasMultiplier)

	txBytes, err = c.prepareBytes(txBuilder, *cliContext, signerData, gas, c.Codec, version)
	if err != nil {
		return "", 65536, err
	}

	res, err := serviceClient.BroadcastTx(ctx, &txtypes.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    txtypes.BroadcastMode_BROADCAST_MODE_BLOCK,
	})

	switch {
	case err != nil:
		//log.Err(err).Msg("Transaction error")
		return "", 65536, err
	case res.TxResponse.Code > 0:
		//log.Error().Msgf("Transaction failed: %v", res.TxResponse)
		return "err", res.TxResponse.Code, fmt.Errorf("transaction failed: %v", res.TxResponse)
	default:
		//log.Info().Msgf("Transaction broadcast successfully: %s", res.TxResponse.TxHash)
		return res.TxResponse.TxHash, res.TxResponse.Code, nil
	}
}
