package types

import (
	"crypto/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GenerateAccAddressForTest() sdk.AccAddress {
	size := 32 // change the length of the generated random string here

	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		panic(err)
	}

	return sdk.AccAddress(rb)
}
