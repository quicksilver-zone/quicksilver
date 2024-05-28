package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/ingenuity-build/xcclookup/pkgs/claims"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func GetAssetsHandler(
	ctx context.Context,
	cfg types.Config,
	cacheMgr *types.CacheManager,
	heights map[string]int64,
	outputFunc func(http.ResponseWriter, *types.Response, map[string]error),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		errors := make(map[string]error)
		vars := mux.Vars(req)

		response := &types.Response{
			Messages: make([]prewards.MsgSubmitClaim, 0),
			Assets:   make(map[string][]types.Asset),
		}

		wg := sync.WaitGroup{}

		// ensure a neatly formatted JSON response
		defer outputFunc(w, response, errors)

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
			errors["MappedAddresses"] = err
		}

		fmt.Println("check config for osmosis chain id...")
		if len(types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr)) == 0 {
			errors["OsmosisConfig"] = fmt.Errorf("osmosis params not set")
		} else {
			chain = types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr)[0].ChainID

			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println("fetch osmosis claim for ", vars["address"])
				messages, assets, err := claims.OsmosisClaim(ctx, cfg, cacheMgr, vars["address"], vars["address"], chain, heights[chain])
				if err != nil {
					errors["OsmosisClaim"] = err
				}
				response.Update(messages, assets, "osmosispool")
			}()

			if mappedAddress, ok := mappedAddresses[chain]; ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					fmt.Println("fetch osmosis claim for mapped account", mappedAddress)
					messages, assets, err := claims.OsmosisClaim(ctx, cfg, cacheMgr, mappedAddress, vars["address"], chain, heights[chain])
					if err != nil {
						errors["OsmosisClaim"] = err
					}
					response.Update(messages, assets, "osmosispool")
				}()
			}
		}

		// umee claim
		fmt.Println("check config for umee chain id...")
		if len(types.GetCache[prewards.UmeeParamsProtocolData](ctx, cacheMgr)) == 0 {
			errors["UmeeConfig"] = fmt.Errorf("umee params not set")
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				chain = types.GetCache[prewards.UmeeParamsProtocolData](ctx, cacheMgr)[0].ChainID

				fmt.Println("fetch umee claim for", vars["address"])
				messages, assets, err := claims.UmeeClaim(ctx, cfg, cacheMgr, vars["address"], vars["address"], chain, heights[chain])
				if err != nil {
					errors["UmeeClaim"] = err
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
						errors["UmeeClaim"] = err
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
					errors[fmt.Sprintf("LiquidClaim:%s", con.ChainID)] = err
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
						errors[fmt.Sprintf("LiquidClaim:%s", con.ChainID)] = err
						return
					}
					response.Update(messages, assets, "liquid")
				}(con)
			}
		}

		wg.Wait()
	}
}
