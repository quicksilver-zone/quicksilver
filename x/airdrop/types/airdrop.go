package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func (zd ZoneDrop) ValidateBasic() error {
	errors := make(map[string]error)

	// must be defined
	if len(zd.ChainId) == 0 {
		errors["ChainId"] = ErrUndefinedAttribute
	}

	// must be greater or equal to 0
	//
	// - equal will result in immediate decay;
	// - greater specifies period of full claim;
	//
	// specific bounds can be applied via proposal process
	if zd.Duration.Microseconds() < 0 {
		errors["Duration"] = fmt.Errorf("%w, must not be negative", ErrInvalidDuration)
	}

	// must be greater or equal to 0
	//
	// - equal will result in a full airdrop reward with immediate cut off on
	//   expiery;
	// - greater will result in a proportionally discounted airdrop reward over
	//   the duration of decay;
	//
	// specific bounds can be applied via proposal process
	if zd.Decay.Microseconds() < 0 {
		errors["Decay"] = fmt.Errorf("%w, must not be negative", ErrInvalidDuration)
	}

	// sum of Duration and Decay may not be zero as this will result in
	// immediate expiery of the zone airdrop
	if zd.Duration.Microseconds()+zd.Decay.Microseconds() == 0 {
		if _, exists := errors["Duration"]; !exists {
			errors["Duration"] = fmt.Errorf("%w, sum of Duration and Decay must not be zero", ErrInvalidDuration)
		}
		if _, exists := errors["Decay"]; !exists {
			errors["Decay"] = fmt.Errorf("%w, sum of Duration and Decay must not be zero", ErrInvalidDuration)
		}
	}

	// must be positive value
	if zd.Allocation == 0 {
		errors["Allocation"] = ErrUndefinedAttribute
	}

	// must have at least one defined
	if len(zd.Actions) == 0 {
		errors["Actions"] = ErrUndefinedAttribute
	} else {
		// may not exceed defined types.Action bounds
		// * (-1) to account for enum: 0 = undefined (protobuf3 spec)
		if len(zd.Actions) > len(Action_name)-1 {
			errors["Action"] = fmt.Errorf("exceeds number of defined actions (%d)", len(Action_name)-1)
		} else {
			wsum := sdk.ZeroDec()
			for _, aw := range zd.Actions {
				wsum = wsum.Add(aw)
			}
			if !wsum.Equal(sdk.OneDec()) {
				errors["Actions"] = fmt.Errorf("%w, got %s", ErrActionWeights, wsum)
			}
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
	if len(cr.ChainId) == 0 {
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
	// action enum, completed action
	sum := uint64(0)
	i := 0
	// map[actionEnum]*CompletedAction
	for ae, ca := range cr.ActionsCompleted {
		// action enum (+1 protobuf3 enum spec)
		// check enum bounds
		kstr := fmt.Sprintf("ActionsCompleted[%d]", i)
		if !Action(ae).InBounds() {
			errors[kstr+": enum:"] = fmt.Errorf("%w, got %d", ErrActionOutOfBounds, ae)
		}
		// calc sum
		sum += ca.ClaimAmount

		// check claim amount
		if ca.ClaimAmount > cr.MaxAllocation {
			errors[kstr+": ClaimAmount:"] = fmt.Errorf("exceeds max allocation of %d", cr.MaxAllocation)
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

	if cr.BaseValue == 0 {
		errors["BaseValue"] = ErrUndefinedAttribute
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
	if ae < 1 || ae >= len(Action_name) {
		return false
	}

	return true
}
