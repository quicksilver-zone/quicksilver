package claims

import (
	"context"
	"fmt"
	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	liquiditytypes "github.com/ingenuity-build/quicksilver/crescent-types/liquidity/types"
	lpfarmtypes "github.com/ingenuity-build/quicksilver/crescent-types/lpfarm"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	"github.com/ingenuity-build/xcclookup/pkgs/failsim"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"time"
)

func CrescentClaim(
	ctx context.Context,
	cfg types.Config,
	poolsManager *types.CacheManager[prewards.CrescentPoolProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
	address string,
	chain string,
	height int64,
) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	// simFailure hooks: 0-8
	simFailures := failsim.FailuresFromContext(ctx)
	failures := make(map[uint8]struct{})
	if UmeeClaimFailures, ok := simFailures[2]; ok {
		fmt.Println("liquid sim failures")
		failures = UmeeClaimFailures
	}
	fmt.Println("simulate failures:", failures)

	_, addrBytes, err := bech32.DecodeNoLimit(address)
	// 0:
	err = failsim.FailureHook(failures, 0, err, "failure decoding bech32 address")
	if err != nil {
		return nil, nil, err
	}
	crescentAddress, err := bech32.Encode("cre", addrBytes)
	if err != nil {
		return nil, nil, err
	}
	// 1:
	err = failsim.FailureHook(failures, 1, err, "failure encoding crescent address")
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("valid crescent address encoding...")

	host, ok := cfg.Chains[chain]
	if !ok {
		err = fmt.Errorf("no endpoint in config for %s", chain)
	}
	// 2:
	err = failsim.FailureHook(failures, 2, err, fmt.Sprintf("no endpoint in config for %s", chain))
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("found %q in config for %q...\n", host, chain)

	client, err := types.NewRPCClient(host, 30*time.Second)
	// 3:

	err = failsim.FailureHook(failures, 4, err, fmt.Sprintf("failure connecting to host %q", host))
	if err != nil {
		return nil, nil, err
	}
	// fetch timestamp of block
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	banktypes.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	positionsQuery := lpfarmtypes.QueryPositionsRequest{Farmer: crescentAddress, Pagination: &query.PageRequest{Offset: 0, Limit: 1000}}
	bytes := marshaler.MustMarshal(&positionsQuery)
	// query for AllBalances; then iterate, match against accepted balances and requery with proof.
	abciquery, err := client.ABCIQueryWithOptions(
		context.Background(),
		"/crescent/lpfarm/v1beta1/positions/{farmer}",
		bytes,
		rpcclient.ABCIQueryOptions{Height: height},
	)
	// 4:
	err = failsim.FailureHook(failures, 6, err, "ABCIQuery: QueryPositions")
	if err != nil {
		return nil, nil, err
	}
	positionsQueryResponse := lpfarmtypes.QueryPositionsResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &positionsQueryResponse)

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
	fmt.Println("got relevant pools...")

	var errors map[string]error

