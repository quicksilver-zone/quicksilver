package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ingenuity-build/xcclookup/internal/multierror"
	"github.com/ingenuity-build/xcclookup/pkgs/claims"
	"github.com/ingenuity-build/xcclookup/pkgs/failsim"
	"github.com/ingenuity-build/xcclookup/pkgs/types"

	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func GetEpochHandler(
	cfg types.Config,
	connectionManager *types.CacheManager[prewards.ConnectionProtocolData],
	poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	osmosisParamsManager *types.CacheManager[prewards.OsmosisParamsProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()

		// simFailure hooks: 0-1
		if useHooks := req.Header.Get("x-simulate-failure"); useHooks != "" {
			fctx, err := failsim.SetFailureContext(ctx, useHooks)
			if err != nil {
				str := fmt.Sprintf("SetFailureContext error: %s\n", err)
				fmt.Println(str)
				fmt.Fprint(w, str)
				return
			}
			ctx = fctx
		}

		simFailures := failsim.FailuresFromContext(ctx)
		failAt := make(map[uint8]struct{})
		if failures, ok := simFailures[0]; ok {
			failAt = failures
		}

		errors := make(map[string]error)
		vars := mux.Vars(req)
		chain := osmosisParamsManager.Get()[0].ChainID
		response := types.Response{
			Messages: make([]prewards.MsgSubmitClaim, 0),
			Assets:   make(map[string][]types.Asset),
		}

		// ensure a neatly formatted JSON response
		defer func() {
			fmt.Println("check for errors...")
			if len(errors) > 0 {
				fmt.Printf("found %d error(s)\n", len(errors))
				response.Errors = multierror.New(errors)
				fmt.Println(response.Errors)
			}

			fmt.Println("prepare JSON...")
			jsonOut, err := json.Marshal(response)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Fprintf(w, "Error: %s", err)
				return
			}
			fmt.Printf("%s\n", jsonOut)
			fmt.Fprint(w, string(jsonOut))
		}()

		fmt.Println("check config for osmosis chain id...")
		// simulate failure hook 0:
		if _, failHere := failAt[0]; failHere {
			chain = ""
		}
		if chain == "" {
			errors["config"] = fmt.Errorf("osmosis chain ID is not set")
			return
		}

		fmt.Println("check osmosis last epoch height...")
		var height int64
		unfilteredConnections := connectionManager.Get()
		connections := []prewards.ConnectionProtocolData{}
		for _, ufc := range unfilteredConnections {
			if ufc.LastEpoch > 0 {
				connections = append(connections, ufc)
			}
		}

		for _, con := range connections {
			if con.ChainID == osmosisParamsManager.Get()[0].ChainID {
				height = con.LastEpoch
				break
			}
		}
		// simulate failure hook 1:
		if _, failHere := failAt[1]; failHere {
			height = 0
		}
		if height == 0 {
			errors["height"] = fmt.Errorf("fetched height is 0")
			return
		}

		fmt.Println("fetch osmosis claim...")
		messages, assets, err := claims.OsmosisClaim(ctx, cfg, poolsManager, tokensManager, vars["address"], chain, height)
		if err != nil {
			errors["OsmosisClaim"] = err
		}

		for _, message := range messages {
			response.Messages = append(response.Messages, message)
		}

		for chainID, asset := range assets {
			response.Assets[chainID] = []types.Asset{{Type: "osmosispool", Amount: asset}}
		}

		// liquid for all zones; config should hold osmosis chainid.
		fmt.Println("fetch liquid claims...")
		for _, con := range connections {
			liquidMessages, assets, err := claims.LiquidClaim(ctx, cfg, tokensManager, vars["address"], con, con.LastEpoch)
			if err != nil {
				errors[fmt.Sprintf("LiquidClaim:%s", con.ChainID)] = err
				continue
			}
			for _, message := range liquidMessages {
				response.Messages = append(response.Messages, message)
			}
			for chainID, asset := range assets {
				response.Assets[chainID] = append(response.Assets[chainID], types.Asset{Type: "liquid", Amount: asset})
			}
		}
	}
}
