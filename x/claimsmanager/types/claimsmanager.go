package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

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

	if c.Amount <= 0 {
		errs["Amount"] = ErrNotPositive
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}
