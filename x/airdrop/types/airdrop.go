package types

import (
	"fmt"

	"github.com/ingenuity-build/multierror"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (zd *ZoneDrop) ValidateBasic() error {
	errs := make(map[string]error)

	validateChainID(zd, errs)
	validateDuration(zd, errs)
	validateDecay(zd, errs)
	validateSumDurationDecay(zd, errs)
	validateAllocation(zd, errs)
	validateActions(zd, errs)

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func validateChainID(zd *ZoneDrop, errs map[string]error) {
	if zd.ChainId == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}
}

func validateDuration(zd *ZoneDrop, errs map[string]error) {
	if zd.Duration.Microseconds() < 0 {
		errs["Duration"] = fmt.Errorf("%w, must not be negative", ErrInvalidDuration)
	}
}

func validateDecay(zd *ZoneDrop, errs map[string]error) {
	if zd.Decay.Microseconds() < 0 {
		errs["Decay"] = fmt.Errorf("%w, must not be negative", ErrInvalidDuration)
	}
}

func validateSumDurationDecay(zd *ZoneDrop, errs map[string]error) {
	if zd.Duration.Microseconds()+zd.Decay.Microseconds() == 0 {
		errs["Duration"] = fmt.Errorf("%w, sum of Duration and Decay must not be zero", ErrInvalidDuration)
		errs["Decay"] = fmt.Errorf("%w, sum of Duration and Decay must not be zero", ErrInvalidDuration)
	}
}

func validateAllocation(zd *ZoneDrop, errs map[string]error) {
	if zd.Allocation == 0 {
		errs["Allocation"] = ErrUndefinedAttribute
	}
}

func validateActions(zd *ZoneDrop, errs map[string]error) {
	if len(zd.Actions) == 0 {
		errs["Actions"] = ErrUndefinedAttribute
		return
	}
	if len(zd.Actions) > len(Action_name)-1 {
		errs["Action"] = fmt.Errorf("exceeds number of defined actions (%d)", len(Action_name)-1)
		return
	}
	validateActionWeights(zd, errs)
}

func validateActionWeights(zd *ZoneDrop, errs map[string]error) {
	weightSum := sdk.ZeroDec()
	for _, aw := range zd.Actions {
		weightSum = weightSum.Add(aw)
	}
	if !weightSum.Equal(sdk.OneDec()) {
		errs["Actions"] = fmt.Errorf("%w, got %s", ErrActionWeights, weightSum)
	}
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
