package types

import (
	"go.uber.org/multierr"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils"
)

func NewClaim(address, chainID string, module ClaimType, srcChainID string, amount math.Int) Claim {
	return Claim{UserAddress: address, ChainId: chainID, Module: module, SourceChainId: srcChainID, Amount: amount}
}

// ValidateBasic performs stateless validation of a Claim.
func (c *Claim) ValidateBasic() error {
	errs := make(map[string]error)

	_, err := sdk.AccAddressFromBech32(c.UserAddress)
	if err != nil {
		errs["userAddress"] = err
	}

	if c.ChainId == "" {
		errs["chainID"] = ErrUndefinedAttribute
	}

	if c.Amount.IsNil() || !c.Amount.IsPositive() {
		errs["amount"] = ErrNotPositive
	}

	if len(errs) > 0 {
		return multierr.Combine(utils.ErrorMapToSlice(errs)...)
	}

	return nil
}
