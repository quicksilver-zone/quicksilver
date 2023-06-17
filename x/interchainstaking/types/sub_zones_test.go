package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	testZoneID    = "test-zone-1"
	testSubzoneID = "test-subzone-1"
)

var (
	zone = types.Zone{
		ConnectionId: testZoneID,
		ChainId:      testZoneID,
		SubzoneInfo:  nil,
	}

	subzone = types.Zone{
		ConnectionId: testSubzoneID,
		ChainId:      testSubzoneID,
		SubzoneInfo: &types.SubzoneInfo{
			Authority:   "test",
			BaseChainID: testZoneID,
		},
	}
)

func TestZone_IsSubzone(t *testing.T) {
	require.False(t, zone.IsSubzone())
	require.True(t, subzone.IsSubzone())
}

func TestZone_ChainID(t *testing.T) {
	require.Equal(t, testZoneID, zone.ChainID())
	require.Equal(t, testZoneID, subzone.ChainID())
}

func TestZone_ID(t *testing.T) {
	require.Equal(t, testZoneID, zone.ID())
	require.Equal(t, testSubzoneID, subzone.ID())
}
