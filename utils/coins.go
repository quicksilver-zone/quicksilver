package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/quicksilver-zone/quicksilver/v7/utils/bankutils"
)

func DenomFromRequestKey(query []byte, accAddr sdk.AccAddress) (string, error) {
	balancesStore := query[1:]
	gotAccAddress, denom, err := bankutils.AddressAndDenomFromBalancesStore(balancesStore)
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

func DeriveIbcDenom(port, channel, denom string) string {
	return DeriveIbcDenomTrace(port, channel, denom).IBCDenom()
}

func DeriveIbcDenomTrace(port, channel, denom string) ibctransfertypes.DenomTrace {
	// generate denomination prefix
	sourcePrefix := ibctransfertypes.GetDenomPrefix(port, channel)
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + denom

	// construct the denomination trace from the full raw denomination
	denomTrace := ibctransfertypes.ParseDenomTrace(prefixedDenom)

	return denomTrace
}
