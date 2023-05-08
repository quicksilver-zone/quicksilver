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

func (lpd *LiquidAllowedDenomProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if lpd.ChainID == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	if len(strings.Split(lpd.ChainID, "-")) < 2 {
		errs["ChainID"] = ErrInvalidChainID
	}

	if lpd.RegisteredZoneChainID == "" {
		errs["RegisteredZoneChainID"] = ErrUndefinedAttribute
	}

	if len(strings.Split(lpd.RegisteredZoneChainID, "-")) < 2 {
		errs["RegisteredZoneChainID"] = ErrInvalidChainID
	}

	if lpd.QAssetDenom == "" {
		errs["QAssetDenom"] = ErrUndefinedAttribute
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}
