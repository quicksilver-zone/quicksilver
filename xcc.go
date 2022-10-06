package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	osmolockup "github.com/ingenuity-build/quicksilver/osmosis-types/lockup"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

var connectionManager CacheManager[prewards.ConnectionProtocolData]
var poolsManager CacheManager[prewards.OsmosisPoolProtocolData]

type Response struct {
	Messages []prewards.MsgSubmitClaim `json:"messages"`
}

type Data[T any] struct {
	Data []T
}

type CacheManager[T any] struct {
	url         string
	cache       []T
	lastUpdated time.Time
}

func (c *CacheManager[T]) Init(url string, updateTime time.Duration) {
	c.url = url
	c.Fetch()
}

func (c *CacheManager[T]) Fetch() {
	response, err := http.Get(c.url)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	fmt.Println(responseData)
	if err != nil {
		panic(err)
	}

	data := Data[T]{}

	err = json.Unmarshal(responseData, &data)
	if err != nil {
		panic(err)
	}
	c.cache = data.Data
}

func (c CacheManager[T]) Get() []T {
	if time.Now().After(c.lastUpdated.Add(time.Minute * 5)) {
		c.Fetch()
	}
	return c.cache
}

func handle(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var err error

	response := Response{}
	var height int64
	connections := connectionManager.Get()
	for _, con := range connections {
		if con.ChainID == "osmosis-1" {
			height = con.LastEpoch
			break
		}
		if con.ChainID == "quickosmo-1" {
			height = con.LastEpoch
			break
		}
		if con.ChainID == "osmo-test-4" {
			height = con.LastEpoch
			break
		}
	}
	if height == 0 {
		return
	}
	// liquid for all zones; config should hold osmosis chainid.

	// query zones and last heights

	response.Messages, err = handleOsmo(vars["address"], height)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	jsonOut, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprint(w, string(jsonOut))

}

type poolMap map[string][]int64

func handleOsmo(address string, height int64) ([]prewards.MsgSubmitClaim, error) {
	_, addrBytes, err := bech32.Decode(address, 51)
	if err != nil {
		return nil, err
	}
	osmoAddress, err := bech32.Encode("osmo", addrBytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("b")

	client, err := NewRPCClient("http://dev:23101", 30*time.Second)
	if err != nil {
		return nil, err
	}
	// fetch timestamp of block

	query := osmolockup.AccountLockedPastTimeRequest{Owner: osmoAddress, Timestamp: time.Unix(timestamp, 0)}
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	osmolockup.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	bytes := marshaler.MustMarshal(&query)
	// how to pass header for height!
	abciquery, err := client.ABCIQuery(context.Background(), "/osmosis.lockup.Query/AccountLockedPastTime", bytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("done 1")
	queryResponse := osmolockup.AccountLockedPastTimeResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &queryResponse)
	if err != nil {
		return nil, err
	}

	poolIds := poolMap{}
	// filter by pool id - query this from Quicksilver (and cache hourly)
	for _, pool := range poolsManager.Get() {
		for chain := range pool.Zones {
			if _, ok := poolIds[chain]; !ok {
				poolIds[chain] = make([]int64, 0)
			}
			poolIds[chain] = append(poolIds[chain], int64(pool.PoolID))
		}
	}

	fmt.Println(poolIds)

	msg := map[string]prewards.MsgSubmitClaim{}

OUTER:
	for _, lockup := range queryResponse.Locks {
		fmt.Println("a")
		for chainID, chainPools := range poolIds {
			fmt.Println("b")

			for _, p := range chainPools {
				fmt.Println("c")

				if fmt.Sprintf("gamm/pool/%d", p) == lockup.Coins.GetDenomByIndex(0) {
					fmt.Println("matched")
					if _, ok := msg[chainID]; !ok {
						msg[chainID] = prewards.MsgSubmitClaim{
							UserAddress: address,
							Zone:        chainID,
							SrcZone:     "quickosmo-1",
							ClaimType:   prewards.ClaimTypeOsmosisPool,
							Proofs:      make([]*prewards.Proof, 0),
						}
					}

					abciquery, err := client.ABCIQueryWithOptions(
						context.Background(), "/store/lockup/key",
						append(osmolockup.KeyPrefixPeriodLock, append(osmolockup.KeyIndexSeparator, sdk.Uint64ToBigEndian(lockup.ID)...)...),
						rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
					)
					if err != nil {
						return nil, err
					}
					lockupResponse := osmolockup.PeriodLock{}
					err = marshaler.Unmarshal(abciquery.Response.Value, &lockupResponse)
					if err != nil {
						return nil, err
					}
					chainMsg := msg[chainID]

					proof := prewards.Proof{
						Data:     abciquery.Response.Value,
						Key:      abciquery.Response.Key,
						ProofOps: abciquery.Response.ProofOps,
						Height:   abciquery.Response.Height,
					}

					chainMsg.Proofs = append(chainMsg.Proofs, &proof)

					msg[chainID] = chainMsg
					continue OUTER

				}

			}
		}
	}

	out := []prewards.MsgSubmitClaim{}
	for _, m := range msg {
		out = append(out, m)
	}

	return out, nil

}

func main() {
	connectionManager.Init("https://lcd.dev.quicksilver.zone/quicksilver/participationrewards/v1/protocoldata/connection", time.Minute*5)
	poolsManager.Init("https://lcd.dev.quicksilver.zone/quicksilver/participationrewards/v1/protocoldata/osmosispools", time.Minute*5)
	r := mux.NewRouter()
	r.HandleFunc("/{address}/{epoch}", handle)
	http.Handle("/", r)
	http.ListenAndServe(":8090", nil)
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
