package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// HandleAddProtocolDataProposal is a handler for executing a passed add protocol data proposal
func HandleAddProtocolDataProposal(ctx sdk.Context, k Keeper, p *types.AddProtocolDataProposal) error {
	protocolData := NewProtocolData(p.Type, p.Protocol, p.Data)

	_, err := UnmarshalProtocolData(p.Type, p.Data)
	if err != nil {
		return err
	}

	k.SetProtocolData(ctx, p.Key, protocolData)

	return nil
}

func UnmarshalProtocolData(datatype string, data json.RawMessage) (IProtocolData, error) {
	switch datatype {
	case types.ClaimTypes[types.ClaimTypeOsmosisPool]:
		pd := types.OsmosisPoolProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return nil, err
		}
		return pd, nil
	case "connection":
		pd := ConnectionProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return nil, err
		}
		return pd, nil
	case types.ClaimTypes[types.ClaimTypeLiquidToken]:
		pd := types.LiquidAllowedDenomProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return nil, err
		}
		return pd, nil
	default:
		return nil, fmt.Errorf("unsupported protocol %s", datatype)
	}
}

type IProtocolData interface{}

type ConnectionProtocolData struct {
	ConnectionID string
	ChainID      string
}

var (
	_ IProtocolData = &ConnectionProtocolData{}
	_ IProtocolData = &types.OsmosisPoolProtocolData{}
	_ IProtocolData = &types.LiquidAllowedDenomProtocolData{}
)
