package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func DenomFromRequestKey(query []byte, accAddr sdk.AccAddress) (string, error) {
	balancesStore := query[1:]
	gotAccAddress, denom, err := banktypes.AddressAndDenomFromBalancesStore(balancesStore)
	if err != nil {
		return "", err
	}

	if denom == "" {
		return "", fmt.Errorf("key contained no denom")
	}

	if !gotAccAddress.Equals(accAddr) {
		return "", fmt.Errorf("account mismatch; expected %s, got %s", accAddr.String(), gotAccAddress.String())
	}

	return denom, nil
}
