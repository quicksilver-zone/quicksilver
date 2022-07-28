package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var ()

// ParamTable for participationrewards module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new ics Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams default ics params
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		//paramtypes.NewParamSetPair(Key, &p.Attribute, validateAttrib),
	}
}

// validate params.
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
