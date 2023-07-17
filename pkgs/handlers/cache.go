package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ingenuity-build/xcclookup/pkgs/types"

	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type CacheOutput struct {
	Connections   []prewards.ConnectionProtocolData
	OsmosisPools  []prewards.OsmosisPoolProtocolData
	OsmosisParams []prewards.OsmosisParamsProtocolData
	Tokens        []prewards.LiquidAllowedDenomProtocolData
}

func GetCacheHandler(
	Config types.Config,
	connectionManager *types.CacheManager[prewards.ConnectionProtocolData],
	poolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	osmosisParamsManager *types.CacheManager[prewards.OsmosisParamsProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		out := CacheOutput{Connections: connectionManager.Get(), OsmosisPools: poolsManager.Get(), OsmosisParams: osmosisParamsManager.Get(), Tokens: tokensManager.Get()}

		jsonOut, err := json.Marshal(out)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprint(w, string(jsonOut))
	}
}
