package ica

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/keeper"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

type TxKeeper struct {
	Txs []icaTx
}

func (i *TxKeeper) Append(tx icaTx) {
	i.Txs = append(i.Txs, tx)
	fmt.Println("append tx")
}

func (i *TxKeeper) Dump() {
	fmt.Println(i.Txs)
}

type icaTx struct {
	Msgs    []sdk.Msg
	Memo    string
	Account *types.ICAAccount
}

func GetTestSubmitTxFn(txk *TxKeeper) keeper.TxSubmitFn {
	return func(ctx sdk.Context, k *keeper.Keeper, msgs []sdk.Msg, account *types.ICAAccount, memo string, messagesPerTx int64) error {
		var newTx icaTx
		newTx.Msgs = msgs
		newTx.Account = account
		newTx.Memo = memo
		txk.Append(newTx)
		return nil
	}
}
