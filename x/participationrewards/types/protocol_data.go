package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func NewProtocolData(datatype string, data json.RawMessage) *ProtocolData {
	return &ProtocolData{Type: datatype, Data: data}
}

func unmarshalProtocolData[V ProtocolDataI](data json.RawMessage) (ProtocolDataI, error) {
	var cpd, blank V
	json.Unmarshal([]byte(`{}`), &blank)
	err := json.Unmarshal(data, &cpd)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal intermediary %s: %w", reflect.TypeOf(cpd).Name(), err)
	}
	if reflect.DeepEqual(cpd, blank) {
		return nil, fmt.Errorf("unable to unmarshal %s from empty JSON object", reflect.TypeOf(cpd).Name())
	}
	return cpd, nil
}

func UnmarshalProtocolData(datatype ProtocolDataType, data json.RawMessage) (ProtocolDataI, error) {
	switch datatype {
	case ProtocolDataTypeConnection:
		return unmarshalProtocolData[*ConnectionProtocolData](data)
	case ProtocolDataTypeOsmosisParams:
		return unmarshalProtocolData[*OsmosisParamsProtocolData](data)
	case ProtocolDataTypeLiquidToken:
		return unmarshalProtocolData[*LiquidAllowedDenomProtocolData](data)
	case ProtocolDataTypeOsmosisPool:
		return unmarshalProtocolData[*OsmosisPoolProtocolData](data)
	case ProtocolDataTypeUmeeParams:
		return unmarshalProtocolData[*UmeeParamsProtocolData](data)
	case ProtocolDataTypeUmeeReserves:
		return unmarshalProtocolData[*UmeeReservesProtocolData](data)
	case ProtocolDataTypeUmeeUTokenSupply:
		return unmarshalProtocolData[*UmeeUTokenSupplyProtocolData](data)
	case ProtocolDataTypeUmeeTotalBorrows:
		return unmarshalProtocolData[*UmeeTotalBorrowsProtocolData](data)
	case ProtocolDataTypeUmeeInterestScalar:
		return unmarshalProtocolData[*UmeeInterestScalarProtocolData](data)
	case ProtocolDataTypeUmeeLeverageModuleBalance:
		return unmarshalProtocolData[*UmeeLeverageModuleBalanceProtocolData](data)
	case ProtocolDataTypeCrescentParams:
		return unmarshalProtocolData[*CrescentParamsProtocolData](data)
	case ProtocolDataTypeCrescentPool:
		return unmarshalProtocolData[*CrescentPoolProtocolData](data)
	case ProtocolDataTypeCrescentPoolCoinSupply:
		return unmarshalProtocolData[*CrescentPoolCoinSupplyProtocolData](data)
	case ProtocolDataTypeCrescentReserveAddressBalance:
		return unmarshalProtocolData[*CrescentReserveAddressBalanceProtocolData](data)
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
	_ ProtocolDataI = &CrescentParamsProtocolData{}
	_ ProtocolDataI = &CrescentPoolProtocolData{}
	_ ProtocolDataI = &CrescentReserveAddressBalanceProtocolData{}
	_ ProtocolDataI = &CrescentPoolCoinSupplyProtocolData{}
)
