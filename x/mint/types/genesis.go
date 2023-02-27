package types

// NewGenesis creates a new Genesis object.
func NewGenesis(minter Minter, params Params, reductionStartedEpoch int64) *GenesisState {
	return &GenesisState{
		Minter:                minter,
		Params:                params,
		ReductionStartedEpoch: reductionStartedEpoch,
	}
}

// DefaultGenesis creates a default GenesisState object.
func DefaultGenesis() *GenesisState {
	return NewGenesis(
		DefaultInitialMinter(),
		DefaultParams(),
		0,
	)
}

// Validate validates the provided genesis state to ensure the
// expected invariants holds.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	return gs.Minter.Validate()
}
