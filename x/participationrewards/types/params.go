package types

import (
	"errors"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var (
	KeyDistributionProportions = []byte("DistributionProportions")

	DefaultValidatorSelectionAllocation = sdk.NewDecWithPrec(33, 2)
	DefaultPariticpationAllocation      = sdk.NewDecWithPrec(33, 2)
	DefaultLockupAllocation             = sdk.NewDecWithPrec(34, 2)
)

// ParamTable for participationrewards module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new ics Params instance
func NewParams(
	validator_selection_allocation sdk.Dec,
	participation_allocation sdk.Dec,
	lockup_allocation sdk.Dec,
) Params {
	return Params{
		DistributionProportions: DistributionProportions{
			ValidatorSelectionAllocation: validator_selection_allocation,
			PariticpationAllocation:      participation_allocation,
			LockupAllocation:             lockup_allocation,
		},
	}
}

// DefaultParams default ics params
func DefaultParams() Params {
	return NewParams(
		DefaultValidatorSelectionAllocation,
		DefaultPariticpationAllocation,
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

	if v.ValidatorSelectionAllocation.IsNegative() {
		return errors.New("ValidatorSelectionAllocation distribution ratio should not be negative")
	}

	if v.PariticpationAllocation.IsNegative() {
		return errors.New("PariticpationAllocation distribution ratio should not be negative")
	}

	if v.LockupAllocation.IsNegative() {
		return errors.New("LockupAllocation distribution ratio should not be negative")
	}

	totalProportions := v.ValidatorSelectionAllocation.Add(v.PariticpationAllocation).Add(v.LockupAllocation)

	if !totalProportions.Equal(sdk.NewDec(1)) {
		return errors.New("total distributions ratio should be 1")
	}

	return nil
}

// validate params.
func (p Params) Validate() error {
	if err := validateDistributionProportions(p.DistributionProportions); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
