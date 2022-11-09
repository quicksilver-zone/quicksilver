package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ingenuity-build/xcclookup/pkgs/claims"
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
		vars := mux.Vars(req)
		var err error
		chain := osmosisParamsManager.Get()[0].ChainID

		response := types.Response{Messages: make([]prewards.MsgSubmitClaim, 0), Assets: make(map[string][]types.Asset)}
		// osmosis
		var height int64
		connections := connectionManager.Get()
		for _, con := range connections {
			if con.ChainID == osmosisParamsManager.Get()[0].ChainID {
				height = con.LastEpoch
				break
			}
		}
		if chain == "" {
			fmt.Fprintf(w, "Error: osmosis chain ID unset")
			return
		}

		if height == 0 {
			fmt.Fprintf(w, "Error: fetched height is 0")
			return
		}

		fmt.Println("1")
		messages, assets, err := claims.OsmosisClaim(cfg, poolsManager, tokensManager, vars["address"], chain, height)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		for _, message := range messages {
			response.Messages = append(response.Messages, message)
		}

		for chainID, asset := range assets {
			response.Assets[chainID] = []types.Asset{{Type: "osmosispool", Amount: asset}}
		}

		//liquid for all zones; config should hold osmosis chainid.
		for _, con := range connections {
			liquidMessages, assets, err := claims.LiquidClaim(cfg, tokensManager, vars["address"], con, con.LastEpoch)
			if err != nil {
				fmt.Fprintf(w, "Error: %s", err)
				return
			}
			for _, message := range liquidMessages {
				response.Messages = append(response.Messages, message)
			}
			for chainID, asset := range assets {
				response.Assets[chainID] = append(response.Assets[chainID], types.Asset{Type: "liquid", Amount: asset})
			}
		}

		jsonOut, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprint(w, string(jsonOut))
	}
}
