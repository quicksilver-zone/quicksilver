package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccFromVal(t *testing.T) {
	operatorAddress := "cosmosvaloper1c4k24jzduc365kywrsvf5ujz4ya6mwympnc4en"
	accAddr := ValAddrToAccAddr(operatorAddress)
	require.Equal(t, "cosmos1c4k24jzduc365kywrsvf5ujz4ya6mwymy8vq4q", accAddr)
}
