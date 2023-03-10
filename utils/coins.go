package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// coinFromRequestKey parses
func CoinFromRequestKey(query []byte, accAddr sdk.AccAddress) (sdk.Coin, error) {
	denom, err := DenomFromRequestKey(query, accAddr)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(denom, sdk.ZeroInt()), nil
}

func DenomFromRequestKey(query []byte, accAddr sdk.AccAddress) (string, error) {
	balancesStore := query[1:]
	accAddr2, denom, err := banktypes.AddressAndDenomFromBalancesStore(balancesStore)
	if err != nil {
		return "", err
	}

	if !accAddr2.Equals(accAddr) {
		return "", fmt.Errorf("account mismatch; expected %s, got %s", accAddr.String(), accAddr2.String())
	}

	return denom, nil
}
