package lockup

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

// Parameter store keys.
var (
	KeyForceUnlockAllowedAddresses = []byte("ForceUnlockAllowedAddresses")

	_ paramtypes.ParamSet = &Params{}
)

// ParamTable for lockup module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(forceUnlockAllowedAddresses []string) Params {
	return Params{
		ForceUnlockAllowedAddresses: forceUnlockAllowedAddresses,
	}
}

// DefaultParams returns default lockup module parameters.
func DefaultParams() Params {
	return Params{
		ForceUnlockAllowedAddresses: []string{},
	}
}

// validate params.
func (p Params) Validate() error {
	if err := validateAddresses(p.ForceUnlockAllowedAddresses); err != nil {
		return err
	}
	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyForceUnlockAllowedAddresses, &p.ForceUnlockAllowedAddresses, validateAddresses),
	}
}

func validateAddresses(i interface{}) error {
	addresses, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, address := range addresses {
		_, err := addressutils.AccAddressFromBech32(address, "")
		if err != nil {
			return err
		}
	}

	return nil
}
