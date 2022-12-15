package utils

import (
	"bytes"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	idx := bytes.Index(query, accAddr)
	if idx == -1 {
		return "", errors.New("AccountBalanceCallback: invalid request query")
	}
	denom := string(query[idx+len(accAddr):])
	if err := sdk.ValidateDenom(denom); err != nil {
		return "", err
	}

	return denom, nil
}
