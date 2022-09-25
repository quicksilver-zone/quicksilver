package types

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// NewDelegation creates a new delegation object
func NewDelegation(delegatorAddr string, validatorAddr string, amount sdk.Coin) Delegation {
	return Delegation{
		DelegationAddress: delegatorAddr,
		ValidatorAddress:  validatorAddr,
		Amount:            amount,
		Height:            0,
		RedelegationEnd:   0,
	}
}

// MustMarshalDelegation returns the delegation bytes.
// This function will panic on failure.
func MustMarshalDelegation(cdc codec.BinaryCodec, delegation Delegation) []byte {
	return cdc.MustMarshal(&delegation)
}

// MustUnmarshalDelegation return the unmarshaled delegation from bytes.
// This function will panic on failure.
func MustUnmarshalDelegation(cdc codec.BinaryCodec, value []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}

	return delegation
}

// return the delegation
func UnmarshalDelegation(cdc codec.BinaryCodec, value []byte) (delegation Delegation, err error) {
	if bytes.Equal(value, []byte("")) {
		return Delegation{}, fmt.Errorf("unable to unmarshal zero-length byte slice")
	}
	err = cdc.Unmarshal(value, &delegation)
	return delegation, err
}

// This function will panic on failure.
func (d Delegation) GetDelegatorAddr() sdk.AccAddress {
	_, delAddr, err := bech32.DecodeAndConvert(d.DelegationAddress)
	if err != nil {
		panic(err)
	}
	return delAddr
}

// This function will panic on failure.
func (d Delegation) GetValidatorAddr() sdk.ValAddress {
	_, valAddr, err := bech32.DecodeAndConvert(d.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return valAddr
}

type ValidatorIntents []*ValidatorIntent

func (vi ValidatorIntents) Sort() ValidatorIntents {
	sort.SliceStable(vi, func(i, j int) bool {
		return vi[i].ValoperAddress < vi[j].ValoperAddress
	})
	return vi
}

func (vi ValidatorIntents) GetForValoper(valoper string) (*ValidatorIntent, bool) {
	for _, i := range vi.Sort() {
		if i.ValoperAddress == valoper {
			return i, true
		}
	}
	return nil, false
}

func (vi ValidatorIntents) SetForValoper(valoper string, intent *ValidatorIntent) ValidatorIntents {
	for idx, i := range vi.Sort() {
		if i.ValoperAddress == valoper {
			var part []*ValidatorIntent
			if idx != 0 {
				part = vi[:idx-1]
			}
			vi = append(part, vi[idx:]...)
		}
	}
	vi = append(vi, intent)

	return vi.Sort()
}

func (vi ValidatorIntents) MustGetForValoper(valoper string) *ValidatorIntent {
	intent, found := vi.GetForValoper(valoper)
	if !found {
		panic("could not find intent for valoper")
	}
	return intent
}
