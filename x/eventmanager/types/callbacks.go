package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EventCallbacks interface {
	AddCallback(id string, fn interface{}) EventCallbacks
	RegisterCallbacks() EventCallbacks
	Call(ctx sdk.Context, id string, args []byte) error
	Has(id string) bool
}
