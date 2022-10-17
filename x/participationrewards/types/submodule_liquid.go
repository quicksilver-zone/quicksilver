package types

import "github.com/ingenuity-build/quicksilver/internal/multierror"

type LiquidAllowedDenomProtocolData struct {
	ChainID       string
	OriginChainID string
	Denom         string
	LocalDenom    string
}

func (lpd LiquidAllowedDenomProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	if len(lpd.ChainID) == 0 {
		errors["ChainID"] = ErrUndefinedAttribute
	}

	if len(lpd.OriginChainID) == 0 {
		errors["OriginChainID"] = ErrUndefinedAttribute
	}

	if len(lpd.Denom) == 0 {
		errors["Denom"] = ErrUndefinedAttribute
	}

	if len(lpd.LocalDenom) == 0 {
		errors["LocalDenom"] = ErrUndefinedAttribute
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
