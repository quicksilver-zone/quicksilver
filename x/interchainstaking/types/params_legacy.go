package types

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Default ics params.
var (
	DefaultDepositInterval      uint64  = 20
	DefaultValidatorSetInterval uint64  = 200
	DefaultCommissionRate       sdk.Dec = sdk.MustNewDecFromStr("0.025")
	DefaultUnbondingEnabled             = false

	// KeyDepositInterval is store's key for the DepositInterval option.
	KeyDepositInterval = []byte("DepositInterval")
	// KeyValidatorSetInterval is store's key for the ValidatorSetInterval option.
	KeyValidatorSetInterval = []byte("ValidatorSetInterval")
	// KeyCommissionRate is store's key for the CommissionRate option.
	KeyCommissionRate = []byte("CommissionRate")
	// KeyUnbondingEnabled is a global flag to indicated whether unbonding txs are permitted.
	KeyUnbondingEnabled = []byte("UnbondingEnabled")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// MustUnmarshalParams unmarshals the current interchainstaking params value from store key or panic.
func MustUnmarshalParams(cdc *codec.LegacyAmino, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}

	return params
}

// UnmarshalParams unmarshals the current interchainstaking params value from store key.
func UnmarshalParams(cdc *codec.LegacyAmino, value []byte) (params Params, err error) {
	if len(value) == 0 {
		return params, errors.New("unable to unmarshal empty byte slice")
	}
	err = cdc.Unmarshal(value, &params)
	return params, err
}

// NewParams creates a new ics Params instance.
func NewParams(
	depositInterval uint64,
	valsetInterval uint64,
	commissionRate sdk.Dec,
	unbondingEnabled bool,
) Params {
	return Params{
		DepositInterval:      depositInterval,
		ValidatorsetInterval: valsetInterval,
		CommissionRate:       commissionRate,
		UnbondingEnabled:     unbondingEnabled,
	}
}

// DefaultParams default ics params.
func DefaultParams() Params {
	return NewParams(
		DefaultDepositInterval,
		DefaultValidatorSetInterval,
		DefaultCommissionRate,
		DefaultUnbondingEnabled,
	)
}

// Validate validates params.
func (p Params) Validate() error {
	if err := validatePositiveInt(p.DepositInterval); err != nil {
		return fmt.Errorf("invalid deposit interval: %w", err)
	}

	if err := validatePositiveInt(p.ValidatorsetInterval); err != nil {
		return fmt.Errorf("invalid valset interval: %w", err)
	}

	if err := validateNonNegativeDec(p.CommissionRate); err != nil {
		return fmt.Errorf("invalid commission rate: %w", err)
	}

	return nil
}

// ParamKeyTable for ics module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDepositInterval, &p.DepositInterval, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyValidatorSetInterval, &p.ValidatorsetInterval, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyCommissionRate, &p.CommissionRate, validateNonNegativeDec),
		paramtypes.NewParamSetPair(KeyUnbondingEnabled, &p.UnbondingEnabled, validateBoolean),
	}
}

func (p ParamsV1) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDepositInterval, &p.DepositInterval, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyValidatorSetInterval, &p.ValidatorsetInterval, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyCommissionRate, &p.CommissionRate, validateNonNegativeDec),
	}
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p) //nolint:errcheck not needed
	return string(out)
}

// String implements the Stringer interface.
func (p ParamsV1) String() string {
	out, _ := yaml.Marshal(p) //nolint:errcheck not needed
	return string(out)
}

func validateBoolean(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
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
	dec, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if dec.IsNil() {
		return fmt.Errorf("invalid (nil) parameter value")
	}

	if dec.IsNegative() {
		return fmt.Errorf("invalid (negative) parameter value: %s", dec.String())
	}
	return nil
}
