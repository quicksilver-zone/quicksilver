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
		totalProportions := dp.ValidatorSelectionAllocation.Add(dp.HoldingsAllocation).Add(dp.LockupAllocation)

		if !totalProportions.Equal(sdk.OneDec()) {
			errors["TotalProportions"] = ErrInvalidTotalProportions
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

func (dp DistributionProportions) TotalProportions() sdk.Dec {
	return dp.ValidatorSelectionAllocation.Add(dp.HoldingsAllocation).Add(dp.LockupAllocation)
}

func (p Params) ValidateBasic() error {
	return p.DistributionProportions.ValidateBasic()
}

func (c Claim) ValidateBasic() error {
	errors := make(map[string]error)

	_, err := sdk.AccAddressFromBech32(c.UserAddress)
	if err != nil {
		errors["UserAddress"] = err
	}

	if len(c.ChainId) == 0 {
		errors["ChainId"] = ErrUndefinedAttribute
	}

	if c.Amount <= 0 {
		errors["Amount"] = ErrNotPositive
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
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

	if len(pd.Protocol) == 0 {
		errors["Protocol"] = ErrUndefinedAttribute
	}

	// type enumerator
	var te ProtocolDataType
	if len(pd.Type) == 0 {
		errors["Type"] = ErrUndefinedAttribute
	} else {
		if tv, exists := ProtocolDataType_value[pd.Type]; !exists {
			errors["Type"] = fmt.Errorf("%w: %s", ErrUnknownProtocolDataType, pd.Type)
		} else {
			// capture enum value to validate protocol data according to type
			te = tv
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
	case ProtocolDataLiquidToken:
		pd := LiquidAllowedDenomProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return err
		}
		pdi = &pd
	case ProtocolDataOsmosisPool:
		pd := OsmosisPoolProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return err
		}
		pdi = &pd
	case ProtocolDataCrescentPool:
		return ErrUnimplementedProtocolDataType
	case ProtocolDataSifchainPool:
		return ErrUnimplementedProtocolDataType
	default:
		return ErrUnknownProtocolDataType
	}

	return pdi.ValidateBasic()
}
