package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func (zd ZoneDrop) ValidateBasic() error {
	errors := make(map[string]error)

	// must be defined
	if zd.ChainId == "" {
		errors["ChainId"] = ErrUndefinedAttribute
	}

	// must be greater than 0
	if zd.Duration.Microseconds() <= 0 {
		errors["Duration"] = ErrInvalidDuration
	}

	// must be greater or equal to 0
	// - equal will result in a full airdrop reward with immediate cut off on
	//   expiery;
	// - greater will result in a proportionally discounted airdrop reward over
	//   the duration of decay;
	if zd.Decay.Microseconds() < 0 {
		errors["Decay"] = ErrInvalidDuration
	}

	// must be positive value
	if zd.Allocation == 0 {
		errors["Allocation"] = ErrUndefinedAttribute
	}

	// must have at least one defined
	if zd.Actions == nil || len(zd.Actions) == 0 {
		errors["Actions"] = ErrUndefinedAttribute
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

func (cr ClaimRecord) ValidateBasic() error {
	errors := make(map[string]error)

	// must be defined
	if cr.ChainId == "" {
		errors["ChainId"] = ErrUndefinedAttribute
	}

	// must be valid bech32
	if _, err := sdk.AccAddressFromBech32(cr.Address); err != nil {
		errors["Address"] = err
	}

	// must be positive value
	if cr.MaxAllocation == 0 {
		errors["MaxAllocation"] = ErrUndefinedAttribute
	}

	// check action enum in bounds and sum <= max allocation
	if cr.ActionsCompleted != nil {
		// action enum, completed action
		i := 0
		sum := uint64(0)
		for ae, ca := range cr.ActionsCompleted {
			if int(ae) >= len(Action_name) {
				kstr := fmt.Sprintf("ActionsCompleted[%d]", i)
				errors[kstr] = fmt.Errorf("enum out of bounds, expects [0-%d), got %d", len(Action_name), ae)
			}
			sum += ca.ClaimAmount
			i++
		}

		if sum > cr.MaxAllocation {
			errors["ActionsCompleted"] = fmt.Errorf(
				"sum of claims exceed max allocation, expected %d got %d",
				cr.MaxAllocation,
				sum,
			)
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
