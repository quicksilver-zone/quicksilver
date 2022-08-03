package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	osmolockup "github.com/osmosis-labs/osmosis/v10/x/lockup/types"
	tmcrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type Response struct {
	Osmosis *OsmosisResponse `json:"osmosis"`
}

type OsmosisResponse struct {
	Locks []LockWithProof `json:"locks"`
}

type LockWithProof struct {
	Lock   []byte            `json:"lock"`
	Key    []byte            `json:"key"`
	Proof  tmcrypto.ProofOps `json:"proof"`
	Height int64             `json:"height"`
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
	// validate ts
	response.Osmosis, err = handleOsmo(vars["address"], ts)
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

func handleOsmo(address string, timestamp int64) (*OsmosisResponse, error) {
	client, err := NewRPCClient("https://osmosis.chorus.one:443", 5*time.Second)
	if err != nil {
		return nil, err
	}
	query := osmolockup.AccountLockedPastTimeRequest{Owner: address, Timestamp: time.Unix(timestamp, 0)}
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	osmolockup.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	bytes := marshaler.MustMarshal(&query)

	abciquery, err := client.ABCIQuery(context.Background(), "/osmosis.lockup.Query/AccountLockedPastTime", bytes)
	if err != nil {
		return nil, err
	}
	queryResponse := osmolockup.AccountLockedPastTimeResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &queryResponse)
	if err != nil {
		return nil, err
	}

	// filter by pool id - query this from Quicksilver (and cache hourly)
	poolIds := []int64{3, 625, 465}

	o := &OsmosisResponse{}

OUTER:
	for _, lockup := range queryResponse.Locks {
		for _, p := range poolIds {
			if fmt.Sprintf("gamm/pool/%d", p) == lockup.Coins.GetDenomByIndex(0) {
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
				o.Locks = append(o.Locks,
					LockWithProof{
						Lock:   abciquery.Response.Value,
						Proof:  *abciquery.Response.ProofOps,
						Height: abciquery.Response.Height,
						Key:    abciquery.Response.Key,
					})

				continue OUTER
			}
		}
	}

	return o, nil

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
