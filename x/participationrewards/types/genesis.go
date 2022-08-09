package types

import (
	fmt "fmt"

	"github.com/ingenuity-build/quicksilver/utils"
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
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, claim := range data.Claims {
		// check user address
		_, err := utils.AccAddressFromBech32(claim.UserAddress, "")
		if err != nil {
			return err
		}

		// check value is valid
		if claim.HeldAmount <= 0 {
			return fmt.Errorf("claim contains a non-positive value")
		}

	}

	// TODO: validate protocol data is valid
OUTER:
	for _, pd := range data.ProtocolData {
		for _, claimType := range ClaimTypes {
			if claimType == pd.ProtocolData.Type {
				continue OUTER
			}
		}
		return fmt.Errorf("invalid protocol data type: %s", pd.ProtocolData.Type)
	}

	return nil
}
