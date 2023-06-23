package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	testZoneID        = "test-zone-1"
	testSubzoneID     = "test-zone-1|1"
	invalidSubZoneID1 = "test-subzone-1"
	invalidSubZoneID2 = "test-subzone-1|1"
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

func TestValidateSubzoneID(t *testing.T) {
	require.NoError(t, types.ValidateSubzoneID(testSubzoneID, testZoneID))
	require.Error(t, types.ValidateSubzoneID(invalidSubZoneID1, testZoneID))
	require.Error(t, types.ValidateSubzoneID(invalidSubZoneID2, testZoneID))
}

func TestValidateSubzoneForBasezone(t *testing.T) {
	baseZone := types.Zone{
		ConnectionId:    "test-connection",
		ChainId:         "test-chain",
		AccountPrefix:   "testprefix",
		LocalDenom:      "qdenom",
		BaseDenom:       "denom",
		MultiSend:       true,
		LiquidityModule: false,
		Decimals:        18,
		Is_118:          true,
		SubzoneInfo:     nil,
	}

	// set up identical zone with subzone info populated
	invalidBaseZone := baseZone
	invalidBaseZone.SubzoneInfo = &types.SubzoneInfo{
		Authority:   "test",
		BaseChainID: "test",
		ChainID:     "test",
	}

	testCases := []struct {
		name     string
		subzone  types.Zone
		baseZone types.Zone
		valid    bool
	}{
		{
			name: "valid",
			subzone: types.Zone{
				ConnectionId:    "test-connection",
				AccountPrefix:   "testprefix",
				LocalDenom:      "subqdenom",
				BaseDenom:       "denom",
				MultiSend:       true,
				LiquidityModule: false,
				Decimals:        18,
				Is_118:          true,
				SubzoneInfo: &types.SubzoneInfo{
					Authority:   "testauth",
					BaseChainID: "test-chain",
					ChainID:     "test-chain|1234",
				},
			},
			baseZone: baseZone,
			valid:    true,
		},
		{
			name: "invalid subzoneID",
			subzone: types.Zone{
				ConnectionId:    "test-connection",
				AccountPrefix:   "testprefix",
				LocalDenom:      "subqdenom",
				BaseDenom:       "denom",
				MultiSend:       true,
				LiquidityModule: false,
				Decimals:        18,
				Is_118:          true,
				SubzoneInfo: &types.SubzoneInfo{
					Authority:   "testauth",
					BaseChainID: "test-chain",
					ChainID:     "test-chain-subzone", // not using delimiter
				},
			},
			baseZone: baseZone,
			valid:    false,
		},
		{
			name: "invalid connection",
			subzone: types.Zone{
				ConnectionId:    "test-connection-invalid",
				AccountPrefix:   "testprefix",
				LocalDenom:      "subqdenom",
				BaseDenom:       "denom",
				MultiSend:       true,
				LiquidityModule: false,
				Decimals:        18,
				Is_118:          true,
				SubzoneInfo: &types.SubzoneInfo{
					Authority:   "testauth",
					BaseChainID: "test-chain",
					ChainID:     "test-chain|1234",
				},
			},
			baseZone: baseZone,
			valid:    false,
		},
		{
			name: "invalid account prefix",
			subzone: types.Zone{
				ConnectionId:    "test-connection",
				AccountPrefix:   "invalid",
				LocalDenom:      "subqdenom",
				BaseDenom:       "denom",
				MultiSend:       true,
				LiquidityModule: false,
				Decimals:        18,
				Is_118:          true,
				SubzoneInfo: &types.SubzoneInfo{
					Authority:   "testauth",
					BaseChainID: "test-chain",
					ChainID:     "test-chain|1234",
				},
			},
			baseZone: baseZone,
			valid:    false,
		},
		{
			name: "invalid duplicate local denom",
			subzone: types.Zone{
				ConnectionId:    "test-connection",
				AccountPrefix:   "testprefix",
				LocalDenom:      "qdenom",
				BaseDenom:       "denom",
				MultiSend:       true,
				LiquidityModule: false,
				Decimals:        18,
				Is_118:          true,
				SubzoneInfo: &types.SubzoneInfo{
					Authority:   "testauth",
					BaseChainID: "test-chain",
					ChainID:     "test-chain|1234",
				},
			},
			baseZone: baseZone,
			valid:    false,
		},
		{
			name: "invalid mismatch base denom",
			subzone: types.Zone{
				ConnectionId:    "test-connection",
				AccountPrefix:   "testprefix",
				LocalDenom:      "subqdenom",
				BaseDenom:       "invalid",
				MultiSend:       true,
				LiquidityModule: false,
				Decimals:        18,
				Is_118:          true,
				SubzoneInfo: &types.SubzoneInfo{
					Authority:   "testauth",
					BaseChainID: "test-chain",
					ChainID:     "test-chain|1234",
				},
			},
			baseZone: baseZone,
			valid:    false,
		},

		{
			name: "invalid capability fields",
			subzone: types.Zone{
				ConnectionId:    "test-connection",
				AccountPrefix:   "testprefix",
				LocalDenom:      "subqdenom",
				BaseDenom:       "denom",
				MultiSend:       false,
				LiquidityModule: true,
				Decimals:        18,
				Is_118:          false,
				SubzoneInfo: &types.SubzoneInfo{
					Authority:   "testauth",
					BaseChainID: "test-chain",
					ChainID:     "test-chain|1234",
				},
			},
			baseZone: baseZone,
			valid:    false,
		},
		{
			name: "invalid no subzone info",
			subzone: types.Zone{
				ConnectionId:    "test-connection",
				AccountPrefix:   "testprefix",
				LocalDenom:      "subqdenom",
				BaseDenom:       "denom",
				MultiSend:       true,
				LiquidityModule: false,
				Decimals:        18,
				Is_118:          true,
				SubzoneInfo:     nil,
			},
			baseZone: baseZone,
			valid:    false,
		},
		{
			name: "invalid base zone is subzone",
			subzone: types.Zone{
				ConnectionId:    "test-connection",
				AccountPrefix:   "testprefix",
				LocalDenom:      "subqdenom",
				BaseDenom:       "denom",
				MultiSend:       true,
				LiquidityModule: false,
				Decimals:        18,
				Is_118:          true,
				SubzoneInfo: &types.SubzoneInfo{
					Authority:   "testauth",
					BaseChainID: "test-chain",
					ChainID:     "test-chain|1234",
				},
			},
			baseZone: invalidBaseZone,
			valid:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.valid {
				require.NoError(t, types.ValidateSubzoneForBasezone(tc.subzone, tc.baseZone))
				return
			}
			require.Error(t, types.ValidateSubzoneForBasezone(tc.subzone, tc.baseZone))
		})
	}
}
