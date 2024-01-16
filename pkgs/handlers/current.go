package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	"github.com/gorilla/mux"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/ingenuity-build/xcclookup/pkgs/claims"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

func GetCurrentHandler(
	ctx context.Context,
	cfg types.Config,
	connectionManager *types.CacheManager[prewards.ConnectionProtocolData],
	osmosisPoolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	crescentPoolsManager *types.CacheManager[prewards.CrescentPoolProtocolData],
	osmosisParamsManager *types.CacheManager[prewards.OsmosisParamsProtocolData],
	umeeParamsManager *types.CacheManager[prewards.UmeeParamsProtocolData],
	crescentParamsManager *types.CacheManager[prewards.CrescentParamsProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
	zonesManager *types.CacheManager[icstypes.Zone],
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		var err error
		response := types.Response{Messages: make([]prewards.MsgSubmitClaim, 0), Assets: make(map[string][]types.Asset)}
		if len(osmosisParamsManager.Get(ctx)) > 0 {
			chain := osmosisParamsManager.Get(ctx)[0].ChainID

			// osmosis

			if chain == "" {
				fmt.Fprintf(w, "Error: osmosis chain ID unset")
				return
			}

			_, assets, err := claims.OsmosisClaim(
				ctx,
				cfg,
				osmosisPoolsManager,
				tokensManager,
				zonesManager,
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

		if len(umeeParamsManager.Get(ctx)) > 0 {
			_, assets, err := claims.UmeeClaim(
				ctx,
				cfg,
				tokensManager,
				zonesManager,
				vars["address"],
				umeeParamsManager.Get(ctx)[0].ChainID,
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

		if len(crescentParamsManager.Get(ctx)) > 0 {
			// crescent claim
			_, assets, err := claims.CrescentClaim(
				ctx,
				cfg,
				crescentPoolsManager,
				tokensManager,
				zonesManager,
				vars["address"],
				crescentParamsManager.Get(ctx)[0].ChainID,
				0,
			)
			if err != nil {
				fmt.Fprintf(w, "Error: %s", err)
				return
			}

			for chainID, asset := range assets {
				response.Assets[chainID] = append(response.Assets[chainID], types.Asset{Type: "crescentpool", Amount: asset})
			}

		}

		connections := connectionManager.Get(ctx)
		// liquid for all zones; config should hold osmosis chainid.
		for _, con := range connections {
			_, assets, err := claims.LiquidClaim(
				ctx,
				cfg,
				tokensManager,
				zonesManager,
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
