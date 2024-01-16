package claims

import (
	"context"
	"fmt"
	"time"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	leverage "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage"
	leveragetypes "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage/types"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/ingenuity-build/xcclookup/pkgs/failsim"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

func UmeeClaim(
	ctx context.Context,
	cfg types.Config,
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
	zonesManager *types.CacheManager[icstypes.Zone],
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

	addrBytes, err := addressutils.AccAddressFromBech32(address, "")
	// 0:
	err = failsim.FailureHook(failures, 0, err, "failure decoding bech32 address")
	if err != nil {
		return nil, nil, err
	}
	umeeAddress, err := addressutils.EncodeAddressToBech32("umee", addrBytes)
	if err != nil {
		return nil, nil, err
	}

	// 1:
	err = failsim.FailureHook(failures, 1, err, "failure encoding umee address")
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("valid umee address encoding...")

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

	// we need the prefix
	bankquery := banktypes.QueryAllBalancesRequest{Address: umeeAddress}
	bytes := marshaler.MustMarshal(&bankquery)

	// query for AllBalances; then iterate, match against accepted balances and requery with proof.
	abciquery, err := client.ABCIQueryWithOptions(
		ctx,
		"/cosmos.bank.v1beta1.Query/AllBalances",
		bytes,
		rpcclient.ABCIQueryOptions{Height: height},
	)
	// 4:
	err = failsim.FailureHook(failures, 5, err, "ABCIQuery: AllBalances")
	if err != nil {
		return nil, nil, err
	}
	bankQueryResponse := banktypes.QueryAllBalancesResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &bankQueryResponse)
	if err != nil {
		return nil, nil, err
	}

	leveragequery := leverage.QueryAccountBalances{Address: umeeAddress}
	bytes = marshaler.MustMarshal(&leveragequery)
	// query for AllBalances; then iterate, match against accepted balances and requery with proof.
	leverageaccountbalancesquery, err := client.ABCIQueryWithOptions(
		ctx,
		"/umee.leverage.v1.Query/AccountBalances",
		bytes,
		rpcclient.ABCIQueryOptions{Height: height},
	)
	// 5:
	err = failsim.FailureHook(failures, 6, err, "ABCIQuery: QueryAccountBalancesRequest")
	if err != nil {
		return nil, nil, err
	}
	leverageQueryResponse := leverage.QueryAccountBalancesResponse{}
	err = marshaler.Unmarshal(leverageaccountbalancesquery.Response.Value, &leverageQueryResponse)
	if err != nil {
		return nil, nil, err
	}

	// add GetFiltered to CacheManager, to allow filtered lookups on a single field == value
	tokens := GetTokenMap(tokensManager.Get(ctx), zonesManager.Get(ctx), chain, leveragetypes.UTokenPrefix)

	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}

	// bank balance
	for _, coin := range bankQueryResponse.Balances {
		if len(coin.GetDenom()) < 2 || coin.GetDenom()[0:2] != leveragetypes.UTokenPrefix {
			continue
		}
		tuple, ok := tokens[coin.GetDenom()]
		if !ok {
			fmt.Println("not dealing with token for chain", chain, coin.GetDenom())
			// token is not present in list of allowed tokens, ignore.
			continue
		}

		if _, ok := msg[tuple.chain]; !ok {
			msg[tuple.chain] = prewards.MsgSubmitClaim{
				UserAddress: address,
				Zone:        tuple.chain,
				SrcZone:     chain,
				ClaimType:   cmtypes.ClaimTypeUmeeToken,
				Proofs:      make([]*cmtypes.Proof, 0),
			}
		}

		accountPrefix := banktypes.CreateAccountBalancesPrefix(addrBytes)
		lookupKey := append(accountPrefix, []byte(coin.GetDenom())...)
		bankquery, err := client.ABCIQueryWithOptions(
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

		amount, err := bankkeeper.UnmarshalBalanceCompat(marshaler, bankquery.Response.Value, tuple.denom)
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
			Data:      bankquery.Response.Value,
			Key:       bankquery.Response.Key,
			ProofOps:  bankquery.Response.ProofOps,
			Height:    bankquery.Response.Height,
			ProofType: prewards.ProofTypeBank, // module name of proof.
		}

		chainMsg.Proofs = append(chainMsg.Proofs, &proof)

		msg[tuple.chain] = chainMsg
	}

	// leverage account balance
	for _, coin := range leverageQueryResponse.Collateral {
		if len(coin.GetDenom()) < 2 || coin.GetDenom()[0:2] != leveragetypes.UTokenPrefix {
			continue
		}
		tuple, ok := tokens[coin.GetDenom()]
		if !ok {
			fmt.Println("not dealing with token for chain", chain, coin.GetDenom())
			// token is not present in list of allowed tokens, ignore.
			continue
		}

		if _, ok := msg[tuple.chain]; !ok {
			msg[tuple.chain] = prewards.MsgSubmitClaim{
				UserAddress: address,
				Zone:        tuple.chain,
				SrcZone:     chain,
				ClaimType:   cmtypes.ClaimTypeUmeeToken,
				Proofs:      make([]*cmtypes.Proof, 0),
			}
		}

		lookupKey := leveragetypes.KeyCollateralAmount(addrBytes, coin.GetDenom())
		leveragequery, err := client.ABCIQueryWithOptions(
			ctx,
			"/store/leverage/key",
			lookupKey,
			rpcclient.ABCIQueryOptions{Height: leverageaccountbalancesquery.Response.Height, Prove: true},
		)
		fmt.Println("Querying for value", "prefix", lookupKey) // debug?
		// 7:
		err = failsim.FailureHook(failures, 7, err, fmt.Sprintf("unable to query for value of denom %q on %q", tuple.denom, chain))
		if err != nil {
			return nil, nil, err
		}

		amount, err := bankkeeper.UnmarshalBalanceCompat(marshaler, leveragequery.Response.Value, tuple.denom)
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
			Data:      leveragequery.Response.Value,
			Key:       leveragequery.Response.Key,
			ProofOps:  leveragequery.Response.ProofOps,
			Height:    leveragequery.Response.Height,
			ProofType: prewards.ProofTypeLeverage, // module name of proof.
		}

		chainMsg.Proofs = append(chainMsg.Proofs, &proof)

		msg[tuple.chain] = chainMsg
	}

	return msg, assets, nil
}
