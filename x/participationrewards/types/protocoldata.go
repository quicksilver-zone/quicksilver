package types

import (
	"encoding/json"
	fmt "fmt"
	"reflect"
)

func UnmarshalProtocolData(datatype ProtocolDataType, data json.RawMessage) (IProtocolData, error) {
	switch datatype {
	case ProtocolDataOsmosisPool:
		{
			pd := OsmosisPoolProtocolData{}
			err := json.Unmarshal(data, &pd)
			if err != nil {
				return nil, err
			}
			var blank OsmosisPoolProtocolData
			if reflect.DeepEqual(pd, blank) {
				return nil, fmt.Errorf("unable to unmarshal osmosispool protocol data from empty JSON object")
			}
			return pd, nil
		}
	case ProtocolDataConnection:
		{
			pd := ConnectionProtocolData{}
			err := json.Unmarshal(data, &pd)
			if err != nil {
				return nil, err
			}
			var blank ConnectionProtocolData
			if reflect.DeepEqual(pd, blank) {
				return nil, fmt.Errorf("unable to unmarshal connection protocol data from empty JSON object")
			}
			return pd, nil
		}
	case ProtocolDataLiquidToken:
		{
			pd := LiquidAllowedDenomProtocolData{}
			err := json.Unmarshal(data, &pd)
			if err != nil {
				return nil, err
			}
			var blank LiquidAllowedDenomProtocolData
			if reflect.DeepEqual(pd, blank) {
				return nil, fmt.Errorf("unable to unmarshal liquid protocol data from empty JSON object")
			}
			return pd, nil
		}
	default:
		return nil, ErrUnknownProtocolDataType
	}
}

type IProtocolData interface{}

type ConnectionProtocolData struct {
	ConnectionID string
	ChainID      string
}

var (
	_ IProtocolData = &ConnectionProtocolData{}
	_ IProtocolData = &OsmosisPoolProtocolData{}
	_ IProtocolData = &LiquidAllowedDenomProtocolData{}
)
