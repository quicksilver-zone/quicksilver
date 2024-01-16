package claims

import (
	"context"
	"fmt"
	"time"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	osmolockup "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/lockup"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/ingenuity-build/xcclookup/pkgs/failsim"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

type TokenTuple struct {
	denom string
	chain string
}

func GetTokenMap(in []prewards.LiquidAllowedDenomProtocolData, zones []icstypes.Zone, chain, keyPrefix string) map[string]TokenTuple {
	out := make(map[string]TokenTuple)
	for _, i := range in {
		if i.ChainID == chain && ZoneOnboarded(zones, i) {
			out[keyPrefix+i.IbcDenom] = TokenTuple{denom: i.QAssetDenom, chain: i.RegisteredZoneChainID}
		} else {
			fmt.Printf("Zone not found: %s for LiquidToken: %s\n", i.RegisteredZoneChainID, i.IbcDenom)
		}
	}
	return out
}

func ZoneOnboarded(zones []icstypes.Zone, token prewards.LiquidAllowedDenomProtocolData) bool {
	for _, zone := range zones {
		if zone.ChainId == token.RegisteredZoneChainID {
			return true
		}
	}
	return false
}

func LiquidClaim(
	ctx context.Context,
	cfg types.Config,
	// poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
	zonesManager *types.CacheManager[icstypes.Zone],
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
	err = failsim.FailureHook(failures, 0, err, "failure decoding bech32 address")
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
		ctx,
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
	tokens := GetTokenMap(tokensManager.Get(ctx), zonesManager.Get(ctx), chain, "")

	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}

	for _, coin := range queryResponse.Balances {
		tuple, ok := tokens[coin.Denom]
		if !ok {
			fmt.Println("not dealing with token for chain", chain, coin.Denom)
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
			ctx,
			"/store/bank/key",
			lookupKey,
			rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
		)
		fmt.Println("Querying for value", "prefix", accountPrefix, "denom", tuple.denom) // debug?
		// 7:
		err = failsim.FailureHook(failures, 7, err, fmt.Sprintf("unable to query for value of denom %q on %q", tuple.denom, chain))
		if err != nil {
			return nil, nil, err
		}

		amount, err := bankkeeper.UnmarshalBalanceCompat(marshaler, abciquery.Response.Value, tuple.denom)
		if err != nil {
			return nil, nil, err
		}
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

	return msg, assets, nil
}
