package addressutils_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
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
