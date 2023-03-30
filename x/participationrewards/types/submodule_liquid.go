package types

import (
	"strings"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

// LiquidAllowedDenomProtocolData defines protocol state to track off-chain
// liquid qAssets.
type LiquidAllowedDenomProtocolData struct {
	// The chain on which the qAssets reside currently.
	ChainID string
	// The chain for which the qAssets were issued.
	RegisteredZoneChainID string
	// The IBC denom.
	IbcDenom string
	// The qAsset denom.
	QAssetDenom string
}

func (lpd LiquidAllowedDenomProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	if lpd.ChainID == "" {
		errors["ChainID"] = ErrUndefinedAttribute
	}

	if len(strings.Split(lpd.ChainID, "-")) < 2 {
		errors["ChainID"] = ErrInvalidChainID
	}

	if lpd.RegisteredZoneChainID == "" {
		errors["RegisteredZoneChainID"] = ErrUndefinedAttribute
	}

	if len(strings.Split(lpd.RegisteredZoneChainID, "-")) < 2 {
		errors["RegisteredZoneChainID"] = ErrInvalidChainID
	}

	if lpd.QAssetDenom == "" {
		errors["QAssetDenom"] = ErrUndefinedAttribute
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
