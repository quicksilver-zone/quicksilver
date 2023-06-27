package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	MsgTypeWithdrawal  = "withdrawal"
	MsgTypeRebalance   = "rebalance"
	MsgTypeUnbondSend  = "unbondSend"
	MsgTypePerformance = "perf"
	// TransferPort is the portID for ibc transfer module.
	TransferPort = "transfer"
)

var (
	ErrUnexpectedEpochMsgMemo = errors.New("unexpected epoch memo format")
	ErrUnexpectedTxMsgMemo    = errors.New("unexpected tx memo format")
)

func ParseEpochMsgMemo(memo, msgType string) (epochNumber int64, err error) {
	parts := strings.Split(memo, "/")
	if len(parts) != 2 || parts[0] != msgType {
		return 0, fmt.Errorf("msg type %s: %w", msgType, ErrUnexpectedEpochMsgMemo)
	}

	epochNumber, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("msg type %s: %w", msgType, err)
	}

	return epochNumber, err
}

func ParseTxMsgMemo(memo, msgType string) (txHash string, err error) {
	parts := strings.Split(memo, "/")
	if len(parts) != 2 || parts[0] != msgType {
		return "", fmt.Errorf("msg type %s: %w", msgType, ErrUnexpectedTxMsgMemo)
	}

	return parts[1], err
}

func EpochMsgMemo(msgType string, epoch int64) string {
	return fmt.Sprintf("%s/%d", msgType, epoch)
}

func EpochRebalanceMemo(epoch int64) string {
	return EpochMsgMemo(MsgTypeRebalance, epoch)
}

func EpochWithdrawalMemo(epoch int64) string {
	return EpochMsgMemo(MsgTypeWithdrawal, epoch)
}

func TxUnbondSendMemo(hash string) string {
	return fmt.Sprintf("%s/%s", MsgTypeUnbondSend, hash)
}
