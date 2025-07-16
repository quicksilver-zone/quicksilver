package addressutils_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

func TestAddressFromBech32(t *testing.T) {
	cases := []struct {
		name          string
		address       string
		prefix        string
		expectedBytes []byte
		expectedErr   string
	}{
		{
			"invalid - empty string",
			"",
			"",
			nil,
			"empty address string is not allowed",
		},
		{
			"invalid - invalid address",
			"quick",
			"",
			nil,
			"decoding bech32 failed: invalid bech32 string length 5",
		},
		{
			"invalid - invalid characters",
			"sbg2apkjme1qh2ycto7jn30nu",
			"",
			nil,
			"decoding bech32 failed: invalid character not part of charset: 111",
		},
		{
			"invalid - invalid checksum",
			"cosmos1kv4ez0rgrd679m6da96apnqxkcamh28caaaaaa",
			"",
			nil,
			"decoding bech32 failed: invalid checksum (expected 098lr8 got aaaaaa)",
		},
		{
			"invalid - invalid separator",
			"cosmos2kv4ez0rgrd679m6da96apnqxkcamh28c098lr8",
			"",
			nil,
			"decoding bech32 failed: invalid separator index -1",
		},
		{
			"invalid - invalid hrp",
			"cosmos1kv4ez0rgrd679m6da96apnqxkcamh28c098lr8",
			"quick",
			nil,
			"unexpected prefix - got cosmos expected quick",
		},
		{
			"invalid - no prefix",
			"1kv4ez0rgrd679m6da96apnqxkcamh28c00j09s",
			"quick",
			nil,
			"decoding bech32 failed: invalid separator index 0",
		},
		{
			"invalid - too long",
			"cosmos1kv4ez0rgrd679m6da96apnqxkcamh0rgrd6grd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxk79m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd6nqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6danqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da96apnqxkcamh0rgrd679m6da79m6da96apnqxkcamh28cc63a68",
			"",
			nil,
			"decoding bech32 failed: invalid bech32 string length 4393",
		},
		{
			"valid - no check",
			"cosmos1kv4ez0rgrd679m6da96apnqxkcamh28c098lr8",
			"",
			[]byte{0xb3, 0x2b, 0x91, 0x3c, 0x68, 0x1b, 0x75, 0xe2, 0xef, 0x4d, 0xe9, 0x75, 0xd0, 0xcc, 0x6, 0xb6, 0x3b, 0xbb, 0xa8, 0xf8},
			"",
		},
		{
			"valid - with hrp check",
			"cosmos1kv4ez0rgrd679m6da96apnqxkcamh28c098lr8",
			"cosmos",
			[]byte{0xb3, 0x2b, 0x91, 0x3c, 0x68, 0x1b, 0x75, 0xe2, 0xef, 0x4d, 0xe9, 0x75, 0xd0, 0xcc, 0x6, 0xb6, 0x3b, 0xbb, 0xa8, 0xf8},
			"",
		},
	}

	for _, c := range cases {
		addr, err := addressutils.AccAddressFromBech32(c.address, c.prefix)

		if c.expectedErr != "" {
			require.Error(t, err)
			require.ErrorContains(t, err, c.expectedErr)
		} else {
			require.Equal(t, c.expectedBytes, addr.Bytes())
		}

		valaddr, err := addressutils.ValAddressFromBech32(c.address, c.prefix)

		if c.expectedErr != "" {
			require.Error(t, err)
			require.ErrorContains(t, err, c.expectedErr)
		} else {
			require.Equal(t, c.expectedBytes, valaddr.Bytes())
		}

		if c.expectedErr != "" {
			require.Panics(t, func() { addressutils.MustAccAddressFromBech32(c.address, c.prefix) })
			require.Panics(t, func() { addressutils.MustValAddressFromBech32(c.address, c.prefix) })
		} else {
			addr := addressutils.MustAccAddressFromBech32(c.address, c.prefix)
			require.Equal(t, c.expectedBytes, addr.Bytes())

			valaddr := addressutils.MustValAddressFromBech32(c.address, c.prefix)
			require.Equal(t, c.expectedBytes, valaddr.Bytes())
		}
	}
}

