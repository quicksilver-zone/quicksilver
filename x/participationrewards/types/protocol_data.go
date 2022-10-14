package types

import (
	"encoding/json"
	"errors"
	fmt "fmt"
	"reflect"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func UnmarshalProtocolData(datatype ProtocolDataType, data json.RawMessage) (ProtocolDataI, error) {
	switch datatype {
	case ProtocolDataTypeConnection:
		{
			pd := ConnectionProtocolData{}
			err := json.Unmarshal(data, &pd)
			if err != nil {
				return nil, err
			}
			var blank ConnectionProtocolData
			if reflect.DeepEqual(pd, blank) {
				return nil, errors.New("unable to unmarshal connection protocol data from empty JSON object")
			}
			return pd, nil
		}
	case ProtocolDataTypeOsmosisParams:
		{
			pd := OsmosisParamsProtocolData{}
			err := json.Unmarshal(data, &pd)
			if err != nil {
				return nil, err
			}
			var blank OsmosisParamsProtocolData
			if reflect.DeepEqual(pd, blank) {
				return nil, fmt.Errorf("unable to unmarshal osmosisparams protocol data from empty JSON object")
			}
			return pd, nil
		}
	case ProtocolDataTypeLiquidToken:
		{
			pd := LiquidAllowedDenomProtocolData{}
			err := json.Unmarshal(data, &pd)
			if err != nil {
				return nil, err
			}
			var blank LiquidAllowedDenomProtocolData
			if reflect.DeepEqual(pd, blank) {
				return nil, errors.New("unable to unmarshal liquid protocol data from empty JSON object")
			}
			return pd, nil
		}
	case ProtocolDataTypeOsmosisPool:
		{
			pd := OsmosisPoolProtocolData{}
			err := json.Unmarshal(data, &pd)
			if err != nil {
				// TODO: PoolData fails here... needs custom Marshal/Unmarshal on OsmosisPoolProtocolData
				return nil, err
			}
			var blank OsmosisPoolProtocolData
			if reflect.DeepEqual(pd, blank) {
				return nil, fmt.Errorf("unable to unmarshal osmosispool protocol data from empty JSON object")
			}
			return pd, nil
		}
	default:
		return nil, ErrUnknownProtocolDataType
	}
}

type ProtocolDataI interface {
	ValidateBasic() error
}

type ConnectionProtocolData struct {
	ConnectionID string
	ChainID      string
	LastEpoch    int64
}

// ValidateBasic satisfies ProtocolDataI and validates basic stateless data.
func (cpd ConnectionProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	if len(cpd.ConnectionID) == 0 {
		errors["ConnectionID"] = ErrUndefinedAttribute
	}

	if len(cpd.ChainID) == 0 {
		errors["ChainID"] = ErrUndefinedAttribute
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

var (
	_ ProtocolDataI = &ConnectionProtocolData{}
	_ ProtocolDataI = &OsmosisPoolProtocolData{}
	_ ProtocolDataI = &OsmosisParamsProtocolData{}
	_ ProtocolDataI = &LiquidAllowedDenomProtocolData{}
)
