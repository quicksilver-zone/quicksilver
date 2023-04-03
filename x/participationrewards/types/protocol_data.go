package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func UnmarshalProtocolData(datatype ProtocolDataType, data json.RawMessage) (ProtocolDataI, error) {
	switch datatype {
	case ProtocolDataTypeConnection:
		cpd := ConnectionProtocolData{}
		err := json.Unmarshal(data, &cpd)
		if err != nil {
			return nil, err
		}
		var blank ConnectionProtocolData
		if reflect.DeepEqual(cpd, blank) {
			return nil, errors.New("unable to unmarshal connection protocol data from empty JSON object")
		}
		return &cpd, nil
	case ProtocolDataTypeOsmosisParams:
		oppd := OsmosisParamsProtocolData{}
		err := json.Unmarshal(data, &oppd)
		if err != nil {
			return nil, err
		}
		var blank OsmosisParamsProtocolData
		if reflect.DeepEqual(oppd, blank) {
			return nil, fmt.Errorf("unable to unmarshal osmosisparams protocol data from empty JSON object")
		}
		return &oppd, nil
	case ProtocolDataTypeLiquidToken:
		ladpd := LiquidAllowedDenomProtocolData{}
		err := json.Unmarshal(data, &ladpd)
		if err != nil {
			return nil, err
		}
		var blank LiquidAllowedDenomProtocolData
		if reflect.DeepEqual(ladpd, blank) {
			return nil, errors.New("unable to unmarshal liquid protocol data from empty JSON object")
		}
		return &ladpd, nil
	case ProtocolDataTypeOsmosisPool:
		oppd := OsmosisPoolProtocolData{}
		err := json.Unmarshal(data, &oppd)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal intermediary osmosisPoolProtocolData: %w", err)
		}
		var blank OsmosisPoolProtocolData
		if reflect.DeepEqual(oppd, blank) {
			return nil, fmt.Errorf("unable to unmarshal osmosispool protocol data from empty JSON object")
		}

		return &oppd, nil
	default:
		return nil, ErrUnknownProtocolDataType
	}
}

type ProtocolDataI interface {
	ValidateBasic() error
}

// ConnectionProtocolData defines state for connection tracking.
type ConnectionProtocolData struct {
	ConnectionID string
	ChainID      string
	LastEpoch    int64
	Prefix       string
}

// ValidateBasic satisfies ProtocolDataI and validates basic stateless data.
func (cpd *ConnectionProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if cpd.ConnectionID == "" {
		errs["ConnectionID"] = ErrUndefinedAttribute
	}

	if cpd.ChainID == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	if cpd.Prefix == "" {
		errs["Prefix"] = ErrUndefinedAttribute
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

var (
	_ ProtocolDataI = &ConnectionProtocolData{}
	_ ProtocolDataI = &OsmosisPoolProtocolData{}
	_ ProtocolDataI = &OsmosisParamsProtocolData{}
	_ ProtocolDataI = &LiquidAllowedDenomProtocolData{}
)
