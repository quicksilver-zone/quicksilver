package types_test

import (
	"bytes"
	"compress/zlib"
	"testing"

	"github.com/quicksilver-zone/quicksilver/x/airdrop/types"
	"github.com/stretchr/testify/require"
)

// TODO test

func TestDecompress(t *testing.T) {
	testString := "hello, world\n"
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte(testString))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	tests := []struct {
		name     string
		data     []byte
		expected []byte
		wantErr  bool
	}{
		{
			"no data",
			nil,
			nil,
			true,
		},
		{
			"no data",
			[]byte{0, 0, 0},
			nil,
			true,
		},

		{
			"valid data",
			b.Bytes(),
			[]byte(testString),
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := types.Decompress(tc.data)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, got)
		})
	}
}
