package types_test

import (
	"strconv"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestAttributesToMap(t *testing.T) {
	tests := []struct {
		name   string
		events []abcitypes.EventAttribute
		want   map[string]string
	}{
		{
			name: "parse valid",
			events: []abcitypes.EventAttribute{
				{
					Key:   "sender",
					Value: "sender",
					Index: false,
				},
				{
					Key:   "recipient",
					Value: "recipient",
					Index: false,
				},
				{
					Key:   "amount",
					Value: strconv.Itoa(100),
					Index: false,
				},
			},
			want: map[string]string{
				"sender":    "sender",
				"recipient": "recipient",
				"amount":    strconv.Itoa(100),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := types.AttributesToMap(tc.events)
			require.Equal(t, tc.want, actual)
		})
	}
}
