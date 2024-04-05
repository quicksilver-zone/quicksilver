package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

type EMKeeper interface {
	IteratePrefixedEvents(sdk.Context, []byte, func(int64, Event) (stop bool))
	GetCodec() codec.Codec
}

type ConditionI interface {
	proto.Message
	Resolve(ctx sdk.Context, k EMKeeper) (bool, error)
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

func (c ConditionAll) Resolve(ctx sdk.Context, k EMKeeper) (bool, error) {
	out := false
	var err error
	k.IteratePrefixedEvents(ctx, nil, func(index int64, event Event) (stop bool) {
		res, err := event.ResolveAllFieldValues(c.Fields)
		if err != nil {
			return true
		}
		if res {
			out = true
			return true
		}
		return false
	})
	return c.Negate != out, err
}

func NewConditionAll(ctx sdk.Context, fields []*FieldValue, negate bool) (*ConditionAll, error) {
	return &ConditionAll{fields, negate}, nil
}

func (c ConditionAnd) Resolve(ctx sdk.Context, k EMKeeper) (bool, error) {
	var condition1 ConditionI
	var condition2 ConditionI
	err := k.GetCodec().UnpackAny(c.Condition1, &condition1)
	if err != nil {
		return false, err
	}
	err = k.GetCodec().UnpackAny(c.Condition2, &condition2)
	if err != nil {
		return false, err
	}

	res1, err := condition1.Resolve(ctx, k)
	if err != nil {
		return false, err
	}
	res2, err := condition2.Resolve(ctx, k)
	if err != nil {
		return false, err
	}
	return res1 && res2, nil
}

func NewConditionAnd(ctx sdk.Context, condition1, condition2 ConditionI) (*ConditionAnd, error) {
	anyc1, err := codectypes.NewAnyWithValue(condition1)
	if err != nil {
		return nil, err
	}
	anyc2, err := codectypes.NewAnyWithValue(condition2)
	if err != nil {
		return nil, err
	}
	return &ConditionAnd{anyc1, anyc2}, nil
}

func (c ConditionOr) Resolve(ctx sdk.Context, k EMKeeper) (bool, error) {
	var condition1 ConditionI
	var condition2 ConditionI
	err := k.GetCodec().UnpackAny(c.Condition1, &condition1)
	if err != nil {
		return false, err
	}
	err = k.GetCodec().UnpackAny(c.Condition2, &condition2)
	if err != nil {
		return false, err
	}

	res1, err := condition1.Resolve(ctx, k)
	if err != nil {
		return false, err
	}
	res2, err := condition2.Resolve(ctx, k)
	if err != nil {
		return false, err
	}
	return res1 || res2, nil
}

func NewConditionOr(ctx sdk.Context, condition1, condition2 ConditionI) (*ConditionOr, error) {
	anyc1, err := codectypes.NewAnyWithValue(condition1)
	if err != nil {
		return nil, err
	}
	anyc2, err := codectypes.NewAnyWithValue(condition2)
	if err != nil {
		return nil, err
	}
	return &ConditionOr{anyc1, anyc2}, nil
}

// CanExecute determines whether a
func (e *Event) CanExecute(ctx sdk.Context, k EMKeeper) (bool, error) {
	fmt.Print(e.ExecuteCondition)
	if e.ExecuteCondition == nil {
		return true, nil
	}
	var condition ConditionI
	err := k.GetCodec().UnpackAny(e.ExecuteCondition, &condition)
	if err != nil {
		return false, err
	}
	fmt.Print(condition)

	return condition.Resolve(ctx, k)
}
