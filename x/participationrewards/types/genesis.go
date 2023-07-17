package types

import (
	"fmt"

	"github.com/ingenuity-build/multierror"
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
	errors := make(map[string]error)

	if err := gs.Params.Validate(); err != nil {
		errors["Params"] = err
	}

	for i, kpd := range gs.ProtocolData {
		if err := kpd.ValidateBasic(); err != nil {
			el := fmt.Sprintf("ProtocolData[%d]", i)
			errors[el] = err
			continue
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
