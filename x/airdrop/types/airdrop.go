package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/ingenuity-build/multierror"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (zd *ZoneDrop) ValidateBasic() error {
	errs := make(map[string]error)

	// must be defined
	if zd.ChainId == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	// must be greater or equal to 0
	//
	// - equal will result in immediate decay;
	// - greater specifies period of full claim;
	//
	// specific bounds can be applied via proposal process
	if zd.Duration.Microseconds() < 0 {
		errs["Duration"] = fmt.Errorf("%w, must not be negative", ErrInvalidDuration)
	}

	// must be greater or equal to 0
	//
	// - equal will result in a full airdrop reward with immediate cut off on
	//   expiry;
	// - greater will result in a proportionally discounted airdrop reward over
	//   the duration of decay;
	//
	// specific bounds can be applied via proposal process
	if zd.Decay.Microseconds() < 0 {
		errs["Decay"] = fmt.Errorf("%w, must not be negative", ErrInvalidDuration)
	}

	// sum of Duration and Decay may not be zero as this will result in
	// immediate expiry of the zone airdrop
	if zd.Duration.Microseconds()+zd.Decay.Microseconds() == 0 {
		if _, exists := errs["Duration"]; !exists {
			errs["Duration"] = fmt.Errorf("%w, sum of Duration and Decay must not be zero", ErrInvalidDuration)
		}
		if _, exists := errs["Decay"]; !exists {
			errs["Decay"] = fmt.Errorf("%w, sum of Duration and Decay must not be zero", ErrInvalidDuration)
		}
	}

	// must be positive value
	if zd.Allocation == 0 {
		errs["Allocation"] = ErrUndefinedAttribute
	}

	// must have at least one defined
	if len(zd.Actions) == 0 {
		errs["Actions"] = ErrUndefinedAttribute
	} else {
		// may not exceed defined types.Action bounds
		// * (-1) to account for enum: 0 = undefined (protobuf3 spec)
		if len(zd.Actions) > len(Action_name)-1 {
			errs["Action"] = fmt.Errorf("exceeds number of defined actions (%d)", len(Action_name)-1)
		} else {
			weightSum := sdkmath.LegacyZeroDec()
			for _, aw := range zd.Actions {
				weightSum = weightSum.Add(aw)
			}
			if !weightSum.Equal(sdkmath.LegacyOneDec()) {
				errs["Actions"] = fmt.Errorf("%w, got %s", ErrActionWeights, weightSum)
			}
		}
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func (cr *ClaimRecord) ValidateBasic() error {
	errs := make(map[string]error)

	// must be defined
	if cr.ChainId == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	// must be valid bech32
	if _, _, err := bech32.DecodeAndConvert(cr.Address); err != nil {
		errs["Address"] = err
	}

	// must be positive value
	if cr.MaxAllocation == 0 {
		errs["MaxAllocation"] = ErrUndefinedAttribute
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
			errs[kstr+": enum:"] = fmt.Errorf("%w, got %d", ErrActionOutOfBounds, ae)
		}
		// calc sum
		sum += ca.ClaimAmount

		// check claim amount
		if ca.ClaimAmount > cr.MaxAllocation {
			errs[kstr+": ClaimAmount:"] = fmt.Errorf("exceeds max allocation of %d", cr.MaxAllocation)
		}
		i++
	}

	if sum > cr.MaxAllocation {
		errs["ActionsCompleted"] = fmt.Errorf(
			"sum of claims exceed max allocation, expected %d got %d",
			cr.MaxAllocation,
			sum,
		)
	}

	if cr.BaseValue == 0 {
		errs["BaseValue"] = ErrUndefinedAttribute
	}

	if len(errs) > 0 {
		return multierror.New(errs)
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
