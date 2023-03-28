package types

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MsgTypeWithdrawal = "withdrawal"
	MsgTypeRebalance  = "rebalance"
	// TransferPort is the portID for ibc transfer module.
	TransferPort = "transfer"
)

func ParseMsgMemo(memo, msgType string) (epochNumber int64, err error) {
	parts := strings.Split(memo, "/")
	if len(parts) != 2 || parts[0] != msgType {
		return 0, fmt.Errorf("unexpected epoch %s memo format", msgType)
	}

	epochNumber, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unexpected epoch %s memo format: %w", msgType, err)
	}

	return epochNumber, err
}