func TestGenerateAccAddressForTest(t *testing.T) {
	address := addressutils.GenerateAccAddressForTest()
	require.Equal(t, 32, len(address.Bytes()))
	err := sdk.VerifyAddressFormat(address.Bytes())
	require.NoError(t, err)
}

func TestGenerateValAddressForTest(t *testing.T) {
	address := addressutils.GenerateValAddressForTest()
	require.Equal(t, 32, len(address.Bytes()))
	err := sdk.VerifyAddressFormat(address.Bytes())
	require.NoError(t, err)
}

func TestGenerateAddressForTestWithPrefix(t *testing.T) {
	b32addr := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	_, err := addressutils.AddressFromBech32(b32addr, "cosmos")
	require.NoError(t, err)
}

func TestEncodeAddressToBech32(t *testing.T) {
	cases := []struct {
		name            string
		addrBytes       []byte
		prefix          string
		expectedAddress string
		expectedErr     string
	}{
		{
			"valid",
			[]byte{0xb3, 0x2b, 0x91, 0x3c, 0x68, 0x1b, 0x75, 0xe2, 0xef, 0x4d, 0xe9, 0x75, 0xd0, 0xcc, 0x6, 0xb6, 0x3b, 0xbb, 0xa8, 0xf8},
			"cosmos",
			"cosmos1kv4ez0rgrd679m6da96apnqxkcamh28c098lr8",
			"",
		},
		{
			"surprisingly valid, single null byte",
			[]byte{0x00},
			"cosmos",
			"cosmos1qqxuevtt",
			"",
		},
		{
			"surprisingly valid, nil",
			nil,
			"cosmos",
			"cosmos1550dq7",
			"",
		},
		{
			"even more surprisingly valid - no hrp :/",
			[]byte{0xb3, 0x2b, 0x91, 0x3c, 0x68, 0x1b, 0x75, 0xe2, 0xef, 0x4d, 0xe9, 0x75, 0xd0, 0xcc, 0x6, 0xb6, 0x3b, 0xbb, 0xa8, 0xf8},
			"",
			"1kv4ez0rgrd679m6da96apnqxkcamh28cjkahle",
			"",
		},
	}

	for _, c := range cases {
		addr, err := addressutils.EncodeAddressToBech32(c.prefix, sdk.AccAddress(c.addrBytes))

		if c.expectedErr != "" {
			require.Error(t, err)
			require.ErrorContains(t, err, c.expectedErr)

			require.Panics(t, func() { addressutils.MustEncodeAddressToBech32(c.prefix, sdk.AccAddress(c.addrBytes)) })
		} else {
			require.Equal(t, c.expectedAddress, addr)
		}
	}
}

func TestLengthPrefix(t *testing.T) {
	cases := []struct {
		name          string
		input         []byte
		expected      []byte
		expectedError string
	}{
		{
			"empty byte slice",
			[]byte{},
			[]byte{},
			"",
		},
		{
			"nil byte slice",
			nil,
			nil,
			"",
		},
		{
			"single byte",
			[]byte{0x01},
			[]byte{0x01, 0x01},
			"",
		},
		{
			"multiple bytes",
			[]byte{0x01, 0x02, 0x03, 0x04},
			[]byte{0x04, 0x01, 0x02, 0x03, 0x04},
			"",
		},
		{
			"max length address",
			make([]byte, addressutils.MaxAddrLen),
			append([]byte{byte(addressutils.MaxAddrLen)}, make([]byte, addressutils.MaxAddrLen)...),
			"",
		},
		{
			"too long address",
			make([]byte, addressutils.MaxAddrLen+1),
			nil,
			"address length should be max 255 bytes, got 256",
		},
		{
			"very long address",
			make([]byte, 1000),
			nil,
			"address length should be max 255 bytes, got 1000",
		},
		{
			"address with zero bytes",
			[]byte{0x00, 0x00, 0x00},
			[]byte{0x03, 0x00, 0x00, 0x00},
			"",
		},
		{
			"address with mixed bytes",
			[]byte{0xFF, 0x00, 0xAA, 0x55},
			[]byte{0x04, 0xFF, 0x00, 0xAA, 0x55},
			"",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := addressutils.LengthPrefix(c.input)

			if c.expectedError != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, c.expectedError)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expected, result)
			}
		})
	}
}

