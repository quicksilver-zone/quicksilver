package types

import (
	sdkioerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default Capability genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		FactoryDenoms: []GenesisDenom{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.Params.Validate()
	if err != nil {
		return err
	}

	seenDenoms := map[string]bool{}

	for _, denom := range gs.GetFactoryDenoms() {
		if seenDenoms[denom.GetDenom()] {
			return sdkioerrors.Wrapf(ErrInvalidGenesis, "duplicate denom: %s", denom.GetDenom())
		}

		seenDenoms[denom.GetDenom()] = true

		_, _, err := DeconstructDenom(denom.GetDenom())
		if err != nil {
			return err
		}

		if denom.AuthorityMetadata.Admin != "" {
			_, err = sdk.AccAddressFromBech32(denom.AuthorityMetadata.Admin)
			if err != nil {
				return sdkioerrors.Wrapf(ErrInvalidAuthorityMetadata, "Invalid admin address (%s)", err)
			}
		}
	}

	return nil
}
