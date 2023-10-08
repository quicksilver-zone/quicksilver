package keeper_test

import (
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestHandleUpdateZoneProposal() {
	tests := []struct {
		name      string
		changes   []*icstypes.UpdateZoneValue
		expectErr string
	}{
		{
			name: "valid - all changes except connection",
			expectErr: "",
		},
		{
			name: "valid - connection",
			expectErr: "",
		},
		{
			name: "invalid zone",
			expectErr: "unable to get registered zone for chain id",
		},
		{
			name: "invalid change key",
			expectErr: "unexpected key",
		},
		{
			name: "invalid - base_denom not valid",
			expectErr: "invalid denom",
		},
		{
			name: "invalid - zone has assets minted",
			expectErr: "zone has assets minted",
		},
		{
			name: "invalid - parse bool",
			expectErr: "ParseBool",
		},
		{
			name: "invalid - atoi",
			expectErr: "parsing",
		},
		{
			name: "invalid - messages_per_tx",
			expectErr: "invalid value for messages_per_tx",
		},
		{
			name: "invalid - connection format",
			expectErr: "unexpected connection format",
		},
		{
			name: "invalid - zone intialised",
			expectErr: "zone already intialised, cannot update connection_id",
		},
		{
			name: "invalid - unable to fetch",
			expectErr: "unable to fetch",
		},
		{
			name: "invalid - unmarshaling",
			expectErr: "error unmarshaling client state",
		},
	}
}