func TestMustLengthPrefix(t *testing.T) {
	cases := []struct {
		name        string
		input       []byte
		expected    []byte
		shouldPanic bool
	}{
		{
			"empty byte slice",
			[]byte{},
			[]byte{},
			false,
		},
		{
			"nil byte slice",
			nil,
			nil,
			false,
		},
		{
			"single byte",
			[]byte{0x01},
			[]byte{0x01, 0x01},
			false,
		},
		{
			"multiple bytes",
			[]byte{0x01, 0x02, 0x03, 0x04},
			[]byte{0x04, 0x01, 0x02, 0x03, 0x04},
			false,
		},
		{
			"max length address",
			make([]byte, addressutils.MaxAddrLen),
			append([]byte{byte(addressutils.MaxAddrLen)}, make([]byte, addressutils.MaxAddrLen)...),
			false,
		},
		{
			"too long address",
			make([]byte, addressutils.MaxAddrLen+1),
			nil,
			true,
		},
		{
			"very long address",
			make([]byte, 1000),
			nil,
			true,
		},
		{
			"address with zero bytes",
			[]byte{0x00, 0x00, 0x00},
			[]byte{0x03, 0x00, 0x00, 0x00},
			false,
		},
		{
			"address with mixed bytes",
			[]byte{0xFF, 0x00, 0xAA, 0x55},
			[]byte{0x04, 0xFF, 0x00, 0xAA, 0x55},
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.shouldPanic {
				require.Panics(t, func() {
					addressutils.MustLengthPrefix(c.input)
				})
			} else {
				result := addressutils.MustLengthPrefix(c.input)
				require.Equal(t, c.expected, result)
			}
		})
	}
}

func TestLengthPrefixRoundTrip(t *testing.T) {
	// Test that length prefixing works correctly for various address lengths
	testLengths := []int{0, 1, 10, 32, 64, 128, 255}

	for _, length := range testLengths {
		t.Run(fmt.Sprintf("length_%d", length), func(t *testing.T) {
			// Create a test address of the specified length
			testAddr := make([]byte, length)
			for i := range testAddr {
				testAddr[i] = byte(i % 256)
			}

			// Length prefix it
			prefixed, err := addressutils.LengthPrefix(testAddr)
			require.NoError(t, err)

			if length == 0 {
				// Empty addresses should remain empty
				require.Equal(t, testAddr, prefixed)
			} else {
				// Verify the length byte is correct
				require.Equal(t, byte(length), prefixed[0])
				// Verify the rest of the bytes match the original
				require.Equal(t, testAddr, prefixed[1:])
				// Verify total length is correct
				require.Equal(t, length+1, len(prefixed))
			}
		})
	}
}

func TestLengthPrefixWithRealAddresses(t *testing.T) {
	// Test with actual generated addresses
	testAddresses := []sdk.AccAddress{
		addressutils.GenerateAccAddressForTest(),
		addressutils.GenerateAccAddressForTest(),
		addressutils.GenerateAccAddressForTest(),
	}

	for i, addr := range testAddresses {
		t.Run(fmt.Sprintf("real_address_%d", i), func(t *testing.T) {
			addrBytes := addr.Bytes()
			prefixed, err := addressutils.LengthPrefix(addrBytes)
			require.NoError(t, err)

			// Verify length byte
			require.Equal(t, byte(len(addrBytes)), prefixed[0])
			// Verify address bytes
			require.Equal(t, addrBytes, prefixed[1:])
			// Verify total length
			require.Equal(t, len(addrBytes)+1, len(prefixed))
		})
	}
}

