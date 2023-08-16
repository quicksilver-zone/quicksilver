package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gorilla/mux"

	"github.com/ingenuity-build/multierror"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	"github.com/ingenuity-build/xcclookup/pkgs/claims"
	"github.com/ingenuity-build/xcclookup/pkgs/failsim"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

func GetEpochHandler(
	ctx context.Context,
	cfg types.Config,
	connectionManager *types.CacheManager[prewards.ConnectionProtocolData],
	osmosisPoolsManager *types.CacheManager[prewards.OsmosisPoolProtocolData],
	crescentPoolsManager *types.CacheManager[prewards.CrescentPoolProtocolData],
	osmosisParamsManager *types.CacheManager[prewards.OsmosisParamsProtocolData],
	umeeParamsManager *types.CacheManager[prewards.UmeeParamsProtocolData],
	crescentParamsManager *types.CacheManager[prewards.CrescentParamsProtocolData],
	tokensManager *types.CacheManager[prewards.LiquidAllowedDenomProtocolData],
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
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

		var messages map[string]prewards.MsgSubmitClaim
		var assets map[string]sdk.Coins
		var err error
		var connections []prewards.ConnectionProtocolData
		var height int64
		var chain string

		unfilteredConnections := connectionManager.Get(ctx)
		for _, ufc := range unfilteredConnections {
			if ufc.LastEpoch > 0 {
				connections = append(connections, ufc)
			}
		}

		fmt.Println("check config for osmosis chain id...")
		if len(osmosisParamsManager.Get(ctx)) == 0 {
			errors["OsmosisConfig"] = fmt.Errorf("osmosis params not set")
		} else {
			chain = osmosisParamsManager.Get(ctx)[0].ChainID
			if err = ValidateChainConfig("Osmosis", chain, failAt); err != nil {
				errors["OsmosisConfig"] = err
			} else {
				fmt.Println("check osmosis last epoch height...")
				if height, err = ValidateHeight(connections, chain, failAt); err != nil {
					errors["OsmosisHeight"] = err
				} else {
					fmt.Println("fetch osmosis claim...")
					messages, assets, err = claims.OsmosisClaim(ctx, cfg, osmosisPoolsManager, tokensManager, vars["address"], chain, height)
					if err != nil {
						errors["OsmosisClaim"] = err
					}
					response = UpdateResponse(response, messages, assets, "osmosispool")
				}
			}
		}

		// umee claim
		fmt.Println("check config for umee chain id...")
		if len(umeeParamsManager.Get(ctx)) == 0 {
			errors["UmeeConfig"] = fmt.Errorf("umee params not set")
		} else {
			chain = umeeParamsManager.Get(ctx)[0].ChainID
			if err = ValidateChainConfig("Umee", chain, failAt); err != nil {
				errors["UmeeConfig"] = err
			} else {
				fmt.Println("check umee last epoch height...")
				if height, err = ValidateHeight(connections, chain, failAt); err != nil {
					errors["UmeeHeight"] = err
				} else {
					fmt.Println("fetch umee claim...")
					messages, assets, err = claims.UmeeClaim(ctx, cfg, tokensManager, vars["address"], chain, height)
					if err != nil {
						errors["UmeeClaim"] = err
					}
					response = UpdateResponse(response, messages, assets, "liquid")
				}
			}
		}

		// crescent claim
		chain = crescentParamsManager.Get(ctx)[0].ChainID
		fmt.Println("check config for crescent chain id...")
		if len(crescentParamsManager.Get(ctx)) == 0 {
			errors["CrescentConfig"] = fmt.Errorf("crescent params not set")
		} else {
			chain = crescentParamsManager.Get(ctx)[0].ChainID
			if err = ValidateChainConfig("Crescent", chain, failAt); err != nil {
				errors["CrescentConfig"] = err
			} else {
				fmt.Println("check crescent last epoch height...")
				if height, err = ValidateHeight(connections, chain, failAt); err != nil {
					errors["UmeeHeight"] = err
				} else {
					fmt.Println("fetch crescent claim...")
					messages, assets, err = claims.CrescentClaim(ctx, cfg, crescentPoolsManager, tokensManager, vars["address"], chain, height)
					if err != nil {
						errors["CrescentClaim"] = err
					}
					response = UpdateResponse(response, messages, assets, "crescentpool")
				}
			}
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

func ValidateChainConfig(name, chain string, failAt map[uint8]struct{}) error {
	// simulate failure hook 0:
	if _, failHere := failAt[0]; failHere {
		chain = ""
	}
	if chain == "" {
		return fmt.Errorf("%s chain ID is not set", name)
	}
	return nil
}

func ValidateHeight(connections []prewards.ConnectionProtocolData, chain string, failAt map[uint8]struct{}) (int64, error) {
	var height int64
	for _, con := range connections {
		if con.ChainID == chain {
			height = con.LastEpoch
			break
		}
	}
	// simulate failure hook 1:
	if _, failHere := failAt[1]; failHere {
		height = 0
	}

	if height == 0 {
		return height, fmt.Errorf("fetched height is 0")
	}
	return height, nil
}

func UpdateResponse(response types.Response, messages map[string]prewards.MsgSubmitClaim, assets map[string]sdk.Coins, assetType string) types.Response {
	for _, message := range messages {
		response.Messages = append(response.Messages, message)
	}

	for chainID, asset := range assets {
		response.Assets[chainID] = append(response.Assets[chainID], types.Asset{Type: assetType, Amount: asset})
	}
	return response
}
