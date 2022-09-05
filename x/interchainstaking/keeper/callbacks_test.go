package keeper

import (
	"testing"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/stretchr/testify/require"
)

func TestCoinFromRequestKey(t *testing.T) {
	accAddr := utils.GenerateAccAddressForTest()
	prefix := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
	query := append(prefix, []byte("denom")...)

	coin, err := coinFromRequestKey(query, accAddr)
	require.NoError(t, err)
	require.Equal(t, "denom", coin.Denom)
}
