package types

import (
	"encoding/json"
	"fmt"

	"go.uber.org/multierr"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
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
		return multierr.Combine(utils.ErrorMapToSlice(errs)...)
	}

	return nil
}

func (dp *DistributionProportions) Total() sdk.Dec {
	return dp.ValidatorSelectionAllocation.Add(dp.HoldingsAllocation).Add(dp.LockupAllocation)
}

func (kpd *KeyedProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if kpd.Key == "" {
		errs["Key"] = ErrUndefinedAttribute
	}

	if kpd.ProtocolData == nil {
		errs["ProtocolData"] = ErrUndefinedAttribute
	} else {
		if err := kpd.ProtocolData.ValidateBasic(); err != nil {
			errs["ProtocolData"] = err
		}
	}

	if len(errs) > 0 {
		return multierr.Combine(utils.ErrorMapToSlice(errs)...)
	}

	return nil
}

func (pd *ProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	// type enumerator
	var te ProtocolDataType
	if pd.Type == "" {
		errs["Type"] = ErrUndefinedAttribute
	} else {
		if tv, exists := ProtocolDataType_value[pd.Type]; !exists {
			errs["Type"] = fmt.Errorf("%w: %s", ErrUnknownProtocolDataType, pd.Type)
		} else {
			// capture enum value to validate protocol data according to type
			te = ProtocolDataType(tv)
		}
	}

	if len(pd.Data) == 0 {
		errs["Data"] = ErrUndefinedAttribute
	} else if te != -1 {
		if err := validateProtocolData(pd.Data, te); err != nil {
			errs["Data"] = err
		}
	}

	if len(errs) > 0 {
		return multierr.Combine(utils.ErrorMapToSlice(errs)...)
	}

	return nil
}

// validateProtocolData unmarshals to appropriate concrete type and validate.
func validateProtocolData(data json.RawMessage, pdt ProtocolDataType) error {
	pdi, err := UnmarshalProtocolData(pdt, data)
	if err != nil {
		return err
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
