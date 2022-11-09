package claims

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	osmogamm "github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
	osmolockup "github.com/ingenuity-build/quicksilver/osmosis-types/lockup"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type poolMap map[string][]osmogamm.PoolI

func OsmosisClaim(
	cfg types.Config,
	poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
	address string,
	chain string,
	height int64,
) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	_, addrBytes, err := bech32.DecodeNoLimit(address)
	if err != nil {
		return nil, nil, err
	}
	osmoAddress, err := bech32.Encode("osmo", addrBytes)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("2")

	host, ok := cfg.Chains[chain]
	if !ok {
		return nil, nil, fmt.Errorf("no endpoint in config for %s", chain)
	}
	fmt.Println("3")

	client, err := types.NewRPCClient(host, 30*time.Second)
	if err != nil {
		return nil, nil, err
	}
	// fetch timestamp of block
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	osmolockup.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	var timestamp time.Time

	if height == 0 {
		blockRequest := tmservice.GetLatestBlockRequest{}
		bytes := marshaler.MustMarshal(&blockRequest)
		abciquery, err := client.ABCIQuery(context.Background(), "/cosmos.base.tendermint.v1beta1.Service/GetLatestBlock", bytes)
		if err != nil {
			return nil, nil, err
		}
		fmt.Println("4")

		blockQueryResponse := tmservice.GetLatestBlockResponse{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &blockQueryResponse)
		if err != nil {
			return nil, nil, err
		}
		emptyBlockResponse := tmservice.GetLatestBlockResponse{}
		if blockQueryResponse == emptyBlockResponse {
			return nil, nil, fmt.Errorf("unable to query height from Osmosis chain")
		}
		//nolint:staticcheck // SA1019 ignore this!
		if blockQueryResponse.Block != nil {
			timestamp = blockQueryResponse.Block.Header.Time
			height = blockQueryResponse.Block.Header.Height
		} else {
			timestamp = blockQueryResponse.SdkBlock.Header.Time
			height = blockQueryResponse.Block.Header.Height
		}
	} else {
		blockRequest := tmservice.GetBlockByHeightRequest{Height: height}
		bytes := marshaler.MustMarshal(&blockRequest)
		abciquery, err := client.ABCIQuery(context.Background(), "/cosmos.base.tendermint.v1beta1.Service/GetBlockByHeight", bytes)
		if err != nil {
			return nil, nil, err
		}
		fmt.Println("4")

		blockQueryResponse := tmservice.GetBlockByHeightResponse{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &blockQueryResponse)
		if err != nil {
			return nil, nil, err
		}
		emptyBlockResponse := tmservice.GetBlockByHeightResponse{}
		if blockQueryResponse == emptyBlockResponse {
			return nil, nil, fmt.Errorf("unable to query height from Osmosis chain")
		}
		if blockQueryResponse.Block != nil { //nolint:staticcheck // SA1019 ignore this!
			timestamp = blockQueryResponse.Block.Header.Time //nolint:staticcheck // SA1019 ignore this!
		} else {
			timestamp = blockQueryResponse.SdkBlock.Header.Time
		}
	}
	fmt.Println("5")

	query := osmolockup.AccountLockedPastTimeRequest{Owner: osmoAddress, Timestamp: timestamp}
	bytes := marshaler.MustMarshal(&query)
	fmt.Println("6")

	abciquery, err := client.ABCIQueryWithOptions(context.Background(), "/osmosis.lockup.Query/AccountLockedPastTime", bytes, rpcclient.ABCIQueryOptions{Height: height})
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("7")

	queryResponse := osmolockup.AccountLockedPastTimeResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &queryResponse)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("8")

	// add GetFiltered to CacheManager, to allow filtered lookups on a single field == value
	tokens := func(in []prewards.LiquidAllowedDenomProtocolData) map[string]TokenTuple {
		out := make(map[string]TokenTuple)
		for _, i := range in {
			if i.ChainID == chain {
				out[i.IbcDenom] = TokenTuple{denom: i.QAssetDenom, chain: i.RegisteredZoneChainID}
			}
		}
		return out
	}(tokensManager.Get())
	fmt.Println("9")

	pools := poolMap{}
	for _, pool := range poolsManager.Get() {
		for chain := range pool.Zones {
			if _, ok := pools[chain]; !ok {
				pools[chain] = make([]osmogamm.PoolI, 0)
			}
			poolData, err := pool.GetPool()
			if err != nil {
				return nil, nil, err
			}
			pools[chain] = append(pools[chain], poolData)
		}
	}

	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}
	fmt.Println("10", queryResponse, pools)

OUTER:
	for _, lockup := range queryResponse.Locks { // for each lock in response
		for chainID, chainPools := range pools { // iterate over chains - are we doing all chains?
			for _, p := range chainPools { // iterate over the pools for this chain
				fmt.Println("PoolID", p.GetId())
				if fmt.Sprintf("gamm/pool/%d", p.GetId()) == lockup.Coins.GetDenomByIndex(0) {
					// perhaps counter intuitively, we want to group messages by chainID - the chain we are claiming for
					// and assets by chain - the chain on which they are located.
					fmt.Println("10a", chain, chainID)
					if _, ok := msg[chain]; !ok {
						msg[chainID] = prewards.MsgSubmitClaim{
							UserAddress: address,
							Zone:        chainID,
							SrcZone:     chain,
							ClaimType:   cmtypes.ClaimTypeOsmosisPool,
							Proofs:      make([]*cmtypes.Proof, 0),
						}
					}

					if _, ok := assets[chain]; !ok {
						assets[chain] = sdk.Coins{}
					}

					lockupKey := append(osmolockup.KeyPrefixPeriodLock, append(osmolockup.KeyIndexSeparator, sdk.Uint64ToBigEndian(lockup.ID)...)...)

					abciquery, err := client.ABCIQueryWithOptions(
						context.Background(), "/store/lockup/key",
						lockupKey,
						rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
					)
					if err != nil {
						return nil, nil, err
					}
					fmt.Println("10b")
					lockupResponse := osmolockup.PeriodLock{}
					err = marshaler.Unmarshal(abciquery.Response.Value, &lockupResponse)
					if err != nil {
						return nil, nil, err
					}
					fmt.Println("10c")
					gammCoins := lockupResponse.Coins
					gammShares := gammCoins.AmountOf("gamm/pool/" + strconv.Itoa(int(p.GetId())))

					exitedCoins, err := p.CalcExitPoolCoinsFromShares(sdk.Context{}, gammShares, sdk.ZeroDec())
					if err != nil {
						return nil, nil, err
					}
					fmt.Println("10d")

					for _, exitToken := range exitedCoins {
						tuple, ok := tokens[exitToken.Denom]
						if ok {
							exitToken.Denom = tuple.denom
							assets[chain] = assets[chain].Add(exitToken)
						}
					}

					chainMsg := msg[chainID]

					proof := cmtypes.Proof{
						Data:      abciquery.Response.Value,
						Key:       abciquery.Response.Key,
						ProofOps:  abciquery.Response.ProofOps,
						Height:    abciquery.Response.Height,
						ProofType: "lockup", // module name of proof.
					}

					chainMsg.Proofs = append(chainMsg.Proofs, &proof)
					fmt.Println("10e")
					msg[chainID] = chainMsg
					continue OUTER

				}
			}
		}
	}

	// fmt.Printf("Msg: %+v\n", msg)
	// fmt.Printf("Lockup Assets: %+v\n", assets)
	return msg, assets, nil
}
