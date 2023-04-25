package types

import (
	"fmt"
	"time"
)

func NewGenesisState(epochs []EpochInfo) *GenesisState {
	return &GenesisState{Epochs: epochs}
}

func NewGenesisEpochInfo(identifier string, duration time.Duration) EpochInfo {
	return EpochInfo{
		Identifier:              identifier,
		StartTime:               time.Time{},
		Duration:                duration,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		CurrentEpochStartTime:   time.Time{},
		EpochCountingStarted:    false,
	}
}

// DefaultGenesis returns the default Capability genesis state.
func DefaultGenesis() *GenesisState {
	epochs := []EpochInfo{
		NewGenesisEpochInfo("week", time.Hour*24*7),
		NewGenesisEpochInfo("day", time.Hour*24),
		NewGenesisEpochInfo("epoch", time.Second*240),
	}
	return NewGenesisState(epochs)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs *GenesisState) Validate() error {
	epochIdentifiers := map[string]bool{}
	for i, epoch := range gs.Epochs {
		if epoch.Identifier == "" {
			return fmt.Errorf("value #%d: epoch identifier should NOT be empty", i+1)
		}
		if epochIdentifiers[epoch.Identifier] {
			return fmt.Errorf("value #%d: epoch identifier should be unique, got duplicate %q", i+1, epoch.Identifier)
		}
		if epoch.Duration <= 0 {
			return fmt.Errorf("value #%d, Identifier: %q: epoch duration should be >0", i+1, epoch.Identifier)
		}
		epochIdentifiers[epoch.Identifier] = true
	}
	return nil
}
