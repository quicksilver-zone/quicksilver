package utils

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DecodeCwNamespacedKey decodes a cw namespaced key into an address and parts.
// The key is expected to be in the format of a cw namespaced key, which is
// a prefix (1 byte), a 32 byte address, a null terminator (1 byte), 1..n parts
// (each part is a length-prefixed byte string), followed by the final key (not
// length-prefixed).
// As specified here: https://github.com/webmaster128/key-namespacing#nesting
//
// As we are unable to determine the length of the final key, we must pass in the
// number of parts we expect out (including the final key).
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
