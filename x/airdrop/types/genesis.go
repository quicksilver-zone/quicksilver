package types

func NewGenesisState(params Params, zoneDrops []*ZoneDrop, claimRecords []*ClaimRecord) *GenesisState {
	return &GenesisState{
		Params:       params,
		ZoneDrops:    zoneDrops,
		ClaimRecords: claimRecords,
	}
}

// DefaultGenesis returns the default ics genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		ZoneDrops:    make([]*ZoneDrop, 0),
		ClaimRecords: make([]*ClaimRecord, 0),
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants hold.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return nil
}
