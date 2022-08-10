package types

import (
	"fmt"
	time "time"

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
	} else {
		wsum := sdk.ZeroDec()
		for _, aw := range zd.Actions {
			wsum = wsum.Add(aw)
		}
		if !wsum.Equal(sdk.OneDec()) {
			errors["Actions"] = fmt.Errorf("%w, got %s", ErrActionWeights, wsum)
		}
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
			// check enum bounds
			kstr := fmt.Sprintf("ActionsCompleted[%d]", i)
			if int(ae) >= len(Action_name) {
				errors[kstr+" Enum"] = fmt.Errorf("%w, got %d", ErrActionOutOfBounds, ae)
			}
			// calc sum
			sum += ca.ClaimAmount
			// check completed time
			if ca.CompleteTime.After(time.Now()) {
				errors[kstr+" CompleteTime"] = fmt.Errorf("invalid spacetime continuum, time is in the future")
			}
			// check claim amount
			if ca.ClaimAmount > cr.MaxAllocation {
				errors[kstr+" ClaimAmount"] = fmt.Errorf("exceeds max allocation")
			}
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

func (a Action) InBounds() bool {
	// get action enum
	ae := int(a)

	// check action enum
	if ae < 0 || ae >= len(Action_name) {
		return false
	}

	return true
}