OUTER:
	for _, position := range positionsQueryResponse.Positions {
		for _, cpd := range poolsManager.Get() {
			for chainID, denomAndZone := range cpd.Denoms {
				if fmt.Sprintf("pool%d", cpd.PoolID) == position.Denom {
					tuple, ok := tokens[chainID]
					if !ok {
						fmt.Println("not dealing with token for chain", chainID, denomAndZone.Denom)
						// token is not present in list of allowed tokens, ignore.
						continue
					}

					if _, ok := msg[tuple.chain]; !ok {
						msg[tuple.chain] = prewards.MsgSubmitClaim{
							UserAddress: address,
							Zone:        tuple.chain,
							SrcZone:     chain,
							ClaimType:   cmtypes.ClaimTypeCrescentPool,
							Proofs:      make([]*cmtypes.Proof, 0),
						}
					}

					if _, ok := assets[chain]; !ok {
						assets[chain] = sdk.Coins{}
					}

					farmerAddr, err := sdk.AccAddressFromBech32(position.Farmer)

					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("invalid farmer address %q: %w", chain, err)
						continue
					}

					positionKey := lpfarmtypes.GetPositionKey(farmerAddr, position.Denom)

					abciquery, err := client.ABCIQueryWithOptions(
						context.Background(), "/store/lpfarm/key",
						positionKey,
						rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
					)
					// 9:
					err = failsim.FailureHook(failures, 9, err, "ABCIQuery: position")
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to account for assets on zone %q: %w", chain, err)
						continue
					}
					fmt.Println("prepared query for position...")
					positionResponse := lpfarmtypes.Position{}
					err = marshaler.Unmarshal(abciquery.Response.Value, &positionResponse)
					// 10:
					err = failsim.FailureHook(failures, 10, err, "ABCIQuery: position response")
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to account for assets on zone %q: %w", chain, err)
						continue
					}

					// query to get pool info
					poolQuery, err := client.ABCIQueryWithOptions(
						context.Background(), "/store/liquidity/key",
						liquiditytypes.GetPoolKey(cpd.PoolID),
						rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
					)
					poolResponse := liquiditytypes.Pool{}
					err = marshaler.Unmarshal(poolQuery.Response.Value, &poolResponse)
					// 11:
					err = failsim.FailureHook(failures, 10, err, "ABCIQuery: pool response")

					// fetch reserveAddress balance
					_, reserveAddrBytes, err := bech32.DecodeNoLimit(poolResponse.ReserveAddress)

					accountPrefix := banktypes.CreateAccountBalancesPrefix(reserveAddrBytes)
					lookupKey := append(accountPrefix, []byte(tuple.denom)...)

					bankQuery, err := client.ABCIQueryWithOptions(
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

					amount, err := bankkeeper.UnmarshalBalanceCompat(marshaler, bankQuery.Response.Value, tuple.denom)
					if err != nil {
						return nil, nil, err
					}
					// 12:
					err = failsim.FailureHook(failures, 8, err, fmt.Sprintf("ABCIQuery: value of denom %q on chain %q", tuple.denom, chain))
					if err != nil {
						return nil, nil, err
					}

					// fetch total poolcoin supply
					supplyQuery, err := client.ABCIQueryWithOptions(
						context.Background(), "/store/bank/key",
						append(banktypes.SupplyKey, []byte(positionResponse.Denom)...),
						rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
					)
					fmt.Println("Querying for poolcoinsupply", "prefix", banktypes.SupplyKey, "denom", positionResponse.Denom) // debug?
					// 7:
					// 13:
					err = failsim.FailureHook(failures, 7, err, fmt.Sprintf("unable to query for value of denom %q on %q", positionResponse.Denom, chain))
					if err != nil {
						return nil, nil, err
					}

					farmingAmount := positionResponse.FarmingAmount
					poolSupply := sdk.ZeroInt()
					err = poolSupply.Unmarshal(supplyQuery.Response.Value)
					if err != nil {
						return nil, nil, err
					}

					uratio := sdk.NewDecFromInt(farmingAmount).QuoInt(poolSupply)

					uAmount := uratio.MulInt(amount.Amount).TruncateInt()

					assets[chain] = assets[chain].Add(sdk.NewCoin(tuple.denom, uAmount))

					chainMsg := msg[chainID]

					proof := cmtypes.Proof{
						Data:      abciquery.Response.Value,
						Key:       abciquery.Response.Key,
						ProofOps:  abciquery.Response.ProofOps,
						Height:    abciquery.Response.Height,
						ProofType: prewards.ProofTypePosition,
					}

					chainMsg.Proofs = append(chainMsg.Proofs, &proof)
					fmt.Println("obtained relevant proofs...")
					msg[chainID] = chainMsg
					continue OUTER
				}
			}
		}
	}

	return msg, assets, nil
}
