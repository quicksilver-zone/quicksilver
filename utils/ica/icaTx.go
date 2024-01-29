package ica

import (
	"fmt"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/keeper"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IcaTxKeeper struct {
	Txs []icaTx
}

func (i *IcaTxKeeper) Append(tx icaTx) {
	i.Txs = append(i.Txs, tx)
	fmt.Println("append tx")
}

func (i *IcaTxKeeper) Dump() {
	fmt.Println(i.Txs)
}

type icaTx struct {
	Msgs    []sdk.Msg
	Memo    string
	Account *types.ICAAccount
}

func GetTestSubmitTxFn(txk *IcaTxKeeper) keeper.TxSubmitFn {

	return func(ctx sdk.Context, k *keeper.Keeper, msgs []sdk.Msg, account *types.ICAAccount, memo string, messagesPerTx int64) error {
		var newTx icaTx
		newTx.Msgs = msgs
		newTx.Account = account
		newTx.Memo = memo
		txk.Append(newTx)
		return nil
	}
}
