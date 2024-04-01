package types

import (
	"github.com/ingenuity-build/multierror"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewClaim(address, chainID string, module ClaimType, srcChainID string, amount math.Int) Claim {
	return Claim{UserAddress: address, ChainId: chainID, Module: module, SourceChainId: srcChainID, Amount: amount}
}

// ValidateBasic performs stateless validation of a Claim.
func (c *Claim) ValidateBasic() error {
	errs := make(map[string]error)

	_, err := sdk.AccAddressFromBech32(c.UserAddress)
	if err != nil {
		errs["UserAddress"] = err
	}

	if c.ChainId == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	if c.Amount.IsNil() || !c.Amount.IsPositive() {
		errs["Amount"] = ErrNotPositive
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}
