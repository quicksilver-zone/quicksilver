package types

import (
	"fmt"
	"gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyClaimsEnabled     = []byte("ClaimsEnabled")
	DefaultClaimsEnabled = false
)

// ParamKeyTable for claimsmanager module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new claimsmanager Params instance.
func NewParams() Params {
	return Params{
		DefaultClaimsEnabled,
	}
}

// DefaultParams default claimsmanager params.
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyClaimsEnabled, &p.ClaimsEnabled, validateBoolean),
	}
}

func validateBoolean(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// Validate validates params.
func (p *Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p *Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
