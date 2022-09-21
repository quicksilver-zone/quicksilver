package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Default ics params
var (
	DefaultDelegateAccountCount uint64  = 100
	DefaultDepositInterval      uint64  = 20
	DefaultValidatorSetInterval uint64  = 200
	DefaultCommissionRate       sdk.Dec = sdk.MustNewDecFromStr("0.025")

	// KeyDelegateAccountCount is store's key for DelegateAccountCount option
	KeyDelegateAccountCount = []byte("DelegateAccountCount")
	// KeyDepositInterval is store's key for the DepositInterval option
	KeyDepositInterval = []byte("DepositInterval")
	// KeyValidatorSetInterval is store's key for the ValidatorSetInterval option
	KeyValidatorSetInterval = []byte("ValidatorSetInterval")
	// KeyCommissionRate is store's key for the CommissionRate option
	KeyCommissionRate = []byte("CommissionRate")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.LegacyAmino, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}

	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.LegacyAmino, value []byte) (params Params, err error) {
	if bytes.Equal(value, []byte("")) {
		return params, fmt.Errorf("unable to unmarshal empty byte slice")
	}
	err = cdc.Unmarshal(value, &params)
	if err != nil {
		return
	}

	return
}

func validateParams(i interface{}) error {
	v, ok := i.(Params)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.DepositInterval <= 0 {
		return fmt.Errorf("deposit interval must be positive: %d", v.DepositInterval)
	}

	if v.ValidatorsetInterval <= 0 {
		return fmt.Errorf("valset interval must be positive: %d", v.ValidatorsetInterval)
	}

	if v.CommissionRate.IsNil() {
		return fmt.Errorf("commission rate must be non-nil")
	}

	if v.CommissionRate.IsNegative() {
		return fmt.Errorf("commission rate must be non-negative: %s", v.CommissionRate.String())
	}
	return nil
}

// NewParams creates a new ics Params instance
func NewParams(
	delegateAccountCount uint64,
	depositInterval uint64,
	valsetInterval uint64,
	commissionRate sdk.Dec,
) Params {
	return Params{
		DepositInterval:      depositInterval,
		ValidatorsetInterval: valsetInterval,
		CommissionRate:       commissionRate,
	}
}

// DefaultParams default ics params
func DefaultParams() Params {
	return NewParams(
		DefaultDelegateAccountCount,
		DefaultDepositInterval,
		DefaultValidatorSetInterval,
		DefaultCommissionRate,
	)
}

// ParamKeyTable for ics module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDepositInterval, &p.DepositInterval, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyValidatorSetInterval, &p.ValidatorsetInterval, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyCommissionRate, &p.CommissionRate, validateNonNegativeDec),
	}
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validatePositiveInt(i interface{}) error {
	intval, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if intval <= 0 {
		return fmt.Errorf("invalid (non-positive) parameter value: %d", intval)
	}
	return nil
}

func validateNonNegativeDec(i interface{}) error {
	intval, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if intval.IsNegative() {
		return fmt.Errorf("invalid (negative) parameter value: %d", intval)
	}
	return nil
}
