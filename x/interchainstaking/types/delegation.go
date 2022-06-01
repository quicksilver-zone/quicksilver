package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// NewDelegation creates a new delegation object
//nolint:interfacer
func NewDelegation(delegatorAddr string, validatorAddr string, amount sdk.Coin) Delegation {
	return Delegation{
		DelegationAddress: delegatorAddr,
		ValidatorAddress:  validatorAddr,
		Amount:            amount,
		Height:            0,
		RedelegationEnd:   0,
	}
}

// MustMarshalDelegation returns the delegation bytes. Panics if fails
func MustMarshalDelegation(cdc codec.BinaryCodec, delegation Delegation) []byte {
	return cdc.MustMarshal(&delegation)
}

// MustUnmarshalDelegation return the unmarshaled delegation from bytes.
// Panics if fails.
func MustUnmarshalDelegation(cdc codec.BinaryCodec, value []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}

	return delegation
}

// return the delegation
func UnmarshalDelegation(cdc codec.BinaryCodec, value []byte) (delegation Delegation, err error) {
	err = cdc.Unmarshal(value, &delegation)
	return delegation, err
}

func (d Delegation) GetDelegatorAddr() sdk.AccAddress {
	_, delAddr, err := bech32.DecodeAndConvert(d.DelegationAddress)
	if err != nil {
		panic(err)
	}
	return delAddr
}

func (d Delegation) GetValidatorAddr() sdk.ValAddress {
	_, valAddr, err := bech32.DecodeAndConvert(d.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return valAddr
}
