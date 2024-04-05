package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

type ConditionI interface {
	Resolve(ctx sdk.Context, k *Keeper) bool
	Marshal(ctx sdk.Context) []byte
}

// type ConditionExistsAny struct {
// 	Fields []types.KV
// 	negate bool
// }

// func (c ConditionExistsAny) Resolve(ctx sdk.Context, k *Keeper) bool {
// 	out := false
// 	k.IterateEvents(ctx, func(index int64, event types.Event) (stop bool) {
// 		if event.ResolveAnyFieldValues(c.Fields) {
// 			out :=
// 		}
// 	}
// 	return negate ^ out
// }

type ConditionExistsAll struct {
	Fields []types.FieldValue
	Negate bool
}

func (c ConditionExistsAll) Resolve(ctx sdk.Context, k *Keeper) bool {
	out := false
	k.IteratePrefixedEvents(ctx, nil, func(index int64, event types.Event) (stop bool) {
		if event.ResolveAllFieldValues(c.Fields) {
			out = true
			return true
		}
		return false
	})
	return c.Negate != out
}

type ConditionAnd struct {
	Condition1 ConditionI
	Condition2 ConditionI
}

func (c ConditionAnd) Resolve(ctx sdk.Context, k *Keeper) bool {
	return c.Condition1.Resolve(ctx, k) && c.Condition2.Resolve(ctx, k)
}

type ConditionOr struct {
	Condition1 ConditionI
	Condition2 ConditionI
}

func (c ConditionOr) Resolve(ctx sdk.Context, k *Keeper) bool {
	return c.Condition1.Resolve(ctx, k) || c.Condition2.Resolve(ctx, k)
}
