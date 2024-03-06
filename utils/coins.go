package utils

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	transfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
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

// DeriveIbcDenom mirrors getDenomForThisChain from the packet-forward-middleware/v5, used under MIT License.
// See: https://github.com/strangelove-ventures/packet-forward-middleware/blob/86f045c12cc48ffc1f016ff122b89a9f6ac8ed63/router/ibc_middleware.go#L104
func DeriveIbcDenom(port, channel, counterpartyPort, counterpartyChannel, denom string) string {
	counterpartyPrefix := transfertypes.GetDenomPrefix(counterpartyPort, counterpartyChannel)
	if strings.HasPrefix(denom, counterpartyPrefix) {
		unwoundDenom := denom[len(counterpartyPrefix):]
		denomTrace := transfertypes.ParseDenomTrace(unwoundDenom)
		if denomTrace.Path == "" {
			return unwoundDenom
		}
		return denomTrace.IBCDenom()
	}
	prefixedDenom := transfertypes.GetDenomPrefix(port, channel) + denom
	return transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom()
}
