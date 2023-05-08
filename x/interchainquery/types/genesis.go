package types

func NewGenesisState(queries []Query) *GenesisState {
	return &GenesisState{Queries: queries}
}

// DefaultGenesisState returns the default Capability genesis state.
func DefaultGenesisState() *GenesisState {
	var queries []Query
	return NewGenesisState(queries)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// TODO: validate genesis state.
	return nil
}
