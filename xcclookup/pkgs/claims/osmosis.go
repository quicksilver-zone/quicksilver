package claims

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	tmhttp "github.com/cometbft/cometbft/rpc/client/http"

	osmotypes "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types"
	osmoclmodelquery "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/client/queryproto"
	osmoclmodel "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/model"
	osmocl "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/types"
	osmogamm "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/types"
	osmolockup "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/lockup/types"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

var (
	// these values are taken from https://github.com/osmosis-labs/osmosis/blob/c6ed71fd59b5d899d8798629e9b9f89e3b0c24e8/x/lockup/types/keys.go#L68;
	// they were inadvertently removed from quicksilver in v1.8.0.
	KeyPrefixPeriodLock = []byte{0x02}
	KeyIndexSeparator   = []byte{0xFF}
)

type (
	poolMap     map[string][]osmogamm.CFMMPoolI
	clPoolMap   map[string][]osmocl.ConcentratedPoolExtension
	OsmosisPool struct {
		Msg    map[string]prewards.MsgSubmitClaim
		Assets map[string]sdk.Coins
		Err    error
	}
	OsmosisClPool struct {
		Msg    map[string]prewards.MsgSubmitClaim
		Assets map[string]sdk.Coins
		Err    error
	}
	OsmosisResult struct {
		Err           error
		OsmosisPool   OsmosisPool
		OsmosisClPool OsmosisClPool
	}
)

func OsmosisClaim(
	ctx context.Context,
	cfg types.Config,
	cacheMgr *types.CacheManager,
	address string,
	submitAddress string,
	chain string,
	height int64,
) OsmosisResult {
	var err error

	addrBytes, err := addressutils.AccAddressFromBech32(address, "")
	if err != nil {
		return OsmosisResult{Err: err}
	}

	osmoAddress, err := addressutils.EncodeAddressToBech32("osmo", addrBytes)
	if err != nil {
		return OsmosisResult{Err: err}
	}

	host, ok := cfg.Chains[chain]
	if !ok {
		err = fmt.Errorf("no endpoint in config for %s", chain)
		return OsmosisResult{Err: err}
	}

	client, err := types.NewRPCClient(host, time.Duration(cfg.Timeout)*time.Second)
	if err != nil {
		return OsmosisResult{Err: err}
	}

	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	cmtypes.RegisterInterfaces(interfaceRegistry)
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
			return OsmosisResult{Err: err}
		}

		blockQueryResponse := tmservice.GetLatestBlockResponse{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &blockQueryResponse)
		if err != nil {
			return OsmosisResult{Err: err}
		}

		emptyBlockResponse := tmservice.GetLatestBlockResponse{}
		if blockQueryResponse == emptyBlockResponse {
			err = errors.New("unable to query height from Osmosis chain")
		}
		if err != nil {
			return OsmosisResult{Err: err}
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
			return OsmosisResult{Err: err}
		}

		blockQueryResponse := tmservice.GetBlockByHeightResponse{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &blockQueryResponse)
		if err != nil {
			return OsmosisResult{Err: err}
		}

		emptyBlockResponse := tmservice.GetBlockByHeightResponse{}
		if blockQueryResponse == emptyBlockResponse {
			err = errors.New("unable to query height from Osmosis chain")
		}
		if err != nil {
			return OsmosisResult{Err: err}
		}
		if blockQueryResponse.Block != nil { //nolint:staticcheck // SA1019 ignore this!
			timestamp = blockQueryResponse.Block.Header.Time //nolint:staticcheck // SA1019 ignore this!
		} else {
			timestamp = blockQueryResponse.SdkBlock.Header.Time
		}
	}

	ignores := cfg.Ignore.GetIgnoresForType(types.IgnoreTypeLiquid)
	// add GetFiltered to CacheManager, to allow filtered lookups on a single field == value
	laCache, err := types.GetCache[prewards.LiquidAllowedDenomProtocolData](ctx, cacheMgr)
	if err != nil {
		return OsmosisResult{Err: err}
	}
	zoneCache, err := types.GetCache[icstypes.Zone](ctx, cacheMgr)
	if err != nil {
		return OsmosisResult{Err: err}
	}
	tokens := GetTokenMap(laCache, zoneCache, chain, "", ignores)

	result := OsmosisResult{}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		msgs, assets, err := GetOsmosisClaim(ctx, cfg, cacheMgr, client, marshaler, addrBytes, osmoAddress, submitAddress, chain, tokens, height, timestamp)
		if err != nil {
			result.OsmosisPool.Err = err
			return
		}
		result.OsmosisPool.Msg = msgs
		result.OsmosisPool.Assets = assets
	}()

	go func() {
		defer wg.Done()
		clmsgs, classets, err := GetOsmosisClClaim(ctx, cfg, cacheMgr, client, marshaler, addrBytes, osmoAddress, submitAddress, chain, tokens, height)
		if err != nil {
			result.OsmosisClPool.Err = err
			return
		}
		result.OsmosisClPool.Msg = clmsgs
		result.OsmosisClPool.Assets = classets
	}()
	wg.Wait()

	return result
}

