package types

import (
	"strconv"
)

func (e Event) ResolveAllFieldValues(fvs []*FieldValue) bool {
	for _, fv := range fvs {
		if !e.resolveFieldValue(fv) {
			return false
		}
	}
	return true
}

func (e Event) ResolveAnyFieldValues(fvs []*FieldValue) bool {
	for _, fv := range fvs {
		if e.resolveFieldValue(fv) {
			return true
		}
	}
	return false
}

func (e Event) resolveFieldValue(fv *FieldValue) bool {

	if fv.Field == FieldEventType {
		v, err := strconv.ParseInt(fv.Value, 10, 32)
		if err != nil {
			return fv.Negate
		}
		if v == int64(e.EventType) {
			return !fv.Negate
		}
		return fv.Negate
	}

	if fv.Field == FieldEventStatus {
		v, err := strconv.ParseInt(fv.Value, 10, 32)
		if err != nil {
			return fv.Negate
		}
		if v == int64(e.EventType) {
			return !fv.Negate
		}
		return fv.Negate
	}

	if fv.Field == FieldModule && fv.Value == e.Module {
		return !fv.Negate
	}
	if fv.Field == FieldIdentifier && fv.Value == e.Identifier {
		return !fv.Negate
	}
	if fv.Field == FieldChainID && fv.Value == e.ChainId {
		return !fv.Negate
	}

	return fv.Negate
}
