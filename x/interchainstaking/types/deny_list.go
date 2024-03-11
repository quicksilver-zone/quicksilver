package types

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
)

func NewValidatorDenyListForZone(chainID string) ValidatorDenyList {
	return ValidatorDenyList{
		ChainId:    chainID,
		DeniedVals: []Validator{},
	}
}

// UnmarshalValidatorDenyList return the deny list from bytes value
func UnmarshalValidatorDenyList(cdc codec.BinaryCodec, value []byte) (ValidatorDenyList, error) {
	denyList := ValidatorDenyList{}
	if len(value) == 0 {
		return ValidatorDenyList{}, errors.New("unable to unmarshal zero-length byte slice")
	}
	err := cdc.Unmarshal(value, &denyList)
	return denyList, err
}

// MustUnmarshalDelegation return the unmarshaled delegation from bytes, panic on failure
func MustUnmarshalDenyList(cdc codec.BinaryCodec, value []byte) ValidatorDenyList {
	denyList, err := UnmarshalValidatorDenyList(cdc, value)
	if err != nil {
		panic(err)
	}
	return denyList
}

func MustMarshalDenyList(cdc codec.BinaryCodec, denyList ValidatorDenyList) []byte {
	return cdc.MustMarshal(&denyList)
}
