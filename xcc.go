package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	osmolockup "github.com/osmosis-labs/osmosis/v10/x/lockup/types"
	tmcrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type Response struct {
	Messages []prewards.MsgSubmitClaim `json:"messages"`
}

func handle(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var err error
	response := Response{}
	ts, err := strconv.ParseInt(vars["timestamp"], 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	// query epoch time

	// query zones and last heights

	response.Messages, err = handleOsmo(vars["address"], ts)
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

func handleOsmo(address string, timestamp int64) ([]prewards.MsgSubmitClaim, error) {
	_, addrBytes, err := bech32.Decode(address, 51)
	if err != nil {
		return nil, err
	}
	osmoAddress, err := bech32.Encode("osmo", addrBytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("b")

	client, err := NewRPCClient("https://rpc-osmosis.whispernode.com:443", 30*time.Second)
	if err != nil {
		return nil, err
	}
	query := osmolockup.AccountLockedPastTimeRequest{Owner: osmoAddress, Timestamp: time.Unix(timestamp, 0)}
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	osmolockup.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	bytes := marshaler.MustMarshal(&query)

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

	// filter by pool id - query this from Quicksilver (and cache hourly)
	poolIds := poolMap{"osmosis-1": []int64{3, 625, 465}, "juno-1": []int64{2, 5, 7}}

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
							ProofType:   0,
							Key:         make([][]byte, 0),
							Data:        make([][]byte, 0),
							ProofOps:    make([]*tmcrypto.ProofOps, 0),
							Height:      abciquery.Response.Height,
						}
					}

					abciquery, err := client.ABCIQueryWithOptions(
						context.Background(), "/store/lockup/key",
						append(osmolockup.KeyPrefixPeriodLock, append(osmolockup.KeyIndexSeparator, sdk.Uint64ToBigEndian(lockup.ID)...)...),
						rpcclient.ABCIQueryOptions{Height: 0, Prove: true},
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

					chainMsg.Data = append(chainMsg.Data, abciquery.Response.Value)
					chainMsg.ProofOps = append(chainMsg.ProofOps, abciquery.Response.ProofOps)
					chainMsg.Key = append(chainMsg.Key, abciquery.Response.Key)
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
	r := mux.NewRouter()
	r.HandleFunc("/{address}/{timestamp}", handle)
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
