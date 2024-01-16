package claims

import (
	"context"
	"fmt"
	"strconv"
	"time"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ingenuity-build/multierror"
	osmogamm "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm"
	osmolockup "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/lockup"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/ingenuity-build/xcclookup/pkgs/failsim"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

type poolMap map[string][]osmogamm.PoolI

func OsmosisClaim(
	ctx context.Context,
	cfg types.Config,
	poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
	zonesManager *types.CacheManager[icstypes.Zone],
	address string,
	chain string,
	height int64,
) (map[string]prewards.MsgSubmitClaim, map[string]sdk.Coins, error) {
	// simFailure hooks: 0-11
	simFailures := failsim.FailuresFromContext(ctx)
	failures := make(map[uint8]struct{})
	if OsmosisClaimFailures, ok := simFailures[1]; ok {
		fmt.Println("osmosis sim failures")
		failures = OsmosisClaimFailures
	}
	fmt.Println("simulate failures:", failures)

	var err error

	_, addrBytes, err := bech32.DecodeNoLimit(address)
	// 0:
	err = failsim.FailureHook(failures, 0, err, "failure decosing bech32 address")
	if err != nil {
		return nil, nil, err
	}
	osmoAddress, err := bech32.Encode("osmo", addrBytes)
	// 1:
	err = failsim.FailureHook(failures, 1, err, "failure encoding osmo address")
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("valid osmosis address encoding...")

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
	err = failsim.FailureHook(failures, 3, err, fmt.Sprintf("failure connecting to host %q", host))
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
		err = failsim.FailureHook(failures, 4, err, "ABCIQuery: GetLatestBlock")
		if err != nil {
			return nil, nil, err
		}
		fmt.Println("height is zero, get latest block height...")

		blockQueryResponse := tmservice.GetLatestBlockResponse{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &blockQueryResponse)
		// 5:
		err = failsim.FailureHook(failures, 5, err, "ABCIQuery: GetLatestBlockResponse")
		if err != nil {
			return nil, nil, err
		}
		emptyBlockResponse := tmservice.GetLatestBlockResponse{}
		if blockQueryResponse == emptyBlockResponse {
			err = fmt.Errorf("unable to query height from Osmosis chain")
		}
		// 6:
		err = failsim.FailureHook(failures, 6, err, "unable to query height from Osmosis chain")
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
		// 4: (hook numbering reset as this is a different execution path)
		err = failsim.FailureHook(failures, 4, err, "ABCIQuery: GetBlockByHeight")
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("height is %d, get block by height...\n", height)

		blockQueryResponse := tmservice.GetBlockByHeightResponse{}
		err = marshaler.Unmarshal(abciquery.Response.Value, &blockQueryResponse)
		// 5:
		err = failsim.FailureHook(failures, 5, err, "ABCIQuery: GetBlockByHeightResponse")
		if err != nil {
			return nil, nil, err
		}
		emptyBlockResponse := tmservice.GetBlockByHeightResponse{}
		// 6:
		if blockQueryResponse == emptyBlockResponse {
			err = fmt.Errorf("unable to query height from Osmosis chain")
		}
		err = failsim.FailureHook(failures, 6, err, "unable to query height from Osmosis chain")
		if err != nil {
			return nil, nil, err
		}
		if blockQueryResponse.Block != nil { //nolint:staticcheck // SA1019 ignore this!
			timestamp = blockQueryResponse.Block.Header.Time //nolint:staticcheck // SA1019 ignore this!
		} else {
			timestamp = blockQueryResponse.SdkBlock.Header.Time
		}
	}
	fmt.Println("got block timestamp...")

	query := osmolockup.AccountLockedPastTimeRequest{Owner: osmoAddress, Timestamp: timestamp}
	bytes := marshaler.MustMarshal(&query)
	fmt.Println("prepared account lockup query...")

	abciquery, err := client.ABCIQueryWithOptions(
		ctx,
		"/osmosis.lockup.Query/AccountLockedPastTime",
		bytes,
		rpcclient.ABCIQueryOptions{Height: height},
	)
	// 7:
	err = failsim.FailureHook(failures, 6, err, "ABCIQuery: AccountLockedPastTime")
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("executed abci query...")

	queryResponse := osmolockup.AccountLockedPastTimeResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &queryResponse)
	// 8:
	err = failsim.FailureHook(failures, 6, err, "ABCIQuery: AccountLockedPastTimeResponse")
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("unmarshalled query response...")

	// add GetFiltered to CacheManager, to allow filtered lookups on a single field == value
	tokens := GetTokenMap(tokensManager.Get(ctx), zonesManager.Get(ctx), chain, "")

	fmt.Println("got relevant tokens...")

	pools := poolMap{}

	for _, pool := range poolsManager.Get(ctx) {
		for _, denom := range pool.Denoms {
			if _, ok := pools[denom.ChainID]; !ok {
				pools[denom.ChainID] = make([]osmogamm.PoolI, 0)
			}
			poolData, err := pool.GetPool()
			if err != nil {
				return nil, nil, err
			}
			pools[denom.ChainID] = append(pools[denom.ChainID], poolData)
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

			bankQuery, err := client.ABCIQueryWithOptions(
				ctx, "/store/bank/key",
				lookupKey,
				rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
			)
			fmt.Println("Querying for value", "prefix", accountPrefix, "denom", poolCoinDenom) // debug?
			// 7:
			err = failsim.FailureHook(failures, 7, err, fmt.Sprintf("unable to query for value of denom %q on %q", poolCoinDenom, chain))
			if err != nil {
				return nil, nil, err
			}

			amount, err := bankkeeper.UnmarshalBalanceCompat(marshaler, bankQuery.Response.Value, poolCoinDenom)
			if err != nil {
				return nil, nil, err
			}
			// 8:
			err = failsim.FailureHook(failures, 8, err, fmt.Sprintf("ABCIQuery: value of denom %q on chain %q", poolCoinDenom, chain))
			if err != nil {
				return nil, nil, err
			}

			if amount.IsZero() {
				fmt.Println("no unbonded tokens found for denom: " + poolCoinDenom)
			} else {
				fmt.Printf("found assets in bank account for zone %q...\n", chainID)
				if _, ok := msg[chainID]; !ok {
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

				assets[chain] = assets[chain].Add(amount)

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
						ctx,
						"/store/lockup/key",
						lockupKey,
						rpcclient.ABCIQueryOptions{Height: abciquery.Response.Height, Prove: true},
					)
					// 9:
					err = failsim.FailureHook(failures, 9, err, "ABCIQuery: lockup")
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
					err = failsim.FailureHook(failures, 10, err, "ABCIQuery: lockupResponse")
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

					exitedCoins, err := p.CalcExitPoolCoinsFromShares(sdk.Context{}, gammShares, sdk.ZeroDec())
					// 11:
					err = failsim.FailureHook(failures, 11, err, "CalcExitPoolCoinsFromShares")
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

	// fmt.Printf("Msg: %+v\n", msg)
	// fmt.Printf("Lockup Assets: %+v\n", assets)
	return msg, assets, nil
}
