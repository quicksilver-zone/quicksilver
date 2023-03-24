package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	BaseCoinUnit = "uqck"
)

// Parameter store keys.
var (
	KeyDenomCreationFee = []byte("DenomCreationFee")
)

// ParamKeyTable returns KeyTable for tokenfactory module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(denomCreationFee sdk.Coins) Params {
	return Params{
		DenomCreationFee: denomCreationFee,
	}
}

// Validate validates params.
func (p Params) Validate() error {
	err := validateDenomCreationFee(p.DenomCreationFee)

	return err
}

// ParamSetPairs implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDenomCreationFee, &p.DenomCreationFee, validateDenomCreationFee),
	}
}

func validateDenomCreationFee(i any) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("invalid denom creation fee: %+v", i)
	}

	return nil
}

func DefaultParams() Params {
	return NewParams(
		sdk.NewCoins(sdk.NewInt64Coin(BaseCoinUnit, 10_000_000)),
	)
}
