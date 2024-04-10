package types

import (
	fmt "fmt"
	"strconv"
	"strings"
)

func (e Event) ResolveAllFieldValues(fvs []*FieldValue) (bool, error) {
	for _, fv := range fvs {
		res, err := e.resolveFieldValue(fv)
		if err != nil {
			return false, err
		}
		if !res {
			return false, nil
		}
	}
	// return true if none of the conditions failed.
	return true, nil
}

func (e Event) ResolveAnyFieldValues(fvs []*FieldValue) (bool, error) {
	for _, fv := range fvs {
		res, err := e.resolveFieldValue(fv)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	// return false if none of the conditons passed.
	return false, nil
}

func (e Event) resolveFieldValue(fv *FieldValue) (bool, error) {

	switch {
	case fv.Field == FieldEventType:
		if fv.Operator != FieldOperator_EQUAL {
			return false, fmt.Errorf("bad operator %d for field %s", fv.Operator, fv.Field)
		}
		v, err := strconv.ParseInt(fv.Value, 10, 32)
		if err != nil {
			return fv.Negate, err
		}
		if v == int64(e.EventType) {
			return !fv.Negate, nil
		}
		return fv.Negate, nil
	case fv.Field == FieldEventStatus:
		if fv.Operator != FieldOperator_EQUAL {
			return false, fmt.Errorf("bad operator %d for field %s", fv.Operator, fv.Field)
		}
		v, err := strconv.ParseInt(fv.Value, 10, 32)
		if err != nil {
			return fv.Negate, err
		}
		if v == int64(e.EventType) {
			return !fv.Negate, nil
		}
		return fv.Negate, nil
	case fv.Field == FieldModule:
		res, err := compare(fv.Operator, fv.Value, e.Module)
		if err != nil {
			return false, nil
		}
		return res != fv.Negate, nil
	case fv.Field == FieldIdentifier:
		res, err := compare(fv.Operator, fv.Value, e.Identifier)
		if err != nil {
			return false, nil
		}
		return res != fv.Negate, nil
	case fv.Field == FieldChainID:
		res, err := compare(fv.Operator, fv.Value, e.ChainId)
		if err != nil {
			return false, nil
		}
		return res != fv.Negate, nil
	case fv.Field == FieldCallback:
		res, err := compare(fv.Operator, fv.Value, e.ChainId)
		if err != nil {
			return false, nil
		}
		return res != fv.Negate, nil
	}

	return fv.Negate, nil
}

func compare(operator FieldOperator, testValue, value string) (bool, error) {
	switch operator {
	case FieldOperator_EQUAL:
		return testValue == value, nil
	case FieldOperator_CONTAINS:
		return strings.Contains(value, testValue), nil
	case FieldOperator_BEGINSWITH:
		return strings.HasPrefix(value, testValue), nil
	case FieldOperator_ENDSWITH:
		return strings.HasSuffix(value, testValue), nil
	default:
		return false, fmt.Errorf("unrecognised operator %d", operator)
	}
}

func NewFieldValues(fields ...*FieldValue) []*FieldValue {
	return fields
}

func FieldEqual(field, value string) *FieldValue {
	return NewFieldValue(field, value, FieldOperator_EQUAL, false)
}

func FieldNotEqual(field, value string) *FieldValue {
	return NewFieldValue(field, value, FieldOperator_EQUAL, true)
}

func FieldBegins(field, value string) *FieldValue {
	return NewFieldValue(field, value, FieldOperator_BEGINSWITH, false)
}

func FieldEnds(field, value string) *FieldValue {
	return NewFieldValue(field, value, FieldOperator_ENDSWITH, false)
}

func NewFieldValue(field, value string, operator FieldOperator, negate bool) *FieldValue {
	return &FieldValue{
		field,
		value,
		operator,
		negate,
	}
}
