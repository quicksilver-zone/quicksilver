package types

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	SelfConnection = "local"
)

func (dp *DistributionProportions) ValidateBasic() error {
	errs := make(map[string]error)

	if dp.ValidatorSelectionAllocation.IsNil() {
		errs["ValidatorSelectionAllocation"] = ErrUndefinedAttribute
	} else if dp.ValidatorSelectionAllocation.IsNegative() {
		errs["ValidatorSelectionAllocation"] = ErrNegativeAttribute
	}

	if dp.HoldingsAllocation.IsNil() {
		errs["HoldingsAllocation"] = ErrUndefinedAttribute
	} else if dp.HoldingsAllocation.IsNegative() {
		errs["HoldingsAllocation"] = ErrNegativeAttribute
	}

	if dp.LockupAllocation.IsNil() {
		errs["LockupAllocation"] = ErrUndefinedAttribute
	} else if dp.LockupAllocation.IsNegative() {
		errs["LockupAllocation"] = ErrNegativeAttribute
	}

	// no errors yet: check total proportions
	if len(errs) == 0 {
		totalProportions := dp.Total()

		if !totalProportions.Equal(sdk.OneDec()) {
			errs["TotalProportions"] = fmt.Errorf("%w, got %v", ErrInvalidTotalProportions, totalProportions)
		}
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func (dp *DistributionProportions) Total() sdk.Dec {
	return dp.ValidatorSelectionAllocation.Add(dp.HoldingsAllocation).Add(dp.LockupAllocation)
}

func (kpd *KeyedProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	if kpd.Key == "" {
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

func (pd *ProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	// type enumerator
	var te ProtocolDataType
	if pd.Type == "" {
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

// validateProtocolData unmarshals to appropriate concrete type and validate.
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
	case ProtocolDataTypeUmeeReserves:
		pd := UmeeReservesProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return err
		}
		pdi = &pd
	case ProtocolDataTypeUmeeInterestScalar:
		pd := UmeeInterestScalarProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return err
		}
		pdi = &pd
	case ProtocolDataTypeUmeeTotalBorrows:
		pd := UmeeTotalBorrowsProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return err
		}
		pdi = &pd
	case ProtocolDataTypeUmeeUTokenSupply:
		pd := UmeeUTokenSupplyProtocolData{}
		err := json.Unmarshal(data, &pd)
		if err != nil {
			return err
		}
		pdi = &pd
	case ProtocolDataTypeUmeeLeverageModuleBalance:
		pd := UmeeLeverageModuleBalanceProtocolData{}
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

// UserAllocation is an internal keeper struct to track transient state for
// rewards distribution. It contains the user address and the coins that are
// allocated to it.
type UserAllocation struct {
	Address string
	Amount  sdk.Coin
}

// ZoneScore is an internal struct to track transient state for the calculation
// of zone scores. It specifically tallies the total zone voting power used in
// calculations to determine validator voting power percentages.
type ZoneScore struct {
	ZoneID           string // chainID
	TotalVotingPower math.Int
	ValidatorScores  map[string]*Validator
}

// Validator is an internal struct to track transient state for the calculation
// of zone scores. It contains all relevant Validator scoring metrics with a
// pointer reference to the actual Validator (embedded).
type Validator struct {
	PowerPercentage   sdk.Dec
	PerformanceScore  sdk.Dec
	DistributionScore sdk.Dec

	*icstypes.Validator
}

// UserScore is an internal struct to track transient state for rewards
// distribution. It contains the user address and individual score.
type UserScore struct {
	Address string
	Score   sdk.Dec
}
