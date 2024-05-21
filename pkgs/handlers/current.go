package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/ingenuity-build/xcclookup/pkgs/claims"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

func GetCurrentHandler(
	ctx context.Context,
	cfg types.Config,
	cacheMgr *types.CacheManager,
	// connectionManager *types.Cache[prewards.ConnectionProtocolData],
	// osmosisPoolsManager *types.Cache[prewards.OsmosisPoolProtocolData],
	// osmosisClPoolsManager *types.Cache[prewards.OsmosisClPoolProtocolData],
	// osmosisParamsManager *types.Cache[prewards.OsmosisParamsProtocolData],
	// umeeParamsManager *types.Cache[prewards.UmeeParamsProtocolData],
	// tokensManager *types.Cache[prewards.LiquidAllowedDenomProtocolData],
	// zonesManager *types.Cache[icstypes.Zone],
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		var err error
		response := types.Response{Messages: make([]prewards.MsgSubmitClaim, 0), Assets: make(map[string][]types.Asset)}
		if len(types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr)) > 0 {
			chain := types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr)[0].ChainID

			// osmosis

			if chain == "" {
				fmt.Fprintf(w, "Error: osmosis chain ID unset")
				return
			}

			_, assets, err := claims.OsmosisClaim(
				ctx,
				cfg,
				cacheMgr,
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
		}

		if len(types.GetCache[prewards.UmeeParamsProtocolData](ctx, cacheMgr)) > 0 {
			_, assets, err := claims.UmeeClaim(
				ctx,
				cfg,
				cacheMgr,
				vars["address"],
				types.GetCache[prewards.UmeeParamsProtocolData](ctx, cacheMgr)[0].ChainID,
				0,
			)
			if err != nil {
				fmt.Fprintf(w, "Error: %s", err)
				return
			}

			for chainID, asset := range assets {
				response.Assets[chainID] = append(response.Assets[chainID], types.Asset{Type: "umeepool", Amount: asset})
			}
		}

		connections := types.GetCache[prewards.ConnectionProtocolData](ctx, cacheMgr)
		// liquid for all zones; config should hold osmosis chainid.
		for _, con := range connections {
			_, assets, err := claims.LiquidClaim(
				ctx,
				cfg,
				cacheMgr,
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
