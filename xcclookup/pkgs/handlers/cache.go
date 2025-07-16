package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

type CacheOutput struct {
	Connections    []prewards.ConnectionProtocolData
	OsmosisPools   []prewards.OsmosisPoolProtocolData
	OsmosisClPools []prewards.OsmosisClPoolProtocolData
	OsmosisParams  []prewards.OsmosisParamsProtocolData
	Tokens         []prewards.LiquidAllowedDenomProtocolData
}

func GetCacheHandler(
	ctx context.Context,
	_ types.Config,
	cacheMgr *types.CacheManager,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		connections, err := types.GetCache[prewards.ConnectionProtocolData](ctx, cacheMgr)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		osmosisPools, err := types.GetCache[prewards.OsmosisPoolProtocolData](ctx, cacheMgr)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		osmosisParams, err := types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		osmosisClPools, err := types.GetCache[prewards.OsmosisClPoolProtocolData](ctx, cacheMgr)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		liquidAllowedDenoms, err := types.GetCache[prewards.LiquidAllowedDenomProtocolData](ctx, cacheMgr)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		out := CacheOutput{
			Connections:    connections,
			OsmosisPools:   osmosisPools,
			OsmosisParams:  osmosisParams,
			OsmosisClPools: osmosisClPools,
			Tokens:         liquidAllowedDenoms,
		}
		jsonOut, err := json.Marshal(out)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprint(w, string(jsonOut))
	}
}