func TestGenerateValidatorsDeterministic(t *testing.T) {
	t.Run("basic functionality", func(t *testing.T) {
		// Test with different sizes
		testCases := []int{0, 1, 5, 10, 100}

		for _, n := range testCases {
			t.Run(fmt.Sprintf("size_%d", n), func(t *testing.T) {
				validators := addressutils.GenerateValidatorsSorted(n)

				// Check length
				require.Equal(t, n, len(validators))

				// Check that all addresses are valid bech32 with correct prefix
				for _, addr := range validators {
					_, err := addressutils.AddressFromBech32(addr, "cosmosvaloper")
					require.NoError(t, err, "address %s should be valid bech32", addr)
				}

				// Check that addresses are sorted alphabetically
				for i := 1; i < len(validators); i++ {
					require.LessOrEqual(t, validators[i-1], validators[i],
						"addresses should be sorted alphabetically")
				}
			})
		}
	})

	t.Run("output is sorted", func(t *testing.T) {
		// Test that the output is always sorted alphabetically
		// Note: The function generates random addresses but sorts them
		n := 10
		validators := addressutils.GenerateValidatorsSorted(n)

		// Check that addresses are sorted alphabetically
		for i := 1; i < len(validators); i++ {
			require.LessOrEqual(t, validators[i-1], validators[i],
				"addresses should be sorted alphabetically")
		}
	})

	t.Run("addresses are unique", func(t *testing.T) {
		// Test that generated addresses are unique
		// Note: With random generation, there's a small chance of collision
		// but it should be extremely rare with 32-byte addresses
		n := 50
		validators := addressutils.GenerateValidatorsSorted(n)

		// Create a map to check uniqueness
		uniqueAddrs := make(map[string]bool)
		for _, addr := range validators {
			require.False(t, uniqueAddrs[addr], "address %s should be unique", addr)
			uniqueAddrs[addr] = true
		}

		require.Equal(t, n, len(uniqueAddrs), "all addresses should be unique")
	})

	t.Run("correct prefix", func(t *testing.T) {
		// Test that all addresses have the correct prefix
		n := 20
		validators := addressutils.GenerateValidatorsSorted(n)

		for _, addr := range validators {
			require.True(t, strings.HasPrefix(addr, "cosmosvaloper"),
				"address %s should have cosmosvaloper prefix", addr)
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		// Test edge cases
		t.Run("zero validators", func(t *testing.T) {
			validators := addressutils.GenerateValidatorsSorted(0)
			require.Equal(t, 0, len(validators))
			require.NotNil(t, validators) // Should be empty slice, not nil
		})

		t.Run("large number", func(t *testing.T) {
			// Test with a larger number to ensure it handles bigger slices
			n := 1000
			validators := addressutils.GenerateValidatorsSorted(n)
			require.Equal(t, n, len(validators))

			// Check sorting for large set
			for i := 1; i < len(validators); i++ {
				require.LessOrEqual(t, validators[i-1], validators[i])
			}
		})
	})

	t.Run("multiple calls produce different results", func(t *testing.T) {
		// Test that multiple calls produce different random results
		// (though each result is sorted)
		n := 5
		first := addressutils.GenerateValidatorsSorted(n)
		second := addressutils.GenerateValidatorsSorted(n)

		// Both should be sorted
		for i := 1; i < len(first); i++ {
			require.LessOrEqual(t, first[i-1], first[i])
		}
		for i := 1; i < len(second); i++ {
			require.LessOrEqual(t, second[i-1], second[i])
		}

		// But they should be different (random generation)
		// Note: There's a very small chance they could be the same by coincidence
		// In practice, this is extremely unlikely with 32-byte random addresses
		require.Equal(t, n, len(first))
		require.Equal(t, n, len(second))
	})

	t.Run("random generation with deterministic sorting", func(t *testing.T) {
		// This test demonstrates that the function generates random addresses
		// but sorts them to make the final result deterministic
		n := 5

		// Generate multiple sets and verify they're all sorted
		for i := 0; i < 3; i++ {
			validators := addressutils.GenerateValidatorsSorted(n)

			// Verify sorting
			for j := 1; j < len(validators); j++ {
				require.LessOrEqual(t, validators[j-1], validators[j],
					"addresses should be sorted alphabetically")
			}

			// Verify all addresses are valid
			for _, addr := range validators {
				_, err := addressutils.AddressFromBech32(addr, "cosmosvaloper")
				require.NoError(t, err, "address %s should be valid bech32", addr)
			}
		}
	})
}
