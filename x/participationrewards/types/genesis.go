package types

import (
	fmt "fmt"
)

func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{Params: params}
}

// DefaultGenesis returns the default ics genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func (gs GenesisState) Validate() error {
	errors := make(map[string]error)

	if err := gs.Params.Validate(); err != nil {
		errors["Params"] = err
	}

	for i, claim := range gs.Claims {
		if err := claim.ValidateBasic(); err != nil {
			el := fmt.Sprintf("Claim[%d]", i)
			errors[el] = err
		}
	}

	for i, kpd := range gs.ProtocolData {
		if err := kpd.ValidateBasic(); err != nil {
			el := fmt.Sprintf("ProtocolData[%d]", i)
			errors[el] = err
			continue
		}
	}

	return nil
}
