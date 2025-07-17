package types

import (
	fmt "fmt"

	"go.uber.org/multierr"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewClaim(address, chainID string, module ClaimType, srcChainID string, amount math.Int) Claim {
	return Claim{UserAddress: address, ChainId: chainID, Module: module, SourceChainId: srcChainID, Amount: amount}
}

// ValidateBasic performs stateless validation of a Claim.
func (c *Claim) ValidateBasic() error {
	var errs error

	_, err := sdk.AccAddressFromBech32(c.UserAddress)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("UserAddress: %w", err))
	}

	if c.ChainId == "" {
		errs = multierr.Append(errs, fmt.Errorf("ChainID: %w", ErrUndefinedAttribute))
	}

	if c.Amount.IsNil() || !c.Amount.IsPositive() {
		errs = multierr.Append(errs, fmt.Errorf("Amount: %w", ErrNotPositive))
	}

	return errs
}
