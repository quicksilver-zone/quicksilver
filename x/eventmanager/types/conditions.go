package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

type EMKeeper interface {
	IteratePrefixedEvents(sdk.Context, []byte, func(int64, Event) (stop bool))
}

type ConditionI interface {
	proto.Message
	Resolve(ctx sdk.Context, k EMKeeper) bool
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

func (c ConditionAll) Resolve(ctx sdk.Context, k EMKeeper) bool {
	out := false
	k.IteratePrefixedEvents(ctx, nil, func(index int64, event Event) (stop bool) {
		if event.ResolveAllFieldValues(c.Fields) {
			out = true
			return true
		}
		return false
	})
	return c.Negate != out
}

func (c ConditionAnd) Resolve(ctx sdk.Context, k EMKeeper) bool {
	var condition1 ConditionI
	var condition2 ConditionI
	_ = ModuleCdc.UnpackAny(c.Condition1, &condition1)
	_ = ModuleCdc.UnpackAny(c.Condition1, &condition1)

	return condition1.Resolve(ctx, k) && condition2.Resolve(ctx, k)
}

func (c ConditionOr) Resolve(ctx sdk.Context, k EMKeeper) bool {
	var condition1 ConditionI
	var condition2 ConditionI
	_ = ModuleCdc.UnpackAny(c.Condition1, &condition1)
	_ = ModuleCdc.UnpackAny(c.Condition1, &condition1)

	return condition1.Resolve(ctx, k) || condition2.Resolve(ctx, k)
}
