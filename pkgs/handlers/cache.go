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
	Connections   []prewards.ConnectionProtocolData
	OsmosisPools  []prewards.OsmosisPoolProtocolData
	OsmosisParams []prewards.OsmosisParamsProtocolData
	Tokens        []prewards.LiquidAllowedDenomProtocolData
}

func GetCacheHandler(
	ctx context.Context,
	_ types.Config,
	connectionManager *types.CacheManager[prewards.ConnectionProtocolData],
	poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	osmosisParamsManager *types.CacheManager[prewards.OsmosisParamsProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		out := CacheOutput{Connections: connectionManager.Get(ctx), OsmosisPools: poolsManager.Get(ctx), OsmosisParams: osmosisParamsManager.Get(ctx), Tokens: tokensManager.Get(ctx)}

		jsonOut, err := json.Marshal(out)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprint(w, string(jsonOut))
	}
}
