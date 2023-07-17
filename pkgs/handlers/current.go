package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"

	"github.com/ingenuity-build/xcclookup/pkgs/claims"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

func GetCurrentHandler(
	ctx context.Context,
	cfg types.Config,
	connectionManager *types.CacheManager[prewards.ConnectionProtocolData],
	poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	osmosisParamsManager *types.CacheManager[prewards.OsmosisParamsProtocolData],
	umeeParamsManager *types.CacheManager[prewards.UmeeParamsProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		var err error
		chain := osmosisParamsManager.Get(ctx)[0].ChainID

		response := types.Response{Messages: make([]prewards.MsgSubmitClaim, 0), Assets: make(map[string][]types.Asset)}
		// osmosis
		connections := connectionManager.Get(ctx)

		if chain == "" {
			fmt.Fprintf(w, "Error: osmosis chain ID unset")
			return
		}

		_, assets, err := claims.OsmosisClaim(
			ctx,
			cfg,
			poolsManager,
			tokensManager,
			vars["address"],
			chain,
			0,
		)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		for chainID, asset := range assets {
			response.Assets[chainID] = []types.Asset{{Type: "osmosispool", Amount: asset}}
		}

		// umee claim
		_, assets, err = claims.UmeeClaim(
			ctx,
			cfg,
			tokensManager,
			vars["address"],
			umeeParamsManager.Get(ctx)[0].ChainID,
			0,
		)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		for chainID, asset := range assets {
			response.Assets[chainID] = []types.Asset{{Type: "liquid", Amount: asset}}
		}

		// liquid for all zones; config should hold osmosis chainid.
		for _, con := range connections {
			_, assets, err := claims.LiquidClaim(
				ctx,
				cfg,
				tokensManager,
				vars["address"],
				con,
				0,
			)
			if err != nil {
				fmt.Fprintf(w, "Error: %s", err)
				return
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
