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

func GetAirdropHandler(
	ctx context.Context,
	cfg types.Config,
	connectionManager *types.CacheManager[prewards.ConnectionProtocolData],
	poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	osmosisParamsManager *types.CacheManager[prewards.OsmosisParamsProtocolData],
	umeeParamsManager *types.CacheManager[prewards.UmeeParamsProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
	zonesManager *types.CacheManager[icstypes.Zone],
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

		messages, assets, err := claims.OsmosisClaim(
			ctx,
			cfg,
			poolsManager,
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

		for _, message := range messages {
			response.Messages = append(response.Messages, message)
		}

		for chainID, asset := range assets {
			response.Assets[chainID] = []types.Asset{{Type: "osmosispool", Amount: asset}}
		}

		// umee
		messages, assets, err = claims.UmeeClaim(
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

		for _, message := range messages {
			response.Messages = append(response.Messages, message)
		}

		for chainID, asset := range assets {
			response.Assets[chainID] = []types.Asset{{Type: "liquid", Amount: asset}}
		}

		// liquid for all zones; config should hold osmosis chainid.
		for _, con := range connections {
			liquidMessages, assets, err := claims.LiquidClaim(
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
			for _, message := range liquidMessages {
				response.Messages = append(response.Messages, message)
			}
			for chainID, asset := range assets {
				response.Assets[chainID] = append(response.Assets[chainID], types.Asset{Type: "liquid", Amount: asset})
			}
		}

		// add update client message for each chain.

		jsonOut, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprint(w, string(jsonOut))
	}
}
