package types

func NewGenesisState(zones []RegisteredZone) *GenesisState {
	return &GenesisState{Zones: zones}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	epochs := []RegisteredZone{}
	return NewGenesisState(epochs)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// TODO: validate genesis state.
	return nil
}
