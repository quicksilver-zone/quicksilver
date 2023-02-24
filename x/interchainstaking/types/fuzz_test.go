package types_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func FuzzParseMsgMemo(f *testing.F) {
	// seed corpus
	tests := []struct {
		memo    string
		msgType string
	}{
		{
			memo:    types.MsgTypeRebalance + "/" + "10",
			msgType: types.MsgTypeRebalance,
		},
		{
			memo:    types.MsgTypeWithdrawal + "/" + "10",
			msgType: types.MsgTypeWithdrawal,
		},
		{
			memo:    "invalid" + "/" + "10",
			msgType: types.MsgTypeWithdrawal,
		},
		{
			memo:    types.MsgTypeWithdrawal + "/" + "A",
			msgType: types.MsgTypeWithdrawal,
		},
	}
	for _, tt := range tests {
		f.Add(tt.memo, tt.msgType)
	}

	f.Fuzz(func(t *testing.T, memo, msgType string) {
		epochNumber, err := types.ParseMsgMemo(memo, msgType)
		if err != nil {
			return
		}

		newMemo := msgType + "/" + strconv.FormatInt(epochNumber, 10)
		epochNumber2, err := types.ParseMsgMemo(newMemo, msgType)
		require.NoError(t, err)
		require.Equal(t, epochNumber, epochNumber2)
	})
}
