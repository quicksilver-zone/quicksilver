package claims

import (
	"context"
	"fmt"
	"time"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/xcclookup/pkgs/types"
)

type TokenTuple struct {
	denom string
	chain string
}

func GetTokenMap(in []prewards.LiquidAllowedDenomProtocolData, zones []icstypes.Zone, chain, keyPrefix string, ignores types.Ignores) map[string]TokenTuple {
	out := make(map[string]TokenTuple)
	for _, i := range in {
		if ignores.Contains(i.QAssetDenom) {
			continue
		}
		if i.ChainID == chain && ZoneOnboarded(zones, i) {
			out[keyPrefix+i.IbcDenom] = TokenTuple{denom: i.QAssetDenom, chain: i.RegisteredZoneChainID}
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
	cacheMgr *types.CacheManager,
	address string,
	submitAddress string,
	connection prewards.ConnectionProtocolData,
	height int64,
) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	chain := connection.ChainID
	prefix := connection.Prefix

	// explitly don't use quick prefix here, as mapped accounts may have a different prefix
	addrBytes, err := addressutils.AccAddressFromBech32(address, "")
	if err != nil {
		return nil, nil, fmt.Errorf("%w [addressutils.AddressFromBech32]", err)
	}

	chainAddress, err := addressutils.EncodeAddressToBech32(prefix, addrBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("%w [bech32.Encode]", err)
	}

	host, ok := cfg.Chains[chain]
	if !ok {
		err = fmt.Errorf("unable to find endpoint for %s", chain)
		return nil, nil, nil
	}

	client, err := types.NewRPCClient(host, time.Duration(cfg.Timeout)*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("%w [NewRPCClient]", err)
	}

	// fetch timestamp of block
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	banktypes.RegisterInterfaces(interfaceRegistry)
	cmtypes.RegisterInterfaces(interfaceRegistry)
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
	if err != nil {
		return nil, nil, fmt.Errorf("%w [ABCIQueryWithOptions/AllBalances]", err)
	}

	queryResponse := banktypes.QueryAllBalancesResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &queryResponse)
	if err != nil {
		return nil, nil, fmt.Errorf("%w [unmarshalling query response]", err)
	}

	ignores := cfg.Ignore.GetIgnoresForType(types.IgnoreTypeLiquid)

	// add GetFiltered to CacheManager, to allow filtered lookups on a single field == value
	tokens := GetTokenMap(types.GetCache[prewards.LiquidAllowedDenomProtocolData](ctx, cacheMgr), types.GetCache[icstypes.Zone](ctx, cacheMgr), chain, "", ignores)

	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}

	for _, coin := range queryResponse.Balances {
		tuple, ok := tokens[coin.Denom]
		if !ok {
			continue
		}

		if _, ok := msg[tuple.chain]; !ok {
			msg[tuple.chain] = prewards.MsgSubmitClaim{
				UserAddress: submitAddress,
				Zone:        tuple.chain,
				SrcZone:     chain,
				ClaimType:   cmtypes.ClaimTypeLiquidToken,
				Proofs:      make([]*cmtypes.Proof, 0),
			}
		}

		accountPrefix := banktypes.CreateAccountBalancesPrefix(addrBytes)
		lookupKey := append(accountPrefix, []byte(coin.Denom)...)
		abciquery, err := client.ABCIQueryWithOptions(
			ctx,
			"/store/bank/key",
			lookupKey,
			rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
		)
		fmt.Println("Querying for value (liquid tokens)", "chain", chain, "prefix", accountPrefix, "denom", tuple.denom) // debug?
		// 7:
		if err != nil {
			return nil, nil, fmt.Errorf("%w [ABCIQueryWithOptions/gamm_tokens]", err)
		}

		amount, err := bankkeeper.UnmarshalBalanceCompat(marshaler, abciquery.Response.Value, tuple.denom)
		if err != nil {
			return nil, nil, fmt.Errorf("%w [UnmarshalBalanceCompat]", err)
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

		msg[tuple.chain] = chainMsg
	}

	return msg, assets, nil
}
