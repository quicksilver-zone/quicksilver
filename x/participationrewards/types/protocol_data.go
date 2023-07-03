package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func NewProtocolData(datatype string, data json.RawMessage) *ProtocolData {
	return &ProtocolData{Type: datatype, Data: data}
}

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
	case ProtocolDataTypeUmeeParams:
		uppd := UmeeParamsProtocolData{}
		err := json.Unmarshal(data, &uppd)
		if err != nil {
			return nil, err
		}
		var blank UmeeParamsProtocolData
		if reflect.DeepEqual(uppd, blank) {
			return nil, fmt.Errorf("unable to unmarshal umeeparams protocol data from empty JSON object")
		}
		return &uppd, nil
	case ProtocolDataTypeUmeeReserves:
		upd := UmeeReservesProtocolData{}
		err := json.Unmarshal(data, &upd)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal intermediary UmeeReservesProtocolData: %w", err)
		}
		var blank UmeeReservesProtocolData
		if reflect.DeepEqual(upd, blank) {
			return nil, fmt.Errorf("unable to unmarshal UmeeReservesProtocolData from empty JSON object")
		}

		return &upd, nil
	case ProtocolDataTypeUmeeUTokenSupply:
		upd := UmeeUTokenSupplyProtocolData{}
		err := json.Unmarshal(data, &upd)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal intermediary UmeeUTokenSupplyProtocolData: %w", err)
		}
		var blank UmeeUTokenSupplyProtocolData
		if reflect.DeepEqual(upd, blank) {
			return nil, fmt.Errorf("unable to unmarshal UmeeUTokenSupplyProtocolData from empty JSON object")
		}

		return &upd, nil
	case ProtocolDataTypeUmeeTotalBorrows:
		upd := UmeeTotalBorrowsProtocolData{}
		err := json.Unmarshal(data, &upd)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal intermediary UmeeTotalBorrowsProtocolData: %w", err)
		}
		var blank UmeeTotalBorrowsProtocolData
		if reflect.DeepEqual(upd, blank) {
			return nil, fmt.Errorf("unable to unmarshal UmeeTotalBorrowsProtocolData from empty JSON object")
		}

		return &upd, nil
	case ProtocolDataTypeUmeeInterestScalar:
		upd := UmeeInterestScalarProtocolData{}
		err := json.Unmarshal(data, &upd)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal intermediary UmeeInterestScalarProtocolData: %w", err)
		}
		var blank UmeeInterestScalarProtocolData
		if reflect.DeepEqual(upd, blank) {
			return nil, fmt.Errorf("unable to unmarshal UmeeInterestScalarProtocolData from empty JSON object")
		}

		return &upd, nil
	case ProtocolDataTypeUmeeLeverageModuleBalance:
		upd := UmeeLeverageModuleBalanceProtocolData{}
		err := json.Unmarshal(data, &upd)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal intermediary UmeeLeverageModuleBalanceProtocolData: %w", err)
		}
		var blank UmeeLeverageModuleBalanceProtocolData
		if reflect.DeepEqual(upd, blank) {
			return nil, fmt.Errorf("unable to unmarshal UmeeLeverageModuleBalanceProtocolData from empty JSON object")
		}

		return &upd, nil
	default:
		return nil, ErrUnknownProtocolDataType
	}
}

type ProtocolDataI interface {
	ValidateBasic() error
	GenerateKey() []byte
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

func (cpd ConnectionProtocolData) GenerateKey() []byte {
	return []byte(cpd.ChainID)
}

var (
	_ ProtocolDataI = &ConnectionProtocolData{}
	_ ProtocolDataI = &OsmosisPoolProtocolData{}
	_ ProtocolDataI = &OsmosisParamsProtocolData{}
	_ ProtocolDataI = &LiquidAllowedDenomProtocolData{}
	_ ProtocolDataI = &UmeeProtocolData{}
	_ ProtocolDataI = &UmeeParamsProtocolData{}
)
