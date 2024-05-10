package utils

import sdk "github.com/cosmos/cosmos-sdk/types"

func ValAddrToAccAddr(operatorAddress string) string {
	addrBytes, err := sdk.ValAddressFromBech32(operatorAddress)
	if err != nil {
		panic(err)
	}
	return sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), addrBytes)
}
