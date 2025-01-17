package types

import (
	"github.com/ingenuity-build/multierror"
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

	if lpd.RegisteredZoneChainID == "" {
		errs["RegisteredZoneChainID"] = ErrUndefinedAttribute
	}

	if lpd.QAssetDenom == "" {
		errs["QAssetDenom"] = ErrUndefinedAttribute
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func (lpd *LiquidAllowedDenomProtocolData) GenerateKey() []byte {
	return []byte(lpd.ChainID + "_" + lpd.IbcDenom)
}
