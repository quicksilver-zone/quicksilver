package keeper

import (
	"testing"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	"github.com/stretchr/testify/require"
)

// tests that {} is an invalid string, and that an error is thrown when unmarshalled.
// see: https://github.com/ingenuity-build/quicksilver/issues/214
func TestUnmarshalProtocolDataRejectsZeroLengthJson(t *testing.T) {
	_, err := UnmarshalProtocolData(types.ProtocolDataOsmosisPool, []byte("{}"))
	require.Error(t, err)
}
