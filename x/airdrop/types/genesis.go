package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

func NewGenesisState(params Params, zoneDrops []*ZoneDrop, claimRecords []*ClaimRecord) *GenesisState {
	return &GenesisState{
		Params:       params,
		ZoneDrops:    zoneDrops,
		ClaimRecords: claimRecords,
	}
}

// DefaultGenesisState returns the default ics genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		ZoneDrops:    make([]*ZoneDrop, 0),
		ClaimRecords: make([]*ClaimRecord, 0),
	}
}

// Validate validates the provided genesis state to ensure the
// expected invariants hold.
func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	zoneMap := make(map[string]int)
	sumMap := make(map[string]uint64)
	for i, zd := range gs.ZoneDrops {
		// check for duplicate zone chain id
		if zdi, exists := zoneMap[zd.ChainId]; exists {
			return fmt.Errorf("%w, [%d] %s already used for zone drop [%d]", ErrDuplicateZoneDrop, i, zd.ChainId, zdi)
		}
		// validate zone drop
		if err := zd.ValidateBasic(); err != nil {
			return err
		}
		// add to lookup map
		zoneMap[zd.ChainId] = i
		sumMap[zd.ChainId] = 0
	}

	claimMap := make(map[string]int)
	for i, cr := range gs.ClaimRecords {
		// check for duplicate
		key := cr.ChainId + "." + cr.Address
		if cmi, exists := claimMap[key]; exists {
			return fmt.Errorf("%w, [%d] %s already used for zone drop [%d]", ErrDuplicateClaimRecord, i, key, cmi)
		}
		// validate claim record
		if err := cr.ValidateBasic(); err != nil {
			return err
		}
		// check corresponding zone drop exists
		if _, exists := zoneMap[cr.ChainId]; !exists {
			return fmt.Errorf("%w, %s for claim record [%d]", ErrZoneDropNotFound, cr.ChainId, i)
		}
		// sum MaxAllocations per zone
		sumMap[cr.ChainId] += cr.MaxAllocation
		// add to lookup map
		claimMap[key] = i
	}

	for i, zd := range gs.ZoneDrops {
		if zd.Allocation < sumMap[zd.ChainId] {
			return fmt.Errorf("%w, zone drop [%d], max %v, got %v", ErrAllocationExceeded, i, zd.Allocation, sumMap[zd.ChainId])
		}
		if sumMap[zd.ChainId] == 0 {
			return fmt.Errorf("%w, %s [%d]", ErrNoClaimRecords, zd.ChainId, i)
		}
	}

	return nil
}

// GetGenesisStateFromAppState returns x/airdrop GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
