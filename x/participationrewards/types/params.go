package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ingenuity-build/quicksilver/internal/multierror"
	"gopkg.in/yaml.v2"
)

var (
	KeyDistributionProportions = []byte("DistributionProportions")

	DefaultValidatorSelectionAllocation = sdk.NewDecWithPrec(34, 2)
	DefaultHoldingsAllocation           = sdk.NewDecWithPrec(33, 2)
	DefaultLockupAllocation             = sdk.NewDecWithPrec(33, 2)
)

// ParamTable for participationrewards module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new ics Params instance
func NewParams(
	validatorSelectionAllocation sdk.Dec,
	holdingsAllocation sdk.Dec,
	lockupAllocation sdk.Dec,
) Params {
	return Params{
		DistributionProportions: DistributionProportions{
			ValidatorSelectionAllocation: validatorSelectionAllocation,
			HoldingsAllocation:           holdingsAllocation,
			LockupAllocation:             lockupAllocation,
		},
	}
}

// DefaultParams default ics params
func DefaultParams() Params {
	return NewParams(
		DefaultValidatorSelectionAllocation,
		DefaultHoldingsAllocation,
		DefaultLockupAllocation,
	)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistributionProportions, &p.DistributionProportions, validateDistributionProportions),
	}
}

func validateDistributionProportions(i interface{}) error {
	v, ok := i.(DistributionProportions)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	errors := make(map[string]error)

	if v.ValidatorSelectionAllocation.IsNil() {
		errors["ValidatorSelectionAllocation"] = ErrUndefinedAttribute
	} else if v.ValidatorSelectionAllocation.IsNegative() {
		errors["ValidatorSelectionAllocation"] = ErrNegativeDistributionRatio
	}

	if v.HoldingsAllocation.IsNil() {
		errors["HoldingsAllocation"] = ErrUndefinedAttribute
	} else if v.HoldingsAllocation.IsNegative() {
		errors["HoldingsAllocation"] = ErrNegativeDistributionRatio
	}

	if v.LockupAllocation.IsNil() {
		errors["LockupAllocation"] = ErrUndefinedAttribute
	} else if v.LockupAllocation.IsNegative() {
		errors["LockupAllocation"] = ErrNegativeDistributionRatio
	}

	// no errors yet: check total proportions
	if len(errors) == 0 {
		totalProportions := v.ValidatorSelectionAllocation.Add(v.HoldingsAllocation).Add(v.LockupAllocation)

		if !totalProportions.Equal(sdk.NewDec(1)) {
			errors["TotalProportions"] = ErrInvalidTotalProportions
		}
	}

	// all checks done, count errors
	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

// validate params.
func (p Params) Validate() error {
	return validateDistributionProportions(p.DistributionProportions)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
