package types

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
)

func MustMarshalValidator(cdc codec.BinaryCodec, validator Validator) []byte {
	return cdc.MustMarshal(&validator)
}

func UnmarshalValidator(cdc codec.BinaryCodec, value []byte) (Validator, error) {
	if len(value) == 0 {
		return Validator{}, errors.New("unable to unmarshal zero-length byte slice")
	}
	validator := Validator{}
	err := cdc.Unmarshal(value, &validator)
	return validator, err
}

func MustUnmarshalValidator(cdc codec.BinaryCodec, value []byte) Validator {
	validator, err := UnmarshalValidator(cdc, value)
	if err != nil {
		panic(err)
	}
	return validator
}