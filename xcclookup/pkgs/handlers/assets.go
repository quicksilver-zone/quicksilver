package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/claims"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

func GetAssetsHandler(
	ctx context.Context,
	cfg types.Config,
	cacheMgr *types.CacheManager,
	heights map[string]int64,
	outputFunc func(http.ResponseWriter, *types.Response, map[string]error),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		errs := make(map[string]error)
		vars := mux.Vars(req)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		response := &types.Response{
			Messages: make([]prewards.MsgSubmitClaim, 0),
			Assets:   make(map[string][]types.Asset),
		}

		wg := sync.WaitGroup{}

		// ensure a neatly formatted JSON response
		defer outputFunc(w, response, errs)

		var connections []prewards.ConnectionProtocolData
		var chain string

		unfilteredConnections := types.GetCache[prewards.ConnectionProtocolData](ctx, cacheMgr)
		for _, ufc := range unfilteredConnections {
			if ufc.LastEpoch > 0 {
				connections = append(connections, ufc)
			}
		}

		mappedAddresses, err := types.GetMappedAddresses(ctx, vars["address"], unfilteredConnections, &cfg)
		if err != nil {
			errs["MappedAddresses"] = err
		}

		fmt.Println("check config for osmosis chain id...")
		if len(types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr)) == 0 {
			errs["OsmosisConfig"] = errors.New("osmosis params not set")
		} else {
			chain = types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr)[0].ChainID

			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println("fetch osmosis claim for ", vars["address"])
				result := claims.OsmosisClaim(ctx, cfg, cacheMgr, vars["address"], vars["address"], chain, heights[chain]) // return OsmosisResult{OsmosisPool{msg, assets, err}, OsmosisClPool{msg, assets, err}}
				if result.Err != nil {
					errs["OsmosisClaim"] = result.Err
				}
				if result.OsmosisPool.Err != nil {
					errs["OsmosisPoolClaim"] = result.OsmosisPool.Err
				}
				if result.OsmosisClPool.Err != nil {
					errs["OsmosisClPoolClaim"] = result.OsmosisClPool.Err
				}
				if result.OsmosisPool.Msg != nil {
					response.Update(result.OsmosisPool.Msg, result.OsmosisPool.Assets, "osmosispool")
				}
				if result.OsmosisClPool.Msg != nil {
					response.Update(result.OsmosisClPool.Msg, result.OsmosisClPool.Assets, "osmosisclpool")
				}
			}()

			if mappedAddress, ok := mappedAddresses[chain]; ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					fmt.Println("fetch osmosis claim for mapped account", mappedAddress)
					result := claims.OsmosisClaim(ctx, cfg, cacheMgr, mappedAddress, vars["address"], chain, heights[chain])
					if result.Err != nil {
						errs["OsmosisClaim"] = result.Err
					}
					if result.OsmosisPool.Err != nil {
						errs["OsmosisPoolClaim"] = result.OsmosisPool.Err
					}
					if result.OsmosisClPool.Err != nil {
						errs["OsmosisClPoolClaim"] = result.OsmosisClPool.Err
					}
					if result.OsmosisPool.Msg != nil {
						response.Update(result.OsmosisPool.Msg, result.OsmosisPool.Assets, "osmosispool")
					}
					if result.OsmosisClPool.Msg != nil {
						response.Update(result.OsmosisClPool.Msg, result.OsmosisClPool.Assets, "osmosisclpool")
					}
				}()
			}
		}

		// umee claim
		fmt.Println("check config for umee chain id...")
		if len(types.GetCache[prewards.UmeeParamsProtocolData](ctx, cacheMgr)) == 0 {
			errs["UmeeConfig"] = errors.New("umee params not set")
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				chain = types.GetCache[prewards.UmeeParamsProtocolData](ctx, cacheMgr)[0].ChainID

				fmt.Println("fetch umee claim for", vars["address"])
				messages, assets, err := claims.UmeeClaim(ctx, cfg, cacheMgr, vars["address"], vars["address"], chain, heights[chain])
				if err != nil {
					errs["UmeeClaim"] = err
				}
				response.Update(messages, assets, "umee")
			}()
			if mappedAddress, ok := mappedAddresses[chain]; ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					fmt.Println("fetch umee claim for mapped account", mappedAddress)
					messages, assets, err := claims.UmeeClaim(ctx, cfg, cacheMgr, mappedAddress, vars["address"], chain, heights[chain])
					if err != nil {
						errs["UmeeClaim"] = err
					}
					response.Update(messages, assets, "umee")
				}()
			}
		}

		// liquid for all zones; config should hold osmosis chainid.
		fmt.Println("fetch liquid claims...")
		for _, con := range connections {
			wg.Add(1)
			go func(con prewards.ConnectionProtocolData) {
				defer wg.Done()
				messages, assets, err := claims.LiquidClaim(ctx, cfg, cacheMgr, vars["address"], vars["address"], con, heights[con.ChainID])
				if err != nil {
					errs[fmt.Sprintf("LiquidClaim:%s", con.ChainID)] = err
					return
				}
				response.Update(messages, assets, "liquid")
			}(con)

			if mappedAddress, ok := mappedAddresses[con.ChainID]; ok {
				wg.Add(1)
				go func(con prewards.ConnectionProtocolData) {
					defer wg.Done()
					messages, assets, err := claims.LiquidClaim(ctx, cfg, cacheMgr, mappedAddress, vars["address"], con, heights[con.ChainID])
					if err != nil {
						errs[fmt.Sprintf("LiquidClaim:%s", con.ChainID)] = err
						return
					}
					response.Update(messages, assets, "liquid")
				}(con)
			}
		}

		wg.Wait()
	}
}