func GetOsmosisClaim(ctx context.Context, cfg types.Config, cacheMgr *types.CacheManager, client *tmhttp.HTTP, marshaler *codec.ProtoCodec, addrBytes []byte, osmoAddress, submitAddress, chain string, tokens map[string]TokenTuple, height int64, timestamp time.Time) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	// fetch timestamp of block
	var errors map[string]error
	query := osmolockup.AccountLockedPastTimeRequest{Owner: osmoAddress, Timestamp: timestamp}
	bytes := marshaler.MustMarshal(&query)

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

	ignores := cfg.Ignore.GetIgnoresForType(types.IgnoreTypeOsmosisPool)

	pools := poolMap{}
	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}

	poolsCache, err := types.GetCache[prewards.OsmosisPoolProtocolData](ctx, cacheMgr)
	if err != nil {
		return nil, nil, err
	}

	for _, pool := range poolsCache {
		if ignores.Contains(strconv.FormatUint(pool.PoolID, 10)) {
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

	keyCache := make(map[string]bool)

	for chainID, chainPools := range pools { // iterate over chains - are we doing all chains?
		for _, p := range chainPools { // iterate over the pools for this chain
			// fetching unbonded gamm tokens from account

			poolCoinDenom := fmt.Sprintf("gamm/pool/%d", p.GetId())

			lookupKey := banktypes.CreateAccountBalancesPrefix(addrBytes)
			lookupKey = append(lookupKey, []byte(poolCoinDenom)...)

			if keyCache[string(lookupKey)] {
				continue
			}

			bankQuery, err := client.ABCIQueryWithOptions(
				ctx, "/store/bank/key",
				lookupKey,
				rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
			)
			if err != nil {
				return nil, nil, err
			}

			keyCache[string(lookupKey)] = true

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

				if _, ok := msg[chainID]; !ok {
					msg[chainID] = prewards.MsgSubmitClaim{
						UserAddress: submitAddress,
						Zone:        chainID,
						SrcZone:     chain,
						ClaimType:   cmtypes.ClaimTypeOsmosisPool,
						Proofs:      make([]*cmtypes.Proof, 0),
					}
				}

				for _, exitToken := range exitedCoins {
					tuple, ok := tokens[exitToken.Denom]
					if ok {
						exitToken.Denom = tuple.denom
						assets[chain] = assets[chain].Add(exitToken)
						break
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

					lockupKey := KeyPrefixPeriodLock
					lockupKey = append(lockupKey, KeyIndexSeparator...)
					lockupKey = append(lockupKey, sdk.Uint64ToBigEndian(lockup.ID)...)

					abciquery, err := client.ABCIQueryWithOptions(
						ctx,
						"/store/lockup/key",
						lockupKey,
						rpcclient.ABCIQueryOptions{Height: height, Prove: true},
					)
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to account for assets on zone %q: %w", chain, err)
						continue
					}

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

					gammCoins := lockupResponse.Coins

					gammShares := gammCoins.AmountOf("gamm/pool/" + strconv.FormatUint(p.GetId(), 10))

					exitedCoins, err := p.CalcExitPoolCoinsFromShares(sdk.Context{}, gammShares, math.LegacyZeroDec())
					// 11:
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to account for assets on zone %q: %w", chain, err)
						continue
					}

					for _, exitToken := range exitedCoins {
						tuple, ok := tokens[exitToken.Denom]
						if ok {
							exitToken.Denom = tuple.denom
							assets[chain] = assets[chain].Add(exitToken)
							break
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
					msg[chainID] = chainMsg
					//nolint:staticcheck // SA1019 ignore this!
					break
				}
			}
		}
	}
	return msg, assets, nil
}

func GetOsmosisClClaim(ctx context.Context, cfg types.Config, cacheMgr *types.CacheManager, client *tmhttp.HTTP, marshaler *codec.ProtoCodec, addrBytes []byte, osmoAddress, submitAddress, chain string, tokens map[string]TokenTuple, height int64) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	errors := make(map[string]error)
	ignores := cfg.Ignore.GetIgnoresForType(types.IgnoreTypeOsmosisCLPool)

	clpools := clPoolMap{}
	msg := map[string]prewards.MsgSubmitClaim{}
	assets := map[string]sdk.Coins{}

	poolsCache, err := types.GetCache[prewards.OsmosisClPoolProtocolData](ctx, cacheMgr)
	if err != nil {
		return nil, nil, err
	}

	for _, clpool := range poolsCache {
		if ignores.Contains(strconv.FormatUint(clpool.PoolID, 10)) {
			continue
		}
		if clpool.IsIncentivized {
			for _, denom := range clpool.Denoms {
				if _, ok := clpools[denom.ChainID]; !ok {
					clpools[denom.ChainID] = make([]osmocl.ConcentratedPoolExtension, 0)
				}
				poolData, err := clpool.GetPool()
				if err != nil {
					return nil, nil, err
				}
				for _, pool := range clpools[denom.ChainID] {
					if pool.GetId() == poolData.GetId() {
						continue
					}
				}
				clpools[denom.ChainID] = append(clpools[denom.ChainID], poolData)

			}
		}
	}
	clquery := osmoclmodelquery.UserPositionsRequest{Address: osmoAddress}
	clbytes := marshaler.MustMarshal(&clquery)
	clabciquery, err := client.ABCIQueryWithOptions(
		ctx,
		"/osmosis.concentratedliquidity.v1beta1.Query/UserPositions",
		clbytes,
		rpcclient.ABCIQueryOptions{Height: height},
	)
	if err != nil {
		return nil, nil, err
	}

	clqueryResponse := osmoclmodelquery.UserPositionsResponse{}
	err = marshaler.Unmarshal(clabciquery.Response.Value, &clqueryResponse)
	if err != nil {
		return nil, nil, err
	}

	keyCache := make(map[string]bool)

	for chainID, chainclPools := range clpools { // iterate over chains - are we doing all chains?
		for _, p := range chainclPools { // iterate over the pools for this chain
			for _, position := range clqueryResponse.Positions { // for each position in response
				if position.Position.PoolId == p.GetId() {

					// perhaps counter intuitively, we want to group messages by chainID - the chain we are claiming for
					// and assets by chain - the chain on which they are located.
					if _, ok := msg[chainID]; !ok {
						msg[chainID] = prewards.MsgSubmitClaim{
							UserAddress: submitAddress,
							Zone:        chainID,
							SrcZone:     chain,
							ClaimType:   cmtypes.ClaimTypeOsmosisCLPool,
							Proofs:      make([]*cmtypes.Proof, 0),
						}
					}

					if _, ok := assets[chain]; !ok {
						assets[chain] = sdk.Coins{}
					}

					positionKey := osmocl.KeyPositionId(position.Position.PositionId)

					if keyCache[string(positionKey)] {
						continue
					}

					abciquery, err := client.ABCIQueryWithOptions(
						ctx,
						"/store/concentratedliquidity/key",
						positionKey,
						rpcclient.ABCIQueryOptions{Height: height, Prove: true},
					)
					if err != nil {
						errors[chain] = fmt.Errorf("unable to query position for pool %d on chain %q: %w", p.GetId(), chain, err)
						continue
					}

					if abciquery.Response.Value == nil {
						continue
					}

					positionResponse := osmoclmodel.Position{}
					err = marshaler.Unmarshal(abciquery.Response.Value, &positionResponse)
					// 10:
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to unmarshal position for pool %d on chain %q: %w", p.GetId(), chain, err)
						continue
					}

					asset0, asset1, err := osmotypes.CalculateUnderlyingAssetsFromPosition(sdk.Context{}, positionResponse, p)
					if err != nil {
						if errors == nil {
							errors = make(map[string]error)
						}
						errors[chain] = fmt.Errorf("unable to calculate underlying assets for position %d on chain %q: %w", p.GetId(), chain, err)
						continue
					}

					for _, exitToken := range sdk.NewCoins(asset0, asset1) {
						tuple, ok := tokens[exitToken.Denom]
						if ok {
							exitToken.Denom = tuple.denom
							assets[chain] = assets[chain].Add(exitToken)
							break
						}
					}

					chainMsg := msg[chainID]

					proof := cmtypes.Proof{
						Data:      abciquery.Response.Value,
						Key:       abciquery.Response.Key,
						ProofOps:  abciquery.Response.ProofOps,
						Height:    abciquery.Response.Height,
						ProofType: "concentratedliquidity", // module name of proof.
					}

					chainMsg.Proofs = append(chainMsg.Proofs, &proof)
					keyCache[string(positionKey)] = true
					msg[chainID] = chainMsg
					//nolint:staticcheck // SA1019 ignore this!
					break
				}
			}
		}
	}
	return msg, assets, nil
}
