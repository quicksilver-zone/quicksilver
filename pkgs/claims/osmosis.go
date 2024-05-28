package claims

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ingenuity-build/multierror"
	osmocl "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity"
	osmogamm "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm"
	osmolockup "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/lockup"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

type (
	poolMap   map[string][]osmogamm.CFMMPoolI
	clPoolMap map[string][]osmocl.ConcentratedPoolExtension
)

func OsmosisClaim(
	ctx context.Context,
	cfg types.Config,
	cacheMgr *types.CacheManager,
	address string,
	submitAddress string,
	chain string,
	height int64,
) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	var err error

	addrBytes, err := addressutils.AccAddressFromBech32(address, "")
	if err != nil {
		return nil, nil, err
	}

	osmoAddress, err := addressutils.EncodeAddressToBech32("osmo", addrBytes)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("valid osmosis address encoding...")

	host, ok := cfg.Chains[chain]
	if !ok {
		err = fmt.Errorf("no endpoint in config for %s", chain)
		return nil, nil, err
	}

	fmt.Printf("found %q in config for %q...\n", host, chain)

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
		abciquery, err := client.ABCIQuery(
			ctx,
			"/cosmos.base.tendermint.v1beta1.Service/GetLatestBlock",
			bytes,
		)
		// 4:
		if err != nil {
			return nil, nil, err
		}
		fmt.Println("height is zero, get latest block height...")

		blockQueryResponse := tmservice.GetLatestBlockResponse{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &blockQueryResponse)
		if err != nil {
			return nil, nil, err
		}

		emptyBlockResponse := tmservice.GetLatestBlockResponse{}
		if blockQueryResponse == emptyBlockResponse {
			err = fmt.Errorf("unable to query height from Osmosis chain")
		}
		if err != nil {
			return nil, nil, err
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
		abciquery, err := client.ABCIQuery(
			ctx,
			"/cosmos.base.tendermint.v1beta1.Service/GetBlockByHeight",
			bytes,
		)
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("height is %d, get block by height...\n", height)

		blockQueryResponse := tmservice.GetBlockByHeightResponse{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &blockQueryResponse)
		if err != nil {
			return nil, nil, err
		}

		emptyBlockResponse := tmservice.GetBlockByHeightResponse{}
		if blockQueryResponse == emptyBlockResponse {
			err = fmt.Errorf("unable to query height from Osmosis chain")
		}
		if err != nil {
			return nil, nil, err
		}
		if blockQueryResponse.Block != nil { //nolint:staticcheck // SA1019 ignore this!
			timestamp = blockQueryResponse.Block.Header.Time //nolint:staticcheck // SA1019 ignore this!
		} else {
			timestamp = blockQueryResponse.SdkBlock.Header.Time
		}
	}
	fmt.Println("got block timestamp...", timestamp)

	query := osmolockup.AccountLockedPastTimeRequest{Owner: osmoAddress, Timestamp: timestamp}
	bytes := marshaler.MustMarshal(&query)
	fmt.Println("prepared account lockup query...")

	abciquery, err := client.ABCIQueryWithOptions(
		ctx,
		"/osmosis.lockup.Query/AccountLockedPastTime",
		bytes,
		rpcclient.ABCIQueryOptions{Height: height},
	)
	if err != nil {
		return nil, nil, err
	}

	queryResponse := osmolockup.AccountLockedPastTimeResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &queryResponse)
	if err != nil {
		return nil, nil, err
	}

	ignores := cfg.Ignore.GetIgnoresForType(types.IgnoreTypeLiquid)
	// add GetFiltered to CacheManager, to allow filtered lookups on a single field == value
	tokens := GetTokenMap(types.GetCache[prewards.LiquidAllowedDenomProtocolData](ctx, cacheMgr), types.GetCache[icstypes.Zone](ctx, cacheMgr), chain, "", ignores)

	pools := poolMap{}
	clpools := clPoolMap{}

	ignores = cfg.Ignore.GetIgnoresForType(types.IgnoreTypeOsmosisPool)

	for _, pool := range types.GetCache[prewards.OsmosisPoolProtocolData](ctx, cacheMgr) {
		if ignores.Contains(strconv.Itoa(int(pool.PoolID))) {
			continue
		}
		if pool.IsIncentivized {
			for _, denom := range pool.Denoms {
				if _, ok := pools[denom.ChainID]; !ok {
					pools[denom.ChainID] = make([]osmogamm.CFMMPoolI, 0)
				}
				poolData, err := pool.GetPool()
				if err != nil {
					return nil, nil, err
				}
				pools[denom.ChainID] = append(pools[denom.ChainID], poolData)
			}
		}
	}

	ignores = cfg.Ignore.GetIgnoresForType(types.IgnoreTypeOsmosisCLPool)

	for _, clpool := range types.GetCache[prewards.OsmosisClPoolProtocolData](ctx, cacheMgr) {
		if ignores.Contains(strconv.Itoa(int(clpool.PoolID))) {
			continue
		}
		if clpool.IsIncentivized {
			for _, denom := range clpool.Denoms {
				if _, ok := pools[denom.ChainID]; !ok {
					clpools[denom.ChainID] = make([]osmocl.ConcentratedPoolExtension, 0)
				}
				poolData, err := clpool.GetPool()
				if err != nil {
					return nil, nil, err
				}
				clpools[denom.ChainID] = append(clpools[denom.ChainID], poolData)
			}
		}
	}

	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}
	fmt.Println("got relevant pools...")

	var errors map[string]error
	var poolCoinDenom string

	for chainID, chainPools := range pools { // iterate over chains - are we doing all chains?
		for _, p := range chainPools { // iterate over the pools for this chain
			// fetching unbonded gamm tokens from account

			poolCoinDenom = fmt.Sprintf("gamm/pool/%d", p.GetId())

			accountPrefix := banktypes.CreateAccountBalancesPrefix(addrBytes)
			lookupKey := append(accountPrefix, []byte(poolCoinDenom)...)

			fmt.Println("Querying for value (liquid gamm)", "prefix", accountPrefix, "denom", poolCoinDenom) // debug?
			bankQuery, err := client.ABCIQueryWithOptions(
				ctx, "/store/bank/key",
				lookupKey,
				rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
			)
			if err != nil {
				return nil, nil, err
			}

			amount, err := bankkeeper.UnmarshalBalanceCompat(marshaler, bankQuery.Response.Value, poolCoinDenom)
			if err != nil {
				return nil, nil, err
			}

			if !amount.IsZero() {

				if _, ok := assets[chain]; !ok {
					assets[chain] = sdk.Coins{}
				}

				var exitedCoins sdk.Coins
				// if this user is the sole position in the pool, sub one share to avoid CalcExitPoolCoinsFromShares erroring.
				if amount.Amount.Equal(p.GetTotalShares()) {
					exitedCoins, err = p.CalcExitPoolCoinsFromShares(sdk.Context{}, amount.Amount.Sub(math.OneInt()), math.LegacyZeroDec())
				} else {
					exitedCoins, err = p.CalcExitPoolCoinsFromShares(sdk.Context{}, amount.Amount, math.LegacyZeroDec())
				}
				if err != nil {
					if errors == nil {
						errors = make(map[string]error)
					}
					errors[chain] = fmt.Errorf("unable to account for assets on zone %q: %w", chain, err)
					continue
				}

				fmt.Println("exitedCoins: ", exitedCoins)
				fmt.Println("tokens: ", tokens)

				for _, exitToken := range exitedCoins {
					tuple, ok := tokens[exitToken.Denom]
					if ok {
						exitToken.Denom = tuple.denom
						assets[chain] = assets[chain].Add(exitToken)

						fmt.Printf("gamm tokens %s -> token %s (%s), on %s...\n", amount, exitToken, chainID, chain)
						if _, ok := msg[chainID]; !ok {
							msg[chainID] = prewards.MsgSubmitClaim{
								UserAddress: submitAddress,
								Zone:        chainID,
								SrcZone:     chain,
								ClaimType:   cmtypes.ClaimTypeOsmosisPool,
								Proofs:      make([]*cmtypes.Proof, 0),
							}
						}
					}
				}

				chainMsg := msg[chainID]

				proof := cmtypes.Proof{
					Data:      bankQuery.Response.Value,
					Key:       bankQuery.Response.Key,
					ProofOps:  bankQuery.Response.ProofOps,
					Height:    bankQuery.Response.Height,
					ProofType: prewards.ProofTypeBank,
				}

				chainMsg.Proofs = append(chainMsg.Proofs, &proof)

				msg[chainID] = chainMsg
			}
			for _, lockup := range queryResponse.Locks { // for each lock in response
				// checking locked coins
				if poolCoinDenom == lockup.Coins.GetDenomByIndex(0) {
					// perhaps counter intuitively, we want to group messages by chainID - the chain we are claiming for
					// and assets by chain - the chain on which they are located.
					fmt.Printf("found assets for zone %q...\n", chainID)
					if _, ok := msg[chainID]; !ok {
						msg[chainID] = prewards.MsgSubmitClaim{
							UserAddress: submitAddress,
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
						ctx,
						"/store/lockup/key",
						lockupKey,
						rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
					)
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to account for assets on zone %q: %w", chain, err)
						continue
					}

					fmt.Println("prepared query for locked assets...")
					lockupResponse := osmolockup.PeriodLock{}
					err = marshaler.Unmarshal(abciquery.Response.Value, &lockupResponse)
					// 10:
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to account for assets on zone %q: %w", chain, err)
						continue
					}

					fmt.Println("got lockup response...")
					gammCoins := lockupResponse.Coins
					gammShares := gammCoins.AmountOf("gamm/pool/" + strconv.Itoa(int(p.GetId())))

					exitedCoins, err := p.CalcExitPoolCoinsFromShares(sdk.Context{}, gammShares, math.LegacyZeroDec())
					// 11:
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to account for assets on zone %q: %w", chain, err)
						continue
					}
					fmt.Println("calculated exit shares...")

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
					fmt.Println("obtained relevant proofs...")
					msg[chainID] = chainMsg
					//nolint:staticcheck // SA1019 ignore this!
					break
				}
			}
		}
	}

	if len(errors) > 0 {
		return msg, assets, multierror.New(errors)
	}

	fmt.Printf("Msg: %+v\n", msg)
	fmt.Printf("Lockup Assets: %+v\n", assets)
	return msg, assets, nil
}
