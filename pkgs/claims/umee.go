package claims

//revive:disable:redundant-import-alias
import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	rpcclient "github.com/cometbft/cometbft/rpc/client"

	leverage "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage"
	leveragetypes "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage/types"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/xcclookup/pkgs/types"
)

func UmeeClaim(
	ctx context.Context,
	cfg types.Config,
	cacheMgr *types.CacheManager,
	address string,
	submitAddress string,
	chain string,
	height int64,
) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	addrBytes, err := addressutils.AccAddressFromBech32(address, "")
	if err != nil {
		return nil, nil, err
	}
	umeeAddress, err := addressutils.EncodeAddressToBech32("umee", addrBytes)
	if err != nil {
		return nil, nil, err
	}

	if err != nil {
		return nil, nil, err
	}
	fmt.Println("valid umee address encoding...")

	host, ok := cfg.Chains[chain]
	if !ok {
		err = fmt.Errorf("no endpoint in config for %s", chain)
	}
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("found %q in config for %q...\n", host, chain)

	client, err := types.NewRPCClient(host, time.Duration(cfg.Timeout)*time.Second)
	if err != nil {
		return nil, nil, err
	}
	// fetch timestamp of block
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	// banktypes.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	leveragequery := leverage.QueryAccountBalances{Address: umeeAddress}
	bytes := marshaler.MustMarshal(&leveragequery)
	// query for AllBalances; then iterate, match against accepted balances and requery with proof.
	leverageaccountbalancesquery, err := client.ABCIQueryWithOptions(
		ctx,
		"/umee.leverage.v1.Query/AccountBalances",
		bytes,
		rpcclient.ABCIQueryOptions{Height: height},
	)
	if err != nil {
		return nil, nil, err
	}
	leverageQueryResponse := leverage.QueryAccountBalancesResponse{}
	err = marshaler.Unmarshal(leverageaccountbalancesquery.Response.Value, &leverageQueryResponse)
	if err != nil {
		return nil, nil, err
	}

	ignores := cfg.Ignore.GetIgnoresForType(types.IgnoreTypeLiquid)

	// add GetFiltered to CacheManager, to allow filtered lookups on a single field == value
	tokens := GetTokenMap(types.GetCache[prewards.LiquidAllowedDenomProtocolData](ctx, cacheMgr), types.GetCache[icstypes.Zone](ctx, cacheMgr), chain, leveragetypes.UTokenPrefix, ignores)

	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}

	// leverage account balance
	for _, coin := range leverageQueryResponse.Collateral {
		if len(coin.GetDenom()) < 2 || coin.GetDenom()[0:2] != leveragetypes.UTokenPrefix {
			continue
		}
		tuple, ok := tokens[coin.GetDenom()]
		if !ok {
			continue
		}

		if _, ok := msg[tuple.chain]; !ok {
			msg[tuple.chain] = prewards.MsgSubmitClaim{
				UserAddress: submitAddress,
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
		fmt.Println("Querying for value (umee - leverage)", "prefix", string(lookupKey)) // debug?
		if err != nil {
			return nil, nil, err
		}

		amount, err := bankkeeper.UnmarshalBalanceCompat(marshaler, leveragequery.Response.Value, tuple.denom)
		if err != nil {
			return nil, nil, err
		}
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
