package utils

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DecodeCwNamespacedKey(key []byte, numParts int) (sdk.AccAddress, [][]byte, error) {

	if len(key) < 37 { // prefix (1 byte) + 32 byte address + 1 byte null terminator + len (1 byte) + min 1 byte namespace + min 1 byte key
		return nil, nil, errors.New("invalid key length")
	}

	if key[0] != 0x03 {
		return nil, nil, errors.New("invalid prefix")
	}

	addressBytes := key[1:33]

	if key[33] != 0x00 {
		return nil, nil, errors.New("expected null terminator after address")
	}

	address := sdk.AccAddress(addressBytes)
	parts := [][]byte{}
	pointer := 34
	for pointer < len(key) && numParts-1 > len(parts) {
		length := int(key[pointer])
		pointer++
		part := key[pointer : pointer+length]
		pointer += length
		parts = append(parts, part)
	}
	parts = append(parts, key[pointer:])

	return address, parts, nil
}
