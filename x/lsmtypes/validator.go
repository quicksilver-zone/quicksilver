package lsmtypes

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"cosmossdk.io/errors"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// String implements the Stringer interface for a Validator object.
func (v Validator) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

func (v Validator) GetConsAddr() ([]byte, error) {
	pk, ok := v.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		fmt.Println("got", pk)
		return nil, errors.Wrapf(sdkerrors.ErrInvalidType, "expecting cryptotypes.PubKey, got %T", pk)
	}

	return pk.Address().Bytes(), nil
}

func (v Validator) GetCommission() sdkmath.LegacyDec { return v.Commission.Rate }
func (v Validator) IsJailed() bool                   { return v.Jailed }
