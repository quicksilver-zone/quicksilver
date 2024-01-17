package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func addressAndDenomFromBalancesStore(key []byte) (sdk.AccAddress, string, error) {
	if len(key) == 0 {
		return nil, "", banktypes.ErrInvalidKey
	}

	kv.AssertKeyAtLeastLength(key, 1)

	addrBound := int(key[0])

	if len(key)-1 < addrBound {
		return nil, "", banktypes.ErrInvalidKey
	}

	return key[1 : addrBound+1], string(key[addrBound+1:]), nil
}

func DenomFromRequestKey(query []byte, accAddr sdk.AccAddress) (string, error) {
	balancesStore := query[1:]

	gotAccAddress, denom, err := addressAndDenomFromBalancesStore(balancesStore)
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
