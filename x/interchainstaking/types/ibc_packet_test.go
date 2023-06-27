package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestParseMsgMemo(t *testing.T) {
	tests := []struct {
		name                string
		memo                string
		msgType             string
		wantErr             bool
		expectedEpochNumber int64
	}{
		{
			name:                "valid rebalance",
			memo:                types.EpochRebalanceMemo(10),
			msgType:             types.MsgTypeRebalance,
			wantErr:             false,
			expectedEpochNumber: 10,
		},
		{
			name:                "valid withdrawal",
			memo:                types.EpochWithdrawalMemo(10),
			msgType:             types.MsgTypeWithdrawal,
			wantErr:             false,
			expectedEpochNumber: 10,
		},
		{
			name:    "invalid msg type",
			memo:    "invalid" + "/" + "10",
			msgType: types.MsgTypeWithdrawal,
			wantErr: true,
		},
		{
			name:    "invalid epoch number",
			memo:    types.MsgTypeWithdrawal + "/" + "A",
			msgType: types.MsgTypeWithdrawal,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			epochNumber, err := types.ParseEpochMsgMemo(tt.memo, tt.msgType)
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedEpochNumber, epochNumber)
		})
	}
}
