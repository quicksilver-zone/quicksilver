package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Default ics params
var (
	DefaultDelegateAccountCount uint64  = 100
	DefaultDelegateAccountSplit uint64  = 10
	DefaultDepositInterval      uint64  = 50
	DefaultDelegateInterval     uint64  = 100
	DefaultDelegationsInterval  uint64  = 200
	DefaultValidatorSetInterval uint64  = 200
	DefaultCommissionRate       sdk.Dec = func() sdk.Dec { v, _ := sdk.NewDecFromStr("0.02"); return v }()

	// KeyDelegateAccountCount is store's key for DelegateAccountCount option
	KeyDelegateAccountCount = []byte("DelegateAccountCount")
	// KeyDelegateAccountSplit is store's key for the DelegateAccountSplit option
	KeyDelegateAccountSplit = []byte("DelegateAccountSplit")
	// KeyDepositInterval is store's key for the DepositInterval option
	KeyDepositInterval = []byte("DepositInterval")
	// KeyDelegateInterval is store's key for the DelegateInterval option
	KeyDelegateInterval = []byte("DelegateInterval")
	// KeyDelegationsInterval is store's key for the DelegationsInterval option
	KeyDelegationsInterval = []byte("DelegationsInterval")
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

	if v.DelegationAccountCount <= 0 {
		return fmt.Errorf("delegate account count must be positive: %d", v.DelegationAccountCount)
	}

	if v.DelegationAccountSplit <= 0 {
		return fmt.Errorf("delegate account split must be positive: %d", v.DelegationAccountSplit)
	}

	if v.DelegationAccountSplit > v.DelegationAccountCount {
		return fmt.Errorf("delegate account split must be less than or equal to delegate account count: %d", v.DelegationAccountCount)
	}

	if v.DepositInterval <= 0 {
		return fmt.Errorf("deposit interval must be positive: %d", v.DepositInterval)
	}

	if v.DelegateInterval <= 0 {
		return fmt.Errorf("delegate interval must be positive: %d", v.DelegateInterval)
	}

	if v.DelegationsInterval <= 0 {
		return fmt.Errorf("delegations interval must be positive: %d", v.DelegationsInterval)
	}

	if v.ValidatorsetInterval <= 0 {
		return fmt.Errorf("valset interval must be positive: %d", v.ValidatorsetInterval)
	}

	if v.CommissionRate.IsNegative() {
		return fmt.Errorf("commission rate must be non-negative: %s", v.CommissionRate.String())
	}
	return nil
}

// NewParams creates a new ics Params instance
func NewParams(
	delegate_account_count uint64,
	delegate_account_split uint64,
	deposit_interval uint64,
	delegate_interval uint64,
	delegations_interval uint64,
	valset_interval uint64,
	commission_rate sdk.Dec,
) Params {
	return Params{
		DelegationAccountCount: delegate_account_count,
		DelegationAccountSplit: delegate_account_split,
		DepositInterval:        deposit_interval,
		DelegateInterval:       delegate_interval,
		DelegationsInterval:    delegations_interval,
		ValidatorsetInterval:   valset_interval,
		CommissionRate:         commission_rate,
	}
}

// DefaultParams default ics params
func DefaultParams() Params {
	return NewParams(
		DefaultDelegateAccountCount,
		DefaultDelegateAccountSplit,
		DefaultDelegateInterval,
		DefaultDelegateInterval,
		DefaultDelegationsInterval,
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
		paramtypes.NewParamSetPair(KeyDelegateAccountCount, &p.DelegationAccountCount, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyDelegateAccountSplit, &p.DelegationAccountSplit, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyDepositInterval, &p.DepositInterval, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyDelegateInterval, &p.DelegateInterval, validatePositiveInt),
		paramtypes.NewParamSetPair(KeyDelegationsInterval, &p.DelegationsInterval, validatePositiveInt),
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
		return fmt.Errorf("invalid (non-positve) parameter value: %d", intval)
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
