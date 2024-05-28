package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ingenuity-build/xcclookup/pkgs/types"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
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
		out := CacheOutput{
			Connections:   types.GetCache[prewards.ConnectionProtocolData](ctx, cacheMgr),
			OsmosisPools:  types.GetCache[prewards.OsmosisPoolProtocolData](ctx, cacheMgr),
			OsmosisParams: types.GetCache[prewards.OsmosisParamsProtocolData](ctx, cacheMgr),
			Tokens:        types.GetCache[prewards.LiquidAllowedDenomProtocolData](ctx, cacheMgr),
		}
		jsonOut, err := json.Marshal(out)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprint(w, string(jsonOut))
	}
}
