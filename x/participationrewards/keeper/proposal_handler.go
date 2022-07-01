package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// HandleAddProtocolDataProposal is a handler for executing a passed add protocol data proposal
func HandleAddProtocolDataProposal(ctx sdk.Context, k Keeper, p *types.AddProtocolDataProposal) error {
	protocolData := k.NewProtocolData(ctx, p.Type, p.Protocol, p.Data)

	_, err := UnmarshalProtocolData(p.Type, p.Data)
	if err != nil {
		return err
	}

	k.SetProtocolData(ctx, p.Key, protocolData)

	return nil
}

func UnmarshalProtocolData(datatype string, data json.RawMessage) (IProtocolData, error) {
	switch datatype {
	case "osmosispool":
		pd := OsmosisPoolProtocolData{}
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
	case "liquidalloweddenoms":
		pd := LiquidAllowedDenomProtocolData{}
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
	ConnectionId string
	ChainId      string
}

type OsmosisPoolProtocolData struct {
	PoolId            uint64
	PoolName          string
	IbcToken          string
	LocalToken        string
	IbcTokenBalance   int64
	LocalTokenBalance int64
}

type LiquidAllowedDenomProtocolData struct {
	ChainId    string
	Denom      string
	LocalDenom string
}

var (
	_ IProtocolData = &ConnectionProtocolData{}
	_ IProtocolData = &OsmosisPoolProtocolData{}
	_ IProtocolData = &LiquidAllowedDenomProtocolData{}
)
