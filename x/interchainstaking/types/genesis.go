package types

func NewGenesisState(params Params, zones []Zone) *GenesisState {
	return &GenesisState{Params: params, Zones: zones}
}

// DefaultGenesis returns the default ics genesis state.
func DefaultGenesis() *GenesisState {
	zones := []Zone{}
	return NewGenesisState(DefaultParams(), zones)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// TODO: validate genesis state.
	return gs.Params.Validate()
}
