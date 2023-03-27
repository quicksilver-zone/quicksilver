package types

import (
	"encoding/json"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func (dp DistributionProportions) ValidateBasic() error {
	errors := make(map[string]error)

	if dp.ValidatorSelectionAllocation.IsNil() {
		errors["ValidatorSelectionAllocation"] = ErrUndefinedAttribute
	} else if dp.ValidatorSelectionAllocation.IsNegative() {
		errors["ValidatorSelectionAllocation"] = ErrNegativeAttribute
	}

	if dp.HoldingsAllocation.IsNil() {
		errors["HoldingsAllocation"] = ErrUndefinedAttribute
	} else if dp.HoldingsAllocation.IsNegative() {
		errors["HoldingsAllocation"] = ErrNegativeAttribute
	}

	if dp.LockupAllocation.IsNil() {
		errors["LockupAllocation"] = ErrUndefinedAttribute
	} else if dp.LockupAllocation.IsNegative() {
		errors["LockupAllocation"] = ErrNegativeAttribute
	}

	// no errors yet: check total proportions
	if len(errors) == 0 {
		totalProportions := dp.Total()

		if !totalProportions.Equal(sdk.OneDec()) {
			errors["TotalProportions"] = fmt.Errorf("%w, got %v", ErrInvalidTotalProportions, totalProportions)
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

func (dp DistributionProportions) Total() sdk.Dec {
	return dp.ValidatorSelectionAllocation.Add(dp.HoldingsAllocation).Add(dp.LockupAllocation)
}

func (kpd KeyedProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	if len(kpd.Key) == 0 {
		errors["Key"] = ErrUndefinedAttribute
	}

	if kpd.ProtocolData == nil {
		errors["ProtocolData"] = ErrUndefinedAttribute
	} else {
		if err := kpd.ProtocolData.ValidateBasic(); err != nil {
			errors["ProtocolData"] = err
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

func (pd ProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	// type enumerator
	var te ProtocolDataType
	if len(pd.Type) == 0 {
		errors["Type"] = ErrUndefinedAttribute
	} else {
		if tv, exists := ProtocolDataType_value[pd.Type]; !exists {
			errors["Type"] = fmt.Errorf("%w: %s", ErrUnknownProtocolDataType, pd.Type)
		} else {
			// capture enum value to validate protocol data according to type
			te = ProtocolDataType(tv)
		}
	}

	if len(pd.Data) == 0 {
		errors["Data"] = ErrUndefinedAttribute
	} else if te != -1 {
		if err := validateProtocolData(pd.Data, te); err != nil {
			errors["Data"] = err
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

// unmarshal to appropriate concrete type and validate
func validateProtocolData(data json.RawMessage, pdt ProtocolDataType) error {
	var pdi ProtocolDataI
	switch pdt {
	case ProtocolDataTypeLiquidToken:
		pd := LiquidAllowedDenomProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return err
		}
		pdi = &pd
	case ProtocolDataTypeOsmosisPool:
		pd := OsmosisPoolProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return err
		}
		pdi = &pd
	case ProtocolDataTypeCrescentPool:
		return ErrUnimplementedProtocolDataType
	case ProtocolDataTypeSifchainPool:
		return ErrUnimplementedProtocolDataType
	default:
		return ErrUnknownProtocolDataType
	}

	return pdi.ValidateBasic()
}
