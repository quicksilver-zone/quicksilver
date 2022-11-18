package claims

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	osmolockup "github.com/ingenuity-build/quicksilver/osmosis-types/lockup"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	"github.com/ingenuity-build/xcclookup/pkgs/failsim"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type TokenTuple struct {
	denom string
	chain string
}

func LiquidClaim(
	ctx context.Context,
	cfg types.Config,
	// poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
	address string,
	connection prewards.ConnectionProtocolData,
	height int64,
) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	// simFailure hooks: 0-8
	simFailures := failsim.FailuresFromContext(ctx)
	failures := make(map[uint8]struct{})
	if LiquidClaimFailures, ok := simFailures[2]; ok {
		fmt.Println("liquid sim failures")
		failures = LiquidClaimFailures
	}
	fmt.Println("simulate failures:", failures)

	chain := connection.ChainID
	prefix := connection.Prefix

	_, addrBytes, err := bech32.Decode(address, 51)
	// 0:
	err = failsim.FailureHook(failures, 0, err, "failure decosing bech32 address")
	if err != nil {
		return nil, nil, err
	}
	sdkAddr, err := bech32.ConvertBits(addrBytes, 5, 8, true)
	// 1:
	err = failsim.FailureHook(failures, 1, err, "failure converting sdk address")
	if err != nil {
		return nil, nil, err
	}

	chainAddress, err := bech32.Encode(prefix, addrBytes)
	// 2:
	err = failsim.FailureHook(failures, 2, err, "failure encoding chain address")
	if err != nil {
		return nil, nil, err
	}

	host, ok := cfg.Chains[chain]
	if !ok {
		err = fmt.Errorf("unable to find endpoint for %s", chain)
	}
	// 3:
	err = failsim.FailureHook(failures, 3, err, fmt.Sprintf("no endpoint in config for %s", chain))
	if err != nil {
		return nil, nil, err
	}
	client, err := types.NewRPCClient(host, 30*time.Second)
	// 4:
	err = failsim.FailureHook(failures, 4, err, fmt.Sprintf("failure connecting to host %q", host))
	if err != nil {
		return nil, nil, err
	}
	// fetch timestamp of block
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	banktypes.RegisterInterfaces(interfaceRegistry)
	osmolockup.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	// we need the prefix
	query := banktypes.QueryAllBalancesRequest{Address: chainAddress}
	bytes := marshaler.MustMarshal(&query)

	// query for AllBalances; then iterate, match against accepted balances and requery with proof.
	abciquery, err := client.ABCIQueryWithOptions(
		context.Background(),
		"/cosmos.bank.v1beta1.Query/AllBalances",
		bytes,
		rpcclient.ABCIQueryOptions{Height: height},
	)
	// 5:
	err = failsim.FailureHook(failures, 5, err, "ABCIQuery: AllBalances")
	if err != nil {
		return nil, nil, err
	}
	queryResponse := banktypes.QueryAllBalancesResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &queryResponse)
	// 6:
	err = failsim.FailureHook(failures, 6, err, "ABCIQuery: QueryAllBalancesResponse")
	if err != nil {
		return nil, nil, err
	}

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

	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}

	for _, coin := range queryResponse.Balances {
		tuple, ok := tokens[coin.Denom]
		if !ok {
			fmt.Println("not dealing with token for chain", chain, coin.Denom)
			// token is not present in list of allowed tokens, ignore.
			// TODO: handle gamm tokens here, if chain is osmosis
			continue
		}

		if _, ok := msg[tuple.chain]; !ok {
			msg[tuple.chain] = prewards.MsgSubmitClaim{
				UserAddress: address,
				Zone:        tuple.chain,
				SrcZone:     chain,
				ClaimType:   cmtypes.ClaimTypeLiquidToken,
				Proofs:      make([]*cmtypes.Proof, 0),
			}
		}

		accountPrefix := banktypes.CreateAccountBalancesPrefix(sdkAddr)
		lookupKey := append(accountPrefix, []byte(coin.Denom)...)
		abciquery, err := client.ABCIQueryWithOptions(
			context.Background(), "/store/bank/key",
			lookupKey,
			rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
		)
		fmt.Println("Querying for value", "prefix", accountPrefix, "denom", tuple.denom) // debug?
		// 7:
		err = failsim.FailureHook(failures, 7, err, fmt.Sprintf("unable to query for value of denom %q on %q", tuple.denom, chain))
		if err != nil {
			return nil, nil, err
		}

		amount := sdk.Coin{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &amount)
		// 8:
		err = failsim.FailureHook(failures, 8, err, fmt.Sprintf("ABCIQuery: value of denom %q on chain %q", tuple.denom, chain))
		if err != nil {
			return nil, nil, err
		}
		amount.Denom = tuple.denom

		assets[chain] = assets[chain].Add(amount)

		chainMsg := msg[tuple.chain]

		proof := cmtypes.Proof{
			Data:      abciquery.Response.Value,
			Key:       abciquery.Response.Key,
			ProofOps:  abciquery.Response.ProofOps,
			Height:    abciquery.Response.Height,
			ProofType: "bank", // module name of proof.
		}

		chainMsg.Proofs = append(chainMsg.Proofs, &proof)

		// fmt.Printf("Liquid Assets: %+v\n", assets)

		msg[tuple.chain] = chainMsg
	}

	// 	pools := poolMap{}
	// 	// filter by pool id - query this from Quicksilver (and cache hourly)
	// 	for _, pool := range poolsManager.Get() {
	// 		for chain := range pool.Zones {
	// 			if _, ok := pools[chain]; !ok {
	// 				pools[chain] = make([]osmogamm.PoolI, 0)
	// 			}
	// 			poolData, err := pool.GetPool()
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			pools[chain] = append(pools[chain], poolData)
	// 		}
	// 	}

	// OUTER:
	// 	for _, lockup := range queryResponse.Locks { // for each lock in response
	// 		for chainID, chainPools := range pools { // iterate over chains - are we doing all chains?
	// 			for _, p := range chainPools { // iterate over the pools for this chain
	// 				if fmt.Sprintf("gamm/pool/%d", p.GetId()) == lockup.Coins.GetDenomByIndex(0) {
	// 					if _, ok := msg[chainID]; !ok {
	// 						msg[chainID] = prewards.MsgSubmitClaim{
	// 							UserAddress: address,
	// 							Zone:        chainID,
	// 							SrcZone:     chain,
	// 							ClaimType:   prewards.ClaimTypeOsmosisPool,
	// 							Proofs:      make([]*prewards.Proof, 0),
	// 						}
	// 					}

	// 					if _, ok := msg[chainID]; !ok {
	// 						assets[chainID] = sdk.Coins{}
	// 					}

	// 					abciquery, err := client.ABCIQueryWithOptions(
	// 						context.Background(), "/store/lockup/key",
	// 						append(osmolockup.KeyPrefixPeriodLock, append(osmolockup.KeyIndexSeparator, sdk.Uint64ToBigEndian(lockup.ID)...)...),
	// 						rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
	// 					)
	// 					if err != nil {
	// 						return nil, nil, err
	// 					}
	// 					lockupResponse := osmolockup.PeriodLock{}
	// 					err = marshaler.Unmarshal(abciquery.Response.Value, &lockupResponse)
	// 					if err != nil {
	// 						return nil, nil, err
	// 					}
	// 					gammCoins := lockupResponse.Coins
	// 					gammShares := gammCoins.AmountOf("gamm/pool/" + strconv.Itoa(int(p.GetId())))

	// 					exitedCoins, err := p.CalcExitPoolCoinsFromShares(sdk.Context{}, gammShares, sdk.ZeroDec())
	// 					if err != nil {
	// 						return nil, nil, err
	// 					}

	// 					assets[chainID] = assets[chainID].Add(exitedCoins...)

	// 					chainMsg := msg[chainID]

	// 					proof := prewards.Proof{
	// 						Data:     abciquery.Response.Value,
	// 						Key:      abciquery.Response.Key,
	// 						ProofOps: abciquery.Response.ProofOps,
	// 						Height:   abciquery.Response.Height,
	// 					}

	// 					chainMsg.Proofs = append(chainMsg.Proofs, &proof)

	// 					msg[chainID] = chainMsg
	// 					continue OUTER

	// 				}

	// 			}
	// 		}
	// 	}

	return msg, assets, nil
}
