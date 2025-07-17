package types

import (
	"fmt"

	"go.uber.org/multierr"
)

func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{Params: params}
}

// DefaultGenesisState returns the default ics genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate validates the provided genesis state to ensure the
// expected invariants holds.
func (gs *GenesisState) Validate() error {
	errs := make(map[string]error)

	if err := gs.Params.Validate(); err != nil {
		errs["Params"] = err
	}

	for i, kpd := range gs.ProtocolData {
		if err := kpd.ValidateBasic(); err != nil {
			el := fmt.Sprintf("ProtocolData[%d]", i)
			errs[el] = err
		}
	}

	if len(errs) > 0 {
		var errList []error
		for _, err := range errs {
			errList = append(errList, err)
		}
		return multierr.Combine(errList...)
	}

	return nil
}
